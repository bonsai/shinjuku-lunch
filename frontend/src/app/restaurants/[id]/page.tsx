import { getRestaurant } from "@/lib/seed-data"
import RestaurantDetailClient from "./client"
import { notFound } from "next/navigation"

export default async function RestaurantDetailPage(props: PageProps<"/restaurants/[id]">) {
  const { id } = await props.params
  const numId = Number(id)
  if (isNaN(numId)) notFound()

  const data = getRestaurant(numId)
  if (!data) notFound()

  return <RestaurantDetailClient data={data} />
}
