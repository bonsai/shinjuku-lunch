import { getRestaurant } from "@/lib/api"
import RestaurantDetailClient from "./client"
import { notFound } from "next/navigation"

export default async function RestaurantDetailPage(props: PageProps<"/restaurants/[id]">) {
  const { id } = await props.params
  const numId = Number(id)
  if (isNaN(numId)) notFound()

  let data
  try {
    data = await getRestaurant(numId)
  } catch {
    notFound()
  }

  return <RestaurantDetailClient data={data} />
}
