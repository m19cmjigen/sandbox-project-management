/** Format a date string for display. Returns '—' for null/empty. */
export function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  return dateStr
}

/**
 * Returns true when the due date is in the past and the issue is not Done.
 * Compares against midnight of the current day (local time).
 */
export function isDueDateOverdue(dateStr: string | null, statusCategory: string): boolean {
  if (!dateStr || statusCategory === 'Done') return false
  return new Date(dateStr) < new Date(new Date().toDateString())
}
