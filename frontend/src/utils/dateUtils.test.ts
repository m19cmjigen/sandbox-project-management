import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { formatDate, isDueDateOverdue } from './dateUtils'

describe('formatDate', () => {
  it('returns the date string as-is when provided', () => {
    expect(formatDate('2026-03-15')).toBe('2026-03-15')
  })

  it('returns em dash for null', () => {
    expect(formatDate(null)).toBe('—')
  })

  it('returns em dash for empty string', () => {
    expect(formatDate('')).toBe('—')
  })
})

describe('isDueDateOverdue', () => {
  // Fix "today" to 2026-02-23 for deterministic tests
  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-02-23T12:00:00'))
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('returns true when due date is in the past and status is not Done', () => {
    expect(isDueDateOverdue('2026-02-20', 'To Do')).toBe(true)
  })

  it('returns false when due date is today', () => {
    // Today's date is not strictly less than today's midnight
    expect(isDueDateOverdue('2026-02-23', 'In Progress')).toBe(false)
  })

  it('returns false when due date is in the future', () => {
    expect(isDueDateOverdue('2026-03-01', 'To Do')).toBe(false)
  })

  it('returns false when dateStr is null', () => {
    expect(isDueDateOverdue(null, 'To Do')).toBe(false)
  })

  it('returns false when status is Done regardless of date', () => {
    expect(isDueDateOverdue('2026-01-01', 'Done')).toBe(false)
  })
})
