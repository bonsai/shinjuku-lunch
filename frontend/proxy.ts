import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

const API_URL = process.env.API_URL ?? "http://localhost:8080"

export function proxy(request: NextRequest) {
  if (request.nextUrl.pathname.startsWith("/api")) {
    const dest = new URL(request.nextUrl.pathname + request.nextUrl.search, API_URL)
    return NextResponse.rewrite(dest)
  }
}

export const config = {
  matcher: "/api/:path*",
}
