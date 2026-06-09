"use client"

import { useState, useEffect, useCallback } from "react"
import type { Area, Genre } from "@/lib/types"

type Props = {
  areas: Area[]
  genres: Genre[]
  onChange: (params: { area?: string; genre?: string; price_max?: number }) => void
}

export default function FilterBar({ areas, genres, onChange }: Props) {
  const [area, setArea] = useState("")
  const [genre, setGenre] = useState("")
  const [priceMax, setPriceMax] = useState("")

  const apply = useCallback(() => {
    onChange({
      area: area || undefined,
      genre: genre || undefined,
      price_max: priceMax ? Number(priceMax) : undefined,
    })
  }, [area, genre, priceMax, onChange])

  return (
    <div className="flex flex-wrap gap-3 items-end">
      <div className="flex flex-col gap-1">
        <label className="text-xs text-zinc-400 font-medium">エリア</label>
        <select
          value={area}
          onChange={(e) => { setArea(e.target.value); setTimeout(apply, 0) }}
          className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-white"
        >
          <option value="">すべて</option>
          {areas.map((a) => (
            <option key={a.id} value={a.name}>{a.name}</option>
          ))}
        </select>
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs text-zinc-400 font-medium">ジャンル</label>
        <select
          value={genre}
          onChange={(e) => { setGenre(e.target.value); setTimeout(apply, 0) }}
          className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-white"
        >
          <option value="">すべて</option>
          {genres.map((g) => (
            <option key={g.id} value={g.name}>{g.name}</option>
          ))}
        </select>
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs text-zinc-400 font-medium">価格上限</label>
        <input
          type="number"
          placeholder="例: 1000"
          value={priceMax}
          onChange={(e) => setPriceMax(e.target.value)}
          className="w-28 rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-sm text-white"
        />
      </div>
      <button
        onClick={apply}
        className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-500 transition-colors"
      >
        絞り込み
      </button>
    </div>
  )
}
