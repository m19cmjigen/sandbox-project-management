# E2E Testing with Playwright

このディレクトリには、Playwrightを使用したエンドツーエンド（E2E）テストが含まれています。

## 概要

Playwrightを使用して、アプリケーションの主要なユーザーフローをテストします：

- **認証フロー** (`auth.spec.ts`)
  - ログイン/ログアウト
  - 無効な認証情報のエラーハンドリング
  - 保護されたルートへのアクセス制御
  - ユーザーメニューとロール表示

- **ダッシュボード** (`dashboard.spec.ts`)
  - サマリー統計の表示
  - プロジェクトヒートマップ
  - ナビゲーション

- **プロジェクト管理** (`projects.spec.ts`)
  - プロジェクト一覧表示
  - 検索・フィルタリング機能
  - プロジェクトカード表示

- **チケット管理** (`issues.spec.ts`)
  - チケット一覧表示
  - テーブル表示
  - ステータスバッジ
  - 複数条件でのフィルタリング

- **組織管理** (`organizations.spec.ts`)
  - 組織ツリー表示
  - 階層構造の展開/折りたたみ
  - 組織情報の表示

- **管理画面** (`admin.spec.ts`)
  - タブベースのインターフェース
  - 組織管理
  - プロジェクト割り当て
  - 同期管理と履歴
  - アクセス制御

## セットアップ

### 前提条件

1. Node.js 18以上
2. バックエンドAPIが実行中（`http://localhost:8080`）
3. フロントエンド開発サーバーまたはビルド済みアプリ

### インストール

```bash
# 依存関係のインストール（初回のみ）
npm install

# Playwrightブラウザのインストール（初回のみ）
npx playwright install
```

## テストの実行

### すべてのテストを実行

```bash
npm run test:e2e
```

### UIモードで実行（推奨）

```bash
npm run test:e2e:ui
```

UIモードでは、テストの実行状況をビジュアルに確認でき、デバッグが容易です。

### ヘッドモードで実行（ブラウザを表示）

```bash
npm run test:e2e:headed
```

### デバッグモードで実行

```bash
npm run test:e2e:debug
```

### 特定のテストファイルを実行

```bash
npx playwright test auth.spec.ts
```

### 特定のブラウザで実行

```bash
npx playwright test --project=chromium
npx playwright test --project=firefox
npx playwright test --project=webkit
```

### テストレポートの表示

```bash
npm run test:e2e:report
```

## テストの書き方

### 基本的なテスト構造

```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test.beforeEach(async ({ page }) => {
    // テスト前のセットアップ
    await page.goto('/login');
  });

  test('should do something', async ({ page }) => {
    // テストコード
    await expect(page.getByText('Expected Text')).toBeVisible();
  });
});
```

### 認証が必要なテスト

認証が必要なページをテストする場合は、`beforeEach`でログインします：

```typescript
test.beforeEach(async ({ page }) => {
  // ログイン
  await page.goto('/login');
  await page.getByLabel('ユーザー名').fill('admin');
  await page.getByLabel('パスワード').fill('admin123');
  await page.getByRole('button', { name: 'ログイン' }).click();
  await page.waitForURL('/');
});
```

### ベストプラクティス

1. **明確なセレクタを使用**
   - 優先順位: `getByRole` > `getByLabel` > `getByText` > `getByTestId` > CSSセレクタ

2. **待機を適切に使用**
   - `waitForURL()` - URLの変更を待つ
   - `waitForSelector()` - 要素の表示を待つ
   - `waitForTimeout()` - 時間ベースの待機（最終手段）

3. **アサーションは明確に**
   - `toBeVisible()` - 要素が表示されている
   - `toHaveText()` - テキストが一致
   - `toHaveValue()` - フォーム値が一致

4. **テストは独立させる**
   - 各テストは他のテストに依存しない
   - `beforeEach`で必要な初期化を行う

5. **データ存在チェック**
   - データがない場合の表示も考慮
   - 条件分岐で柔軟に対応

## CI/CDでの実行

GitHub Actionsでは自動的にE2Eテストが実行されます：

```yaml
- name: Run E2E tests
  run: npm run test:e2e
```

## トラブルシューティング

### テストが失敗する場合

1. **バックエンドAPIが起動しているか確認**
   ```bash
   curl http://localhost:8080/health
   ```

2. **フロントエンドが起動しているか確認**
   ```bash
   curl http://localhost:5173
   ```

3. **スクリーンショットを確認**
   テスト失敗時のスクリーンショットは `test-results/` ディレクトリに保存されます。

4. **トレースを確認**
   ```bash
   npx playwright show-trace test-results/*/trace.zip
   ```

### よくある問題

- **タイムアウトエラー**: `playwright.config.ts`でタイムアウト値を調整
- **要素が見つからない**: セレクタが正しいか確認、`page.locator().count()`でデバッグ
- **ログインできない**: デフォルトの管理者アカウント（admin/admin123）が存在するか確認

## 参考資料

- [Playwright Documentation](https://playwright.dev/)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)
