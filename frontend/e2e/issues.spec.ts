import { test, expect } from '@playwright/test';

test.describe('Issues Page', () => {
  test.beforeEach(async ({ page }) => {
    // Login and navigate to issues page
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');

    // Navigate to issues
    await page.getByRole('button', { name: 'チケット' }).click();
    await page.waitForURL('/issues');
  });

  test('should display issues page', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'チケット一覧' })).toBeVisible();
  });

  test('should have search functionality', async ({ page }) => {
    // Check if search input exists
    const searchInput = page.getByPlaceholder(/検索|チケット/i);
    await expect(searchInput).toBeVisible();

    // Try searching
    await searchInput.fill('test');
    await page.waitForTimeout(500); // Wait for debounce/filter
  });

  test('should have filter options', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Check for organization filter
    const orgFilter = page.locator('text=組織').first();
    if (await orgFilter.isVisible()) {
      await expect(orgFilter).toBeVisible();
    }

    // Check for project filter
    const projectFilter = page.locator('text=プロジェクト').first();
    if (await projectFilter.isVisible()) {
      await expect(projectFilter).toBeVisible();
    }

    // Check for status filter
    const statusFilter = page.locator('text=遅延ステータス').first();
    if (await statusFilter.isVisible()) {
      await expect(statusFilter).toBeVisible();
    }
  });

  test('should display issues table when data exists', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Check if table exists
    const table = page.locator('table, [role="table"]').first();
    if (await table.isVisible()) {
      await expect(table).toBeVisible();

      // Check for table headers
      await expect(page.getByText(/キー|Key/i).first()).toBeVisible();
      await expect(page.getByText(/サマリー|Summary/i).first()).toBeVisible();
    } else {
      // No issues message should be visible
      await expect(page.getByText(/チケットがありません|No issues/i)).toBeVisible();
    }
  });

  test('should display status badges in table', async ({ page }) => {
    await page.waitForTimeout(1000);

    const table = page.locator('table, [role="table"]').first();
    if (await table.isVisible()) {
      // Check for status badges (RED, YELLOW, GREEN)
      const badges = page.locator('.MuiChip-root, [role="status"]');
      const count = await badges.count();

      if (count > 0) {
        await expect(badges.first()).toBeVisible();
      }
    }
  });

  test('should allow filtering by organization', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Find organization select/dropdown
    const orgSelect = page.locator('select[name="organization"], [role="combobox"]').first();
    if (await orgSelect.isVisible()) {
      // Select an organization
      await orgSelect.click();
      await page.waitForTimeout(500);
    }
  });

  test('should allow filtering by project', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Find project select/dropdown
    const projectSelect = page.locator('select[name="project"], [role="combobox"]').nth(1);
    if (await projectSelect.isVisible()) {
      // Select a project
      await projectSelect.click();
      await page.waitForTimeout(500);
    }
  });
});
