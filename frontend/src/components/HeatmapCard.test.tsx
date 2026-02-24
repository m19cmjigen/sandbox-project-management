import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import HeatmapCard from './HeatmapCard'
import type { DashboardOrgNode } from '../types/dashboard'

// react-router-dom navigate mock
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

const makeNode = (partial: Partial<DashboardOrgNode> = {}): DashboardOrgNode => ({
  id: 1,
  name: 'テスト組織',
  parent_id: null,
  level: 0,
  total_projects: 5,
  red_projects: 2,
  yellow_projects: 1,
  green_projects: 2,
  delay_status: 'RED',
  delay_rate: 0.4,
  children: [],
  ...partial,
})

const renderCard = (node: DashboardOrgNode, isChild = false) =>
  render(
    <MemoryRouter>
      <HeatmapCard node={node} isChild={isChild} />
    </MemoryRouter>
  )

describe('HeatmapCard', () => {
  it('renders organization name', () => {
    renderCard(makeNode())
    expect(screen.getByText('テスト組織')).toBeInTheDocument()
  })

  it('renders delay status chip label', () => {
    renderCard(makeNode({ delay_status: 'RED' }))
    expect(screen.getByText('遅延あり')).toBeInTheDocument()
  })

  it('renders YELLOW status chip label', () => {
    renderCard(makeNode({ delay_status: 'YELLOW' }))
    expect(screen.getByText('注意')).toBeInTheDocument()
  })

  it('renders GREEN status chip label', () => {
    renderCard(makeNode({ delay_status: 'GREEN' }))
    expect(screen.getByText('正常')).toBeInTheDocument()
  })

  it('navigates to projects page with org id on click', () => {
    renderCard(makeNode({ id: 42 }))
    fireEvent.click(screen.getByText('テスト組織'))
    expect(mockNavigate).toHaveBeenCalledWith('/projects?organization_id=42')
  })

  it('shows project count in child mode', () => {
    renderCard(makeNode({ total_projects: 7 }), true)
    expect(screen.getByText(/7件/)).toBeInTheDocument()
  })

  it('shows red project count in non-child mode', () => {
    renderCard(makeNode({ red_projects: 3 }), false)
    // The redesigned component renders the count as a standalone number next to a colored dot
    expect(screen.getByText('3')).toBeInTheDocument()
  })
})
