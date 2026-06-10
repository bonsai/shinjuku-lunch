import { NextResponse } from "next/server"
import { getRestaurants } from "@/lib/seed-data"

export function GET(req: Request) {
  const { searchParams } = new URL(req.url)
  const area = searchParams.get("area") ?? undefined
  const genre = searchParams.get("genre") ?? undefined
  const priceMax = searchParams.get("price_max") ? Number(searchParams.get("price_max")) : undefined
  return NextResponse.json(getRestaurants({ area, genre, price_max: priceMax }))
}
