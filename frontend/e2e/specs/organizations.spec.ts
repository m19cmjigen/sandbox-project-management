import { test, expect } from '@playwright/test'
import { setupApiMocks, mockOrganizationsResponse } from '../fixtures/api-mocks'

test.describe('Organization management page', () => {
  test.beforeEach(async ({ page }) => {
    await setupApiMocks(page)

    // Add method-aware handler for organizations endpoint
    await page.route('**/api/v1/organizations', async (route) => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          json: {
            id: 10,
            name: '新規組織',
            parent_id: null,
            path: '/10/',
            level: 0,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
            total_projects: 0,
            red_projects: 0,
            yellow_projects: 0,
            green_projects: 0,
            delay_status: 'GREEN',
          },
        })
      } else {
        await route.fulfill({ json: mockOrganizationsResponse })
      }
    })

    await page.goto('/organizations/manage')
  })

  test('displays organization tree', async ({ page }) => {
    await expect(page.getByText('開発本部')).toBeVisible()
    await expect(page.getByText('第一開発部')).toBeVisible()
  })

  test('shows add root organization button', async ({ page }) => {
    // Button label is "本部を追加"
    await expect(page.getByRole('button', { name: /本部を追加/ })).toBeVisible()
  })

  test('opens create dialog when add button clicked', async ({ page }) => {
    await page.getByRole('button', { name: /本部を追加/ }).click()
    // Dialog title is "組織を追加"
    await expect(page.getByRole('dialog')).toBeVisible()
    await expect(page.getByText('組織を追加')).toBeVisible()
  })

  test('create dialog has name input field', async ({ page }) => {
    await page.getByRole('button', { name: /本部を追加/ }).click()
    await expect(page.getByLabel('組織名')).toBeVisible()
  })

  test('shows unassigned projects panel', async ({ page }) => {
    // Right panel heading is "未分類プロジェクト"
    await expect(page.getByText('未分類プロジェクト')).toBeVisible()
  })

  test('clicking cancel closes the dialog', async ({ page }) => {
    await page.getByRole('button', { name: /本部を追加/ }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()

    await page.getByRole('button', { name: /キャンセル/ }).click()
    await expect(dialog).not.toBeVisible()
  })
})
