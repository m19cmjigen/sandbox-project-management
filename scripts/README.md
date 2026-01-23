# GitHub Issue作成スクリプト

このディレクトリには、`tickets/`ディレクトリ内のMarkdownファイルをGitHub Issueとして一括作成するスクリプトが含まれています。

## 方法1: GitHub CLI (gh) を使用【推奨】

最も簡単で確実な方法です。

### 前提条件

GitHub CLI (`gh`) がインストールされている必要があります。

#### インストール

**macOS:**
```bash
brew install gh
```

**Ubuntu/Debian:**
```bash
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh
```

**Windows:**
```bash
winget install --id GitHub.cli
```

または: https://github.com/cli/cli#installation

### 使用方法

1. **GitHub認証**
```bash
gh auth login
```

ブラウザが開くので、指示に従って認証してください。

2. **スクリプト実行**
```bash
cd /path/to/sandbox-project-management
./scripts/create-issues-with-gh.sh
```

3. **確認**
```bash
gh issue list
```

または、ブラウザで確認:
https://github.com/m19cmjigen/sandbox-project-management/issues

## 方法2: GitHub API (curl) を使用

`gh` CLIが使えない場合は、GitHub APIを直接使用します。

### 前提条件

GitHub Personal Access Tokenが必要です。

#### Personal Access Token の取得

1. GitHubにログイン
2. Settings → Developer settings → Personal access tokens → Tokens (classic)
3. "Generate new token (classic)" をクリック
4. 以下の権限を選択:
   - `repo` (すべて)
5. トークンを生成してコピー

### 使用方法

1. **環境変数にトークンを設定**
```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxx"
```

2. **スクリプト実行**
```bash
cd /path/to/sandbox-project-management
./scripts/create-github-issues.sh
```

3. **確認**

ブラウザで確認:
https://github.com/m19cmjigen/sandbox-project-management/issues

## 作成されるIssueの情報

### タイトル

各Markdownファイルの最初の `# ` 行がタイトルになります。

例: `# INFRA-001: AWS アカウント・環境セットアップ`

### 本文

Markdownファイルの全内容が本文として使用されます。

### ラベル

自動的に以下のラベルが付与されます:

#### カテゴリラベル
- `INFRA-*` → `infrastructure`
- `DB-*` → `database`
- `BACK-*` → `backend`
- `BATCH-*` → `batch-worker`
- `FRONT-*` → `frontend`
- `SEC-*` → `security`
- `TEST-*` → `testing`
- `DEPLOY-*` → `deployment`
- `DOC-*` → `documentation`

#### 優先度ラベル
- 優先度: 高 → `priority: high`
- 優先度: 中 → `priority: medium`
- 優先度: 低 → `priority: low`

#### 共通ラベル
- すべてのIssueに `ticket` ラベルが付与されます

## トラブルシューティング

### gh コマンドが見つからない

```bash
which gh
```

見つからない場合は、上記のインストール手順に従ってください。

### 認証エラー

```bash
gh auth status
```

認証されていない場合:
```bash
gh auth login
```

### API Rate Limit エラー

GitHubのAPI制限に達した場合は、1時間待ってから再実行してください。

認証済みユーザーの場合、1時間あたり5000リクエストまで可能です。

### 権限エラー

Personal Access Tokenに `repo` 権限があることを確認してください。

## 手動でIssueを作成する場合

スクリプトを使用せず、手動で作成する場合:

1. https://github.com/m19cmjigen/sandbox-project-management/issues/new にアクセス
2. `tickets/` ディレクトリ内の各Markdownファイルの内容をコピー＆ペースト
3. 適切なラベルを付与
4. "Submit new issue" をクリック

37個のチケットがあるため、手動作成は非常に時間がかかります。

## 一括削除（必要な場合）

誤って作成したIssueを一括削除する場合:

```bash
# すべてのオープンIssueを表示
gh issue list --limit 100

# 特定のラベルのIssueをクローズ
gh issue list --label "ticket" --json number --jq '.[].number' | \
  xargs -I {} gh issue close {}
```

## チケット一覧

全37チケットの一覧と詳細は `tickets/README.md` を参照してください。

## 関連ドキュメント

- [tickets/README.md](../tickets/README.md) - チケット一覧
- [GitHub CLI Documentation](https://cli.github.com/manual/)
- [GitHub REST API - Issues](https://docs.github.com/en/rest/issues/issues)
