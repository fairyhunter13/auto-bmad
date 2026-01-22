/**
 * JSON-RPC 2.0 Client
 *
 * Implements the JSON-RPC 2.0 protocol over length-prefixed framing.
 * Frame format: [4 bytes: big-endian length][N bytes: JSON payload][1 byte: newline]
 *
 * This matches the framing protocol implemented in the Golang backend
 * (apps/core/internal/server/framing.go).
 */

import { EventEmitter } from 'events'
import { Readable, Writable } from 'stream'

/** JSON-RPC 2.0 Request */
export interface JsonRpcRequest {
  jsonrpc: '2.0'
  method: string
  params?: unknown
  id: number
}

/** JSON-RPC 2.0 Response */
export interface JsonRpcResponse {
  jsonrpc: '2.0'
  result?: unknown
  error?: JsonRpcError
  id: number | null
}

/** JSON-RPC 2.0 Error */
export interface JsonRpcError {
  code: number
  message: string
  data?: unknown
}

/** Pending request tracking */
interface PendingRequest {
  resolve: (result: unknown) => void
  reject: (error: Error) => void
  timer: NodeJS.Timeout
}

/** Configuration for RpcClient */
export interface RpcClientConfig {
  /** Default timeout for requests (ms) */
  defaultTimeout: number
}

const DEFAULT_CONFIG: RpcClientConfig = {
  defaultTimeout: 30000 // 30 seconds
}

/**
 * MessageWriter handles writing length-prefixed JSON-RPC messages.
 */
class MessageWriter {
  constructor(private stream: Writable) {}

  write(message: unknown): boolean {
    const payload = JSON.stringify(message)
    const payloadBuffer = Buffer.from(payload, 'utf8')
    const length = payloadBuffer.length

    // Create frame: [4-byte length][payload][newline]
    const header = Buffer.alloc(4)
    header.writeUInt32BE(length, 0)

    // Write header
    if (!this.stream.write(header)) {
      return false
    }

    // Write payload
    if (!this.stream.write(payloadBuffer)) {
      return false
    }

    // Write trailing newline
    return this.stream.write('\n')
  }
}

/**
 * MessageReader handles reading length-prefixed JSON-RPC messages.
 * Emits 'message' events when complete messages are received.
 */
class MessageReader extends EventEmitter {
  private buffer: Buffer = Buffer.alloc(0)

  constructor(stream: Readable) {
    super()
    stream.on('data', (chunk: Buffer) => this.onData(chunk))
    stream.on('error', (err) => this.emit('error', err))
    stream.on('close', () => this.emit('close'))
  }

  private onData(chunk: Buffer): void {
    // Append new data to buffer
    this.buffer = Buffer.concat([this.buffer, chunk])
    this.processBuffer()
  }

  private processBuffer(): void {
    // Process all complete messages in buffer
    while (this.buffer.length >= 4) {
      // Read 4-byte length prefix (big-endian)
      const length = this.buffer.readUInt32BE(0)
      const totalLength = 4 + length + 1 // header + payload + newline

      // Check if we have the complete message
      if (this.buffer.length < totalLength) {
        break
      }

      // Extract payload (skip header)
      const payload = this.buffer.subarray(4, 4 + length).toString('utf8')

      // Remove processed message from buffer (including newline)
      this.buffer = this.buffer.subarray(totalLength)

      try {
        const message = JSON.parse(payload)
        this.emit('message', message)
      } catch (err) {
        console.error('[RPC] Parse error:', err)
        this.emit('parseError', err)
      }
    }
  }
}

/**
 * RpcClient provides a Promise-based interface for JSON-RPC 2.0 calls.
 */
export class RpcClient extends EventEmitter {
  private nextId = 1
  private pending = new Map<number, PendingRequest>()
  private reader: MessageReader | null = null
  private writer: MessageWriter | null = null
  private config: RpcClientConfig
  private connected = false

  constructor(config: Partial<RpcClientConfig> = {}) {
    super()
    this.config = { ...DEFAULT_CONFIG, ...config }
  }

  /**
   * Connect the client to the backend streams.
   */
  connect(stdin: Writable, stdout: Readable): void {
    if (this.connected) {
      throw new Error('RpcClient already connected')
    }

    this.writer = new MessageWriter(stdin)
    this.reader = new MessageReader(stdout)

    this.reader.on('message', (response: JsonRpcResponse) => {
      this.handleResponse(response)
    })

    this.reader.on('error', (err) => {
      console.error('[RPC] Reader error:', err)
      this.emit('error', err)
    })

    this.reader.on('close', () => {
      console.log('[RPC] Connection closed')
      this.rejectAllPending(new Error('Connection closed'))
      this.emit('close')
    })

    this.connected = true
    this.emit('connected')
  }

  /**
   * Disconnect the client and reject all pending requests.
   */
  disconnect(): void {
    if (!this.connected) return

    this.rejectAllPending(new Error('Client disconnected'))
    this.reader = null
    this.writer = null
    this.connected = false
    this.emit('disconnected')
  }

  /**
   * Check if the client is connected.
   */
  isConnected(): boolean {
    return this.connected
  }

  /**
   * Call a JSON-RPC method and wait for the response.
   */
  async call<T = unknown>(method: string, params?: unknown, timeout?: number): Promise<T> {
    if (!this.connected || !this.writer) {
      throw new Error('RpcClient not connected')
    }

    const id = this.nextId++
    const request: JsonRpcRequest = {
      jsonrpc: '2.0',
      method,
      id
    }

    if (params !== undefined) {
      request.params = params
    }

    return new Promise<T>((resolve, reject) => {
      // Set up timeout
      const timeoutMs = timeout ?? this.config.defaultTimeout
      const timer = setTimeout(() => {
        if (this.pending.has(id)) {
          this.pending.delete(id)
          reject(new Error(`Request ${method} (id=${id}) timed out after ${timeoutMs}ms`))
        }
      }, timeoutMs)

      // Track pending request
      this.pending.set(id, {
        resolve: resolve as (result: unknown) => void,
        reject,
        timer
      })

      // Send request
      const success = this.writer!.write(request)
      if (!success) {
        // Stream buffer is full, but message is queued
        console.warn(`[RPC] Write buffer full for request ${method} (id=${id})`)
      }
    })
  }

  /**
   * Send a notification (no response expected).
   */
  notify(method: string, params?: unknown): void {
    if (!this.connected || !this.writer) {
      throw new Error('RpcClient not connected')
    }

    // Notifications don't have an id
    const notification = {
      jsonrpc: '2.0' as const,
      method,
      params
    }

    this.writer.write(notification)
  }

  private handleResponse(response: JsonRpcResponse): void {
    // Check if this is a response to a pending request
    if (response.id !== null && response.id !== undefined) {
      const pending = this.pending.get(response.id as number)
      if (pending) {
        this.pending.delete(response.id as number)
        clearTimeout(pending.timer)

        if (response.error) {
          const err = new Error(response.error.message)
          ;(err as Error & { code?: number; data?: unknown }).code = response.error.code
          ;(err as Error & { code?: number; data?: unknown }).data = response.error.data
          pending.reject(err)
        } else {
          pending.resolve(response.result)
        }
        return
      }
    }

    // If no pending request found, emit as event (server-initiated notification)
    this.emit('notification', response)
  }

  private rejectAllPending(error: Error): void {
    for (const [id, pending] of this.pending) {
      clearTimeout(pending.timer)
      pending.reject(error)
      this.pending.delete(id)
    }
  }

  /**
   * Get the number of pending requests.
   */
  getPendingCount(): number {
    return this.pending.size
  }
}
