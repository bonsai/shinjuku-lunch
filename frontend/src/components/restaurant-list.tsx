"use client"

import { useState, useMemo } from "react"
import type { Restaurant, Area, Genre, LunchLog } from "@/lib/types"
import RestaurantCard from "./restaurant-card"
import FilterBar from "./filter-bar"

const PRICE_LIMIT = 1500

type Props = {
  initialRestaurants: Restaurant[]
  areas: Area[]
  genres: Genre[]
  logs: LunchLog[]
}

export default function RestaurantList({ initialRestaurants, areas, genres, logs }: Props) {
  const [area, setArea] = useState("")
  const [genre, setGenre] = useState("")
  const [priceMax, setPriceMax] = useState(PRICE_LIMIT)
  const [search, setSearch] = useState("")

  const minPriceByRestaurant = useMemo(() => {
    const map = new Map<number, number>()
    logs.forEach((log) => {
      const cur = map.get(log.restaurant_id)
      if (cur === undefined || log.price < cur) map.set(log.restaurant_id, log.price)
    })
    return map
  }, [logs])

  const filtered = useMemo(() => {
    return initialRestaurants.filter((r) => {
      if (area && r.area !== area) return false
      if (genre && r.genre !== genre) return false
      if (search && !r.name.includes(search)) return false
      if (priceMax < PRICE_LIMIT) {
        const minPrice = minPriceByRestaurant.get(r.id)
        if (minPrice !== undefined && minPrice > priceMax) return false
      }
      return true
    })
  }, [initialRestaurants, area, genre, search, priceMax, minPriceByRestaurant])

  return (
    <div className="space-y-4">
      <FilterBar
        areas={areas}
        genres={genres}
        area={area}
        genre={genre}
        priceMax={priceMax}
        search={search}
        onArea={setArea}
        onGenre={setGenre}
        onPriceMax={setPriceMax}
        onSearch={setSearch}
      />
      <p className="text-xs text-zinc-500">{filtered.length} 件表示</p>
      {filtered.length === 0 ? (
        <div className="text-center py-12 text-zinc-500">該当する店舗がありません</div>
      ) : (
        <div className="grid gap-3 sm:gap-4 grid-cols-1 sm:grid-cols-2">
          {filtered.map((r) => (
            <RestaurantCard key={r.id} r={r} />
          ))}
        </div>
      )}
    </div>
  )
}
