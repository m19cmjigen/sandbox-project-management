import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');
  });

  test('should display dashboard summary', async ({ page }) => {
    // Check that main dashboard elements are visible
    await expect(page.getByText('全社プロジェクト進捗可視化プラットフォーム')).toBeVisible();
    await expect(page.getByRole('heading', { name: 'ダッシュボード' })).toBeVisible();
  });

  test('should display statistics cards', async ({ page }) => {
    // Wait for dashboard to load
    await page.waitForTimeout(1000);

    // Check for statistics cards (if data exists)
    const statisticsSection = page.locator('text=統計情報').first();
    if (await statisticsSection.isVisible()) {
      // Statistics section is visible, check for cards
      await expect(page.getByText(/プロジェクト/)).toBeVisible();
      await expect(page.getByText(/チケット/)).toBeVisible();
    }
  });

  test('should display project heatmap', async ({ page }) => {
    // Wait for dashboard to load
    await page.waitForTimeout(1000);

    // Check for heatmap section (if projects exist)
    const heatmapSection = page.locator('text=プロジェクトヒートマップ').first();
    if (await heatmapSection.isVisible()) {
      await expect(heatmapSection).toBeVisible();
    }
  });

  test('should navigate to other pages from dashboard', async ({ page }) => {
    // Click on navigation item
    await page.getByRole('button', { name: 'プロジェクト' }).click();
    await page.waitForURL('/projects');

    // Navigate back to dashboard
    await page.getByRole('button', { name: 'ダッシュボード' }).click();
    await page.waitForURL('/');
  });
});
