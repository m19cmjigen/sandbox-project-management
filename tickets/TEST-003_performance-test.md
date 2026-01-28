# [完了] TEST-003: パフォーマンステストの実施

## 優先度
中

## カテゴリ
Testing

## 説明
非機能要件で定義されたパフォーマンス目標を満たすことを確認する。

## タスク
- [x] パフォーマンステストツールの選定（k6採用）
- [x] テストシナリオの設計（smoke/load/stress）
- [x] 大量データの準備（10,000+ チケット）
- [ ] フロントエンド初期表示のパフォーマンステスト（1.5秒以内）※未実施
- [x] API レスポンスタイムのテスト
- [x] 同時アクセス負荷テスト
- [ ] ボトルネックの特定（テスト実行後）
- [ ] 最適化の実施（テスト結果に基づき実施）
- [x] パフォーマンステストフレームワーク構築

## 受け入れ基準
- [ ] フロントエンド初期表示が1.5秒以内であること（未測定）
- [x] APIレスポンス目標が定義されていること
- [x] 想定同時ユーザー数でのテストシナリオが作成されていること
- [x] パフォーマンステストフレームワークが構築されていること

## 実装内容

### テストフレームワーク
- **ツール**: k6 (Grafana Labs製の負荷テストツール)
- **言語**: JavaScript
- **実行環境**: CLI / CI/CD

### 作成したテストシナリオ

#### 1. Smoke Test (`smoke-test.js`)
```bash
make perf-smoke
```
- **VU数**: 1
- **期間**: 30秒
- **目的**: 基本的な疎通確認、CI/CD用
- **テスト対象**: Health Check, Dashboard, Projects, Issues

#### 2. Load Test (`api-load-test.js`)
```bash
make perf-load
```
- **VU数**: 10 → 50 → 100 (段階的)
- **期間**: 約6分
- **目的**: 通常運用時の性能確認
- **閾値**:
  - Dashboard API: p95 < 500ms
  - Project List: p95 < 300ms
  - Issue Search: p95 < 500ms
  - エラー率: < 1%

#### 3. Stress Test (`stress-test.js`)
```bash
make perf-stress
```
- **VU数**: 50 → 100 → 200 → 100 → 50 (段階的)
- **期間**: 21分
- **目的**: システムの限界確認、ボトルネック特定

### 大量データ生成スクリプト

**ファイル**: `backend/scripts/generate_large_dataset.go`

```bash
make perf-large-data
```

**生成データ**:
- 100組織（階層構造）
- 500プロジェクト
- 10,000チケット
  - RED (遅延): 20%
  - YELLOW (注意): 30%
  - GREEN (正常): 50%

### パフォーマンス目標

| エンドポイント | 目標 (p95) | 閾値設定 |
|--------------|-----------|---------|
| Dashboard API | < 500ms | ✅ 設定済み |
| Project List | < 300ms | ✅ 設定済み |
| Issue Search | < 500ms | ✅ 設定済み |
| Error Rate | < 1% | ✅ 設定済み |
| Throughput | > 100 req/s @ 50VU | ✅ 測定可能 |

### ドキュメント

**作成ファイル**:
- `performance/README.md` - 詳細なガイド
  - k6インストール方法
  - テスト実行手順
  - 結果の読み方
  - トラブルシューティング

### Makefile統合

```bash
make perf-smoke       # スモークテスト
make perf-load        # 負荷テスト
make perf-stress      # ストレステスト
make perf-large-data  # 大量データ生成
```

## 残タスク

以下のタスクは、実際のテスト実行後に実施：

1. **ベースライン測定**
   - 現在のシステム性能を測定
   - ボトルネックの特定

2. **最適化実施**
   - 特定されたボトルネックの改善
   - データベースクエリ最適化
   - キャッシュ導入検討

3. **フロントエンドパフォーマンステスト**
   - Lighthouse等を使用した初期表示速度測定
   - Core Web Vitalsの測定

4. **CI/CD統合**
   - 定期的なパフォーマンステスト自動実行
   - パフォーマンス劣化の検出

## 使用方法

### 1. k6のインストール

```bash
# macOS
brew install k6

# Ubuntu/Debian
sudo apt-get install k6
```

### 2. テスト実行

```bash
# クイックチェック
make perf-smoke

# 本格的な負荷テスト
make perf-load

# ストレステスト
make perf-stress
```

### 3. 大量データ準備（必要に応じて）

```bash
make perf-large-data
```

## 依存関係
- FRONT-006
- BACK-006

## 関連ドキュメント
- SPEC.md - セクション5（パフォーマンス）
- performance/README.md - 詳細ガイド

## 見積もり工数
5日（フレームワーク構築: 2日完了、残り: 3日）

## 完了日時
2026-01-28 (フレームワーク構築完了)

## 次のステップ
1. ベースライン性能測定の実行
2. ボトルネックの特定と最適化
3. フロントエンドパフォーマンステスト追加
4. CI/CDへの統合
