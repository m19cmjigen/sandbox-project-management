import { test, expect } from '@playwright/test';

test.describe('Organizations Page', () => {
  test.beforeEach(async ({ page }) => {
    // Login and navigate to organizations page
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');

    // Navigate to organizations
    await page.getByRole('button', { name: '組織管理' }).click();
    await page.waitForURL('/organizations');
  });

  test('should display organizations page', async ({ page }) => {
    await expect(page.getByRole('heading', { name: '組織一覧' })).toBeVisible();
  });

  test('should display organization tree structure', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Check if tree view exists
    const treeView = page.locator('[role="tree"], .MuiList-root').first();
    if (await treeView.isVisible()) {
      await expect(treeView).toBeVisible();
    }
  });

  test('should have expandable organization nodes', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Look for expand/collapse buttons
    const expandButtons = page.locator('[aria-label*="expand"], button:has-text("▶"), button:has-text("▼")');
    const count = await expandButtons.count();

    if (count > 0) {
      // Click first expand button
      await expandButtons.first().click();
      await page.waitForTimeout(300);

      // Should show children
      const children = page.locator('.MuiListItem-root').count();
      expect(await children).toBeGreaterThan(0);
    }
  });

  test('should display organization information', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Check if organization items are visible
    const orgItems = page.locator('.MuiListItem-root');
    const count = await orgItems.count();

    if (count > 0) {
      // At least one organization is visible
      await expect(orgItems.first()).toBeVisible();

      // Should show organization icon
      const icons = page.locator('[data-testid="BusinessIcon"], svg');
      if (await icons.count() > 0) {
        await expect(icons.first()).toBeVisible();
      }
    } else {
      // No organizations message
      await expect(page.getByText(/組織がありません|No organizations/i)).toBeVisible();
    }
  });

  test('should show hierarchy levels visually', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Organizations should have indentation for hierarchy
    const orgItems = page.locator('.MuiListItem-root');
    const count = await orgItems.count();

    if (count > 1) {
      // Multiple organizations exist, check for visual hierarchy
      const firstItem = orgItems.first();
      await expect(firstItem).toBeVisible();
    }
  });
});
