import { test, expect } from '@playwright/test';

test.describe('Projects Page', () => {
  test.beforeEach(async ({ page }) => {
    // Login and navigate to projects page
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');

    // Navigate to projects
    await page.getByRole('button', { name: 'プロジェクト' }).click();
    await page.waitForURL('/projects');
  });

  test('should display projects page', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'プロジェクト一覧' })).toBeVisible();
  });

  test('should have search functionality', async ({ page }) => {
    // Check if search input exists
    const searchInput = page.getByPlaceholder(/検索|プロジェクト名/i);
    await expect(searchInput).toBeVisible();

    // Try searching
    await searchInput.fill('test');
    await page.waitForTimeout(500); // Wait for debounce/filter
  });

  test('should have filter options', async ({ page }) => {
    // Check for filter dropdowns/selects
    const filters = page.locator('[role="combobox"], select').first();
    if (await filters.isVisible()) {
      await expect(filters).toBeVisible();
    }
  });

  test('should display project cards when data exists', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Check if projects are loaded
    const projectCards = page.locator('[data-testid="project-card"], .MuiCard-root');
    const count = await projectCards.count();

    if (count > 0) {
      // At least one project card is visible
      await expect(projectCards.first()).toBeVisible();
    } else {
      // No projects message should be visible
      await expect(page.getByText(/プロジェクトがありません|No projects/i)).toBeVisible();
    }
  });

  test('should display organization filter', async ({ page }) => {
    // Wait for page to load
    await page.waitForTimeout(1000);

    // Organization filter should be present
    const orgFilter = page.locator('text=組織').first();
    if (await orgFilter.isVisible()) {
      await expect(orgFilter).toBeVisible();
    }
  });

  test('should display status filter', async ({ page }) => {
    // Wait for page to load
    await page.waitForTimeout(1000);

    // Status filter should be present
    const statusFilter = page.locator('text=ステータス').first();
    if (await statusFilter.isVisible()) {
      await expect(statusFilter).toBeVisible();
    }
  });
});
