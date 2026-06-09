export type Restaurant = {
  id: number
  name: string
  area: string
  genre: string
  address?: string
  station?: string
  walk_min?: number
  latitude?: number
  longitude?: number
  business_hours?: string
  url_tabelog?: string
  url_hotpepper?: string
  notes?: string
  created_at: string
}

export type LunchLog = {
  id: number
  restaurant_id: number
  menu: string
  price: number
  rating: number
  comment?: string
  revisit: boolean
  visited_date: string
  created_at: string
}

export type LunchLogInput = {
  restaurant_id: number
  menu: string
  price: number
  rating: number
  comment?: string
  revisit?: boolean
  visited_date?: string
}

export type RestaurantDetail = Restaurant & {
  logs: LunchLog[]
}

export type Area = {
  id: number
  name: string
}

export type Genre = {
  id: number
  name: string
}
