# DB-004: シードデータの作成 設計書

## 概要

開発・テスト用のシードデータを作成する。
組織階層・プロジェクト・チケット（RED/YELLOW/GREEN混在）・同期ログを投入する。

## 作成するファイル

| ファイル | 内容 |
|---------|-----|
| `database/seeds/seed.sql` | シードデータSQL本体 |
| `database/seeds/apply.sh` | 実行シェルスクリプト |
| `Makefile` | `db-seed` ターゲット追加 |
| `tickets/DB-004_seed-data-creation.md` | 完了マーク追加 |

## データ設計

### Organizations (3階層)
- level 0 (本部): 技術本部, 営業本部, 管理本部
- level 1 (部): 開発部, インフラ部, 営業推進部, 人事部
- level 2 (課): Webシステム課, モバイル開発課, クラウド基盤課, 営業企画課, 人事企画課

### Projects (6件)
- 5件: 各課に紐付け
- 1件: 未分類（organization_id = NULL）

### Issues (約50件)
遅延ステータスはトリガーが自動計算（CURRENT_DATE基準）:
- RED: status_category != 'Done' AND due_date < CURRENT_DATE
- YELLOW: status_category != 'Done' AND due_date <= CURRENT_DATE+3 または due_date IS NULL
- GREEN: 完了済み または 十分余裕のある期限

### Sync Logs (3件)
- SUCCESS × 2, FAILURE × 1
