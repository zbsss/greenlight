import client from "./client";
import {
  createQueryHook,
  createImmutableHook,
  createInfiniteHook,
  // createMutateHook,
} from "swr-openapi";
// import { isMatch } from "lodash-es";

const prefix = "movies";

export const useQuery = createQueryHook(client, prefix);
export const useImmutable = createImmutableHook(client, prefix);
export const useInfinite = createInfiniteHook(client, prefix);
// export const useMutate = createMutateHook(client, prefix);
