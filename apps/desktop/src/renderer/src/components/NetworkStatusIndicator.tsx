/**
 * NetworkStatusIndicator Component
 *
 * Displays current network connectivity status in the UI.
 * Shows warning dialog when network goes offline during active operations.
 */

import { useEffect, useState } from 'react'
import { Wifi, WifiOff, Loader2 } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

type NetworkStatus = 'online' | 'offline' | 'checking'

interface NetworkStatusResult {
  status: NetworkStatus
  lastChecked: string
  latency?: number
}

export function NetworkStatusIndicator(): JSX.Element {
  const [status, setStatus] = useState<NetworkStatus>('checking')

  useEffect(() => {
    // Get initial status
    window.api.network
      .getStatus()
      .then((result: NetworkStatusResult) => {
        setStatus(result.status)
      })
      .catch((err) => {
        console.error('Failed to get initial network status:', err)
      })

    // Listen for status changes
    const unsubscribe = window.api.on.networkStatusChanged((event) => {
      console.log('[NetworkStatus] Status changed:', event.previous, '->', event.current)
      setStatus(event.current)

      // Show notification when going offline
      if (event.current === 'offline' && event.previous === 'online') {
        // Could trigger a toast notification here
        console.warn('[NetworkStatus] Network connection lost')
      }

      // Show notification when coming back online
      if (event.current === 'online' && event.previous === 'offline') {
        console.log('[NetworkStatus] Network connection restored')
      }
    })

    return () => {
      unsubscribe()
    }
  }, [])

  // Determine badge variant based on status
  const variant = status === 'online' ? 'default' : status === 'offline' ? 'destructive' : 'secondary'

  return (
    <Badge variant={variant} className="flex items-center gap-1 cursor-default">
      {status === 'online' && (
        <>
          <Wifi className="h-3 w-3" />
          <span>Online</span>
        </>
      )}
      {status === 'offline' && (
        <>
          <WifiOff className="h-3 w-3" />
          <span>Offline</span>
        </>
      )}
      {status === 'checking' && (
        <>
          <Loader2 className="h-3 w-3 animate-spin" />
          <span>Checking</span>
        </>
      )}
    </Badge>
  )
}
