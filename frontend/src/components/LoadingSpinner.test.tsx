import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import LoadingSpinner from './LoadingSpinner'

describe('LoadingSpinner', () => {
  it('renders a circular progress indicator', () => {
    render(<LoadingSpinner />)
    expect(screen.getByRole('progressbar')).toBeInTheDocument()
  })

  it('applies minHeight to container', () => {
    const { container } = render(<LoadingSpinner minHeight={300} />)
    const wrapper = container.firstChild as HTMLElement
    expect(wrapper).toHaveStyle({ minHeight: '300px' })
  })
})
