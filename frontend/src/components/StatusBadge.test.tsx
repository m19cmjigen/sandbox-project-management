import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import StatusBadge from './StatusBadge'

describe('StatusBadge', () => {
  it('renders "遅延" label for RED status', () => {
    render(<StatusBadge status="RED" />)
    expect(screen.getByText('遅延')).toBeInTheDocument()
  })

  it('renders "注意" label for YELLOW status', () => {
    render(<StatusBadge status="YELLOW" />)
    expect(screen.getByText('注意')).toBeInTheDocument()
  })

  it('renders "正常" label for GREEN status', () => {
    render(<StatusBadge status="GREEN" />)
    expect(screen.getByText('正常')).toBeInTheDocument()
  })
})
