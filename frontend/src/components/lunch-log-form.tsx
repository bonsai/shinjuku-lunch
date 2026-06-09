"use client"

import { useState } from "react"
import { postLunchLog } from "@/lib/api"
import Stars from "./stars"

export default function LunchLogForm({
  restaurantId,
  onCreated,
}: {
  restaurantId: number
  onCreated: () => void
}) {
  const [menu, setMenu] = useState("")
  const [price, setPrice] = useState("")
  const [rating, setRating] = useState(3)
  const [comment, setComment] = useState("")
  const [revisit, setRevisit] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!menu || !price) return
    setSubmitting(true)
    setError("")
    try {
      await postLunchLog({
        restaurant_id: restaurantId,
        menu,
        price: Number(price),
        rating,
        comment: comment || undefined,
        revisit,
      })
      setMenu("")
      setPrice("")
      setRating(3)
      setComment("")
      setRevisit(false)
      onCreated()
    } catch (err) {
      setError(err instanceof Error ? err.message : "送信失敗")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <div className="flex flex-col gap-1">
          <label className="text-sm text-zinc-400">メニュー *</label>
          <input
            required
            value={menu}
            onChange={(e) => setMenu(e.target.value)}
            className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-white"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-sm text-zinc-400">価格 (円) *</label>
          <input
            required
            type="number"
            min={0}
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-white"
          />
        </div>
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-sm text-zinc-400">評価</label>
        <div className="flex gap-1">
          {[1, 2, 3, 4, 5].map((n) => (
            <button
              key={n}
              type="button"
              onClick={() => setRating(n)}
              className={`text-2xl transition-colors ${n <= rating ? "text-yellow-400" : "text-zinc-600"}`}
            >
              {n <= rating ? "★" : "☆"}
            </button>
          ))}
        </div>
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-sm text-zinc-400">コメント</label>
        <textarea
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          rows={2}
          className="rounded-lg border border-zinc-700 bg-zinc-800 px-3 py-2 text-white resize-none"
        />
      </div>
      <label className="flex items-center gap-2 text-sm text-zinc-300">
        <input
          type="checkbox"
          checked={revisit}
          onChange={(e) => setRevisit(e.target.checked)}
          className="rounded border-zinc-600"
        />
        また行きたい
      </label>
      {error && <p className="text-sm text-red-400">{error}</p>}
      <button
        type="submit"
        disabled={submitting}
        className="rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-500 disabled:opacity-50 transition-colors"
      >
        {submitting ? "送信中..." : "記録する"}
      </button>
    </form>
  )
}
