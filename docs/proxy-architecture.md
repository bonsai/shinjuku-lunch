# Proxy Architecture: Next.js → Go API のリクエスト転送

## 概要

新宿ランチナビでは **Next.js (frontend)** と **Go API (backend)** の2プロセス構成をとる。
フロントエンドの `proxy.ts` (Next.js 16 の Proxy 機構) が `/api` 以下のリクエストを Go バックエンドに転送する役割を担う。

```
  Browser
    │
    ├─ http://localhost:3000/          → Next.js がSSR/CSRでページを返す
    │
    └─ http://localhost:3000/api/*     → Next.js Proxy が Go に転送
           │
           └─ http://localhost:8080/api/*  (Go API)
                  │
                  └─ PostgreSQL (Neon)
```

## proxy.ts の動作

```ts
const API_URL = process.env.API_URL ?? "http://localhost:8080"

export function proxy(request: NextRequest) {
  if (request.nextUrl.pathname.startsWith("/api")) {
    const dest = new URL(request.nextUrl.pathname + request.nextUrl.search, API_URL)
    return NextResponse.rewrite(dest)
  }
}

export const config = {
  matcher: "/api/:path*",
}
```

### 3つの要素

| 要素 | 役割 |
|------|------|
| `config.matcher` | Proxy が起動するパスを `/api/:path*` に限定。それ以外のリクエストは素通り |
| `proxy()` 関数 | マッチしたリクエストをインターセプトし、rewrite で転送先 URL を差し替え |
| `NextResponse.rewrite()` | **サーバーサイドでの内部リクエスト書き換え**。ブラウザの URL は変更されず、レスポンスだけが差し替わる |

rewrite は HTTP リダイレクト (302) とは異なり、ブラウザから見ると Next.js のオリジン (`localhost:3000`) と通信しているように見える。CORS の問題が発生しない。

## 2つのAPI URL設定

このプロジェクトには API のベース URL を指定する箇所が2箇所ある。

### 1. proxy.ts — `API_URL` (環境変数 `API_URL`)

| 対象 | 値 |
|------|-----|
| 開発時デフォルト | `http://localhost:8080` |
| Render 本番 | 未設定（同一サービス内通信であれば internal URL） |

`proxy.ts` は **サーバーサイド** でのみ動作する。serve
r から Go へのリクエストに使われる。

### 2. api.ts — `NEXT_PUBLIC_API_URL` (環境変数 `NEXT_PUBLIC_API_URL`)

| 対象 | 値 |
|------|-----|
| 開発時デフォルト | `http://localhost:8080` |
| Render 本番 | `shinjuku-lunch-api` の internal host |

`NEXT_PUBLIC_` 接頭辞により、この値は **ブラウザのクライアント JS**
にもバンドルされる。サーバーコンポーネント (SSR) からの fetch は proxy を経由せず直接 Go に向かう。

### 使い分けの理由

| レンダリング方式 | リクエスト送信元 | 経路 |
|----------------|----------------|------|
| **SSR** (Server Component) | Next.js サーバー | `api.ts` → Go に直接 fetch |
| **CSR** (Client Component) | ブラウザ | `fetch("/api/...")` → 同梱 JS → Next.js (port 3000) → `proxy.ts` → Go (port 8080) |

Server Component からのデータ取得は proxy を経由する必要がない。そこで `NEXT_PUBLIC_API_URL`
を使って直接 Go API を叩く。

Client Component (`"use client"`) からの fetch は相対パス `/api/...`
で書かれている。これが Next.js サーバーに届き、proxy.ts が Go に転送する。

## リクエストフロー図解

### 開発時のトップページ表示 (SSR)

```
1. ブラウザ GET /
2. Next.js Server Component がレンダリングを開始
3. getRestaurants() を実行
   → api.ts が NEXT_PUBLIC_API_URL (= http://localhost:8080) を使って
      http://localhost:8080/api/restaurants に直接 fetch
4. Go API が PostgreSQL からデータを取得 → JSON を返却
5. Server Component が HTML を生成 → ブラウザに返却
```

### ブラウザ上のフィルター操作 (CSR)

```
1. ユーザーが「エリア=歌舞伎町」を選択
2. RestaurantList (Client Component) が
   fetch("/api/restaurants?area=歌舞伎町") を実行
3. ブラウザが http://localhost:3000/api/restaurants?area=歌舞伎町 にリクエスト
4. リクエストが Next.js サーバーに到着
5. proxy.ts の config.matcher が /api/:path* にマッチ
6. proxy() 関数が実行される
   → dest = new URL("/api/restaurants?area=歌舞伎町", "http://localhost:8080")
   → NextResponse.rewrite(dest) で内部転送
7. Next.js サーバーが http://localhost:8080/api/restaurants?area=歌舞伎町 にリクエスト
8. Go API がレスポンスを返す
9. Next.js がそのレスポンスをブラウザに中継
10. Client Component が JSON を受け取り、UI を更新
```

## なぜ proxy が必要か

proxy を使わない場合、Client Component から Go API に直接 fetch するにはブラウザが
`http://localhost:8080/api/...` にアクセスすることになる。これは **異なるオリジン**
へのリクエストとなり、CORS の設定が別途必要になる。

proxy によりすべてのリクエストが `http://localhost:3000` (
Next.js のオリジン) に統一されるため、CORS 不要でシンプルに運用できる。

## Render デプロイ時の動作

`render.yaml` では frontend と API が別サービスとしてデプロイされる。

```yaml
services:
  - type: web  # Go API (port 8080)
    name: shinjuku-lunch-api
    ...

  - type: web  # Next.js (port 10000)
    name: shinjuku-lunch-frontend
    ...
    envVars:
      - key: NEXT_PUBLIC_API_URL
        fromService:
          name: shinjuku-lunch-api
          type: web
          property: host  # https://shinjuku-lunch-api.onrender.com
```

Render は `NEXT_PUBLIC_API_URL` に API サービスの内部ホスト名を自動設定する。
Server Component の fetch はこの URL を通じて同じ Render ネットワーク内で API と通信する。

Client Component の fetch (`/api/...`) は Next.js サーバー (port 10000) に向かい、
proxy.ts がこれを API 内部ホストに rewrite する。この rewrite 先は `API_URL`
環境変数で指定する（未指定時は `http://localhost:8080` になるため、
本番では Render のダッシュボードで明示的に設定が必要）。

## 参考: 関連ファイル

| ファイル | 役割 |
|---------|------|
| `frontend/proxy.ts` | Next.js 16 Proxy。リクエスト転送のエントリポイント |
| `frontend/src/lib/api.ts` | API クライアント。fetch ラッパー＋型付き関数群 |
| `frontend/src/lib/types.ts` | Go API のモデルと完全一致する TypeScript 型定義 |
| `frontend/next.config.ts` | Next.js 設定（本Proxyに関しては特別な設定不要） |
| `run_all.ps1` | 開発起動スクリプト（Next.js と Go を同時起動） |
| `render.yaml` | Render デプロイ設定（frontend + API 2サービス構成） |
