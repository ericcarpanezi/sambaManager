import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api";
import type { DirectoryUser } from "../types";

export function useUsers(search: string) {
  return useQuery({
    queryKey: ["users", search],
    queryFn: async () => {
      const { data } = await api.get<{ items: DirectoryUser[] }>("/users", {
        params: { search },
      });
      return data.items;
    },
  });
}
