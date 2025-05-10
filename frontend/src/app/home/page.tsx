import client from "@/lib/api/client";

export default async function Home() {
  const { data, error } = await client.GET("/v1/movies");

  if (error) return `An error occurred: ${JSON.stringify(error)}`;

  if (!data.movies) return "No movies found";

  return <div>{data.movies.map((movie) => movie.title).join(", ")}</div>;
}
