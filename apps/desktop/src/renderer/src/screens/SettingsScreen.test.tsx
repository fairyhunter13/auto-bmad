import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { SettingsScreen } from './SettingsScreen'
import type { Settings } from '../../types/settings'

const mockSettings: Settings = {
  maxRetries: 3,
  retryDelay: 5000,
  desktopNotifications: true,
  soundEnabled: false,
  stepTimeoutDefault: 300000,
  heartbeatInterval: 60000,
  theme: 'system',
  showDebugOutput: false,
  projectProfiles: {},
  recentProjectsMax: 10
}

describe('SettingsScreen', () => {
  beforeEach(() => {
    // Mock the window.api object
    global.window.api = {
      dialog: {
        selectFolder: vi.fn()
      },
      project: {
        scan: vi.fn(),
        getRecent: vi.fn(),
        addRecent: vi.fn(),
        removeRecent: vi.fn(),
        setContext: vi.fn()
      },
      network: {
        getStatus: vi.fn()
      },
      settings: {
        get: vi.fn().mockResolvedValue(mockSettings),
        set: vi.fn().mockResolvedValue(mockSettings),
        reset: vi.fn().mockResolvedValue(mockSettings)
      },
      on: {
        backendCrash: vi.fn(),
        backendStatus: vi.fn(),
        networkStatusChanged: vi.fn()
      }
    } as any
  })

  it('renders loading state initially', () => {
    render(<SettingsScreen />)
    expect(screen.getByText('Loading settings...')).toBeInTheDocument()
  })

  it('loads and displays settings', async () => {
    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })

    // Check that settings values are displayed
    expect(screen.getByLabelText('Maximum Retries')).toHaveValue(3)
    expect(screen.getByLabelText('Retry Delay (ms)')).toHaveValue(5000)
  })

  it('calls settings.get on mount', async () => {
    const getSpy = vi.spyOn(window.api.settings, 'get')
    
    render(<SettingsScreen />)

    await waitFor(() => {
      expect(getSpy).toHaveBeenCalledTimes(1)
    })
  })

  it('updates settings when input changes', async () => {
    const setSpy = vi.spyOn(window.api.settings, 'set').mockResolvedValue({
      ...mockSettings,
      maxRetries: 5
    })

    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })

    const maxRetriesInput = screen.getByLabelText('Maximum Retries') as HTMLInputElement
    
    // Use fireEvent.change to atomically set the value to avoid intermediate states
    fireEvent.change(maxRetriesInput, { target: { value: '5' } })

    // Wait for settings.set to be called with the expected value
    await waitFor(() => {
      expect(setSpy).toHaveBeenCalledWith({ maxRetries: 5 })
    })
  })

  it('updates boolean settings when switch is toggled', async () => {
    const user = userEvent.setup()
    const setSpy = vi.spyOn(window.api.settings, 'set').mockResolvedValue({
      ...mockSettings,
      desktopNotifications: false
    })

    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })

    const notificationSwitch = screen.getByLabelText('Desktop Notifications')
    await user.click(notificationSwitch)

    await waitFor(() => {
      expect(setSpy).toHaveBeenCalledWith({ desktopNotifications: false })
    })
  })

  it('resets settings when reset button is clicked', async () => {
    const user = userEvent.setup()
    const resetSpy = vi.spyOn(window.api.settings, 'reset')

    // Mock window.confirm
    global.confirm = vi.fn().mockReturnValue(true)

    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })

    const resetButton = screen.getByText('Reset to Defaults')
    await user.click(resetButton)

    await waitFor(() => {
      expect(resetSpy).toHaveBeenCalledTimes(1)
    })
  })

  it('does not reset if user cancels confirmation', async () => {
    const user = userEvent.setup()
    const resetSpy = vi.spyOn(window.api.settings, 'reset')

    // Mock window.confirm to return false
    global.confirm = vi.fn().mockReturnValue(false)

    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })

    const resetButton = screen.getByText('Reset to Defaults')
    await user.click(resetButton)

    // Reset should not be called
    expect(resetSpy).not.toHaveBeenCalled()
  })

  it('displays error message when settings fail to load', async () => {
    // Mock get to fail
    global.window.api.settings.get = vi.fn().mockRejectedValue(new Error('Network error'))

    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Failed to load settings')).toBeInTheDocument()
    })
  })

  it('displays error message when settings fail to save', async () => {
    // Mock set to fail
    vi.spyOn(window.api.settings, 'set').mockRejectedValue(new Error('Save failed'))

    render(<SettingsScreen />)

    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })

    const maxRetriesInput = screen.getByLabelText('Maximum Retries')
    
    // Use fireEvent.change to atomically set the value to avoid intermediate states
    fireEvent.change(maxRetriesInput, { target: { value: '7' } })

    // Wait for debounced save (300ms) + time for async operation and state update
    // Field-level error is displayed via role="alert"
    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Save failed')
    }, { timeout: 1000 })
  })

  describe('Field-level error handling', () => {
    it('displays field-specific error message next to the offending field', async () => {
      // Mock set to fail with specific error
      vi.spyOn(window.api.settings, 'set').mockRejectedValue(
        new Error('Value must be between 0 and 10')
      )

      render(<SettingsScreen />)

      await waitFor(() => {
        expect(screen.getByText('Settings')).toBeInTheDocument()
      })

      const maxRetriesInput = screen.getByLabelText('Maximum Retries')
      fireEvent.change(maxRetriesInput, { target: { value: '15' } })

      // Wait for field-level error to appear
      await waitFor(() => {
        const alert = screen.getByRole('alert')
        expect(alert).toHaveTextContent('Value must be between 0 and 10')
      }, { timeout: 1000 })

      // Verify the input has aria-invalid attribute
      await waitFor(() => {
        expect(maxRetriesInput).toHaveAttribute('aria-invalid', 'true')
      }, { timeout: 1000 })
    })

    it('clears field error when field is successfully saved', async () => {
      const setSpy = vi.spyOn(window.api.settings, 'set')
      
      // First call fails
      setSpy.mockRejectedValueOnce(new Error('Invalid value'))
      // Second call succeeds
      setSpy.mockResolvedValueOnce({ ...mockSettings, maxRetries: 5 })

      render(<SettingsScreen />)

      await waitFor(() => {
        expect(screen.getByText('Settings')).toBeInTheDocument()
      })

      const maxRetriesInput = screen.getByLabelText('Maximum Retries')
      
      // First change - will fail
      fireEvent.change(maxRetriesInput, { target: { value: '100' } })

      // Wait for error to appear
      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent('Invalid value')
      }, { timeout: 1000 })

      // Second change - will succeed
      fireEvent.change(maxRetriesInput, { target: { value: '5' } })

      // Wait for error to be cleared
      await waitFor(() => {
        expect(screen.queryByRole('alert')).not.toBeInTheDocument()
      }, { timeout: 1000 })
    })

    it('displays generic error message for non-Error exceptions', async () => {
      // Mock set to throw non-Error value
      vi.spyOn(window.api.settings, 'set').mockRejectedValue('string error')

      render(<SettingsScreen />)

      await waitFor(() => {
        expect(screen.getByText('Settings')).toBeInTheDocument()
      })

      const retryDelayInput = screen.getByLabelText('Retry Delay (ms)')
      fireEvent.change(retryDelayInput, { target: { value: '-1' } })

      // Should show generic "Invalid value" message
      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent('Invalid value')
      }, { timeout: 1000 })
    })

    it('clears all field errors when reset is clicked', async () => {
      const user = userEvent.setup()
      
      // Mock set to fail to create field error
      vi.spyOn(window.api.settings, 'set').mockRejectedValue(new Error('Save failed'))
      global.confirm = vi.fn().mockReturnValue(true)

      render(<SettingsScreen />)

      await waitFor(() => {
        expect(screen.getByText('Settings')).toBeInTheDocument()
      })

      // Trigger an error
      const maxRetriesInput = screen.getByLabelText('Maximum Retries')
      fireEvent.change(maxRetriesInput, { target: { value: '999' } })

      // Wait for error to appear
      await waitFor(() => {
        expect(screen.getByRole('alert')).toBeInTheDocument()
      }, { timeout: 1000 })

      // Click reset
      const resetButton = screen.getByText('Reset to Defaults')
      await user.click(resetButton)

      // Wait for field errors to be cleared
      await waitFor(() => {
        expect(screen.queryByRole('alert')).not.toBeInTheDocument()
      })
    })
  })
})
