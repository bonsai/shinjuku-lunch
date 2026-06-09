# デプロイ診断レポート

**作成日**: 2026-06-05
**プロジェクト**: shinjuku-lunch (新宿ランチナビ)
**依頼者**: 0501JP

## 要件

| 条件 | 詳細 |
|------|------|
| 無料 | 無料枠で運用可能 |
| コールドスタートなし | アイドル時もスリープしない |
| Go対応 | Go 1.26 / chi ビルド可能 |
| API使いやすい | REST API そのまま公開 |
| スマホ認証できる | Web UI がスマホ対応 |
| gh認証 | GitHub CLI (gh) で認証・デプロイ可能 |

## 診断結果

### 推奨: Railway

| 項目 | 評価 |
|------|------|
| 無料枠 | $5/月クレジット — 今回のAPIなら十分 |
| コールドスタート | **なし** (常時稼働) |
| Go ビルド | Dockerfile or Nixpacks 自動検知 |
| API公開 | デフォルトで `*.railway.app` ドメイン |
| スマホ対応 | Web UI 完全対応、GitHub OAuth 連携 |
| gh認証 | GitHub リポジトリ連携、PR push で自動デプロイ |

### 次点: Fly.io

- 無料3VM、常時稼働、SQLite永続化可能
- CLI必須だがスマホからも`flyctl`操作可能
- gh認証は不可（flyctl auth が必要）

### 非推奨: Render

- 15分 idle でスリープ → コールドスタート要件を満たさない

### 非推奨: Cloud Run

- 最小インスタンス0でコールドスタート、1にすると課金発生

## デプロイ手順 (Railway)

```bash
# 1. GitHub リポジトリと連携
gh repo view  # 確認

# 2. Railway ダッシュボード → New Project → Deploy from GitHub repo
#    リポジトリ選択 → 自動で Dockerfile 検出

# 3. 環境変数設定 (任意)
#    DATABASE_URL = （Neon接続文字列、未設定ならseed.jsonモード）

# 4. 自動デプロイ
#    main ブランチに push するたび自動デプロイ
```

## 参考ファイル

- `deploy/deployment-targets.json` — 全PaaS比較データ
- `api/Dockerfile` — 単一ステージ、alpineベース
- `api/.air.toml` — 開発用ホットリロード設定
