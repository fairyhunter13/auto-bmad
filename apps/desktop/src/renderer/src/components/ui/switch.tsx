import * as React from 'react'
import { cn } from '@/lib/utils'

export interface SwitchProps {
  id?: string
  checked?: boolean
  onCheckedChange?: (checked: boolean) => void
  disabled?: boolean
  className?: string
}

/**
 * Simple toggle switch component (custom implementation)
 * TODO: Replace with @radix-ui/react-switch when installed
 */
const Switch = React.forwardRef<HTMLInputElement, SwitchProps>(
  ({ id, checked = false, onCheckedChange, disabled = false, className }, ref) => {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      onCheckedChange?.(e.target.checked)
    }

    return (
      <label
        htmlFor={id}
        className={cn(
          'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
          checked ? 'bg-primary' : 'bg-input',
          disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer',
          className
        )}
      >
        <input
          ref={ref}
          type="checkbox"
          id={id}
          checked={checked}
          onChange={handleChange}
          disabled={disabled}
          className="sr-only"
        />
        <span
          className={cn(
            'inline-block h-5 w-5 transform rounded-full bg-background transition-transform',
            checked ? 'translate-x-5' : 'translate-x-1'
          )}
        />
      </label>
    )
  }
)
Switch.displayName = 'Switch'

export { Switch }
