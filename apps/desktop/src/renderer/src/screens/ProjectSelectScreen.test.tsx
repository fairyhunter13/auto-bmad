/**
 * Tests for ProjectSelectScreen component.
 * Validates project selection UI, recent projects, and folder picker integration.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { ProjectSelectScreen } from './ProjectSelectScreen'

// Mock window.api
const mockApi = {
  dialog: {
    selectFolder: vi.fn()
  },
  project: {
    scan: vi.fn(),
    getRecent: vi.fn(),
    addRecent: vi.fn(),
    removeRecent: vi.fn()
  }
}

declare global {
  interface Window {
    api: typeof mockApi
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  window.api = mockApi
})

describe('ProjectSelectScreen', () => {
  it('renders project selection screen with title', () => {
    render(<ProjectSelectScreen />)
    expect(screen.getByText('Select Project')).toBeInTheDocument()
  })

  it('displays "Select Project Folder" button', () => {
    render(<ProjectSelectScreen />)
    expect(screen.getByRole('button', { name: /select project folder/i })).toBeInTheDocument()
  })

  it('displays recent projects section', () => {
    render(<ProjectSelectScreen />)
    expect(screen.getByText('Recent Projects')).toBeInTheDocument()
  })

  it('shows "No recent projects" when list is empty', async () => {
    mockApi.project.getRecent.mockResolvedValue([])
    
    render(<ProjectSelectScreen />)
    
    await waitFor(() => {
      expect(screen.getByText(/no recent projects/i)).toBeInTheDocument()
    })
  })

  it('calls selectFolder dialog when button is clicked', async () => {
    mockApi.dialog.selectFolder.mockResolvedValue(null)
    
    render(<ProjectSelectScreen />)
    
    const button = screen.getByRole('button', { name: /select project folder/i })
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(mockApi.dialog.selectFolder).toHaveBeenCalledTimes(1)
    })
  })

  it('scans project when folder is selected', async () => {
    const mockPath = '/home/user/my-project'
    const mockScanResult = {
      isBmad: true,
      projectType: 'greenfield' as const,
      bmadCompatible: true,
      minBmadVersion: '6.0.0',
      path: mockPath,
      hasBmadFolder: true,
      hasOutputFolder: false
    }

    mockApi.dialog.selectFolder.mockResolvedValue(mockPath)
    mockApi.project.scan.mockResolvedValue(mockScanResult)
    mockApi.project.addRecent.mockResolvedValue(undefined)
    mockApi.project.getRecent.mockResolvedValue([])

    render(<ProjectSelectScreen />)
    
    const button = screen.getByRole('button', { name: /select project folder/i })
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(mockApi.project.scan).toHaveBeenCalledWith(mockPath)
    })
  })

  it('displays project info after successful scan', async () => {
    const mockPath = '/home/user/my-project'
    const mockScanResult = {
      isBmad: true,
      projectType: 'greenfield' as const,
      bmadCompatible: true,
      minBmadVersion: '6.0.0',
      path: mockPath,
      hasBmadFolder: true,
      hasOutputFolder: false
    }

    mockApi.dialog.selectFolder.mockResolvedValue(mockPath)
    mockApi.project.scan.mockResolvedValue(mockScanResult)
    mockApi.project.addRecent.mockResolvedValue(undefined)
    mockApi.project.getRecent.mockResolvedValue([])

    render(<ProjectSelectScreen />)
    
    const button = screen.getByRole('button', { name: /select project folder/i })
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(screen.getByText('Project Details')).toBeInTheDocument()
    })
  })

  it('adds project to recent list when valid BMAD project is selected', async () => {
    const mockPath = '/home/user/my-project'
    const mockScanResult = {
      isBmad: true,
      projectType: 'greenfield' as const,
      bmadCompatible: true,
      minBmadVersion: '6.0.0',
      path: mockPath,
      hasBmadFolder: true,
      hasOutputFolder: false
    }

    mockApi.dialog.selectFolder.mockResolvedValue(mockPath)
    mockApi.project.scan.mockResolvedValue(mockScanResult)
    mockApi.project.addRecent.mockResolvedValue(undefined)
    mockApi.project.getRecent.mockResolvedValue([])

    render(<ProjectSelectScreen />)
    
    const button = screen.getByRole('button', { name: /select project folder/i })
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(mockApi.project.addRecent).toHaveBeenCalledWith(mockPath)
    })
  })

  it('does not add to recent list when non-BMAD project is selected', async () => {
    const mockPath = '/home/user/not-bmad'
    const mockScanResult = {
      isBmad: false,
      projectType: 'not-bmad' as const,
      bmadCompatible: false,
      minBmadVersion: '6.0.0',
      path: mockPath,
      hasBmadFolder: false,
      hasOutputFolder: false,
      error: 'Not a BMAD project'
    }

    mockApi.dialog.selectFolder.mockResolvedValue(mockPath)
    mockApi.project.scan.mockResolvedValue(mockScanResult)
    mockApi.project.getRecent.mockResolvedValue([])

    render(<ProjectSelectScreen />)
    
    const button = screen.getByRole('button', { name: /select project folder/i })
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(mockApi.project.addRecent).not.toHaveBeenCalled()
    })
  })

  it('displays recent projects when available', async () => {
    const mockRecentProjects = [
      {
        path: '/home/user/project1',
        name: 'project1',
        lastOpened: '2026-01-20T10:00:00Z'
      },
      {
        path: '/home/user/project2',
        name: 'project2',
        lastOpened: '2026-01-19T10:00:00Z'
      }
    ]

    mockApi.project.getRecent.mockResolvedValue(mockRecentProjects)

    render(<ProjectSelectScreen />)
    
    await waitFor(() => {
      expect(screen.getByText('project1')).toBeInTheDocument()
      expect(screen.getByText('project2')).toBeInTheDocument()
    })
  })
})
