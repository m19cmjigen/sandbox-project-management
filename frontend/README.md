# Frontend (React + TypeScript)

全社プロジェクト進捗可視化プラットフォームのフロントエンドアプリケーション。

## 技術スタック

- **React** 18.2 - UIライブラリ
- **TypeScript** 5.3 - 型安全性
- **Vite** 5.1 - ビルドツール
- **Material-UI (MUI)** 5.15 - UIコンポーネントライブラリ
- **React Router** 6.22 - ルーティング
- **Zustand** 4.5 - 状態管理
- **Axios** 1.6 - HTTP通信
- **Recharts** 2.12 - チャート描画
- **ESLint** + **Prettier** - コード品質管理

## プロジェクト構造

```
frontend/
├── src/
│   ├── components/        # 共通コンポーネント
│   │   └── Layout.tsx
│   ├── pages/             # ページコンポーネント
│   │   ├── Dashboard.tsx
│   │   ├── Organizations.tsx
│   │   ├── Projects.tsx
│   │   └── Issues.tsx
│   ├── hooks/             # カスタムフック
│   ├── services/          # API通信サービス
│   ├── store/             # 状態管理（Zustand）
│   ├── types/             # TypeScript型定義
│   ├── utils/             # ユーティリティ関数
│   ├── App.tsx            # アプリケーションルート
│   ├── main.tsx           # エントリーポイント
│   ├── theme.ts           # MUIテーマ設定
│   └── index.css          # グローバルスタイル
├── index.html
├── vite.config.ts         # Vite設定
├── tsconfig.json          # TypeScript設定
├── package.json
└── README.md
```

## セットアップ

### 前提条件

- Node.js 18.0 以上
- npm 9.0 以上

### インストール

```bash
cd frontend
npm install
```

### 開発サーバーの起動

```bash
npm run dev
```

アプリケーションが `http://localhost:3000` で起動します。

## 利用可能なスクリプト

```bash
# 開発サーバー起動
npm run dev

# プロダクションビルド
npm run build

# ビルドされたアプリをプレビュー
npm run preview

# 型チェック
npm run type-check

# リント実行
npm run lint

# リント自動修正
npm run lint:fix

# コードフォーマット
npm run format

# テスト実行
npm run test

# カバレッジ付きテスト
npm run test:coverage

# テストUI起動
npm run test:ui
```

## パス aliases

TypeScriptのパスマッピングにより、以下のようにインポートできます：

```typescript
import Component from '@components/Component'
import useSomething from '@hooks/useSomething'
import { api } from '@services/api'
import { useStore } from '@store/index'
import { User } from '@types/user'
import { formatDate } from '@utils/date'
```

## API連携

バックエンドAPIとの通信は `/api` プレフィックスで行います。
Viteのプロキシ設定により、開発時は `http://localhost:8080` に転送されます。

```typescript
// 例: /api/v1/organizations へのリクエスト
axios.get('/api/v1/organizations')
```

## Material-UI テーマ

`src/theme.ts` でMUIのテーマをカスタマイズしています。

- プライマリカラー: #1976d2 (青)
- セカンダリカラー: #dc004e (ピンク)
- エラー: #f44336 (赤)
- 警告: #ff9800 (オレンジ)
- 成功: #4caf50 (緑)

## ルーティング

React Routerを使用したルーティング構成：

- `/` - ダッシュボード
- `/organizations` - 組織管理
- `/projects` - プロジェクト一覧
- `/issues` - チケット一覧

## 状態管理

Zustandを使用した軽量な状態管理：

```typescript
// store/index.ts
import create from 'zustand'

interface AppState {
  user: User | null
  setUser: (user: User) => void
}

export const useStore = create<AppState>((set) => ({
  user: null,
  setUser: (user) => set({ user }),
}))

// コンポーネントで使用
const { user, setUser } = useStore()
```

## コーディング規約

### TypeScript

- すべてのコンポーネントで型定義を使用
- `any`型の使用を最小限に
- propsインターフェースは明示的に定義

### React

- 関数コンポーネントを使用
- カスタムフックで再利用可能なロジックを抽出
- useEffect の依存配列を適切に管理

### スタイリング

- Material-UIの `sx` prop を優先的に使用
- グローバルスタイルは最小限に
- レスポンシブデザインを考慮

## ビルド最適化

Viteの設定により、以下の最適化を実施：

- コード分割（React、MUI を別チャンクに）
- ソースマップ生成
- Tree shaking
- 圧縮

## デプロイ

### ビルド

```bash
npm run build
```

`dist/` ディレクトリに静的ファイルが生成されます。

### S3 + CloudFront へのデプロイ

```bash
# ビルド
npm run build

# S3バケットにアップロード
aws s3 sync dist/ s3://your-bucket-name/ --delete

# CloudFront キャッシュ無効化
aws cloudfront create-invalidation --distribution-id YOUR_DIST_ID --paths "/*"
```

## 開発中の機能

以下の機能は後続チケットで実装予定：

- **FRONT-002**: 共通コンポーネントの実装
- **FRONT-003**: 組織階層ツリー表示
- **FRONT-004**: プロジェクト一覧表示
- **FRONT-005**: チケット詳細・フィルタ表示
- **FRONT-006**: ダッシュボードHeatmap表示
- **FRONT-007**: 組織管理画面

## トラブルシューティング

### ポートが既に使用されている

```bash
# 別のポートで起動
PORT=3001 npm run dev
```

### 型エラーが発生する

```bash
# node_modules を削除して再インストール
rm -rf node_modules package-lock.json
npm install
```

### ビルドエラー

```bash
# TypeScriptの型チェック
npm run type-check

# ESLintチェック
npm run lint
```

## 関連ドキュメント

- [SPEC.md](../SPEC.md) - プロジェクト要件定義書
- [Material-UI Documentation](https://mui.com/)
- [React Router Documentation](https://reactrouter.com/)
- [Zustand Documentation](https://github.com/pmndrs/zustand)
