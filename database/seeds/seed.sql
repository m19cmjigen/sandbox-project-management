-- ==============================================
-- Development / Test Seed Data
-- DB-004: Seed data creation
--
-- Applying this script multiple times is safe because each INSERT uses
-- ON CONFLICT DO NOTHING. Run via: make db-seed
-- ==============================================

-- ==============================================
-- 1. Organizations (組織階層: 本部 → 部 → 課)
-- ==============================================

INSERT INTO organizations (id, name, parent_id, path, level)
OVERRIDING SYSTEM VALUE VALUES
  -- level 0: 本部
  (1,  '技術本部',       NULL, '/1/',          0),
  (7,  '営業本部',       NULL, '/7/',          0),
  (10, '管理本部',       NULL, '/10/',         0),
  -- level 1: 部
  (2,  '開発部',         1,    '/1/2/',        1),
  (5,  'インフラ部',     1,    '/1/5/',        1),
  (8,  '営業推進部',     7,    '/7/8/',        1),
  (11, '人事部',         10,   '/10/11/',      1),
  -- level 2: 課
  (3,  'Webシステム課',  2,    '/1/2/3/',      2),
  (4,  'モバイル開発課', 2,    '/1/2/4/',      2),
  (6,  'クラウド基盤課', 5,    '/1/5/6/',      2),
  (9,  '営業企画課',     8,    '/7/8/9/',      2),
  (12, '人事企画課',     11,   '/10/11/12/',   2)
ON CONFLICT (id) DO NOTHING;

SELECT setval('organizations_id_seq', GREATEST((SELECT MAX(id) FROM organizations), 12));

-- ==============================================
-- 2. Projects (Jiraプロジェクト)
-- ==============================================

INSERT INTO projects (id, jira_project_id, key, name, lead_email, organization_id, is_active)
OVERRIDING SYSTEM VALUE VALUES
  (1, 'PROJ-10001', 'WEBAPP',  '顧客Webポータル刷新',    'tanaka@example.com',  3,    true),
  (2, 'PROJ-10002', 'MOBILE',  'モバイルアプリv2開発',   'sato@example.com',    4,    true),
  (3, 'PROJ-10003', 'INFRA',   'クラウドインフラ整備',   'yamada@example.com',  6,    true),
  (4, 'PROJ-10004', 'CRM',     '顧客管理システム刷新',   'suzuki@example.com',  9,    true),
  (5, 'PROJ-10005', 'HR',      '人事管理システム導入',   'ito@example.com',     12,   true),
  (6, 'PROJ-10006', 'UNCAT',   '未分類プロジェクト',     'other@example.com',   NULL, true)
ON CONFLICT (id) DO NOTHING;

SELECT setval('projects_id_seq', GREATEST((SELECT MAX(id) FROM projects), 6));

-- ==============================================
-- 3. Issues
--
-- delay_status is auto-calculated by the trigger (calculate_issue_delay_status):
--   RED    : status_category != 'Done' AND due_date < CURRENT_DATE
--   YELLOW : status_category != 'Done' AND (due_date <= CURRENT_DATE+3 OR due_date IS NULL)
--   GREEN  : status_category = 'Done'  OR  due_date > CURRENT_DATE+3
-- All due_date values use CURRENT_DATE offsets so statuses stay accurate
-- regardless of when the seed is applied.
-- ==============================================

