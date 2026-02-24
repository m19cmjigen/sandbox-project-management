# セキュリティ監査レポート

**監査日**: 2026-02-24
**対象バージョン**: main ブランチ（PR #22 / SEC-002 マージ後）
**監査ツール**: govulncheck、コードレビュー

---

## 1. エグゼクティブサマリー

| 項目 | 結果 |
|------|------|
| OWASP Top 10 対策 | 主要項目に対応済み |
| SQLインジェクション | 問題なし（全クエリ parameterized）|
| XSS | 問題なし（React 自動エスケープ + CSP ヘッダー）|
| CSRF | 問題なし（JWT Bearer 方式 → Cookie 非使用）|
| 認証・認可 | JWT + RBAC 実装済み |
| 依存ライブラリ脆弱性 | 1件修正、3件は Go ツールチェーン側の問題 |

---

## 2. OWASP Top 10 チェック結果

### A01 - Broken Access Control ✅
- すべての `/api/v1/**` エンドポイントに JWT 認証ミドルウェアを適用済み
- 書き込み系操作（POST/PUT/DELETE）に `RequireRole("admin")` RBAC ガードを適用
- `/health`, `/ready`, `POST /auth/login` のみ公開エンドポイント

### A02 - Cryptographic Failures ✅
- パスワードは bcrypt (cost=12) でハッシュ化
- JWT は HS256 署名（24時間有効）
- `JWT_SECRET` は環境変数で管理（本番では強力なランダム値を使用）
- HTTPS は AWS ALB/CloudFront でターミネーション

### A03 - Injection ✅
- **SQLインジェクション**: 全 DB クエリが `$N` プレースホルダを使用（sqlx/database/sql）
- JQL（Jira API）: バッチが生成する JQL にはユーザー入力を含まない
- 動的 ORDER BY / WHERE 句はホワイトリスト方式でのみ構築

```go
// 例: sort パラメータはホワイトリストマップで検証済み
validSortCols := map[string]string{
    "due_date": "i.due_date",
    ...
}
sortCol, ok := validSortCols[sortParam]  // 不正値は無視
```

### A04 - Insecure Design ✅
- バッチとAPIを分離（Jira 認証情報は API サーバーには不要）
- 最小権限原則: `viewer` ロールは読み取り専用

### A05 - Security Misconfiguration ✅
- セキュリティヘッダーを全レスポンスに付与:
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`
  - `X-XSS-Protection: 1; mode=block`
  - `Content-Security-Policy: default-src 'none'`
  - `Referrer-Policy: strict-origin-when-cross-origin`
- CORS は `CORS_ALLOWED_ORIGINS` 環境変数で制御（本番は明示的に設定必須）
- リリースモードでデフォルト `JWT_SECRET` を検出した場合は警告ログを出力

### A06 - Vulnerable and Outdated Components ⚠️ (対応済み)
**修正済み:**
- `GO-2025-3595`: `golang.org/x/net` の XSS 脆弱性 → v0.38.0 にアップグレード済み

**ツールチェーン側（アプリコードでは対応不可）:**
- `GO-2026-4341`: `net/url` メモリ枯渇 → Go 1.25.6 以上で修正
- `GO-2026-4340`: `crypto/tls` ハンドシェイク問題 → Go 1.25.6 以上で修正
- `GO-2026-4337`: `crypto/tls` セッション再開問題 → Go 1.25.7 以上で修正

**アクション**: CI/CD で定期的に `govulncheck` を実行し、Go ツールチェーンを最新に保つこと。

### A07 - Identification and Authentication Failures ✅
- JWT アクセストークン（24時間有効）
- タイミング攻撃対策: ユーザーが存在しない場合も `invalid email or password` を返す
- パスワード最小長チェック（8文字以上）
- bcrypt ハッシュで平文パスワードは保存しない
- `is_active` フラグによるアカウント無効化に対応

### A08 - Software and Data Integrity Failures ✅
- JWT 署名検証で改ざん検知
- DB への全入力はバリデーション後に parameterized query 経由でのみ書き込み

### A09 - Security Logging and Monitoring Failures ✅
- 全リクエストをリクエストログミドルウェアで記録
- バッチ実行ログを `sync_logs` テーブルに記録
- CloudWatch EMF によるバッチメトリクス（BATCH-005）
- AWS CloudTrail で Secrets Manager アクセスを記録（docs/secrets-management.md）

### A10 - Server-Side Request Forgery (SSRF) ✅
- 外部 HTTP リクエストはバッチの Jira API 呼び出しのみ
- URL は環境変数 `JIRA_BASE_URL` で固定（ユーザー入力を URL に使用しない）

---

## 3. 追加セキュリティ推奨事項

以下は現時点では未実装だが、本番運用前に対応を推奨する項目。

### 高優先度
| 推奨事項 | 理由 |
|---------|------|
| ログインエンドポイントへのレート制限実装 | ブルートフォース攻撃対策（IP ベースで 5回/分程度）|
| Go ツールチェーンを 1.25.7 以上に更新 | crypto/tls 脆弱性の修正 |
| AWS WAF の導入 | DDoS・SQLi・XSS フィルタリング |

### 中優先度
| 推奨事項 | 理由 |
|---------|------|
| パスワードリセット機能の実装 | 現状はパスワード変更手段がない |
| アクセストークンのブラックリスト（Redis）| ログアウト直後のトークン無効化 |
| 依存ライブラリの定期スキャン（CI に govulncheck 追加）| 新規脆弱性の早期発見 |
| フロントエンド依存パッケージの npm audit | npm 依存の脆弱性確認 |

### 低優先度
| 推奨事項 | 理由 |
|---------|------|
| パスワード複雑度ポリシーの強化 | 現状は最小長のみ |
| セキュリティイベントの専用ログストリーム | 認証失敗・権限エラーの集約監視 |

---

## 4. フロントエンドセキュリティ

| チェック項目 | 結果 |
|------------|------|
| XSS | React 自動エスケープにより問題なし |
| 機密情報の localStorage 保存 | JWTをlocalStorageに保存予定 → httpOnly Cookieへの変更を推奨 |
| 依存パッケージ脆弱性 | `npm audit` で定期確認を推奨 |

---

## 5. 修正済み問題一覧

| ID | 問題 | 対応 | コミット |
|----|------|------|---------|
| VUL-001 | golang.org/x/net XSS (GO-2025-3595) | v0.38.0 にアップグレード | feature/sec-003 |
| SEC-001 | 本番でデフォルトJWT_SECRETの使用 | リリースモードで警告ログ出力 | feature/sec-003 |
| SEC-002 | パスワード最小長チェック欠如 | 8文字未満は認証失敗として処理 | feature/sec-003 |
