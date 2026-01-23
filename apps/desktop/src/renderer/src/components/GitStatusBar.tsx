import { useEffect, useState } from 'react'

/**
 * GitStatusBar - Displays Git repository status for a project
 *
 * Shows:
 * - Branch name when in a Git repo
 * - Warning indicator for uncommitted changes
 * - Hides gracefully when not in a Git repo
 */

interface GitStatus {
  isGitRepo: boolean
  branch?: string
  hasChanges: boolean
  error?: string
}

interface GitStatusBarProps {
  projectPath: string
}

export function GitStatusBar({ projectPath }: GitStatusBarProps): JSX.Element | null {
  const [status, setStatus] = useState<GitStatus | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false

    async function fetchStatus(): Promise<void> {
      try {
        setLoading(true)
        const result = await window.api.git.getRepoStatus(projectPath)
        if (!cancelled) {
          setStatus(result)
        }
      } catch (error) {
        if (!cancelled) {
          setStatus({ isGitRepo: false, hasChanges: false, error: String(error) })
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    if (projectPath) {
      fetchStatus()
    }

    return () => {
      cancelled = true
    }
  }, [projectPath])

  // Don't render anything while loading or if not a git repo
  if (loading || !status?.isGitRepo) {
    return null
  }

  return (
    <div
      className="flex items-center gap-2 text-sm text-muted-foreground"
      data-testid="git-status-bar"
    >
      <svg
        className="h-4 w-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
        />
      </svg>
      <span data-testid="git-branch">{status.branch || 'unknown'}</span>
      {status.hasChanges && (
        <span
          className="text-yellow-500 flex items-center gap-1"
          data-testid="git-changes-indicator"
          title="Uncommitted changes"
        >
          <span className="inline-block w-2 h-2 bg-yellow-500 rounded-full" />
          <span>Modified</span>
        </span>
      )}
    </div>
  )
}

export default GitStatusBar
