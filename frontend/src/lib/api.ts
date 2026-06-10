const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? ""

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`${res.status} ${text}`)
  }
  return res.json()
}

import type { Restaurant, RestaurantDetail, LunchLog, LunchLogInput, Area, Genre } from "./types"

export async function getRestaurants(params?: {
  area?: string
  genre?: string
  price_max?: number
}): Promise<Restaurant[]> {
  const q = new URLSearchParams()
  if (params?.area) q.set("area", params.area)
  if (params?.genre) q.set("genre", params.genre)
  if (params?.price_max) q.set("price_max", String(params.price_max))
  const qs = q.toString()
  return fetchJSON(`/api/restaurants${qs ? "?" + qs : ""}`)
}

export async function getRestaurant(id: number): Promise<RestaurantDetail> {
  return fetchJSON(`/api/restaurants/${id}`)
}

export async function getLunchLogs(restaurantId?: number): Promise<LunchLog[]> {
  const qs = restaurantId ? `?restaurant_id=${restaurantId}` : ""
  return fetchJSON(`/api/lunch-logs${qs}`)
}

export async function postLunchLog(input: LunchLogInput): Promise<LunchLog> {
  return fetchJSON("/api/lunch-logs", {
    method: "POST",
    body: JSON.stringify(input),
  })
}

export async function getAreas(): Promise<Area[]> {
  return fetchJSON("/api/areas")
}

export async function getGenres(): Promise<Genre[]> {
  return fetchJSON("/api/genres")
}
