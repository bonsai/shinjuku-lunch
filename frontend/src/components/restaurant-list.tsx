"use client"

import { useState, useEffect, useCallback } from "react"
import type { Restaurant, Area, Genre } from "@/lib/types"
import { getRestaurants } from "@/lib/api"
import RestaurantCard from "./restaurant-card"
import FilterBar from "./filter-bar"

type Props = {
  initialRestaurants: Restaurant[]
  areas: Area[]
  genres: Genre[]
}

export default function RestaurantList({ initialRestaurants, areas, genres }: Props) {
  const [restaurants, setRestaurants] = useState(initialRestaurants)
  const [loading, setLoading] = useState(false)

  const handleFilter = useCallback(async (params: { area?: string; genre?: string; price_max?: number }) => {
    setLoading(true)
    try {
      const data = await getRestaurants(params)
      setRestaurants(data)
    } catch {
      // keep current list on error
    } finally {
      setLoading(false)
    }
  }, [])

  return (
    <div className="space-y-6">
      <FilterBar areas={areas} genres={genres} onChange={handleFilter} />
      {loading ? (
        <div className="text-center py-12 text-zinc-500">読み込み中...</div>
      ) : restaurants.length === 0 ? (
        <div className="text-center py-12 text-zinc-500">該当する店舗がありません</div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2">
          {restaurants.map((r) => (
            <RestaurantCard key={r.id} r={r} />
          ))}
        </div>
      )}
    </div>
  )
}
