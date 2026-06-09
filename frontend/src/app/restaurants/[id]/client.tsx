"use client"

import { useState } from "react"
import Link from "next/link"
import type { RestaurantDetail } from "@/lib/types"
import Stars from "@/components/stars"
import LunchLogForm from "@/components/lunch-log-form"
import { getRestaurant } from "@/lib/api"

export default function RestaurantDetailClient({ data: initial }: { data: RestaurantDetail }) {
  const [data, setData] = useState(initial)

  async function refresh() {
    try {
      const updated = await getRestaurant(data.id)
      setData(updated)
    } catch { /* ignore */ }
  }

  const logs = data.logs

  return (
    <main className="py-6 sm:py-10">
      <Link href="/" className="text-sm text-blue-400 hover:text-blue-300 mb-4 inline-block">
        ← 一覧に戻る
      </Link>

      <div className="rounded-xl border border-zinc-700 bg-zinc-800/50 p-4 sm:p-6">
        <h1 className="text-xl sm:text-2xl font-bold text-white">{data.name}</h1>
        <div className="mt-2 flex flex-wrap gap-2">
          {data.genre && (
            <span className="rounded-full bg-zinc-700 px-2.5 py-0.5 text-xs sm:text-sm text-zinc-300">{data.genre}</span>
          )}
          {data.area && (
            <span className="rounded-full bg-zinc-700 px-2.5 py-0.5 text-xs sm:text-sm text-zinc-300">{data.area}</span>
          )}
        </div>
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-2 sm:gap-3 text-sm">
          {data.station && (
            <div><span className="text-zinc-500">最寄駅:</span> <span className="text-zinc-300">{data.station}</span></div>
          )}
          {data.walk_min != null && (
            <div><span className="text-zinc-500">徒歩:</span> <span className="text-zinc-300">{data.walk_min}分</span></div>
          )}
          {data.business_hours && (
            <div className="col-span-1 sm:col-span-2"><span className="text-zinc-500">営業時間:</span> <span className="text-zinc-300">{data.business_hours}</span></div>
          )}
          {data.address && (
            <div className="col-span-1 sm:col-span-2"><span className="text-zinc-500">住所:</span> <span className="text-zinc-300">{data.address}</span></div>
          )}
        </div>
        {data.notes && (
          <p className="mt-4 text-sm text-zinc-400">{data.notes}</p>
        )}
        {data.latitude && data.longitude && (
          <a
            href={`https://www.google.com/maps/search/?api=1&query=${data.latitude},${data.longitude}`}
            target="_blank"
            rel="noopener noreferrer"
            className="mt-4 inline-block rounded-lg bg-zinc-700 px-4 py-2 text-sm text-blue-400 hover:bg-zinc-600 active:bg-zinc-500 touch-manipulation"
          >
            Google Maps で見る
          </a>
        )}
        <div className="mt-4 flex flex-wrap gap-3 text-sm">
          {data.url_tabelog && (
            <a href={data.url_tabelog} target="_blank" rel="noopener noreferrer" className="text-orange-400 hover:underline">食べログ</a>
          )}
          {data.url_hotpepper && (
            <a href={data.url_hotpepper} target="_blank" rel="noopener noreferrer" className="text-red-400 hover:underline">ホットペッパー</a>
          )}
        </div>
      </div>

      <section className="mt-6 sm:mt-8">
        <h2 className="text-base sm:text-lg font-semibold text-white mb-3 sm:mb-4">ランチ記録 {logs.length > 0 && `(${logs.length})`}</h2>
        {logs.length === 0 ? (
          <p className="text-zinc-500 text-sm">まだ記録がありません</p>
        ) : (
          <div className="space-y-3">
            {logs.map((log) => (
              <div key={log.id} className="rounded-lg border border-zinc-700 bg-zinc-800/30 p-3 sm:p-4">
                <div className="flex items-center justify-between">
                  <span className="font-medium text-white text-sm sm:text-base">{log.menu}</span>
                  <span className="text-base sm:text-lg text-green-400">¥{log.price.toLocaleString()}</span>
                </div>
                <div className="mt-1 flex items-center flex-wrap gap-2 sm:gap-3 text-sm">
                  <Stars rating={log.rating} size="sm" />
                  <span className="text-zinc-500 text-xs sm:text-sm">{log.visited_date}</span>
                  {log.revisit && (
                    <span className="text-green-400 text-xs">また行きたい</span>
                  )}
                </div>
                {log.comment && (
                  <p className="mt-2 text-sm text-zinc-400">{log.comment}</p>
                )}
              </div>
            ))}
          </div>
        )}
      </section>

      <section className="mt-6 sm:mt-8 rounded-xl border border-zinc-700 bg-zinc-800/50 p-4 sm:p-6">
        <h2 className="text-base sm:text-lg font-semibold text-white mb-4">ランチを記録</h2>
        <LunchLogForm restaurantId={data.id} onCreated={refresh} />
      </section>
    </main>
  )
}
