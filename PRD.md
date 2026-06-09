# PRD: 新宿ランチナビ (Shinjuku Lunch Navigator)

## 1. プロダクト概要

新宿・歌舞伎町・大久保エリアのランチ情報を収集・可視化・探索する PWA。
Next.js によるモダンな Web アプリで店舗を発見・記録でき、Neon DB でデータを管理する。

## 2. ターゲットユーザー

- 新宿勤務・通学で毎日のランチ選びに悩む人
- コスパの良い穴場を探したい人
- エリアやジャンルから直感的に店を選びたい人

## 3. 技術スタック

| レイヤー | 技術 | 役割 |
|---------|------|------|
| **DB** | Neon (PostgreSQL) | 店舗・レビューデータの保存、サーバーレス |
| **バックエンド** | Go (net/http or chi) | REST API、Neon 接続、認証不要・匿名投稿 |
| **フロントエンド** | Next.js (React 19 / App Router) | SSR・CSR・ルーティング・UI |
| **PWA** | Next.js (PWA 対応予定) | オフライン対応・ホーム画面インストール |
| **認証** | なし | 全ユーザー匿名、投稿・編集に認証不要 |

## 4. データモデル

> 認証不要・匿名投稿のため `users` テーブルは持たない。
> すべての投稿は同一 namespace で管理される。

### 4.1 生データのフィールド（既存 raw データより抽出）

ランチ日記_雛形.md から同定した属性:
- `店名` — 店舗名
- `エリア` — エリア分類（新宿 / 歌舞伎町 / 西新宿 / 新宿三丁目 / 大久保 / その他）
- `ジャンル` — 料理ジャンル（韓国料理 / タイ料理 / 和食 / カレー / インド / 牛丼 / etc.）
- `注文` — 注文メニュー
- `価格` — 価格（円）
- `評価` — ★ 評価
- `一言` — コメント
- `再訪` — 再訪フラグ（したい / しない）

kabukicho_thai_research.md から追加で同定した属性:
- `徒歩_分` — 最寄駅からの徒歩分数
- `最寄駅` — 最寄駅（西武新宿駅 / JR新宿駅 / 新宿三丁目駅 / 大久保駅）
- `食べログ評価` — 食べログスコア
- `営業時間` — ランチ営業時間
- `住所` — 所在地
- `座標` — GeoJSON 座標（経度・緯度）
- `URL` — 食べログ / ホットペッパーリンク
- `特徴` — フリーテキストメモ

### 4.2 Neon DB スキーマ (PostgreSQL)

```sql
CREATE TABLE areas (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE          -- 新宿 / 歌舞伎町 / 西新宿 / 新宿三丁目 / 大久保 / その他
);

CREATE TABLE genres (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE          -- 韓国料理 / タイ料理 / 和食 / カレー / インド / 牛丼 / etc.
);

CREATE TABLE restaurants (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,                -- 店名
  area_id INTEGER REFERENCES areas(id),
  genre_id INTEGER REFERENCES genres(id),
  address TEXT,
  station TEXT,                      -- 最寄駅
  walk_min INTEGER,                  -- 徒歩分数
  latitude DOUBLE PRECISION,         -- 緯度
  longitude DOUBLE PRECISION,        -- 経度
  business_hours TEXT,               -- 営業時間
  url_tabelog TEXT,
  url_hotpepper TEXT,
  notes TEXT,                        -- 特徴・メモ
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE lunch_logs (
  id SERIAL PRIMARY KEY,
  restaurant_id INTEGER REFERENCES restaurants(id),
  menu TEXT,                         -- 注文メニュー
  price INTEGER,                     -- 価格（円）
  rating INTEGER CHECK (rating BETWEEN 1 AND 5),  -- ★評価
  comment TEXT,                      -- 一言
  revisit BOOLEAN,                   -- 再訪したい？
  visited_date DATE DEFAULT CURRENT_DATE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## 5. 機能一覧

### MVP (Phase 1)

| # | 機能 | 説明 |
|---|------|------|
| 1 | 店舗一覧 | Next.js SSR でレストラン一覧を表示。エリア・ジャンル・価格フィルター対応 |
| 2 | 店舗詳細 | 店名、ジャンル、価格帯、評価、徒歩分数、マップリンクを表示 |
| 3 | ランチ記録登録 | フォームから行った店・メニュー・評価を記録 → Neon DB に保存 |
| 4 | フィルター | エリア・ジャンル・価格帯で店舗を絞り込み、動的にリスト更新 |
| 5 | レスポンシブ | PC / タブレット / スマホに対応した Tailwind CSS デザイン |

### Phase 2

| # | 機能 | 説明 |
|---|------|------|
| 6 | オフライン対応 | PWA Service Worker でキャッシュ、オフラインでも閲覧可能 |
| 7 | インストール | ホーム画面に追加対応 |
| 8 | おすすめランチ提案 | 今まで行っていないジャンル・エリアを提案 |
| 9 | 統計ダッシュボード | 月間ランチ費、よく行くジャンル、行った店舗マップヒートマップ |

### Phase 3

| # | 機能 | 説明 |
|---|------|------|
| 10 | 写真アップロード | メニュー写真を添付 |
| 11 | AI ランチ提案 | ChatGPT API と連携し「今日の気分」から提案 |

## 6. アーキテクチャ

```
[Next.js Client (SSR/CSR)]
  │ Next.js App Router (React 19)
  │ Proxy: /api → Go API
  │
  ├─→ Go API (chi router)
  │       │
  │       └─→ Neon (PostgreSQL)
  │
  └─→ (offline) Service Worker キャッシュ
