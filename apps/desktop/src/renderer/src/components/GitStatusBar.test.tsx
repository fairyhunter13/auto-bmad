import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { GitStatusBar } from './GitStatusBar'

// Mock the window.api object
const mockGitApi = {
  getRepoStatus: vi.fn()
}

beforeEach(() => {
  // Reset mock
  mockGitApi.getRepoStatus.mockReset()

  // Setup global mock
  global.window.api = {
    git: mockGitApi
  } as unknown as typeof window.api
})

afterEach(() => {
  vi.clearAllMocks()
})

describe('GitStatusBar', () => {
  it('renders nothing when path is empty', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: false,
      hasChanges: false
    })

    const { container } = render(<GitStatusBar projectPath="" />)

    // Wait a bit for any async operations
    await waitFor(() => {
      expect(container.firstChild).toBeNull()
    })
  })

  it('renders nothing when not a git repo', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: false,
      hasChanges: false
    })

    const { container } = render(<GitStatusBar projectPath="/some/path" />)

    await waitFor(() => {
      expect(container.firstChild).toBeNull()
    })

    expect(mockGitApi.getRepoStatus).toHaveBeenCalledWith('/some/path')
  })

  it('displays branch name when in a git repo', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: true,
      branch: 'main',
      hasChanges: false
    })

    render(<GitStatusBar projectPath="/git/repo" />)

    await waitFor(() => {
      expect(screen.getByTestId('git-status-bar')).toBeInTheDocument()
    })

    expect(screen.getByTestId('git-branch')).toHaveTextContent('main')
  })

  it('shows changes indicator when there are uncommitted changes', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: true,
      branch: 'feature-branch',
      hasChanges: true
    })

    render(<GitStatusBar projectPath="/git/repo" />)

    await waitFor(() => {
      expect(screen.getByTestId('git-changes-indicator')).toBeInTheDocument()
    })

    expect(screen.getByTestId('git-changes-indicator')).toHaveTextContent('Modified')
  })

  it('does not show changes indicator when there are no uncommitted changes', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: true,
      branch: 'main',
      hasChanges: false
    })

    render(<GitStatusBar projectPath="/git/repo" />)

    await waitFor(() => {
      expect(screen.getByTestId('git-status-bar')).toBeInTheDocument()
    })

    expect(screen.queryByTestId('git-changes-indicator')).not.toBeInTheDocument()
  })

  it('handles API errors gracefully', async () => {
    mockGitApi.getRepoStatus.mockRejectedValue(new Error('API error'))

    const { container } = render(<GitStatusBar projectPath="/error/path" />)

    await waitFor(() => {
      expect(container.firstChild).toBeNull()
    })
  })

  it('shows "unknown" when branch is not available', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: true,
      branch: undefined,
      hasChanges: false
    })

    render(<GitStatusBar projectPath="/git/repo" />)

    await waitFor(() => {
      expect(screen.getByTestId('git-branch')).toHaveTextContent('unknown')
    })
  })

  it('refetches status when projectPath changes', async () => {
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: true,
      branch: 'main',
      hasChanges: false
    })

    const { rerender } = render(<GitStatusBar projectPath="/path1" />)

    await waitFor(() => {
      expect(mockGitApi.getRepoStatus).toHaveBeenCalledWith('/path1')
    })

    // Change the path
    mockGitApi.getRepoStatus.mockResolvedValue({
      isGitRepo: true,
      branch: 'develop',
      hasChanges: true
    })

    rerender(<GitStatusBar projectPath="/path2" />)

    await waitFor(() => {
      expect(mockGitApi.getRepoStatus).toHaveBeenCalledWith('/path2')
    })
  })
})
