#!/bin/bash

# GitHub Issue作成スクリプト
# 使用方法:
#   1. GitHub Personal Access Tokenを取得（repo権限が必要）
#   2. export GITHUB_TOKEN="your_token_here"
#   3. ./scripts/create-github-issues.sh

set -e

# 設定
REPO_OWNER="m19cmjigen"
REPO_NAME="sandbox-project-management"
TICKETS_DIR="tickets"

# GitHubトークンの確認
if [ -z "$GITHUB_TOKEN" ]; then
    echo "エラー: GITHUB_TOKEN環境変数が設定されていません"
    echo "使用方法:"
    echo "  export GITHUB_TOKEN='your_github_personal_access_token'"
    echo "  ./scripts/create-github-issues.sh"
    exit 1
fi

# チケットファイルのリストを取得
TICKET_FILES=$(find "$TICKETS_DIR" -name "*.md" ! -name "README.md" | sort)

echo "========================================="
echo "GitHub Issue作成スクリプト"
echo "========================================="
echo "リポジトリ: $REPO_OWNER/$REPO_NAME"
echo "チケット数: $(echo "$TICKET_FILES" | wc -l)"
echo "========================================="
echo ""

# カテゴリとラベルのマッピング
declare -A CATEGORY_LABELS
CATEGORY_LABELS["INFRA"]="infrastructure"
CATEGORY_LABELS["DB"]="database"
CATEGORY_LABELS["BACK"]="backend"
CATEGORY_LABELS["BATCH"]="batch-worker"
CATEGORY_LABELS["FRONT"]="frontend"
CATEGORY_LABELS["SEC"]="security"
CATEGORY_LABELS["TEST"]="testing"
CATEGORY_LABELS["DEPLOY"]="deployment"
CATEGORY_LABELS["DOC"]="documentation"

# 優先度ラベルのマッピング
declare -A PRIORITY_LABELS
PRIORITY_LABELS["高"]="priority: high"
PRIORITY_LABELS["中"]="priority: medium"
PRIORITY_LABELS["低"]="priority: low"

# Markdownファイルからタイトルと本文を抽出してIssue作成
create_issue_from_markdown() {
    local file=$1
    local filename=$(basename "$file")

    echo "処理中: $filename"

    # タイトルを抽出（最初の# 行）
    local title=$(grep -m 1 "^# " "$file" | sed 's/^# //')

    # 本文全体を取得
    local body=$(cat "$file")

    # カテゴリを抽出
    local category=$(echo "$filename" | cut -d'-' -f1)
    local category_label="${CATEGORY_LABELS[$category]:-enhancement}"

    # 優先度を抽出
    local priority=$(grep "^## 優先度" -A 1 "$file" | tail -1)
    local priority_label="${PRIORITY_LABELS[$priority]:-priority: medium}"

    # チケット番号を抽出
    local ticket_number=$(echo "$filename" | sed 's/_.*$//' | sed 's/\.md$//')

    # JSONペイロードを作成（jqを使用してエスケープ）
    local json_body=$(jq -n \
        --arg title "$title" \
        --arg body "$body" \
        --arg label1 "$category_label" \
        --arg label2 "$priority_label" \
        --arg label3 "ticket" \
        '{
            title: $title,
            body: $body,
            labels: [$label1, $label2, $label3]
        }')

    # GitHub APIでIssue作成
    local response=$(curl -s -X POST \
        -H "Authorization: token $GITHUB_TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        -H "Content-Type: application/json" \
        "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/issues" \
        -d "$json_body")

    # レスポンスから作成されたIssue番号を取得
    local issue_number=$(echo "$response" | jq -r '.number // empty')

    if [ -n "$issue_number" ]; then
        echo "  ✓ Issue #$issue_number 作成: $title"
    else
        local error_message=$(echo "$response" | jq -r '.message // "不明なエラー"')
        echo "  ✗ エラー: $error_message"
        echo "  レスポンス: $response" >&2
    fi

    # API Rate Limitを回避するため少し待機
    sleep 1
}

# 全てのチケットファイルを処理
for file in $TICKET_FILES; do
    create_issue_from_markdown "$file"
done

echo ""
echo "========================================="
echo "完了しました！"
echo "========================================="
echo "GitHubリポジトリを確認してください:"
echo "https://github.com/$REPO_OWNER/$REPO_NAME/issues"
