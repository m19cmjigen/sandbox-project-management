import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Alert,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  FormControlLabel,
  IconButton,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  People as PeopleIcon,
  LockReset as LockResetIcon,
} from '@mui/icons-material'
import type { SelectChangeEvent } from '@mui/material'
import {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  changePassword,
  type User,
} from '../api/users'
import type { Role } from '../stores/authStore'
import { useAuthStore } from '../stores/authStore'
import { canManageUsers } from '../utils/permissions'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import ConfirmDialog from '../components/ConfirmDialog'

const roleLabel: Record<Role, string> = {
  admin: '管理者',
  project_manager: 'PJ管理者',
  viewer: '閲覧者',
}

const roleChipColor: Record<Role, 'error' | 'warning' | 'default'> = {
  admin: 'error',
  project_manager: 'warning',
  viewer: 'default',
}

// ---- Create Dialog ----

interface CreateDialogProps {
  open: boolean
  onClose: () => void
  onCreated: () => void
}

function CreateUserDialog({ open, onClose, onCreated }: CreateDialogProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState<Role>('viewer')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleClose = () => {
    setEmail('')
    setPassword('')
    setRole('viewer')
    setError(null)
    onClose()
  }

  const handleSubmit = async () => {
    setError(null)
    setSaving(true)
    try {
      await createUser({ email, password, role })
      handleClose()
      onCreated()
    } catch {
      setError('ユーザーの作成に失敗しました。メールアドレスが既に使用されている可能性があります。')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>ユーザー作成</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <TextField
          label="メールアドレス"
          type="email"
          fullWidth
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          sx={{ mt: 1, mb: 2 }}
        />
        <TextField
          label="パスワード (8文字以上)"
          type="password"
          fullWidth
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          sx={{ mb: 2 }}
        />
        <FormControl fullWidth>
          <InputLabel>ロール</InputLabel>
          <Select
            label="ロール"
            value={role}
            onChange={(e: SelectChangeEvent) => setRole(e.target.value as Role)}
          >
            <MenuItem value="admin">管理者</MenuItem>
            <MenuItem value="project_manager">PJ管理者</MenuItem>
            <MenuItem value="viewer">閲覧者</MenuItem>
          </Select>
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>キャンセル</Button>
        <Button onClick={handleSubmit} variant="contained" disabled={saving || !email || !password}>
          {saving ? '作成中...' : '作成'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

// ---- Edit Dialog ----

interface EditDialogProps {
  user: User | null
  onClose: () => void
  onUpdated: () => void
}

function EditUserDialog({ user, onClose, onUpdated }: EditDialogProps) {
  const [role, setRole] = useState<Role>('viewer')
  const [isActive, setIsActive] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // ダイアログが開いたときにユーザー情報をセット
  useEffect(() => {
    if (user) {
      setRole(user.role)
      setIsActive(user.is_active)
      setError(null)
    }
  }, [user])

  const handleSubmit = async () => {
    if (!user) return
    setError(null)
    setSaving(true)
    try {
      await updateUser(user.id, { role, is_active: isActive })
      onClose()
      onUpdated()
    } catch {
      setError('ユーザーの更新に失敗しました。')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={user !== null} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>ユーザー編集</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2, mt: 1 }}>
          {user?.email}
        </Typography>
        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel>ロール</InputLabel>
          <Select
            label="ロール"
            value={role}
            onChange={(e: SelectChangeEvent) => setRole(e.target.value as Role)}
          >
            <MenuItem value="admin">管理者</MenuItem>
            <MenuItem value="project_manager">PJ管理者</MenuItem>
            <MenuItem value="viewer">閲覧者</MenuItem>
          </Select>
        </FormControl>
        <FormControlLabel
          control={<Switch checked={isActive} onChange={(e) => setIsActive(e.target.checked)} />}
          label="アクティブ"
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>キャンセル</Button>
        <Button onClick={handleSubmit} variant="contained" disabled={saving}>
          {saving ? '更新中...' : '更新'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

// ---- Change Password Dialog ----

interface ChangePasswordDialogProps {
  user: User | null
  onClose: () => void
  onDone: () => void
}

function ChangePasswordDialog({ user, onClose, onDone }: ChangePasswordDialogProps) {
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // ダイアログが開いたときに入力をリセット
  useEffect(() => {
    if (user) {
      setNewPassword('')
      setConfirmPassword('')
      setError(null)
    }
  }, [user])

  const handleClose = () => {
    setNewPassword('')
    setConfirmPassword('')
    setError(null)
    onClose()
  }

  const handleSubmit = async () => {
    if (!user) return
    if (newPassword.length < 8) {
      setError('パスワードは8文字以上で入力してください。')
      return
    }
    if (newPassword !== confirmPassword) {
      setError('パスワードが一致しません。')
      return
    }
    setError(null)
    setSaving(true)
    try {
      await changePassword(user.id, newPassword)
      handleClose()
      onDone()
    } catch {
      setError('パスワードの変更に失敗しました。')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={user !== null} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>パスワード変更</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2, mt: 1 }}>
          {user?.email}
        </Typography>
        <TextField
          label="新しいパスワード (8文字以上)"
          type="password"
          fullWidth
          required
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          sx={{ mb: 2 }}
        />
        <TextField
          label="新しいパスワード (確認)"
          type="password"
          fullWidth
          required
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>キャンセル</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={saving || !newPassword || !confirmPassword}
        >
          {saving ? '変更中...' : '変更'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

// ---- Main Page ----

export default function UserManagement() {
  const navigate = useNavigate()
  const currentUser = useAuthStore((s) => s.user)

  // admin 以外はダッシュボードにリダイレクト
  useEffect(() => {
    if (currentUser && !canManageUsers(currentUser.role)) {
      navigate('/', { replace: true })
    }
  }, [currentUser, navigate])

  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<User | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<User | null>(null)
  const [passwordTarget, setPasswordTarget] = useState<User | null>(null)

  const fetchUsers = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getUsers()
      setUsers(data)
    } catch {
      setError('ユーザー一覧の取得に失敗しました。')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  const handleDelete = async () => {
    if (!deleteTarget) return
    try {
      await deleteUser(deleteTarget.id)
      setDeleteTarget(null)
      fetchUsers()
    } catch {
      setError('ユーザーの削除に失敗しました。')
      setDeleteTarget(null)
    }
  }

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage message={error} onRetry={fetchUsers} />

  return (
    <Box>
      <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 3 }}>
        <Stack direction="row" alignItems="center" spacing={1}>
          <PeopleIcon color="primary" />
          <Typography variant="h5" fontWeight={700}>
            ユーザー管理
          </Typography>
        </Stack>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setCreateOpen(true)}>
          ユーザー作成
        </Button>
      </Stack>

      <TableContainer
        component={Paper}
        elevation={0}
        sx={{ border: '1px solid', borderColor: 'divider' }}
      >
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>メールアドレス</TableCell>
              <TableCell>ロール</TableCell>
              <TableCell>状態</TableCell>
              <TableCell align="right">操作</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((u) => (
              <TableRow key={u.id} hover>
                <TableCell>{u.id}</TableCell>
                <TableCell>{u.email}</TableCell>
                <TableCell>
                  <Chip
                    label={roleLabel[u.role]}
                    size="small"
                    color={roleChipColor[u.role]}
                    variant="outlined"
                  />
                </TableCell>
                <TableCell>
                  <Chip
                    label={u.is_active ? 'アクティブ' : '無効'}
                    size="small"
                    color={u.is_active ? 'success' : 'default'}
                    variant="outlined"
                  />
                </TableCell>
                <TableCell align="right">
                  <Tooltip title="編集">
                    <IconButton size="small" onClick={() => setEditTarget(u)}>
                      <EditIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="パスワード変更">
                    <IconButton size="small" onClick={() => setPasswordTarget(u)}>
                      <LockResetIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title={u.id === currentUser?.id ? '自分自身は削除できません' : '削除'}>
                    {/* Tooltipはdisabledなボタンに対して機能しないため、spanでラップ */}
                    <span>
                      <IconButton
                        size="small"
                        color="error"
                        disabled={u.id === currentUser?.id}
                        onClick={() => setDeleteTarget(u)}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </span>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))}
            {users.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} align="center" sx={{ py: 4 }}>
                  <Typography color="text.secondary">ユーザーが存在しません</Typography>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>

      <CreateUserDialog
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={fetchUsers}
      />

      <EditUserDialog
        user={editTarget}
        onClose={() => setEditTarget(null)}
        onUpdated={fetchUsers}
      />

      <ChangePasswordDialog
        user={passwordTarget}
        onClose={() => setPasswordTarget(null)}
        onDone={fetchUsers}
      />

      <ConfirmDialog
        open={deleteTarget !== null}
        title="ユーザー削除"
        message={`${deleteTarget?.email} を削除しますか？この操作は元に戻せません。`}
        confirmLabel="削除"
        destructive
        onConfirm={handleDelete}
        onClose={() => setDeleteTarget(null)}
      />
    </Box>
  )
}
