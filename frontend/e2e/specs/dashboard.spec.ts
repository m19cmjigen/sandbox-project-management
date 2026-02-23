import { test, expect } from '@playwright/test'
import { setupApiMocks } from '../fixtures/api-mocks'

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await setupApiMocks(page)
    await page.goto('/')
  })

  test('displays summary stat cards', async ({ page }) => {
    // Stat card labels as defined in Dashboard.tsx SummaryCard
    await expect(page.getByText('総プロジェクト')).toBeVisible()
    await expect(page.getByText('遅延プロジェクト')).toBeVisible()
    await expect(page.getByText('注意プロジェクト')).toBeVisible()
    await expect(page.getByText('正常プロジェクト')).toBeVisible()
  })

  test('displays project count values', async ({ page }) => {
    // Total: 10 projects
    await expect(page.getByText('10').first()).toBeVisible()
  })

  test('displays org heatmap cards', async ({ page }) => {
    // Root orgs should appear as heatmap cards
    await expect(page.getByText('開発本部')).toBeVisible()
    await expect(page.getByText('営業本部')).toBeVisible()
  })

  test('heatmap card shows RED delay status chip', async ({ page }) => {
    // 遅延あり chip label is used for RED status
    await expect(page.getByText('遅延あり').first()).toBeVisible()
  })

  test('clicking heatmap card navigates to projects page', async ({ page }) => {
    // Click on the 開発本部 card
    const card = page.getByText('開発本部').first()
    await card.click()

    // Should navigate to /projects?organization_id=1
    await expect(page).toHaveURL(/\/projects.*organization_id=1/)
  })
})
