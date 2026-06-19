# ADR-002: 界隈ナビ — ビジョン拡張: リンカー・ラッパー SNS

## Status: Proposed
## Date: 2026-06-18
## Supersedes: ADR-001 (拡張・補完)

---

## Context (背景)

ADR-001 で技術骨格を整理した。さらに以下のようなビジョンが共有された:

| キーワード | 意味 |
|-----------|------|
| **SNS/BBS/記録帳** | SNSというより「個人の記録」が本質。コミュニティはその延長 |
| **自慢したい** | 行った・集めた・知ってる を可視化して自慢するモチベーション |
| **シール帳** | 集めた証 (行った店・イベント) がコレクションになっていく感覚 |
| **小さいものは JS/DB で完結** | 個人用ツールレベルならフルスタック不要 |
| **5人以上 → Rails 発行** | 小規模は軽量で、閾値越えたらスケール |
| **界隈が発展する物語** | データが積み上がるほど町が面白くなる |
| **バックも発展** = バックエンドも共に成長 | 技術も一緒にスケールする設計 |
| **X/Instagram を束ねるリンカー/ラッパー** | 既存SNSを否定しない。ハブとして機能 |

---

## Decision (決定)

### D5: プロダクト本質は「シール帳」である

**これはSNSではない。街のシール帳である。**

SNSは「人の投稿を見る」が中心。違う。
ここは「自分の行ったところを集める」が中心。

```
SNS のタイムライン:    他人の投稿 → 他人の投稿 → 他人の投稿
シール帳のタイムライン: 自分が行った → 自分が行った → 自分が行った
                                           ↳ 他人のも見えるけど、自分のが主語
```

**UI/UX に影響する設計指針:**

- マイページ = シール帳 (行った数、コレクション、バッジ)
- 他人の投稿より **自分の記録が先**
- フォロー機能より **フォロワーはあくまで付随**
- いいねより **「私も行きたい」が主アクション**

### D6: ラッパー・リンカー SNS モデル

**X (Twitter) や Instagram を「否定」しない。それらを「束ねる」。**

```
従来のSNS発信:
  ユーザー → Xに投稿 → Xのタイムライン
  ユーザー → Instagram投稿 → Instagramのフィード

ラッパーモデル:
  ユーザー → X/Instagram に投稿 (普段通り)
                    ↓
              [リンカー層] ← ハッシュタグ or APIで収集
                    ↓
              界隈ナビに自動インポート
                    ↓
              シール帳に「この投稿、この店のこの日」で整理
                    ↓
              地図上にプロット / コレクションに追加
```

**ユーザー体験:**

1. 普段通り X で「新宿のクラフトビール最高だった #tokyocraftbeer」と投稿
2. 界隈ナビがハッシュタグを検出 → 自動でシール帳に追加
3. ユーザーは界隈ナビを開くと「自分の街の記録」が地図に広がっている
4. 他人の投稿も見れる → 新しい発見 → 「行きたい」を登録

**技術的実装:**

```ruby
# Rails でのリンカー実装イメージ
class Linker::XCollector
  # ハッシュタグ監視 → 該当venueに紐づけて保存
  HASHTAGS = %w[tokyocraftbeer takinoiyu tokyodjbar asiafood_shinjuku]
  
  def collect
    HASHTAGS.each do |tag|
      posts = XClient.search("##{tag}", max_results: 100)
      posts.each { |post| import_as_activity(post, tag) }
    end
  end
end

class Linker::InstagramCollector
  # Instagram Basic Display API or oEmbed
  # 同様のハッシュタグ収集
end
```

**収集データは `external_posts` テーブルに格納:**

```sql
CREATE TABLE external_posts (
  id SERIAL PRIMARY KEY,
  platform TEXT NOT NULL,          -- 'x' | 'instagram'
  external_id TEXT NOT NULL UNIQUE,
  author_name TEXT,
  body TEXT,
  media_urls JSONB,                -- [{ url, type }]
  hashtags JSONB,                  -- ['tokyocraftbeer', ...]
  venue_id INTEGER REFERENCES venues(id),  -- 紐付け先 (手動 or 自動)
  posted_at TIMESTAMPTZ,
  collected_at TIMESTAMPTZ DEFAULT NOW()
);
```

### D7: 段階的スケール設計 (JS → Rails)

**「小さいものは JS/DB で完結。5人以上で Rails 発行」**

```
Phase 0: 個人ツール (JS + SQLite/JSON)
  ├── 1人用の「行ったリスト」
  ├── ブラウザ内完結 (localStorage or SQLite WASM)
  ├── 共有はURLパラメータ or JSON export
  └── コスト: ゼロ

Phase 1: 小規模共有 (JS + SQLite + 軽量API)
  ├── 5人まで
  ├── 共有サーバー不要 (GitHub Pages + JSONbin.io 等)
  ├── シール帳の共有はJSONファイルのやりとり
  └── コスト: ほぼゼロ

Phase 2: コミュニティ (Rails + PostgreSQL)
  ├── 5人以上 or 界隈が広がったら
  ├── Rails API 発行
  ├── ユーザー管理、フォロー、フィード
  └── コスト: Render $7~/月

Phase 3: プラットフォーム (Rails + Redis + CDN)
  ├── リアルタイム、通知、メーカーAPI
  └── コスト: スケールに応じて
```

**Phase 0/1 の最小実装:**

