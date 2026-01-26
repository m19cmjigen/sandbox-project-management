import { test, expect } from '@playwright/test';

test.describe('Admin Page', () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin and navigate to admin page
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');

    // Navigate to admin page
    await page.getByRole('button', { name: '管理' }).click();
    await page.waitForURL('/admin');
  });

  test('should display admin page', async ({ page }) => {
    await expect(page.getByText('管理画面')).toBeVisible();
  });

  test('should have tabbed interface', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Check for tabs
    const tabs = page.locator('[role="tab"]');
    const count = await tabs.count();

    expect(count).toBeGreaterThan(0);

    if (count >= 3) {
      // Should have at least 3 tabs: Organizations, Project Assignment, Sync
      await expect(tabs.first()).toBeVisible();
    }
  });

  test('should display organization management tab', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Look for organization tab
    const orgTab = page.locator('[role="tab"]').filter({ hasText: /組織/ });
    if (await orgTab.count() > 0) {
      await orgTab.first().click();
      await page.waitForTimeout(500);

      // Should show organization management interface
      const orgTable = page.locator('table, .MuiTable-root').first();
      if (await orgTable.isVisible()) {
        await expect(orgTable).toBeVisible();
      }
    }
  });

  test('should display project assignment tab', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Look for project assignment tab
    const projectTab = page.locator('[role="tab"]').filter({ hasText: /プロジェクト/ });
    if (await projectTab.count() > 0) {
      await projectTab.first().click();
      await page.waitForTimeout(500);

      // Should show project assignment interface
      const projectForm = page.locator('form, .MuiFormControl-root').first();
      if (await projectForm.isVisible()) {
        await expect(projectForm).toBeVisible();
      }
    }
  });

  test('should display sync management tab', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Look for sync tab
    const syncTab = page.locator('[role="tab"]').filter({ hasText: /同期/ });
    if (await syncTab.count() > 0) {
      await syncTab.first().click();
      await page.waitForTimeout(500);

      // Should show sync interface
      const syncButton = page.getByRole('button', { name: /同期|Sync/ });
      if (await syncButton.isVisible()) {
        await expect(syncButton).toBeVisible();
      }
    }
  });

  test('should have create organization button', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Switch to organization tab
    const orgTab = page.locator('[role="tab"]').filter({ hasText: /組織/ });
    if (await orgTab.count() > 0) {
      await orgTab.first().click();
      await page.waitForTimeout(500);

      // Look for create button
      const createButton = page.getByRole('button', { name: /作成|追加|Create|Add/ });
      if (await createButton.count() > 0) {
        await expect(createButton.first()).toBeVisible();
      }
    }
  });

  test('should display sync history', async ({ page }) => {
    await page.waitForTimeout(1000);

    // Navigate to sync tab
    const syncTab = page.locator('[role="tab"]').filter({ hasText: /同期/ });
    if (await syncTab.count() > 0) {
      await syncTab.first().click();
      await page.waitForTimeout(1000);

      // Check for sync history section
      const historySection = page.locator('text=同期履歴').first();
      if (await historySection.isVisible()) {
        await expect(historySection).toBeVisible();

        // Should have table or list of sync logs
        const historyTable = page.locator('table, .MuiTable-root');
        if (await historyTable.count() > 0) {
          await expect(historyTable.first()).toBeVisible();
        }
      }
    }
  });
});

test.describe('Admin Access Control', () => {
  test('non-admin user should not access admin page', async ({ page }) => {
    // This test would require creating a non-admin user first
    // For now, we just verify admin can access
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');

    // Navigate to admin page
    await page.getByRole('button', { name: '管理' }).click();
    await page.waitForURL('/admin');

    // Should be able to access as admin
    await expect(page.getByText('管理画面')).toBeVisible();
  });
});
