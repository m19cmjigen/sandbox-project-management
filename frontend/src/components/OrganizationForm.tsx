import { useState } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
} from '@mui/material'
import { Organization } from '@/types'

interface OrganizationFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: OrganizationFormData) => Promise<void>
  organizations: Organization[]
  editingOrganization?: Organization | null
}

export interface OrganizationFormData {
  name: string
  parent_id: number | null
}

export default function OrganizationForm({
  open,
  onClose,
  onSubmit,
  organizations,
  editingOrganization,
}: OrganizationFormProps) {
  const [name, setName] = useState(editingOrganization?.name || '')
  const [parentId, setParentId] = useState<number | ''>(
    editingOrganization?.parent_id || ''
  )
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    // Validation
    if (!name.trim()) {
      setError('組織名は必須です')
      return
    }

    // Check for circular reference
    if (editingOrganization && parentId === editingOrganization.id) {
      setError('自分自身を親組織に設定することはできません')
      return
    }

    setError(null)
    setSubmitting(true)

    try {
      await onSubmit({
        name: name.trim(),
        parent_id: parentId === '' ? null : (parentId as number),
      })
      handleClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存に失敗しました')
    } finally {
      setSubmitting(false)
    }
  }

  const handleClose = () => {
    setName('')
    setParentId('')
    setError(null)
    onClose()
  }

  // Filter out current organization and its children from parent options
  const availableParents = organizations.filter((org) => {
    if (!editingOrganization) return true
    if (org.id === editingOrganization.id) return false
    // Check if org is a descendant of editingOrganization
    return !org.path.startsWith(editingOrganization.path + '.')
  })

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {editingOrganization ? '組織を編集' : '新しい組織を作成'}
      </DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <TextField
          autoFocus
          margin="dense"
          label="組織名"
          type="text"
          fullWidth
          variant="outlined"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={submitting}
          required
          sx={{ mt: 2 }}
        />

        <FormControl fullWidth margin="dense" sx={{ mt: 2 }}>
          <InputLabel>親組織</InputLabel>
          <Select
            value={parentId}
            label="親組織"
            onChange={(e) => setParentId(e.target.value as number | '')}
            disabled={submitting}
          >
            <MenuItem value="">なし（トップレベル）</MenuItem>
            {availableParents.map((org) => (
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