```

## 7. 初期シードデータ

既存 raw データから投入する店舗一覧:

### 歌舞伎町エリア
- TOP トッポギ（韓国料理）
- トルコケバブ屋台（¥390）
- 焼肉 ナンバーワン（¥600）
- バンタイ（タイ料理、徒歩2分、食べログ3.60）
- サームロット（タイ料理、徒歩6分）
- ゲウチャイ（タイ料理、徒歩7分）
- ランブータン（タイ料理、徒歩7分）
- バンコクスパイス（タイ料理、徒歩8分）
- モモタイ（タイ料理、徒歩8分、24h）
- チャオサイゴン パリバール（アジア料理）
- マサラステーション（インドカレー）

### 大久保エリア
- でじにらんど 大久保店（韓国料理、¥660）
- うま煮や（和食、¥700）
- ハレルヤ（韓国料理、¥500）
- 小さなカレー家（カレー、¥950）
- すき家 大久保二丁目店（牛丼、¥450〜）

## 8. 非機能要件

| 項目 | 要件 |
|------|------|
| オフライン | Service Worker キャッシュ、直近閲覧データは IndexedDB に保存 |
| 初回ロード | Next.js SSR による高速初期表示 |
| レスポンシブ | PC / タブレット / スマホ に対応（Tailwind CSS） |
| DB 接続 | Neon サーバーレス（Go: lib/pq または pgx） |
| シード | `seed.json` または `raw/` の md から `import_from_raw` で投入 |

## 9. OpenAPI 仕様

認証不要・匿名の REST API。
全エンドポイントは `/api` 配下、JSON 応答。

| Method | Path | 説明 |
|--------|------|------|
| GET | /api/restaurants | 店舗一覧（フィルタ: ?area=&genre=&price_max=） |
| GET | /api/restaurants/:id | 店舗詳細 +  lunch_logs 一覧 |
| POST | /api/lunch-logs | ランチログを投稿（匿名） |
| GET | /api/areas | エリア一覧 |
| GET | /api/genres | ジャンル一覧 |

詳細は `api/openapi.yaml` 参照。

## 10. ディレクトリ構成案

```
shinjuku-lunch/
├── PRD.md
├── neon/
│   ├── schema.sql            # DDL
│   └── seed.json             # 初期データ（raw/md から抽出済みJSON）
├── api/
│   ├── main.go               # エントリ (chi)
│   ├── openapi.yaml          # OpenAPI 3.1 定義
│   ├── handler/
│   │   ├── restaurants.go    # GET /api/restaurants
│   │   ├── lunch_logs.go     # GET/POST /api/lunch-logs
│   │   ├── areas.go          # GET /api/areas
│   │   └── genres.go         # GET /api/genres
│   ├── db/
│   │   └── db.go             # DB 接続管理
│   ├── model/
│   │   └── types.go          # データ構造
│   └── go.mod
├── frontend/
│   ├── proxy.ts              # /api → Go へのプロキシ
│   ├── src/
│   │   ├── app/
│   │   │   ├── page.tsx              # トップ（店舗一覧）
│   │   │   ├── layout.tsx
│   │   │   └── restaurants/[id]/
│   │   │       ├── page.tsx          # 店舗詳細（SSR）
│   │   │       └── client.tsx        # 詳細クライアントコンポーネント
│   │   ├── components/
│   │   │   ├── filter-bar.tsx        # フィルターUI
│   │   │   ├── restaurant-card.tsx   # 店舗カード
│   │   │   ├── restaurant-list.tsx   # 店舗一覧（クライアント）
│   │   │   ├── lunch-log-form.tsx    # ランチ記録フォーム
│   │   │   └── stars.tsx            # ★表示
│   │   └── lib/
│   │       ├── api.ts               # API クライアント
│   │       └── types.ts             # 型定義
│   ├── package.json
│   └── next.config.ts
├── scripts/
│   └── import_from_raw.go    # raw/ データを Neon に投入
├── run_all.ps1               # 開発起動スクリプト
└── render.yaml               # Render デプロイ設定
```

## 11. 今後の検討課題

- PWA 対応（Service Worker + マニフェスト）
- 地図表示（Leaflet / MapLibre の導入検討）
- 匿名投稿のためスパム対策（簡易レートリミット or Cloudflare Turnstile）
