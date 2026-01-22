# 全社プロジェクト進捗可視化プラットフォーム 要件定義書

## 1. プロジェクト概要
目的: 組織ごとに分散しているJiraプロジェクトの進捗状況（特に納期遅延）を一元管理し、経営層・PMO・管理職が早期に対策を打てる状態にする。
ターゲット: 経営層、PMO、部門長、プロジェクトマネージャー。
開発体制: 社内スクラッチ開発。

## 2. システムアーキテクチャ・技術スタック
Jira CloudのAPIレート制限を回避し、かつ高速なレスポンスを実現するため、データを中間DBにキャッシュする**「非同期収集型アーキテクチャ」**を採用する。

### 2.1 技術選定
領域	技術要素	選定理由
Frontend	TypeScript, React	型安全性による品質担保と、SPAによる高速な画面遷移。UIコンポーネントライブラリ（MUIやChakra UI等）の活用推奨。
Backend	Go (Golang)	静的型付けによる堅牢性と、Goroutineによる大量チケットデータの並行取得・処理に最適。
Infra (AWS)	AWS Fargate (ECS) / Lambda	サーバー管理コストの削減。バッチ処理はFargate、APIはLambdaまたはFargate。
Database	Amazon Aurora (PostgreSQL)	複雑な組織階層データと大量のチケットデータを扱うため、リレーショナルDBを採用。
IaC	Terraform / AWS CDK	インフラ構成のコード管理（推奨）。

### 2.2 システム構成図 (概念)
コード スニペット
[User (Browser)] --(HTTPS)--> [CloudFront + S3 (React App)]
       |
       +--(API Request)--> [ALB / API Gateway]
                                |
                          [Backend API (Go)] <--(Read/Write)--> [Aurora DB]
                                |
                          [Batch Worker (Go)] --(REST API)--> [Jira Cloud]

## 3. 機能要件 (Functional Requirements)
### 3.1 データ収集機能 (Batch)
収集対象: Jira Cloud上の全プロジェクト（または指定JQL範囲）。
実行頻度:
Full Sync: 1日1回（深夜）。全データを再取得・整合性チェック。
Delta Sync: 1時間に1回。updated >= -1h の条件で変更分のみ取得・更新。
Go言語による実装ポイント:
各プロジェクトのチケット取得をGoroutineで並列化し、処理時間を短縮すること。
Jira APIのレート制限（Rate Limiting）ヘッダーを監視し、429エラー時はExponential Backoff（指数関数的待機）でリトライするロジックを実装すること。
### 3.2 データ正規化・判定ロジック
収集したデータを以下のルールで統一規格化し、DBへ保存する。

#### ① ステータス正規化
Jiraの statusCategory を利用し、以下の3つに分類。
未着手 (To Do)
進行中 (In Progress)
完了 (Done)
② 遅延 (Delay) 判定ロジック 各チケットに対し、以下の優先順位でステータスを付与する。
🔴 RED (遅延): Status != Done かつ DueDate < Today（納期過ぎ）
🟡 YELLOW (注意): Status != Done かつ Today <= DueDate <= Today + 3days（3日以内期限）
🟡 YELLOW (設定不備): Status != Done かつ DueDate IS NULL（期限未設定）
🟢 GREEN (正常): 上記以外

### 3.3 組織・プロジェクト管理機能
組織マスタ管理:
全社の組織ツリー（本部 - 部 - 課）を作成・編集する機能。
プロジェクト紐付け:
Jiraから新規検出されたプロジェクトを「未分類」リストに表示。
管理者が手動で「組織」にドラッグ＆ドロップ等で紐付けを行う。
支援機能: プロジェクトリーダーのEmailドメイン等から所属組織を推測し、デフォルト値として提示する。

### 3.4 ダッシュボード機能 (Frontend)
① 全社サマリ (Heatmap View)
組織階層をツリー表示し、各組織の「遅延プロジェクト率」を色分け表示（赤・黄・緑）。
クリックで下位組織へドリルダウン。
② プロジェクト一覧
選択した組織配下のプロジェクトカードを表示。
各カードには「遅延チケット数」「期限切れ間近チケット数」を表示。
③ 遅延チケット詳細
フィルタリング機能: 「遅延のみ」「期限未設定のみ」等のフィルタを提供。
Jiraへのディープリンクボタンを設置。

## 4. データベース設計指針 (Schema Concept)
PostgreSQLを想定した主要テーブル構成案。
Organizations (組織マスタ)
id, name, parent_id, path (階層検索用)
Projects (Jiraプロジェクト)
jira_project_id, key, name, lead_account_id, organization_id (FK)
Issues (チケット情報)
jira_issue_id, project_id (FK), summary, status, status_category, due_date, assignee_name, delay_status (Red/Yellow/Green), last_updated_at
SyncLogs (バッチ実行ログ)
executed_at, status, details

## 5. 非機能要件 (Non-Functional Requirements)
セキュリティ:
Jiraとの通信にはOAuth 2.0 (3LO) または、セキュアに管理されたAPI Tokenを使用。
AWS Secrets Managerを利用して認証情報を管理。
パフォーマンス:
フロントエンドの初期表示（First Contentful Paint）は1.5秒以内。
一覧取得APIはページネーションまたは仮想スクロールを実装し、大量データ時もブラウザをフリーズさせない。
拡張性:
将来的に「工数予実（Time Tracking）」や「品質分析」機能を追加できるよう、Issuesテーブルは拡張可能な設計とする。
