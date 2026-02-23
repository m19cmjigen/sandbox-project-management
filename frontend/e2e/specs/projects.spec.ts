import { test, expect } from '@playwright/test'
import { setupApiMocks } from '../fixtures/api-mocks'

test.describe('Projects page', () => {
  test.beforeEach(async ({ page }) => {
    await setupApiMocks(page)
    await page.goto('/projects')
  })

  test('renders project cards', async ({ page }) => {
    await expect(page.getByText('プロジェクトA')).toBeVisible()
    await expect(page.getByText('プロジェクトB')).toBeVisible()
  })

  test('displays project key', async ({ page }) => {
    await expect(page.getByText('PROJ-A')).toBeVisible()
  })

  test('shows delay status chip on RED project', async ({ page }) => {
    // RED status chip label is 遅延あり (from statusConfig in ProjectCard)
    await expect(page.getByText('遅延あり').first()).toBeVisible()
  })

  test('displays issue count in project card', async ({ page }) => {
    // RED project has 3 red issues: "遅延 3"
    await expect(page.getByText('遅延 3')).toBeVisible()
  })

  test('navigates to issues page when clicking ticket button', async ({ page }) => {
    // IconButton with aria-label="チケット一覧"
    const btn = page.getByRole('button', { name: 'チケット一覧' }).first()
    await btn.click()
    await expect(page).toHaveURL(/\/issues.*project_id=1/)
  })

  test('sort select combobox is rendered', async ({ page }) => {
    await expect(page.getByRole('combobox').first()).toBeVisible()
  })
})
