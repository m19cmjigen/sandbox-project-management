import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ErrorMessage from './ErrorMessage'

describe('ErrorMessage', () => {
  it('renders the error message text', () => {
    render(<ErrorMessage message="取得に失敗しました" />)
    expect(screen.getByText('取得に失敗しました')).toBeInTheDocument()
  })

  it('does not show retry button when onRetry is not provided', () => {
    render(<ErrorMessage message="エラー" />)
    expect(screen.queryByRole('button', { name: '再試行' })).not.toBeInTheDocument()
  })

  it('shows retry button when onRetry is provided', () => {
    const onRetry = vi.fn()
    render(<ErrorMessage message="エラー" onRetry={onRetry} />)
    expect(screen.getByRole('button', { name: '再試行' })).toBeInTheDocument()
  })

  it('calls onRetry when retry button is clicked', () => {
    const onRetry = vi.fn()
    render(<ErrorMessage message="エラー" onRetry={onRetry} />)
    fireEvent.click(screen.getByRole('button', { name: '再試行' }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })
})
