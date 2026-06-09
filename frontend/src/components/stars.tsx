export default function Stars({ rating, size = "md" }: { rating: number; size?: "sm" | "md" | "lg" }) {
  const sizeClass = size === "sm" ? "text-sm" : size === "lg" ? "text-xl" : "text-base"
  return (
    <span className={`${sizeClass} text-yellow-400`}>
      {"★".repeat(Math.min(5, Math.max(0, rating)))}
      {"☆".repeat(5 - Math.min(5, Math.max(0, rating)))}
    </span>
  )
}
