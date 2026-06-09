import { describe, it, expect } from "vitest"
import type { Restaurant, LunchLog, LunchLogInput, Area, Genre } from "../types"

describe("Types", () => {
  it("Restaurant type is properly structured", () => {
    const r: Restaurant = {
      id: 1,
      name: "テスト",
      area: "新宿",
      genre: "和食",
      created_at: "2026-01-01",
    }
    expect(r.id).toBe(1)
    expect(r.name).toBe("テスト")
  })

  it("LunchLog type allows optional fields", () => {
    const log: LunchLog = {
      id: 1,
      restaurant_id: 1,
      menu: "定食",
      price: 800,
      rating: 4,
      revisit: false,
      visited_date: "2026-06-01",
      created_at: "2026-06-01T00:00:00Z",
    }
    expect(log.comment).toBeUndefined()
  })

  it("LunchLogInput has required fields", () => {
    const input: LunchLogInput = {
      restaurant_id: 1,
      menu: "定食",
      price: 800,
      rating: 4,
    }
    expect(input.restaurant_id).toBe(1)
    expect(input.revisit).toBeUndefined()
  })

  it("Area and Genre types", () => {
    const area: Area = { id: 1, name: "歌舞伎町" }
    const genre: Genre = { id: 1, name: "タイ料理" }
    expect(area.name).toBe("歌舞伎町")
    expect(genre.name).toBe("タイ料理")
  })
})
