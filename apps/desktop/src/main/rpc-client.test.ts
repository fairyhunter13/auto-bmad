/**
 * Tests for JSON-RPC 2.0 Client
 *
 * Tests the length-prefixed framing protocol and request/response correlation.
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { PassThrough } from 'stream'
import { RpcClient, JsonRpcResponse } from './rpc-client'

describe('RpcClient', () => {
  let client: RpcClient
  let mockStdin: PassThrough
  let mockStdout: PassThrough

  beforeEach(() => {
    // Create mock streams
    mockStdin = new PassThrough()
    mockStdout = new PassThrough()
    client = new RpcClient({ defaultTimeout: 1000 })
  })

  afterEach(() => {
    if (client.isConnected()) {
      client.disconnect()
    }
    mockStdin.destroy()
    mockStdout.destroy()
  })

  describe('connect/disconnect', () => {
    it('should connect to streams', () => {
      client.connect(mockStdin, mockStdout)
      expect(client.isConnected()).toBe(true)
    })

    it('should throw if already connected', () => {
      client.connect(mockStdin, mockStdout)
      expect(() => client.connect(mockStdin, mockStdout)).toThrow('already connected')
    })

    it('should disconnect cleanly', () => {
      client.connect(mockStdin, mockStdout)
      client.disconnect()
      expect(client.isConnected()).toBe(false)
    })

    it('should emit connected event', () => {
      const handler = vi.fn()
      client.on('connected', handler)
      client.connect(mockStdin, mockStdout)
      expect(handler).toHaveBeenCalled()
    })

    it('should emit disconnected event', () => {
      const handler = vi.fn()
      client.on('disconnected', handler)
      client.connect(mockStdin, mockStdout)
      client.disconnect()
      expect(handler).toHaveBeenCalled()
    })
  })

  describe('call', () => {
    it('should throw if not connected', async () => {
      await expect(client.call('test')).rejects.toThrow('not connected')
    })

    it('should write length-prefixed message to stdin', async () => {
      client.connect(mockStdin, mockStdout)

      // Capture what's written to stdin
      const chunks: Buffer[] = []
      mockStdin.on('data', (chunk) => chunks.push(chunk))

      // Start call (will timeout, but we just want to check the write)
      const callPromise = client.call('system.ping')

      // Wait for write to complete
      await new Promise((resolve) => setImmediate(resolve))

      // Verify written data
      const written = Buffer.concat(chunks)
      expect(written.length).toBeGreaterThan(4) // At least header + some payload

      // Read length prefix
      const length = written.readUInt32BE(0)
      expect(length).toBeGreaterThan(0)

      // Read payload
      const payload = written.subarray(4, 4 + length).toString()
      const parsed = JSON.parse(payload)

      expect(parsed.jsonrpc).toBe('2.0')
      expect(parsed.method).toBe('system.ping')
      expect(parsed.id).toBe(1)

      // Verify trailing newline
      expect(written[4 + length]).toBe(10) // '\n'

      // Clean up
      callPromise.catch(() => {}) // Ignore timeout error
    })

    it('should resolve with result on success response', async () => {
      client.connect(mockStdin, mockStdout)

      // Make call
      const callPromise = client.call<string>('system.ping')

      // Wait for request to be sent
      await new Promise((resolve) => setImmediate(resolve))

      // Send response
      const response: JsonRpcResponse = {
        jsonrpc: '2.0',
        result: 'pong',
        id: 1
      }
      sendResponse(mockStdout, response)

      const result = await callPromise
      expect(result).toBe('pong')
    })

    it('should reject with error on error response', async () => {
      client.connect(mockStdin, mockStdout)

      const callPromise = client.call('invalid.method')

      await new Promise((resolve) => setImmediate(resolve))

      const response: JsonRpcResponse = {
        jsonrpc: '2.0',
        error: {
          code: -32601,
          message: 'Method not found'
        },
        id: 1
      }
      sendResponse(mockStdout, response)

      await expect(callPromise).rejects.toThrow('Method not found')
    })

    it('should reject on timeout', async () => {
      client.connect(mockStdin, mockStdout)

      // Use short timeout
      await expect(client.call('slow.method', undefined, 50)).rejects.toThrow('timed out')
    })

    it('should correlate responses by id', async () => {
      client.connect(mockStdin, mockStdout)

      // Start multiple calls
      const call1 = client.call<string>('method1')
      const call2 = client.call<string>('method2')
      const call3 = client.call<string>('method3')

      await new Promise((resolve) => setImmediate(resolve))

      // Send responses out of order
      sendResponse(mockStdout, { jsonrpc: '2.0', result: 'result3', id: 3 })
      sendResponse(mockStdout, { jsonrpc: '2.0', result: 'result1', id: 1 })
      sendResponse(mockStdout, { jsonrpc: '2.0', result: 'result2', id: 2 })

      expect(await call1).toBe('result1')
      expect(await call2).toBe('result2')
      expect(await call3).toBe('result3')
    })

    it('should include params in request', async () => {
      client.connect(mockStdin, mockStdout)

      const chunks: Buffer[] = []
      mockStdin.on('data', (chunk) => chunks.push(chunk))

      const callPromise = client.call('echo', { message: 'hello' })
      await new Promise((resolve) => setImmediate(resolve))

      const written = Buffer.concat(chunks)
      const length = written.readUInt32BE(0)
      const payload = written.subarray(4, 4 + length).toString()
      const parsed = JSON.parse(payload)

      expect(parsed.params).toEqual({ message: 'hello' })

      callPromise.catch(() => {})
    })
  })

  describe('notify', () => {
    it('should throw if not connected', () => {
      expect(() => client.notify('test')).toThrow('not connected')
    })

    it('should write notification without id', () => {
      client.connect(mockStdin, mockStdout)

      const chunks: Buffer[] = []
      mockStdin.on('data', (chunk) => chunks.push(chunk))

      client.notify('log', { level: 'info', message: 'test' })

      const written = Buffer.concat(chunks)
      const length = written.readUInt32BE(0)
      const payload = written.subarray(4, 4 + length).toString()
      const parsed = JSON.parse(payload)

      expect(parsed.jsonrpc).toBe('2.0')
      expect(parsed.method).toBe('log')
      expect(parsed.params).toEqual({ level: 'info', message: 'test' })
      expect(parsed.id).toBeUndefined()
    })
  })

  describe('connection errors', () => {
    it('should reject pending requests on disconnect', async () => {
      client.connect(mockStdin, mockStdout)

      const callPromise = client.call('slow.method', undefined, 10000)
      await new Promise((resolve) => setImmediate(resolve))

      client.disconnect()

      await expect(callPromise).rejects.toThrow('disconnected')
    })

    it('should reject pending requests on stream close', async () => {
      client.connect(mockStdin, mockStdout)

      const callPromise = client.call('slow.method', undefined, 10000)
      await new Promise((resolve) => setImmediate(resolve))

      mockStdout.destroy()

      await expect(callPromise).rejects.toThrow()
    })
  })

  describe('getPendingCount', () => {
    it('should return number of pending requests', async () => {
      client.connect(mockStdin, mockStdout)

      expect(client.getPendingCount()).toBe(0)

      const call1 = client.call('method1', undefined, 10000)
      const call2 = client.call('method2', undefined, 10000)

      await new Promise((resolve) => setImmediate(resolve))

      expect(client.getPendingCount()).toBe(2)

      // Resolve one
      sendResponse(mockStdout, { jsonrpc: '2.0', result: 'ok', id: 1 })
      await call1

      expect(client.getPendingCount()).toBe(1)

      // Cancel the other
      client.disconnect()
      call2.catch(() => {})

      expect(client.getPendingCount()).toBe(0)
    })
  })
})

/**
 * Helper to send a length-prefixed response to the mock stdout.
 */
function sendResponse(stream: PassThrough, response: JsonRpcResponse): void {
  const payload = JSON.stringify(response)
  const payloadBuffer = Buffer.from(payload, 'utf8')
  const header = Buffer.alloc(4)
  header.writeUInt32BE(payloadBuffer.length, 0)

  stream.write(header)
  stream.write(payloadBuffer)
  stream.write('\n')
}
