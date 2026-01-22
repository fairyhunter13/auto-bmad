/**
 * Tests for Backend Process Manager
 *
 * Tests process lifecycle, crash detection, and restart logic.
 * These are unit tests that mock child_process to avoid spawning real processes.
 */

import { describe, it, expect, beforeEach, vi, afterEach, Mock } from 'vitest'
import { EventEmitter } from 'events'
import { spawn } from 'child_process'
import { Writable, Readable } from 'stream'
import { BackendProcess } from './backend'

// Mock child_process
vi.mock('child_process', () => ({
  spawn: vi.fn()
}))

// Mock electron
vi.mock('electron', () => ({
  app: {
    isPackaged: false
  },
  BrowserWindow: vi.fn()
}))

/**
 * Mock ChildProcess implementation for testing.
 */
class MockChildProcess extends EventEmitter {
  killed = false
  stdin: Writable
  stdout: Readable
  stderr: Readable & { push: (chunk: unknown) => boolean }
  pid = 12345
  kill: Mock

  constructor() {
    super()
    this.stdin = new Writable({
      write(_chunk, _encoding, callback) {
        callback()
      }
    })
    this.stdout = new Readable({
      read() {}
    })
    this.stderr = new Readable({
      read() {}
    }) as Readable & { push: (chunk: unknown) => boolean }
    this.kill = vi.fn((signal?: string) => {
      if (signal === 'SIGKILL') {
        this.killed = true
      }
      return true
    })
  }
}

describe('BackendProcess', () => {
  let backend: BackendProcess
  let mockProcess: MockChildProcess

  beforeEach(() => {
    vi.clearAllMocks()

    // Create mock process
    mockProcess = new MockChildProcess()
    ;(spawn as Mock).mockReturnValue(mockProcess)

    // Create backend with short timeouts for testing
    backend = new BackendProcess({
      maxRestarts: 3,
      baseRestartDelay: 10,
      maxRestartDelay: 100,
      shutdownTimeout: 100
    })
  })

  afterEach(() => {
    // Clean up timers
    vi.useRealTimers()
  })

  describe('getBinaryPath', () => {
    it('should return development path when not packaged', () => {
      const path = backend.getBinaryPath()
      expect(path).toContain('resources/bin/autobmad-')
      expect(path).toMatch(/autobmad-(linux|darwin)-(amd64|arm64)/)
    })
  })

  describe('spawn', () => {
    it('should spawn the process with correct options', async () => {
      await backend.spawn()

      expect(spawn).toHaveBeenCalledWith(
        expect.stringContaining('autobmad-'),
        [],
        expect.objectContaining({
          stdio: ['pipe', 'pipe', 'pipe']
        })
      )
    })

    it('should emit spawn event', async () => {
      const handler = vi.fn()
      backend.on('spawn', handler)

      await backend.spawn()

      expect(handler).toHaveBeenCalled()
    })

    it('should set isRunning to true', async () => {
      expect(backend.isRunning()).toBe(false)

      await backend.spawn()

      expect(backend.isRunning()).toBe(true)
    })

    it('should not spawn if already running', async () => {
      await backend.spawn()

      // Second spawn should be a no-op
      await backend.spawn()

      expect(spawn).toHaveBeenCalledTimes(1)
    })

    it('should provide stdin and stdout streams', async () => {
      await backend.spawn()

      expect(backend.stdin).toBeDefined()
      expect(backend.stdout).toBeDefined()
    })
  })

  describe('crash detection', () => {
    it('should emit crash event on process error', async () => {
      await backend.spawn()

      const crashHandler = vi.fn()
      backend.on('crash', crashHandler)

      mockProcess.emit('error', new Error('spawn ENOENT'))

      expect(crashHandler).toHaveBeenCalledWith('spawn ENOENT')
    })

    it('should emit crash event on unexpected exit', async () => {
      await backend.spawn()

      const crashHandler = vi.fn()
      backend.on('crash', crashHandler)

      mockProcess.emit('exit', 1, null)

      expect(crashHandler).toHaveBeenCalledWith(expect.stringContaining('code=1'))
    })

    it('should emit exit event with code and signal', async () => {
      await backend.spawn()

      const exitHandler = vi.fn()
      backend.on('exit', exitHandler)

      mockProcess.emit('exit', null, 'SIGTERM')

      expect(exitHandler).toHaveBeenCalledWith(null, 'SIGTERM')
    })
  })

  describe('shutdown', () => {
    it('should send SIGTERM', async () => {
      await backend.spawn()

      const shutdownPromise = backend.shutdown()

      expect(mockProcess.kill).toHaveBeenCalledWith('SIGTERM')

      // Simulate process exit
      mockProcess.emit('exit', 0, 'SIGTERM')

      await shutdownPromise
    })

    it('should not emit crash event during shutdown', async () => {
      await backend.spawn()

      const crashHandler = vi.fn()
      backend.on('crash', crashHandler)

      const shutdownPromise = backend.shutdown()
      mockProcess.emit('exit', 0, 'SIGTERM')
      await shutdownPromise

      expect(crashHandler).not.toHaveBeenCalled()
    })

    it('should set isRunning to false after shutdown', async () => {
      await backend.spawn()
      expect(backend.isRunning()).toBe(true)

      const shutdownPromise = backend.shutdown()
      mockProcess.emit('exit', 0, null)
      await shutdownPromise

      expect(backend.isRunning()).toBe(false)
    })

    it('should handle shutdown when no process is running', async () => {
      // Should not throw
      await backend.shutdown()
    })
  })

  describe('stderr handling', () => {
    it('should emit stderr event when stderr data is received', async () => {
      await backend.spawn()

      const stderrHandler = vi.fn()
      backend.on('stderr', stderrHandler)

      // Directly emit data event on stderr
      mockProcess.stderr.emit('data', Buffer.from('test log message'))

      expect(stderrHandler).toHaveBeenCalledWith('test log message')
    })

    it('should trim stderr messages', async () => {
      await backend.spawn()

      const stderrHandler = vi.fn()
      backend.on('stderr', stderrHandler)

      mockProcess.stderr.emit('data', Buffer.from('  trimmed message  \n'))

      expect(stderrHandler).toHaveBeenCalledWith('trimmed message')
    })
  })

  describe('restart scheduling', () => {
    it('should schedule restart after crash', async () => {
      await backend.spawn()
      const initialCallCount = (spawn as Mock).mock.calls.length

      // Simulate crash
      mockProcess.emit('exit', 1, null)

      // A restart should be scheduled
      // We can verify by checking that after enough time, spawn is called again
      mockProcess = new MockChildProcess()
      ;(spawn as Mock).mockReturnValue(mockProcess)

      // Wait for restart delay
      await new Promise((resolve) => setTimeout(resolve, 50))

      // Should have called spawn at least once more after the initial spawn
      expect((spawn as Mock).mock.calls.length).toBeGreaterThan(initialCallCount)
    })

    it('should not restart during shutdown', async () => {
      await backend.spawn()

      // Start shutdown
      const shutdownPromise = backend.shutdown()

      // Emit exit during shutdown
      mockProcess.emit('exit', 0, 'SIGTERM')
      await shutdownPromise

      // Wait to make sure no restart is scheduled
      await new Promise((resolve) => setTimeout(resolve, 50))

      // Should only have spawned once
      expect(spawn).toHaveBeenCalledTimes(1)
    })
  })
})
