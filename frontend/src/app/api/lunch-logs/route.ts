import { NextResponse } from "next/server"
import { getLunchLogs } from "@/lib/seed-data"

export function GET(req: Request) {
  const { searchParams } = new URL(req.url)
  const restaurantId = searchParams.get("restaurant_id") ? Number(searchParams.get("restaurant_id")) : undefined
  return NextResponse.json(getLunchLogs(restaurantId))
}

export async function POST(req: Request) {
  const body = await req.json()
  const { restaurant_id, menu, price, rating, comment, revisit, visited_date } = body
  return NextResponse.json({
    id: Date.now(),
    restaurant_id,
    menu,
    price,
    rating,
    comment: comment ?? "",
    revisit: revisit ?? false,
    visited_date: visited_date ?? new Date().toISOString().split("T")[0],
    created_at: new Date().toISOString(),
  }, { status: 201 })
}
