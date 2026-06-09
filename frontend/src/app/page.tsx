import { getRestaurants, getAreas, getGenres } from "@/lib/api"
import RestaurantList from "@/components/restaurant-list"

export const dynamic = "force-dynamic"

export default async function Home() {
  let restaurants: Awaited<ReturnType<typeof getRestaurants>> = []
  let areas: Awaited<ReturnType<typeof getAreas>> = []
  let genres: Awaited<ReturnType<typeof getGenres>> = []

  try {
    [restaurants, areas, genres] = await Promise.all([
      getRestaurants(),
      getAreas(),
      getGenres(),
    ])
  } catch {
    // API not available during build or offline
  }

  return (
    <main className="flex-1 w-full max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-white">新宿ランチナビ</h1>
        <p className="mt-1 text-sm text-zinc-400">
          新宿・歌舞伎町・大久保のランチ情報
        </p>
      </div>
      <RestaurantList initialRestaurants={restaurants} areas={areas} genres={genres} />
    </main>
  )
}