```javascript
// 個人シール帳 — ブラウザ完結
const stickerBook = {
  visited: [],    // [{ venue_id, date, memo, photo }]
  wantToGo: [],   // [{ venue_id, added_at }]
  badges: [],     // ['first_craftbeer', 'onsen_lover', ...]
  
  addVisit(venue) {
    this.visited.push({ ...venue, date: new Date() });
    this.checkBadges();
    this.save();
  },
  
  checkBadges() {
    if (this.visited.length >= 10) this.badges.push('collector_10');
    if (this.visited.filter(v => v.category === 'craftbeer'). >= 5)
      this.badges.push('beer_lover');
  },
  
  save() { localStorage.setItem('stickerbook', JSON.stringify(this)); },
  export() { return JSON.stringify(this, null, 2); }
};
```

### D8: 界隈が発展する物語

**データが積み上がるほど、街が面白くなる設計。**

```
1人目が「tokyo-craftbeer の ○○ビール 飲んだ」を記録
  → 地図に1ピン

5人目が同じ店を記録
  → 店のページに「5人が行った」と表示

10人がクラフトビールを記録
  → 「クラフトビール界隈」の統計が出現
  → 「今月の人気ビール TOP5」が自動生成

50人が記録
  → 地図がヒートマップに変わる
  → 「この界隈、今アツい」が可視化される

100人・1000人
  → 街の物語になる
  → 「2026年のtokyo-craftbeer トレンド」が語れる
  → メーカーが自発的に参入したくなる
```

**これを支えるDB設計:**

```sql
-- 界隈の「物語」を自動生成する集計ビュー
CREATE MATERIALIZED VIEW kaiwai_stats AS
SELECT
  v.category,
  DATE_TRUNC('month', ua.created_at) AS month,
  COUNT(DISTINCT ua.user_id) AS unique_visitors,
  COUNT(*) AS total_actions,
  AVG(r.rating) AS avg_rating,
  json_agg(DISTINCT v.name ORDER BY COUNT(*) DESC) FILTER (
    WHERE ua.action = 'visited'
  ) AS top_venues
FROM user_actions ua
JOIN venues v ON ua.target_id = v.id AND ua.target_type = 'venue'
LEFT JOIN reviews r ON r.target_id = v.id AND r.target_type = 'venue'
GROUP BY v.category, DATE_TRUNC('month', ua.created_at);

-- バッジ自動付与
CREATE OR REPLACE FUNCTION check_badges()
RETURNS TRIGGER AS $$
BEGIN
  -- 10店舗訪問で 'collector_10'
  IF (SELECT COUNT(DISTINCT target_id) FROM user_actions
      WHERE user_id = NEW.user_id AND action = 'visited') >= 10
  THEN
    INSERT INTO user_badges (user_id, badge)
    VALUES (NEW.user_id, 'collector_10')
    ON CONFLICT DO NOTHING;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### D9: バックエンドも発展する設計

**技術スタックも「界隈と共に成長」する。**

```
個人ツール (Phase 0)
  └── JS + localStorage

小規模 (Phase 1)
  └── Go API + SQLite (現状の新宿ランチ)

コミュニティ (Phase 2)
  └── Rails API + PostgreSQL + Redis

プラットフォーム (Phase 3)
  └── Rails + Sidekiq + Redis + CDN + Elasticsearch
```

**Go → Rails 移行のトリガー:**

| 条件 | 移行判断 |
|------|---------|
| ユーザー5人以上 | 検討開始 |
| フォロー/フィード必要 | 移行決定 |
| リアルタイム通知必要 | Redis + ActionCable |
| 検索が重い | Elasticsearch |
| メーカーAPI提供 | GraphQL |

---

## Consequences (影響)

### プラス
- **本質が「シール帳」なので、SNS疲れしない** (他人の投稿に追われない)
- **リンカー設計でX/Instagramの資産を活かせる** (ゼロから集客不要)
- **Phase 0 は今日から作れる** (JS + localStorage で個人ツール)
- **データが街の物語になる** (メーカーが参入したくなる構造)
- **技術も段階的にスケール** (オーバーエンジニアリングしない)

### マイナス/リスク
- **リンカーはX/InstagramのAPI規約に依存** (API有料化リスク)
- **ハッシュタグ収集は精度が完璧でない** (手動紐付けも必要)
- **Phase 0/1 は「ただのリスト」** (面白さが出るのはPhase 2以降)
- **5人集めるための初期モチベーション設計が重要**

### やらないこと
- ❌ X/Instagram の代替 (あくまでハブ)
- ❌ リアルタイムチャット (既存コミュニティに任せる)
- ❌ 汎用SNS機能 (DM、グループ等は作らない)
- ❌ 動画ホスティング (リンクのみ)

---

## 全体アーキテクチャ図

```
  X / Instagram (発信層)
       │
       │ ハッシュタグ / API
       ▼
  ┌─────────────────────────────┐
  │  リンカー層 (Linker)         │
  │  XCollector / InstaCollector │
  │  → external_posts テーブル   │
  └──────────────┬──────────────┘
                 │
                 ▼
  ┌─────────────────────────────┐
  │  コア (Rails API)            │
  │  venues / events / reviews   │
  │  user_actions (シール帳)     │
  │  kaiwai_stats (物語)         │
  └──────────────┬──────────────┘
                 │
       ┌─────────┼─────────┐
       ▼         ▼         ▼
  個人シール帳  界隈マップ  メーカーAPI
  (JS/DB)     (地図UI)    (Phase 3)
```

---

## 次のアクション

1. **Phase 0 の個人シール帳をJSで試作** (1日で作れる)
2. **DBスキーマに `external_posts` テーブル追加** (リンカー対応)
3. **Go API に venues/events テーブル追加** (ADR-001 の拡張)
4. **Rails プロジェクトの雛形作成** (Phase 2 の準備)
