import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ConfirmDialog from './ConfirmDialog'

const defaultProps = {
  open: true,
  title: '削除の確認',
  message: 'この操作は取り消せません。',
  onConfirm: vi.fn(),
  onClose: vi.fn(),
}

describe('ConfirmDialog', () => {
  it('renders title and message when open', () => {
    render(<ConfirmDialog {...defaultProps} />)
    expect(screen.getByText('削除の確認')).toBeInTheDocument()
    expect(screen.getByText('この操作は取り消せません。')).toBeInTheDocument()
  })

  it('does not render when closed', () => {
    render(<ConfirmDialog {...defaultProps} open={false} />)
    expect(screen.queryByText('削除の確認')).not.toBeInTheDocument()
  })

  it('calls onConfirm when confirm button is clicked', () => {
    const onConfirm = vi.fn()
    render(<ConfirmDialog {...defaultProps} onConfirm={onConfirm} />)
    fireEvent.click(screen.getByRole('button', { name: '確認' }))
    expect(onConfirm).toHaveBeenCalledTimes(1)
  })

  it('calls onClose when cancel button is clicked', () => {
    const onClose = vi.fn()
    render(<ConfirmDialog {...defaultProps} onClose={onClose} />)
    fireEvent.click(screen.getByRole('button', { name: 'キャンセル' }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('renders custom confirmLabel', () => {
    render(<ConfirmDialog {...defaultProps} confirmLabel="削除する" />)
    expect(screen.getByRole('button', { name: '削除する' })).toBeInTheDocument()
  })
})
