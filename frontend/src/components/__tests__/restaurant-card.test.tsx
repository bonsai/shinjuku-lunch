import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import RestaurantCard from "../restaurant-card"
import type { Restaurant } from "@/lib/types"

const mockRestaurant: Restaurant = {
  id: 1,
  name: "バンタイ",
  area: "歌舞伎町",
  genre: "タイ料理",
  station: "西武新宿駅",
  walk_min: 2,
  latitude: 35.6958,
  longitude: 139.7012,
  notes: "平日ランチ¥950",
  created_at: "2026-01-01T00:00:00Z",
}

describe("RestaurantCard", () => {
  it("renders restaurant name", () => {
    render(<RestaurantCard r={mockRestaurant} />)
    expect(screen.getByText("バンタイ")).toBeInTheDocument()
  })

  it("renders genre and area badges", () => {
    render(<RestaurantCard r={mockRestaurant} />)
    expect(screen.getByText("タイ料理")).toBeInTheDocument()
    expect(screen.getByText("歌舞伎町")).toBeInTheDocument()
  })

  it("renders walk time", () => {
    render(<RestaurantCard r={mockRestaurant} />)
    expect(screen.getByText("徒歩 2分")).toBeInTheDocument()
  })

  it("renders notes", () => {
    render(<RestaurantCard r={mockRestaurant} />)
    expect(screen.getByText("平日ランチ¥950")).toBeInTheDocument()
  })

  it("renders map link when coordinates exist", () => {
    render(<RestaurantCard r={mockRestaurant} />)
    const mapLink = screen.getByText("Map")
    expect(mapLink).toBeInTheDocument()
    expect(mapLink).toHaveAttribute("href", "https://www.google.com/maps/search/?api=1&query=35.6958,139.7012")
  })
})
