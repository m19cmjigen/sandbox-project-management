# パフォーマンステスト (Performance Tests)

k6を使用したAPIパフォーマンステストスイートです。

## 概要

このディレクトリには、システムの非機能要件を検証するためのパフォーマンステストが含まれています。

### テスト対象

1. **API レスポンスタイム**
   - ダッシュボードAPI: < 500ms
   - プロジェクト一覧API: < 300ms
   - チケット検索API: < 500ms

2. **同時ユーザー負荷**
   - 想定: 100同時ユーザー
   - 目標: 99%のリクエストが正常応答

3. **大量データ処理**
   - 10,000+ チケット環境での動作確認

## 前提条件

### k6のインストール

```bash
# macOS
brew install k6

# Ubuntu/Debian
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Windows
choco install k6

# Docker
docker pull grafana/k6
```

### テストデータの準備

大量のテストデータが必要な場合：

```bash
# 大量データ生成スクリプト実行
cd backend
go run ./scripts/generate_large_dataset.go
```

## テストの実行

### 1. 基本的なAPIテスト

```bash
# シンプルな負荷テスト
k6 run performance/api-load-test.js

# カスタムVU数で実行
k6 run --vus 50 --duration 30s performance/api-load-test.js
```

### 2. 認証付きAPIテスト

```bash
# 認証が必要なエンドポイントのテスト
k6 run performance/api-auth-test.js
```

### 3. ダッシュボードパフォーマンステスト

```bash
# ダッシュボード特化テスト
k6 run performance/dashboard-perf-test.js
```

### 4. ストレステスト

```bash
# 段階的に負荷を上げるテスト
k6 run performance/stress-test.js
```

## テストシナリオ

### Smoke Test（疎通確認）
```bash
k6 run --vus 1 --duration 10s performance/api-load-test.js
```
- VU数: 1
- 期間: 10秒
- 目的: 基本的な動作確認

### Load Test（通常負荷）
```bash
k6 run --vus 50 --duration 5m performance/api-load-test.js
```
- VU数: 50
- 期間: 5分
- 目的: 通常運用時の性能確認

### Stress Test（高負荷）
```bash
k6 run performance/stress-test.js
```
- VU数: 50 → 100 → 200 → 100 → 50
- 期間: 15分
- 目的: システムの限界を確認

### Spike Test（突発負荷）
```bash
k6 run --vus 10 --duration 2m \
       --stage 10s:10 \
       --stage 30s:200 \
       --stage 10s:10 \
       performance/api-load-test.js
```
- 突発的な負荷スパイクへの対応確認

## 結果の読み方

### 主要メトリクス

```
     ✓ status is 200
     ✓ response time < 500ms

     checks.........................: 100.00% ✓ 50000      ✗ 0
     data_received..................: 50 MB   1.7 MB/s
     data_sent......................: 5.0 MB  167 kB/s
     http_req_blocked...............: avg=1.2ms    min=1µs   med=4µs   max=200ms p(90)=7µs   p(95)=10µs
     http_req_connecting............: avg=500µs    min=0s    med=0s    max=100ms p(90)=0s    p(95)=0s
     http_req_duration..............: avg=250ms    min=50ms  med=200ms max=2s    p(90)=400ms p(95)=500ms
       { expected_response:true }...: avg=250ms    min=50ms  med=200ms max=2s    p(90)=400ms p(95)=500ms
     http_req_failed................: 0.00%   ✓ 0          ✗ 50000
     http_req_receiving.............: avg=500µs    min=10µs  med=100µs max=50ms  p(90)=1ms   p(95)=2ms
     http_req_sending...............: avg=50µs     min=5µs   med=20µs  max=10ms  p(90)=100µs p(95)=200µs
     http_req_tls_handshaking.......: avg=0s       min=0s    med=0s    max=0s    p(90)=0s    p(95)=0s
     http_req_waiting...............: avg=249ms    min=49ms  med=199ms max=1.9s  p(90)=399ms p(95)=499ms
     http_reqs......................: 50000   1666.666667/s
     iteration_duration.............: avg=1.25s    min=1s    med=1.2s  max=3s    p(90)=1.4s  p(95)=1.5s
     iterations.....................: 50000   1666.666667/s
     vus............................: 50      min=50       max=50
     vus_max........................: 50      min=50       max=50
```

### 重要な指標

- **http_req_duration**: リクエストの応答時間（目標: p95 < 500ms）
- **http_req_failed**: 失敗率（目標: < 1%）
- **checks**: アサーション成功率（目標: 100%）
- **http_reqs**: 秒間リクエスト数（スループット）

### 合格基準

- ✅ **p95 response time** < 500ms
- ✅ **error rate** < 1%
- ✅ **success rate** > 99%
- ✅ **throughput** > 100 req/s（50 VU時）

## HTML レポート生成

```bash
# JSON形式で結果出力
k6 run --out json=results.json performance/api-load-test.js

# HTMLレポート生成（k6-to-htmlツール使用）
npm install -g k6-to-html
k6-to-html results.json results.html
```

## CI/CD統合

`.github/workflows/performance.yml` でパフォーマンステストが自動実行されます。

```yaml
name: Performance Tests

on:
  schedule:
    - cron: '0 2 * * *'  # 毎日午前2時
  workflow_dispatch:

jobs:
  performance-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install k6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
      - name: Run performance tests
        run: k6 run performance/api-load-test.js
```

## トラブルシューティング

### 高いエラー率

```
http_req_failed: 10.00%
```

**原因**:
- サーバーリソース不足
- データベース接続プール枯渇
- ネットワーク帯域不足

**対策**:
- サーバースペックの確認
- データベースコネクションプールサイズの調整
- キャッシュの導入

### 遅いレスポンスタイム

```
http_req_duration: avg=2s p(95)=5s
```

**原因**:
- N+1クエリ問題
- インデックス不足
- 非効率なクエリ

**対策**:
- SQLのEXPLAIN ANALYZEで確認
- 適切なインデックスの追加
- クエリの最適化

### メモリリーク

長時間テストでメモリ使用量が増加し続ける場合：

**対策**:
- プロファイリングツールで原因特定
- goroutineリークの確認
- データベースコネクションの適切なクローズ

## 参考資料

- [k6 公式ドキュメント](https://k6.io/docs/)
- [k6 Examples](https://k6.io/docs/examples/)
- [Performance Testing Guide](https://k6.io/docs/testing-guides/)
