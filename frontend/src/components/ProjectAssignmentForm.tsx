import { useState } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Typography,
  Box,
} from '@mui/material'
import { ProjectWithStats, Organization } from '@/types'

interface ProjectAssignmentFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (projectId: number, organizationId: number | null) => Promise<void>
  project: ProjectWithStats | null
  organizations: Organization[]
}

export default function ProjectAssignmentForm({
  open,
  onClose,
  onSubmit,
  project,
  organizations,
}: ProjectAssignmentFormProps) {
  const [organizationId, setOrganizationId] = useState<number | ''>(
    project?.organization_id || ''
  )
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    if (!project) return

    setError(null)
    setSubmitting(true)

    try {
      await onSubmit(
        project.id,
        organizationId === '' ? null : (organizationId as number)
      )
      handleClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : '割り当てに失敗しました')
    } finally {
      setSubmitting(false)
    }
  }

  const handleClose = () => {
    setOrganizationId(project?.organization_id || '')
    setError(null)
    onClose()
  }

  if (!project) return null

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>プロジェクトを組織に割り当て</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ mb: 3, mt: 2 }}>
          <Typography variant="body2" color="textSecondary">
            プロジェクト
          </Typography>
          <Typography variant="h6">{project.name}</Typography>
          <Typography variant="body2" color="textSecondary">
            キー: {project.key}
          </Typography>
        </Box>

        <FormControl fullWidth margin="dense">
          <InputLabel>組織</InputLabel>
          <Select
            value={organizationId}
            label="組織"
            onChange={(e) => setOrganizationId(e.target.value as number | '')}
            disabled={submitting}
          >
            <MenuItem value="">未割り当て</MenuItem>
            {organizations.map((org) => (
              <MenuItem key={org.id} value={org.id}>
                {org.name} (レベル {org.level})
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={submitting}>
          キャンセル
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={submitting}
        >
          {submitting ? '保存中...' : '保存'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
