import rawSeed from "@/data/seed.json"

const areas = rawSeed.areas.map((a, i) => ({ id: i + 1, name: a.name }))
const genres = rawSeed.genres.map((g, i) => ({ id: i + 1, name: g.name }))
const areaMap = new Map(areas.map(a => [a.name, a.id]))
const genreMap = new Map(genres.map(g => [g.name, g.id]))

const restaurants = rawSeed.restaurants.map((r, i) => ({
  id: i + 1,
  name: r.name,
  area: r.area,
  genre: r.genre,
  area_id: areaMap.get(r.area) ?? 0,
  genre_id: genreMap.get(r.genre) ?? 0,
  address: r.address ?? undefined,
  station: r.station ?? undefined,
  walk_min: r.walk_min ?? undefined,
  latitude: r.latitude ?? undefined,
  longitude: r.longitude ?? undefined,
  business_hours: r.business_hours ?? undefined,
  url_tabelog: r.url_tabelog ?? undefined,
  url_hotpepper: r.url_hotpepper ?? undefined,
  notes: r.notes ?? undefined,
  created_at: "2026-06-01T00:00:00Z",
}))

const restNameMap = new Map(restaurants.map(r => [r.name, r]))

const lunchLogs = rawSeed.lunch_logs.map((l, i) => ({
  id: i + 1,
  restaurant_id: restNameMap.get(l.restaurant)?.id ?? 0,
  menu: l.menu,
  price: l.price,
  rating: l.rating,
  comment: l.comment ?? "",
  revisit: l.revisit,
  visited_date: l.visited_date,
  created_at: `${l.visited_date}T00:00:00Z`,
}))

export function getAreas() {
  return areas
}

export function getGenres() {
  return genres
}

export function getRestaurants(params?: { area?: string; genre?: string; price_max?: number }) {
  let result = restaurants
  if (params?.area) result = result.filter(r => r.area === params.area)
  if (params?.genre) result = result.filter(r => r.genre === params.genre)
  return result
}

export function getRestaurant(id: number) {
  const r = restaurants.find(r => r.id === id)
  if (!r) return null
  const logs = lunchLogs.filter(l => l.restaurant_id === id)
  return { ...r, logs }
}

export function getLunchLogs(restaurantId?: number) {
  if (restaurantId) return lunchLogs.filter(l => l.restaurant_id === restaurantId)
  return lunchLogs
}

export { lunchLogs }
