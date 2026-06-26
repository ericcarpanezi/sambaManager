import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api";
import type { DashboardSnapshot } from "../types";

export function useDashboard() {
  return useQuery({
    queryKey: ["dashboard"],
    queryFn: async () => {
      const { data } = await api.get<{ snapshot: DashboardSnapshot }>("/dashboard");
      return data.snapshot;
    },
    refetchInterval: 15000,
  });
}
