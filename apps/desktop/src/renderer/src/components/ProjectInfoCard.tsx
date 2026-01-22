/**
 * ProjectInfoCard - Display project detection results
 * 
 * Shows:
 * - Project type (greenfield/brownfield)
 * - BMAD version and compatibility
 * - Existing artifacts (brownfield only)
 * - Project context editor (brownfield only)
 */

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import type { ProjectScanResult } from '../types/project'
import { CheckCircle, XCircle } from 'lucide-react'

interface ProjectInfoCardProps {
  project: ProjectScanResult
}

export function ProjectInfoCard({ project }: ProjectInfoCardProps): JSX.Element {
  const [context, setContext] = useState('')
  const [savingContext, setSavingContext] = useState(false)

  const handleSaveContext = async (): Promise<void> => {
    setSavingContext(true)
    try {
      await window.api.project.setContext(project.path, context)
    } catch (error) {
      console.error('Failed to save context:', error)
    } finally {
      setSavingContext(false)
    }
  }

  // Non-BMAD project error state
  if (!project.isBmad) {
    return (
      <Card className="mt-6 border-destructive">
        <CardContent className="pt-6">
          <p className="text-destructive font-medium">Not a BMAD Project</p>
          <p className="text-muted-foreground mt-2">
            Run <code className="bg-muted px-1 py-0.5 rounded">bmad-init</code> to initialize this folder as a BMAD project.
          </p>
          {project.error && (
            <p className="text-sm text-muted-foreground mt-2">{project.error}</p>
          )}
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="mt-6">
      <CardHeader>
        <CardTitle>Project Details</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Project Type Badge */}
        <div className="flex items-center gap-2">
          <span className="font-medium">Type:</span>
          <Badge variant={project.projectType === 'greenfield' ? 'default' : 'secondary'}>
            {project.projectType}
          </Badge>
        </div>

        {/* BMAD Version */}
        {project.bmadVersion && (
          <div className="flex items-center gap-2">
            <span className="font-medium">BMAD Version:</span>
            <span>{project.bmadVersion}</span>
            {project.bmadCompatible ? (
              <CheckCircle className="h-4 w-4 text-green-500" />
            ) : (
              <XCircle className="h-4 w-4 text-red-500" />
            )}
          </div>
        )}

        {!project.bmadCompatible && (
          <div className="text-sm text-destructive">
            ⚠️ This project requires BMAD {project.minBmadVersion} or higher
          </div>
        )}

        {/* Existing Artifacts (Brownfield) */}
        {project.existingArtifacts && project.existingArtifacts.length > 0 && (
          <div>
            <span className="font-medium">Existing Artifacts:</span>
            <ul className="mt-2 space-y-1">
              {project.existingArtifacts.map((artifact) => (
                <li key={artifact.path} className="text-sm text-muted-foreground">
                  • {artifact.name} ({artifact.type})
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Project Context (Brownfield only) */}
        {project.projectType === 'brownfield' && (
          <div>
            <Label htmlFor="context">Project Context (optional)</Label>
            <Textarea
              id="context"
              placeholder="Describe your project context for better AI assistance..."
              value={context}
              onChange={(e) => setContext(e.target.value)}
              className="mt-2"
              rows={4}
            />
            <Button
              variant="outline"
              size="sm"
              className="mt-2"
              onClick={handleSaveContext}
              disabled={savingContext || !context.trim()}
            >
              {savingContext ? 'Saving...' : 'Save Context'}
            </Button>
          </div>
        )}

        {/* Continue Button */}
        <Button className="w-full mt-4">
          Continue to Dashboard
        </Button>
      </CardContent>
    </Card>
  )
}