INSERT INTO issues (
  jira_issue_id, jira_issue_key, project_id,
  summary, status, status_category,
  due_date, assignee_name, priority, issue_type, last_updated_at
) VALUES

  -- ---- Project 1: WEBAPP (顧客Webポータル刷新) ----
  -- RED: 期限超過
  ('JIRA-W001', 'WEBAPP-1',  1, 'ログイン画面のUI刷新',             '進行中',     'In Progress', CURRENT_DATE - 14, '田中太郎',   'High',   'Story', NOW()),
  ('JIRA-W002', 'WEBAPP-2',  1, '決済フローのリファクタリング',       '進行中',     'In Progress', CURRENT_DATE - 7,  '田中太郎',   'High',   'Task',  NOW()),
  ('JIRA-W003', 'WEBAPP-3',  1, 'セッション管理の脆弱性対応',         'レビュー中', 'In Progress', CURRENT_DATE - 3,  '佐藤花子',   'High',   'Bug',   NOW()),
  -- YELLOW: 期限3日以内
  ('JIRA-W004', 'WEBAPP-4',  1, 'パスワードリセット機能の実装',       '進行中',     'In Progress', CURRENT_DATE + 1,  '佐藤花子',   'Medium', 'Story', NOW()),
  ('JIRA-W005', 'WEBAPP-5',  1, 'メール通知テンプレートの更新',       '未着手',     'To Do',       CURRENT_DATE + 2,  NULL,         'Low',    'Task',  NOW()),
  -- YELLOW: 期限未設定
  ('JIRA-W006', 'WEBAPP-6',  1, '多言語対応（i18n）の調査',           '未着手',     'To Do',       NULL,              '山田次郎',   'Low',    'Task',  NOW()),
  -- GREEN: 完了済み
  ('JIRA-W007', 'WEBAPP-7',  1, 'フロントエンド技術選定',             '完了',       'Done',        CURRENT_DATE - 30, '田中太郎',   'High',   'Task',  NOW()),
  ('JIRA-W008', 'WEBAPP-8',  1, 'APIインターフェース設計',            '完了',       'Done',        CURRENT_DATE - 20, '佐藤花子',   'High',   'Story', NOW()),
  -- GREEN: 期限に余裕あり
  ('JIRA-W009', 'WEBAPP-9',  1, '負荷テストの実施',                   '未着手',     'To Do',       CURRENT_DATE + 14, NULL,         'Medium', 'Task',  NOW()),
  ('JIRA-W010', 'WEBAPP-10', 1, 'リリースドキュメント作成',           '未着手',     'To Do',       CURRENT_DATE + 21, '田中太郎',   'Low',    'Task',  NOW()),

  -- ---- Project 2: MOBILE (モバイルアプリv2開発) ----
  -- RED: 期限超過
  ('JIRA-M001', 'MOBILE-1',  2, 'iOS/Android共通APIクライアント実装', '進行中',     'In Progress', CURRENT_DATE - 10, '鈴木一郎',   'High',   'Story', NOW()),
  ('JIRA-M002', 'MOBILE-2',  2, 'プッシュ通知の実装',                 'レビュー中', 'In Progress', CURRENT_DATE - 5,  '鈴木一郎',   'High',   'Story', NOW()),
  -- YELLOW: 期限3日以内
  ('JIRA-M003', 'MOBILE-3',  2, 'オフラインモードの対応',             '進行中',     'In Progress', CURRENT_DATE + 0,  '伊藤美咲',   'Medium', 'Story', NOW()),
  ('JIRA-M004', 'MOBILE-4',  2, 'アクセシビリティ対応',               '未着手',     'To Do',       CURRENT_DATE + 3,  NULL,         'Medium', 'Task',  NOW()),
  -- YELLOW: 期限未設定
  ('JIRA-M005', 'MOBILE-5',  2, 'ダークモード対応の検討',             '未着手',     'To Do',       NULL,              '伊藤美咲',   'Low',    'Task',  NOW()),
  -- GREEN: 完了済み
  ('JIRA-M006', 'MOBILE-6',  2, '画面設計・ワイヤーフレーム',         '完了',       'Done',        CURRENT_DATE - 40, '鈴木一郎',   'High',   'Story', NOW()),
  ('JIRA-M007', 'MOBILE-7',  2, 'プロトタイプのユーザーテスト',       '完了',       'Done',        CURRENT_DATE - 25, '伊藤美咲',   'High',   'Task',  NOW()),
  -- GREEN: 期限に余裕あり
  ('JIRA-M008', 'MOBILE-8',  2, 'App Store申請準備',                  '未着手',     'To Do',       CURRENT_DATE + 30, NULL,         'High',   'Task',  NOW()),

  -- ---- Project 3: INFRA (クラウドインフラ整備) ----
  -- RED: 期限超過
  ('JIRA-I001', 'INFRA-1',   3, 'Aurora PostgreSQLのパラメータ最適化','進行中',     'In Progress', CURRENT_DATE - 8,  '山田次郎',   'High',   'Task',  NOW()),
  -- YELLOW: 期限3日以内
  ('JIRA-I002', 'INFRA-2',   3, 'CloudWatchアラート設定',             '進行中',     'In Progress', CURRENT_DATE + 1,  '山田次郎',   'Medium', 'Task',  NOW()),
  -- YELLOW: 期限未設定
  ('JIRA-I003', 'INFRA-3',   3, 'セキュリティグループの棚卸し',       '未着手',     'To Do',       NULL,              NULL,         'Medium', 'Task',  NOW()),
  -- GREEN: 完了済み
  ('JIRA-I004', 'INFRA-4',   3, 'VPC設計・構築',                      '完了',       'Done',        CURRENT_DATE - 60, '山田次郎',   'High',   'Story', NOW()),
  ('JIRA-I005', 'INFRA-5',   3, 'ECSクラスター構築',                  '完了',       'Done',        CURRENT_DATE - 45, '渡辺健太',   'High',   'Story', NOW()),
  ('JIRA-I006', 'INFRA-6',   3, 'CI/CDパイプライン構築',              '完了',       'Done',        CURRENT_DATE - 30, '渡辺健太',   'High',   'Task',  NOW()),
  -- GREEN: 期限に余裕あり
  ('JIRA-I007', 'INFRA-7',   3, 'DR（災害復旧）計画策定',             '未着手',     'To Do',       CURRENT_DATE + 45, NULL,         'Medium', 'Story', NOW()),

  -- ---- Project 4: CRM (顧客管理システム刷新) ----
  -- RED: 期限超過
  ('JIRA-C001', 'CRM-1',     4, '顧客データ移行バッチ実装',           '進行中',     'In Progress', CURRENT_DATE - 20, '中村誠',     'High',   'Story', NOW()),
  ('JIRA-C002', 'CRM-2',     4, '重複顧客データのクレンジング',       '進行中',     'In Progress', CURRENT_DATE - 12, '中村誠',     'High',   'Task',  NOW()),
  ('JIRA-C003', 'CRM-3',     4, 'SFAとのAPI連携実装',                 'レビュー中', 'In Progress', CURRENT_DATE - 4,  '小林奈々',   'High',   'Story', NOW()),
  -- YELLOW: 期限3日以内
  ('JIRA-C004', 'CRM-4',     4, '顧客検索機能の高速化',               '進行中',     'In Progress', CURRENT_DATE + 2,  '小林奈々',   'Medium', 'Task',  NOW()),
  -- YELLOW: 期限未設定
  ('JIRA-C005', 'CRM-5',     4, '売上予測レポート機能',               '未着手',     'To Do',       NULL,              NULL,         'Low',    'Story', NOW()),
  -- GREEN: 完了済み
  ('JIRA-C006', 'CRM-6',     4, '要件定義・業務フロー分析',           '完了',       'Done',        CURRENT_DATE - 90, '中村誠',     'High',   'Story', NOW()),
  ('JIRA-C007', 'CRM-7',     4, 'DB設計・ER図作成',                   '完了',       'Done',        CURRENT_DATE - 70, '小林奈々',   'High',   'Task',  NOW()),
  -- GREEN: 期限に余裕あり
  ('JIRA-C008', 'CRM-8',     4, 'ユーザー受け入れテスト（UAT）',      '未着手',     'To Do',       CURRENT_DATE + 20, NULL,         'High',   'Task',  NOW()),
  ('JIRA-C009', 'CRM-9',     4, '本番リリース・移行計画',             '未着手',     'To Do',       CURRENT_DATE + 35, '中村誠',     'High',   'Story', NOW()),

  -- ---- Project 5: HR (人事管理システム導入) ----
  -- RED: 期限超過
  ('JIRA-H001', 'HR-1',      5, '勤怠管理モジュールのカスタマイズ',   '進行中',     'In Progress', CURRENT_DATE - 6,  '加藤さくら', 'High',   'Story', NOW()),
  -- YELLOW: 期限3日以内
  ('JIRA-H002', 'HR-2',      5, '給与計算ロジックの検証',             '進行中',     'In Progress', CURRENT_DATE + 1,  '加藤さくら', 'High',   'Task',  NOW()),
  -- YELLOW: 期限未設定
  ('JIRA-H003', 'HR-3',      5, '人事評価フォームの設計',             '未着手',     'To Do',       NULL,              NULL,         'Medium', 'Story', NOW()),
  -- GREEN: 完了済み
  ('JIRA-H004', 'HR-4',      5, 'パッケージ選定・POC',               '完了',       'Done',        CURRENT_DATE - 50, '加藤さくら', 'High',   'Task',  NOW()),
  ('JIRA-H005', 'HR-5',      5, 'データ移行計画策定',                '完了',       'Done',        CURRENT_DATE - 35, '松本浩二',   'High',   'Task',  NOW()),
  -- GREEN: 期限に余裕あり
  ('JIRA-H006', 'HR-6',      5, '従業員向けトレーニング計画',         '未着手',     'To Do',       CURRENT_DATE + 25, '松本浩二',   'Medium', 'Task',  NOW()),
  ('JIRA-H007', 'HR-7',      5, '本番稼働・並行運用期間',             '未着手',     'To Do',       CURRENT_DATE + 60, NULL,         'High',   'Story', NOW()),

  -- ---- Project 6: UNCAT (未分類プロジェクト) ----
  -- YELLOW: 期限未設定 (組織未割り当てかつ期限なし)
  ('JIRA-U001', 'UNCAT-1',   6, 'レガシーシステム調査',               '未着手',     'To Do',       NULL,              NULL,         'Low',    'Task',  NOW()),
  ('JIRA-U002', 'UNCAT-2',   6, '社内ツール統廃合の検討',             '未着手',     'To Do',       NULL,              NULL,         'Low',    'Task',  NOW()),
  -- GREEN: 完了済み
  ('JIRA-U003', 'UNCAT-3',   6, 'プロジェクト概要ヒアリング',         '完了',       'Done',        CURRENT_DATE - 15, '不明',       'Low',    'Task',  NOW())

