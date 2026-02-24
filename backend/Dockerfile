# マルチステージビルド

# ビルドステージ
FROM golang:1.23-alpine AS builder

# 必要なツールのインストール
RUN apk add --no-cache git make

WORKDIR /app

# Go modulesのダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# 実行ステージ
FROM alpine:latest

# セキュリティ更新とca-certificatesのインストール
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ビルド成果物をコピー
COPY --from=builder /app/api .

# タイムゾーン設定
ENV TZ=Asia/Tokyo

# ポート公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 実行
CMD ["./api"]
