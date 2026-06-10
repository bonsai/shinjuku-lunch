import { NextResponse } from "next/server"
import { getAreas } from "@/lib/seed-data"

export function GET() {
  return NextResponse.json(getAreas())
}
