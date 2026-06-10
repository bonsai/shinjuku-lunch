import { getRestaurants, getAreas, getGenres, getLunchLogs } from "@/lib/seed-data"
import RestaurantList from "@/components/restaurant-list"

export default function Home() {
  const restaurants = getRestaurants()
  const areas = getAreas()
  const genres = getGenres()
  const logs = getLunchLogs()

  return (
    <main className="py-6 sm:py-10">
      <div className="mb-6 sm:mb-8">
        <h1 className="text-xl sm:text-2xl font-bold text-white">新宿ランチナビ</h1>
        <p className="mt-1 text-xs sm:text-sm text-zinc-400">
          新宿・歌舞伎町・大久保のランチ情報
        </p>
      </div>
      <RestaurantList initialRestaurants={restaurants} areas={areas} genres={genres} logs={logs} />
    </main>
  )
}
