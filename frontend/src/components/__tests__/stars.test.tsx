import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import Stars from "../stars"

describe("Stars", () => {
  it("renders correct number of filled stars", () => {
    render(<Stars rating={3} />)
    const text = screen.getByText("★★★☆☆")
    expect(text).toBeInTheDocument()
  })

  it("renders 0 stars for rating 0", () => {
    render(<Stars rating={0} />)
    expect(screen.getByText("☆☆☆☆☆")).toBeInTheDocument()
  })

  it("renders 5 stars for rating 5", () => {
    render(<Stars rating={5} />)
    expect(screen.getByText("★★★★★")).toBeInTheDocument()
  })

  it("clamps rating to 0-5", () => {
    render(<Stars rating={-1} />)
    expect(screen.getByText("☆☆☆☆☆")).toBeInTheDocument()

    render(<Stars rating={10} />)
    expect(screen.getByText("★★★★★")).toBeInTheDocument()
  })
})
