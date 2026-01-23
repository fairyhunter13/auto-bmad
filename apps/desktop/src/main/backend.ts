/**
 * Backend Process Manager
 *
 * Manages the lifecycle of the Golang backend process (autobmad).
 * Responsibilities:
 * - Spawn the binary with stdin/stdout pipes
 * - Detect crashes and emit events
 * - Graceful shutdown (SIGTERM -> wait -> SIGKILL)
 * - Auto-restart with exponential backoff
 */

import { spawn, ChildProcess } from 'child_process'
import { app, BrowserWindow } from 'electron'
import path from 'path'
import os from 'os'
import { EventEmitter } from 'events'
import { Writable, Readable } from 'stream'

/** Events emitted by BackendProcess */
export interface BackendEvents {
  spawn: () => void
  crash: (error: string) => void
  exit: (code: number | null, signal: string | null) => void
  stderr: (data: string) => void
}

/** Configuration for BackendProcess */
export interface BackendConfig {
  /** Maximum restart attempts before giving up */
  maxRestarts: number
  /** Base delay for exponential backoff (ms) */
  baseRestartDelay: number
  /** Maximum delay between restarts (ms) */
  maxRestartDelay: number
  /** Timeout for graceful shutdown (ms) */
  shutdownTimeout: number
}

const DEFAULT_CONFIG: BackendConfig = {
  maxRestarts: 5,
  baseRestartDelay: 1000,
  maxRestartDelay: 30000,
  shutdownTimeout: 5000
}

export class BackendProcess extends EventEmitter {
  private process: ChildProcess | null = null
  private isShuttingDown = false
  private restartCount = 0
  private restartTimer: NodeJS.Timeout | null = null
  private config: BackendConfig
  private mainWindow: BrowserWindow | null = null
  private projectPath: string | null = null

  constructor(config: Partial<BackendConfig> = {}) {
    super()
    this.config = { ...DEFAULT_CONFIG, ...config }
  }

  /** Set the main window for emitting crash events */
  setMainWindow(window: BrowserWindow | null): void {
    this.mainWindow = window
  }

  /**
   * Get the path to the platform-specific backend binary.
   * In development: uses the binary in apps/desktop/resources/bin/
   * In production: uses the binary in process.resourcesPath/bin/
   */
  getBinaryPath(): string {
    const platform = process.platform === 'darwin' ? 'darwin' : 'linux'
    const arch = os.arch() === 'arm64' ? 'arm64' : 'amd64'
    const binaryName = `autobmad-${platform}-${arch}`

    if (app.isPackaged) {
      // Production: binary in resources/bin/
      return path.join(process.resourcesPath, 'bin', binaryName)
    } else {
      // Development: binary in apps/desktop/resources/bin/
      return path.join(__dirname, '../../resources/bin', binaryName)
    }
  }

  /** Spawn the backend process */
  async spawn(projectPath: string): Promise<void> {
    if (this.process) {
      console.warn('[Backend] Process already running')
      return
    }

    const binaryPath = this.getBinaryPath()
    console.log(`[Backend] Spawning: ${binaryPath}`)
    console.log(`[Backend] Project path: ${projectPath}`)

    // Store project path for restarts
    this.projectPath = projectPath

    try {
      this.process = spawn(binaryPath, ['--project-path', projectPath], {
        stdio: ['pipe', 'pipe', 'pipe'],
        env: { ...process.env }
      })

      this.setupProcessHandlers()
      this.emit('spawn')
      console.log(`[Backend] Process spawned with PID: ${this.process.pid}`)

      // Reset restart count on successful spawn
      this.restartCount = 0
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : String(err)
      console.error('[Backend] Spawn error:', errorMsg)
      this.emitCrashEvent(`Spawn failed: ${errorMsg}`)
      throw err
    }
  }

  private setupProcessHandlers(): void {
    if (!this.process) return

    this.process.on('error', (err) => {
      console.error('[Backend] Process error:', err.message)
      this.emitCrashEvent(err.message)
      this.scheduleRestart()
    })

    this.process.on('exit', (code, signal) => {
      console.log(`[Backend] Process exited: code=${code}, signal=${signal}`)

      if (!this.isShuttingDown) {
        const errorMsg = `Process exited unexpectedly: code=${code}, signal=${signal}`
        console.error(`[Backend] ${errorMsg}`)
        this.emitCrashEvent(errorMsg)
        this.emit('exit', code, signal)
        this.process = null
        this.scheduleRestart()
      } else {
        this.emit('exit', code, signal)
        this.process = null
      }
    })

    // Capture stderr for logging
    this.process.stderr?.on('data', (data: Buffer) => {
      const msg = data.toString().trim()
      if (msg) {
        console.log('[Backend]', msg)
        this.emit('stderr', msg)
      }
    })
  }

  private emitCrashEvent(error: string): void {
    this.emit('crash', error)

    // Send to renderer process
    if (this.mainWindow && !this.mainWindow.isDestroyed()) {
      this.mainWindow.webContents.send('backend:crash', error)
    }
  }

  private scheduleRestart(): void {
    if (this.isShuttingDown) return
    if (this.restartCount >= this.config.maxRestarts) {
      console.error('[Backend] Max restarts exceeded, giving up')
      return
    }

    // Calculate delay with exponential backoff
    const delay = Math.min(
      this.config.baseRestartDelay * Math.pow(2, this.restartCount),
      this.config.maxRestartDelay
    )

    this.restartCount++
    console.log(
      `[Backend] Scheduling restart ${this.restartCount}/${this.config.maxRestarts} in ${delay}ms`
    )

    this.restartTimer = setTimeout(async () => {
      this.restartTimer = null
      try {
        if (this.projectPath) {
          await this.spawn(this.projectPath)
        } else {
          console.error('[Backend] Cannot restart: project path not set')
        }
      } catch {
        // Spawn error already handled in spawn()
      }
    }, delay)
  }

  /**
   * Gracefully shutdown the backend process.
   * Sends SIGTERM, waits up to shutdownTimeout, then sends SIGKILL.
   */
  async shutdown(): Promise<void> {
    if (!this.process) {
      console.log('[Backend] No process to shutdown')
      return
    }

    console.log('[Backend] Initiating graceful shutdown')
    this.isShuttingDown = true

    // Cancel any pending restart
    if (this.restartTimer) {
      clearTimeout(this.restartTimer)
      this.restartTimer = null
    }

    // Send SIGTERM
    this.process.kill('SIGTERM')

    // Wait for exit or timeout
    await Promise.race([
      new Promise<void>((resolve) => {
        this.process?.on('exit', () => resolve())
      }),
      new Promise<void>((resolve) =>
        setTimeout(() => {
          if (this.process && !this.process.killed) {
            console.warn('[Backend] Graceful shutdown timeout, sending SIGKILL')
            this.process.kill('SIGKILL')
          }
          resolve()
        }, this.config.shutdownTimeout)
      )
    ])

    console.log('[Backend] Shutdown complete')
    this.process = null
    this.isShuttingDown = false
  }

  /** Check if the backend process is running */
  isRunning(): boolean {
    return this.process !== null && !this.process.killed
  }

  /** Get the stdin stream for writing requests */
  get stdin(): Writable | null {
    return this.process?.stdin ?? null
  }

  /** Get the stdout stream for reading responses */
  get stdout(): Readable | null {
    return this.process?.stdout ?? null
  }
}
