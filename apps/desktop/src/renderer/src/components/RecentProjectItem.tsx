/**
 * RecentProjectItem - Individual recent project list item
 * 
 * Displays project name, path, and last opened date.
 * Allows selecting or removing the project.
 */

import { Button } from '@/components/ui/button'
import type { RecentProject } from '../screens/ProjectSelectScreen'

interface RecentProjectItemProps {
  project: RecentProject
  onSelect: () => void
  onRemove: () => void
}

export function RecentProjectItem({ project, onSelect, onRemove }: RecentProjectItemProps): JSX.Element {
  // Format last opened date
  const formatDate = (isoString: string): string => {
    try {
      const date = new Date(isoString)
      const now = new Date()
      const diffMs = now.getTime() - date.getTime()
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

      if (diffDays === 0) return 'Today'
      if (diffDays === 1) return 'Yesterday'
      if (diffDays < 7) return `${diffDays} days ago`
      return date.toLocaleDateString()
    } catch {
      return 'Unknown'
    }
  }

  return (
    <li className="flex items-center justify-between gap-2 p-2 rounded-md hover:bg-accent">
      <button
        onClick={onSelect}
        className="flex-1 text-left cursor-pointer"
      >
        <div className="font-medium">{project.name}</div>
        <div className="text-xs text-muted-foreground">{project.path}</div>
        <div className="text-xs text-muted-foreground">{formatDate(project.lastOpened)}</div>
      </button>
      <Button
        variant="ghost"
        size="sm"
        onClick={(e) => {
          e.stopPropagation()
          onRemove()
        }}
        aria-label={`Remove ${project.name}`}
      >
        ×
      </Button>
    </li>
  )
}
