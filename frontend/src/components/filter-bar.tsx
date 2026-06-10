"use client"

import type { Area, Genre } from "@/lib/types"

type Props = {
  areas: Area[]
  genres: Genre[]
  area: string
  genre: string
  priceMax: number
  search: string
  onArea: (v: string) => void
  onGenre: (v: string) => void
  onPriceMax: (v: number) => void
  onSearch: (v: string) => void
}

const PRICE_LIMIT = 1500

export default function FilterBar({ areas, genres, area, genre, priceMax, search, onArea, onGenre, onPriceMax, onSearch }: Props) {
  return (
    <div className="space-y-3 rounded-xl border border-zinc-700 bg-zinc-800/30 p-3 sm:p-4">
      <div className="flex flex-wrap gap-2 sm:gap-3">
        <div className="flex flex-col gap-1 flex-1 min-w-[120px] sm:flex-none">
          <label className="text-xs text-zinc-400 font-medium">エリア</label>
          <select
            value={area}
            onChange={(e) => onArea(e.target.value)}
            className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-white w-full"
          >
            <option value="">すべて</option>
            {areas.map((a) => (
              <option key={a.id} value={a.name}>{a.name}</option>
            ))}
          </select>
        </div>
        <div className="flex flex-col gap-1 flex-1 min-w-[120px] sm:flex-none">
          <label className="text-xs text-zinc-400 font-medium">ジャンル</label>
          <select
            value={genre}
            onChange={(e) => onGenre(e.target.value)}
            className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-white w-full"
          >
            <option value="">すべて</option>
            {genres.map((g) => (
              <option key={g.id} value={g.name}>{g.name}</option>
            ))}
          </select>
        </div>
        <div className="flex flex-col gap-1 flex-1 min-w-[140px]">
          <label className="text-xs text-zinc-400 font-medium">店名</label>
          <input
            type="text"
            placeholder="絞り込み..."
            value={search}
            onChange={(e) => onSearch(e.target.value)}
            className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-white w-full"
          />
        </div>
      </div>
      <div className="flex flex-col gap-1.5">
        <div className="flex justify-between items-center">
          <label className="text-xs text-zinc-400 font-medium">価格上限</label>
          <span className="text-sm font-semibold text-blue-400">
            {priceMax >= PRICE_LIMIT ? "上限なし" : `¥${priceMax}`}
          </span>
        </div>
        <input
          type="range"
          min={300}
          max={PRICE_LIMIT}
          step={50}
          value={priceMax}
          onChange={(e) => onPriceMax(Number(e.target.value))}
          className="w-full h-2 rounded-full accent-blue-500 cursor-pointer"
        />
        <div className="flex justify-between text-xs text-zinc-600">
          <span>¥300</span>
          <span>¥750</span>
          <span>¥1,200</span>
          <span>上限なし</span>
        </div>
      </div>
    </div>
  )
}
