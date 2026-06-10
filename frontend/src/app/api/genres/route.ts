import { NextResponse } from "next/server"
import { getGenres } from "@/lib/seed-data"

export function GET() {
  return NextResponse.json(getGenres())
}
