import { test, expect } from '@playwright/test'
import { setupApiMocks, mockIssuesResponse } from '../fixtures/api-mocks'

test.describe('Issues page', () => {
  test.beforeEach(async ({ page }) => {
    await setupApiMocks(page)
    await page.goto('/issues')
  })

  test('renders issue rows in table', async ({ page }) => {
    await expect(page.getByText('遅延チケットのサンプル')).toBeVisible()
    await expect(page.getByText('注意チケットのサンプル')).toBeVisible()
  })

  test('displays project key in project column', async ({ page }) => {
    // Issues table shows project_key (e.g. PROJ-A), not project_name
    await expect(page.getByText('PROJ-A').first()).toBeVisible()
  })

  test('displays assignee name', async ({ page }) => {
    await expect(page.getByText('担当者A')).toBeVisible()
  })

  test('displays delay status chip', async ({ page }) => {
    // RED issue shows 遅延 chip
    await expect(page.getByText('遅延').first()).toBeVisible()
  })

  test('delay filter buttons are rendered', async ({ page }) => {
    // Issues uses ToggleButtonGroup for delay filter
    await expect(page.getByRole('button', { name: 'すべて' })).toBeVisible()
    await expect(page.getByRole('button', { name: '遅延' })).toBeVisible()
    await expect(page.getByRole('button', { name: '注意' })).toBeVisible()
    await expect(page.getByRole('button', { name: '正常' })).toBeVisible()
  })

  test('CSV export button is visible', async ({ page }) => {
    await expect(page.getByRole('button', { name: /CSV出力/ })).toBeVisible()
  })

  test('filtering by RED updates results', async ({ page }) => {
    // Override mock to return only RED issues when filtered
    await page.route('**/api/v1/issues**', async (route) => {
      const url = route.request().url()
      if (url.includes('delay_status=RED')) {
        await route.fulfill({
          json: { data: [mockIssuesResponse.data[0]], pagination: { page: 1, per_page: 25, total: 1, total_pages: 1 } },
        })
      } else {
        await route.fulfill({ json: mockIssuesResponse })
      }
    })

    // Click the 遅延 (RED) toggle button
    await page.getByRole('button', { name: '遅延' }).click()

    // 注意チケット should no longer be visible
    await expect(page.getByText('注意チケットのサンプル')).not.toBeVisible()
  })
})
