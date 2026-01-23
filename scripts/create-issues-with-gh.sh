#!/bin/bash

# GitHub CLI (gh) を使用したIssue作成スクリプト
# 前提条件: gh CLI がインストールされていること
# インストール: https://cli.github.com/

set -e

TICKETS_DIR="tickets"

# GitHub CLIがインストールされているか確認
if ! command -v gh &> /dev/null; then
    echo "エラー: GitHub CLI (gh) がインストールされていません"
    echo ""
    echo "インストール方法:"
    echo "  macOS:   brew install gh"
    echo "  Linux:   https://github.com/cli/cli/blob/trunk/docs/install_linux.md"
    echo "  Windows: https://github.com/cli/cli#windows"
    exit 1
fi

# 認証確認
if ! gh auth status &> /dev/null; then
    echo "GitHub認証が必要です"
    gh auth login
fi

echo "========================================="
echo "GitHub Issue作成スクリプト (gh CLI版)"
echo "========================================="
echo ""

# カテゴリとラベルのマッピング
get_labels() {
    local filename=$1
    local category=$(echo "$filename" | cut -d'-' -f1)
    local file_path="$TICKETS_DIR/$filename"

    # 優先度を取得
    local priority=$(grep "^## 優先度" -A 1 "$file_path" | tail -1 | tr -d '[:space:]')

    local labels="ticket"

    case $category in
        INFRA) labels="$labels,infrastructure" ;;
        DB) labels="$labels,database" ;;
        BACK) labels="$labels,backend" ;;
        BATCH) labels="$labels,batch-worker" ;;
        FRONT) labels="$labels,frontend" ;;
        SEC) labels="$labels,security" ;;
        TEST) labels="$labels,testing" ;;
        DEPLOY) labels="$labels,deployment" ;;
        DOC) labels="$labels,documentation" ;;
    esac

    case $priority in
        高) labels="$labels,priority: high" ;;
        中) labels="$labels,priority: medium" ;;
        低) labels="$labels,priority: low" ;;
    esac

    echo "$labels"
}

# チケットファイルのリストを取得（ソート済み）
TICKET_FILES=$(find "$TICKETS_DIR" -name "*.md" ! -name "README.md" | sort)

echo "チケット数: $(echo "$TICKET_FILES" | wc -l)"
echo "========================================="
echo ""

# 各チケットファイルからIssue作成
for file in $TICKET_FILES; do
    filename=$(basename "$file")

    # タイトルを抽出
    title=$(grep -m 1 "^# " "$file" | sed 's/^# //')

    # ラベルを取得
    labels=$(get_labels "$filename")

    echo "作成中: $filename"
    echo "  タイトル: $title"
    echo "  ラベル: $labels"

    # GitHub Issue作成
    gh issue create \
        --title "$title" \
        --body-file "$file" \
        --label "$labels"

    echo "  ✓ 完了"
    echo ""

    # API Rate Limitを回避
    sleep 1
done

echo "========================================="
echo "全てのIssueが作成されました！"
echo "========================================="
gh issue list --limit 50
