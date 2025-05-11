import client from "@/lib/api/client";

export const dynamic = "force-dynamic";

export default async function Home() {
  const { data, error } = await client.GET("/v1/movies");

  console.log(data, error);

  if (error) return `An error occurred: ${JSON.stringify(error)}`;

  if (!data.movies) return "No movies found";

  return <pre>{JSON.stringify(data, null, 2)}</pre>;
}