ON CONFLICT (jira_issue_id) DO NOTHING;

SELECT setval('issues_id_seq', GREATEST((SELECT MAX(id) FROM issues), 1));

-- ==============================================
-- 4. Sync Logs (バッチ実行履歴)
-- ==============================================

-- Insert only when the table is empty to avoid accumulating duplicates on re-runs.
INSERT INTO sync_logs (sync_type, executed_at, completed_at, status, projects_synced, issues_synced, duration_seconds)
SELECT v.*
FROM (VALUES
  ('FULL' ::VARCHAR, NOW() - INTERVAL '2 days',   NOW() - INTERVAL '2 days'  + INTERVAL '8 minutes',  'SUCCESS'::VARCHAR, 6, 48, 480),
  ('DELTA'::VARCHAR, NOW() - INTERVAL '1 day',    NOW() - INTERVAL '1 day'   + INTERVAL '30 seconds', 'SUCCESS'::VARCHAR, 3, 12, 30),
  ('DELTA'::VARCHAR, NOW() - INTERVAL '12 hours', NULL::TIMESTAMP,                                     'FAILURE'::VARCHAR, 0, 0,  NULL::INTEGER)
) AS v(sync_type, executed_at, completed_at, status, projects_synced, issues_synced, duration_seconds)
WHERE NOT EXISTS (SELECT 1 FROM sync_logs LIMIT 1);
