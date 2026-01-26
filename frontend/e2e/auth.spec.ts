import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Start from the login page
    await page.goto('/login');
  });

  test('should display login page correctly', async ({ page }) => {
    // Check that login page elements are visible
    await expect(page.getByRole('heading', { name: 'ログイン' })).toBeVisible();
    await expect(page.getByLabel('ユーザー名')).toBeVisible();
    await expect(page.getByLabel('パスワード')).toBeVisible();
    await expect(page.getByRole('button', { name: 'ログイン' })).toBeVisible();
  });

  test('should show error for invalid credentials', async ({ page }) => {
    // Fill in invalid credentials
    await page.getByLabel('ユーザー名').fill('invalid_user');
    await page.getByLabel('パスワード').fill('wrong_password');

    // Click login button
    await page.getByRole('button', { name: 'ログイン' }).click();

    // Should show error message
    await expect(page.getByText(/Invalid username or password|ユーザー名またはパスワードが正しくありません/i)).toBeVisible();
  });

  test('should successfully login with valid credentials', async ({ page }) => {
    // Fill in valid credentials (default admin account)
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');

    // Click login button
    await page.getByRole('button', { name: 'ログイン' }).click();

    // Should redirect to dashboard
    await page.waitForURL('/');

    // Should see dashboard content
    await expect(page.getByText('全社プロジェクト進捗可視化プラットフォーム')).toBeVisible();

    // Should see user menu
    await expect(page.getByText('admin')).toBeVisible();
  });

  test('should toggle password visibility', async ({ page }) => {
    const passwordInput = page.getByLabel('パスワード');
    const toggleButton = page.getByLabel('toggle password visibility');

    // Password should be hidden by default
    await expect(passwordInput).toHaveAttribute('type', 'password');

    // Click toggle to show password
    await toggleButton.click();
    await expect(passwordInput).toHaveAttribute('type', 'text');

    // Click toggle again to hide password
    await toggleButton.click();
    await expect(passwordInput).toHaveAttribute('type', 'password');
  });
});

test.describe('Authenticated User Actions', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.getByLabel('ユーザー名').fill('admin');
    await page.getByLabel('パスワード').fill('admin123');
    await page.getByRole('button', { name: 'ログイン' }).click();
    await page.waitForURL('/');
  });

  test('should navigate through menu items', async ({ page }) => {
    // Navigate to Organizations
    await page.getByRole('button', { name: '組織管理' }).click();
    await page.waitForURL('/organizations');
    await expect(page.getByRole('heading', { name: '組織一覧' })).toBeVisible();

    // Navigate to Projects
    await page.getByRole('button', { name: 'プロジェクト' }).click();
    await page.waitForURL('/projects');
    await expect(page.getByRole('heading', { name: 'プロジェクト一覧' })).toBeVisible();

    // Navigate to Issues
    await page.getByRole('button', { name: 'チケット' }).click();
    await page.waitForURL('/issues');
    await expect(page.getByRole('heading', { name: 'チケット一覧' })).toBeVisible();

    // Navigate to Dashboard
    await page.getByRole('button', { name: 'ダッシュボード' }).click();
    await page.waitForURL('/');
  });

  test('should logout successfully', async ({ page }) => {
    // Open user menu
    await page.getByLabel('account of current user').click();

    // Click logout
    await page.getByRole('menuitem', { name: 'ログアウト' }).click();

    // Should redirect to login page
    await page.waitForURL('/login');
    await expect(page.getByRole('heading', { name: 'ログイン' })).toBeVisible();
  });

  test('should display user role badge', async ({ page }) => {
    // Should see admin role badge
    await expect(page.getByText('管理者')).toBeVisible();
  });

  test('should access admin page as admin', async ({ page }) => {
    // Navigate to Admin page
    await page.getByRole('button', { name: '管理' }).click();
    await page.waitForURL('/admin');

    // Should see admin page content
    await expect(page.getByText('管理画面')).toBeVisible();
  });
});

test.describe('Protected Routes', () => {
  test('should redirect to login when accessing protected route without auth', async ({ page }) => {
    // Try to access dashboard without login
    await page.goto('/');

    // Should redirect to login
    await page.waitForURL('/login');
    await expect(page.getByRole('heading', { name: 'ログイン' })).toBeVisible();
  });

  test('should redirect to login when accessing admin page without auth', async ({ page }) => {
    // Try to access admin page without login
    await page.goto('/admin');

    // Should redirect to login
    await page.waitForURL('/login');
  });
});
