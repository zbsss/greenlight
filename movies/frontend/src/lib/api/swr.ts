import client from "./client";
import {
  createQueryHook,
  createImmutableHook,
  createInfiniteHook,
} from "swr-openapi";

const prefix = "movies";

export const useQuery = createQueryHook(client, prefix);
export const useImmutable = createImmutableHook(client, prefix);
export const useInfinite = createInfiniteHook(client, prefix);
