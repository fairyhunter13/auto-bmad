/**
 * SettingsScreen - User preferences and configuration
 * 
 * Allows users to configure:
 * - Retry behavior (max retries, delay)
 * - Notifications (desktop, sound)
 * - Timeouts (step timeout, heartbeat interval)
 * - UI preferences (theme, debug output)
 */

import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import type { Settings } from '../../types/settings'

export function SettingsScreen(): JSX.Element {
  const [settings, setSettings] = useState<Settings | null>(null)
  const [saving, setSaving] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Load settings on mount
  useEffect(() => {
    loadSettings()
  }, [])

  const loadSettings = async (): Promise<void> => {
    setLoading(true)
    setError(null)
    try {
      const loadedSettings = await window.api.settings.get()
      setSettings(loadedSettings as Settings)
    } catch (err) {
      console.error('Failed to load settings:', err)
      setError('Failed to load settings. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const handleChange = async <K extends keyof Settings>(
    key: K,
    value: Settings[K]
  ): Promise<void> => {
    if (!settings) return

    setSaving(true)
    setError(null)
    try {
      const updated = await window.api.settings.set({ [key]: value })
      setSettings(updated as Settings)
    } catch (err) {
      console.error('Failed to save settings:', err)
      setError('Failed to save settings. Please try again.')
    } finally {
      setSaving(false)
    }
  }

  const handleReset = async (): Promise<void> => {
    if (!confirm('Reset all settings to defaults?')) return

    setSaving(true)
    setError(null)
    try {
      const defaults = await window.api.settings.reset()
      setSettings(defaults as Settings)
    } catch (err) {
      console.error('Failed to reset settings:', err)
      setError('Failed to reset settings. Please try again.')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto p-8">
        <div className="flex items-center justify-center h-64">
          <p className="text-muted-foreground">Loading settings...</p>
        </div>
      </div>
    )
  }

  if (!settings) {
    return (
      <div className="container mx-auto p-8">
        <div className="flex flex-col items-center justify-center h-64 gap-4">
          <p className="text-destructive">Failed to load settings</p>
          <Button onClick={loadSettings}>Retry</Button>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-2">
          Configure your Auto-BMAD preferences
        </p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-destructive/10 text-destructive rounded-md">
          {error}
        </div>
      )}

      <div className="space-y-6">
        {/* Retry Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Retry Behavior</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="maxRetries">Maximum Retries</Label>
              <Input
                id="maxRetries"
                type="number"
                value={settings.maxRetries}
                onChange={(e) =>
                  handleChange('maxRetries', parseInt(e.target.value) || 0)
                }
                disabled={saving}
                min="0"
                max="10"
              />
              <p className="text-sm text-muted-foreground">
                Number of times to retry failed steps (0-10)
              </p>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="retryDelay">Retry Delay (ms)</Label>
              <Input
                id="retryDelay"
                type="number"
                value={settings.retryDelay}
                onChange={(e) =>
                  handleChange('retryDelay', parseInt(e.target.value) || 0)
                }
                disabled={saving}
                min="0"
                step="1000"
              />
              <p className="text-sm text-muted-foreground">
                Time to wait before retrying (in milliseconds)
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Notification Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Notifications</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="desktopNotifications">Desktop Notifications</Label>
                <p className="text-sm text-muted-foreground">
                  Show system notifications for important events
                </p>
              </div>
              <Switch
                id="desktopNotifications"
                checked={settings.desktopNotifications}
                onCheckedChange={(v) => handleChange('desktopNotifications', v)}
                disabled={saving}
              />
            </div>
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="soundEnabled">Sound Effects</Label>
                <p className="text-sm text-muted-foreground">
                  Play sounds for notifications
                </p>
              </div>
              <Switch
                id="soundEnabled"
                checked={settings.soundEnabled}
                onCheckedChange={(v) => handleChange('soundEnabled', v)}
                disabled={saving}
              />
            </div>
          </CardContent>
        </Card>

        {/* Timeout Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Timeouts</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="stepTimeoutDefault">Step Timeout (ms)</Label>
              <Input
                id="stepTimeoutDefault"
                type="number"
                value={settings.stepTimeoutDefault}
                onChange={(e) =>
                  handleChange('stepTimeoutDefault', parseInt(e.target.value) || 0)
                }
                disabled={saving}
                min="0"
                step="60000"
              />
              <p className="text-sm text-muted-foreground">
                Maximum time to wait for a step to complete (5 minutes = 300000)
              </p>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="heartbeatInterval">Heartbeat Interval (ms)</Label>
              <Input
                id="heartbeatInterval"
                type="number"
                value={settings.heartbeatInterval}
                onChange={(e) =>
                  handleChange('heartbeatInterval', parseInt(e.target.value) || 0)
                }
                disabled={saving}
                min="0"
                step="30000"
              />
              <p className="text-sm text-muted-foreground">
                How often to check for progress updates (60 seconds = 60000)
              </p>
            </div>
          </CardContent>
        </Card>

        {/* UI Preferences */}
        <Card>
          <CardHeader>
            <CardTitle>User Interface</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="theme">Theme</Label>
              <select
                id="theme"
                value={settings.theme}
                onChange={(e) =>
                  handleChange('theme', e.target.value as 'light' | 'dark' | 'system')
                }
                disabled={saving}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="system">System</option>
                <option value="light">Light</option>
                <option value="dark">Dark</option>
              </select>
              <p className="text-sm text-muted-foreground">
                Choose your preferred color theme
              </p>
            </div>
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="showDebugOutput">Show Debug Output</Label>
                <p className="text-sm text-muted-foreground">
                  Display detailed technical information
                </p>
              </div>
              <Switch
                id="showDebugOutput"
                checked={settings.showDebugOutput}
                onCheckedChange={(v) => handleChange('showDebugOutput', v)}
                disabled={saving}
              />
            </div>
          </CardContent>
        </Card>

        {/* Reset Button */}
        <div className="flex justify-end gap-2">
          <Button variant="destructive" onClick={handleReset} disabled={saving}>
            Reset to Defaults
          </Button>
        </div>

        {saving && (
          <div className="text-sm text-muted-foreground text-center">
            Saving...
          </div>
        )}
      </div>
    </div>
  )
}
