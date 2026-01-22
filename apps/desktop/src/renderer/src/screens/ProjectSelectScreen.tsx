/**
 * ProjectSelectScreen - Main project selection UI
 * 
 * Allows users to:
 * - Select a project folder via native dialog
 * - View recently opened projects
 * - See project detection results
 * - Provide project context for brownfield projects
 */

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ProjectInfoCard } from '../components/ProjectInfoCard'
import { RecentProjectItem } from '../components/RecentProjectItem'
import type { ProjectScanResult } from '../types/project'

export interface RecentProject {
  path: string
  name: string
  lastOpened: string
  context?: string
}

export function ProjectSelectScreen(): JSX.Element {
  const [recentProjects, setRecentProjects] = useState<RecentProject[]>([])
  const [selectedProject, setSelectedProject] = useState<ProjectScanResult | null>(null)
  const [loading, setLoading] = useState(false)

  // Load recent projects on mount
  useEffect(() => {
    loadRecentProjects()
  }, [])

  const loadRecentProjects = async (): Promise<void> => {
    try {
      const recent = await window.api.project.getRecent()
      setRecentProjects(recent)
    } catch (error) {
      console.error('Failed to load recent projects:', error)
      setRecentProjects([])
    }
  }

  const handleSelectFolder = async (): Promise<void> => {
    setLoading(true)
    try {
      const path = await window.api.dialog.selectFolder()
      if (path) {
        // Scan the selected folder
        const result = await window.api.project.scan(path)
        setSelectedProject(result)

        // Add to recent projects if it's a valid BMAD project
        if (result.isBmad) {
          await window.api.project.addRecent(path)
          // Reload recent projects list
          await loadRecentProjects()
        }
      }
    } catch (error) {
      console.error('Failed to select folder:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenRecent = async (path: string): Promise<void> => {
    setLoading(true)
    try {
      const result = await window.api.project.scan(path)
      setSelectedProject(result)

      // Update recent list (moves to top)
      if (result.isBmad) {
        await window.api.project.addRecent(path)
        await loadRecentProjects()
      }
    } catch (error) {
      console.error('Failed to open recent project:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">Select Project</h1>

      <div className="grid gap-6 md:grid-cols-2">
        {/* New Project Selection */}
        <Card>
          <CardHeader>
            <CardTitle>Open Project</CardTitle>
          </CardHeader>
          <CardContent>
            <Button 
              onClick={handleSelectFolder} 
              className="w-full"
              disabled={loading}
            >
              Select Project Folder
            </Button>
          </CardContent>
        </Card>

        {/* Recent Projects */}
        <Card>
          <CardHeader>
            <CardTitle>Recent Projects</CardTitle>
          </CardHeader>
          <CardContent>
            {recentProjects.length === 0 ? (
              <p className="text-muted-foreground">No recent projects</p>
            ) : (
              <ul className="space-y-2">
                {recentProjects.map((project) => (
                  <RecentProjectItem
                    key={project.path}
                    project={project}
                    onSelect={() => handleOpenRecent(project.path)}
                    onRemove={async () => {
                      await window.api.project.removeRecent(project.path)
                      await loadRecentProjects()
                    }}
                  />
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Project Info (shown after selection) */}
      {selectedProject && (
        <ProjectInfoCard project={selectedProject} />
      )}
    </div>
  )
}
