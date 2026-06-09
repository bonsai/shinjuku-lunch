import Link from "next/link"
import type { Restaurant } from "@/lib/types"
import Stars from "./stars"

export default function RestaurantCard({ r }: { r: Restaurant }) {
  return (
    <Link
      href={`/restaurants/${r.id}`}
      className="block rounded-xl border border-zinc-700 bg-zinc-800/50 p-3 sm:p-4 hover:bg-zinc-700/50 active:bg-zinc-700 transition-colors touch-manipulation"
    >
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 flex-1">
          <h3 className="text-base sm:text-lg font-semibold text-white truncate">{r.name}</h3>
          <div className="mt-1 flex flex-wrap gap-1.5 text-sm text-zinc-400">
            {r.genre && (
              <span className="rounded-full bg-zinc-700 px-2 py-0.5 text-xs">{r.genre}</span>
            )}
            {r.area && (
              <span className="rounded-full bg-zinc-700 px-2 py-0.5 text-xs">{r.area}</span>
            )}
          </div>
        </div>
        {r.latitude && r.longitude && (
          <a
            href={`https://www.google.com/maps/search/?api=1&query=${r.latitude},${r.longitude}`}
            target="_blank"
            rel="noopener noreferrer"
            onClick={(e) => e.stopPropagation()}
            className="shrink-0 rounded-lg bg-zinc-700 px-2 py-1 text-xs text-blue-400 hover:bg-zinc-600 active:bg-zinc-500 touch-manipulation"
          >
            Map
          </a>
        )}
      </div>
      <div className="mt-2 flex items-center gap-3 text-xs sm:text-sm">
        {r.walk_min != null && (
          <span className="text-zinc-400">徒歩 {r.walk_min}分</span>
        )}
        {r.station && (
          <span className="text-zinc-500">{r.station}</span>
        )}
      </div>
      {r.notes && (
        <p className="mt-1.5 text-xs sm:text-sm text-zinc-400 line-clamp-2">{r.notes}</p>
      )}
    </Link>
  )
}
