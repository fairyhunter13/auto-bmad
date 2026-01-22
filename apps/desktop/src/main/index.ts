/**
 * Electron Main Process Entry Point
 *
 * Responsibilities:
 * - Create the main BrowserWindow with secure settings
 * - Spawn and manage the Golang backend process
 * - Bridge IPC requests between renderer and backend
 * - Handle graceful shutdown
 *
 * Security settings (NFR-S5):
 * - contextIsolation: true
 * - nodeIntegration: false
 * - sandbox: true
 */

import { app, shell, BrowserWindow, ipcMain } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'
import { BackendProcess } from './backend'
import { RpcClient } from './rpc-client'

// Global instances
let mainWindow: BrowserWindow | null = null
let backend: BackendProcess | null = null
let rpcClient: RpcClient | null = null

/**
 * Create the main application window with secure settings.
 */
function createWindow(): void {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    show: false,
    autoHideMenuBar: true,
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      // Security settings (NFR-S5)
      contextIsolation: true, // REQUIRED: Isolate renderer from preload
      nodeIntegration: false, // REQUIRED: No Node.js in renderer
      sandbox: true // RECOMMENDED: OS-level sandboxing
    }
  })

  mainWindow.on('ready-to-show', () => {
    mainWindow?.show()
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // HMR for renderer based on electron-vite cli
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

/**
 * Initialize the Golang backend process.
 */
async function initBackend(): Promise<void> {
  backend = new BackendProcess({
    maxRestarts: 5,
    baseRestartDelay: 1000,
    maxRestartDelay: 30000,
    shutdownTimeout: 5000
  })

  // Set main window for crash event forwarding
  if (mainWindow) {
    backend.setMainWindow(mainWindow)
  }

  // Listen for backend events
  backend.on('spawn', () => {
    console.log('[Main] Backend spawned successfully')
    connectRpcClient()
    notifyBackendStatus('connected')
  })

  backend.on('crash', (error: string) => {
    console.error('[Main] Backend crashed:', error)
    disconnectRpcClient()
    notifyBackendStatus('disconnected')
  })

  backend.on('stderr', (data: string) => {
    // Backend logs go to console
    console.log('[Backend]', data)
  })

  // Spawn the backend
  await backend.spawn()
}

/**
 * Connect the RPC client to the backend streams.
 */
function connectRpcClient(): void {
  if (!backend?.stdin || !backend?.stdout) {
    console.error('[Main] Cannot connect RPC client: backend streams not available')
    return
  }

  rpcClient = new RpcClient({
    defaultTimeout: 30000
  })

  rpcClient.on('error', (err) => {
    console.error('[Main] RPC client error:', err)
  })

  rpcClient.on('close', () => {
    console.log('[Main] RPC connection closed')
  })

  // Forward server-initiated notifications to renderer
  rpcClient.on('notification', (notification: { method?: string; params?: unknown }) => {
    if (notification.method && mainWindow && !mainWindow.isDestroyed()) {
      console.log('[Main] Forwarding notification:', notification.method)
      mainWindow.webContents.send(`event:${notification.method}`, notification.params)
    }
  })

  rpcClient.connect(backend.stdin, backend.stdout)
  console.log('[Main] RPC client connected')
}

/**
 * Disconnect the RPC client.
 */
function disconnectRpcClient(): void {
  if (rpcClient) {
    rpcClient.disconnect()
    rpcClient = null
  }
}

/**
 * Notify renderer of backend connection status.
 */
function notifyBackendStatus(status: 'connected' | 'disconnected'): void {
  if (mainWindow && !mainWindow.isDestroyed()) {
    mainWindow.webContents.send('backend:status', status)
  }
}

/**
 * Handle RPC calls from renderer via IPC.
 */
function setupIpcHandlers(): void {
  // Handle native folder selection dialog
  ipcMain.handle('dialog:selectFolder', async () => {
    const { dialog } = require('electron')
    const result = await dialog.showOpenDialog({
      properties: ['openDirectory'],
      title: 'Select BMAD Project Folder'
    })

    if (result.canceled || result.filePaths.length === 0) {
      return null
    }

    return result.filePaths[0]
  })

  // Handle JSON-RPC calls from renderer
  ipcMain.handle('rpc:call', async (_event, method: string, params?: unknown) => {
    if (!rpcClient?.isConnected()) {
      throw new Error('Backend not connected')
    }

    try {
      return await rpcClient.call(method, params)
    } catch (err) {
      // Re-throw with a clean error object for IPC
      if (err instanceof Error) {
        const error = new Error(err.message)
        ;(error as Error & { code?: number }).code = (
          err as Error & { code?: number }
        ).code
        throw error
      }
      throw err
    }
  })
}

/**
 * Graceful shutdown of all components.
 */
async function shutdown(): Promise<void> {
  console.log('[Main] Initiating graceful shutdown')

  // Disconnect RPC client
  disconnectRpcClient()

  // Shutdown backend
  if (backend) {
    await backend.shutdown()
    backend = null
  }

  console.log('[Main] Shutdown complete')
}

// Application lifecycle
app.whenReady().then(async () => {
  // Set app user model id for Windows
  electronApp.setAppUserModelId('com.autobmad.app')

  // Default open or close DevTools by F12 in development
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  // Setup IPC handlers before creating window
  setupIpcHandlers()

  // Create window first
  createWindow()

  // Initialize backend
  try {
    await initBackend()
  } catch (err) {
    console.error('[Main] Failed to initialize backend:', err)
    // Continue without backend - user will see error in UI
  }

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

// Quit when all windows are closed, except on macOS
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// Graceful shutdown before quit
app.on('before-quit', async (event) => {
  // Prevent default to allow async shutdown
  event.preventDefault()

  try {
    await shutdown()
  } finally {
    // Now actually quit
    app.exit(0)
  }
})
