import { NextResponse } from "next/server"
import { getRestaurant } from "@/lib/seed-data"

export function GET(_req: Request, { params }: { params: { id: string } }) {
  const id = Number(params.id)
  const data = getRestaurant(id)
  if (!data) return NextResponse.json({ error: "not found" }, { status: 404 })
  return NextResponse.json(data)
}
