import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api";
export function useUsers(search) {
    return useQuery({
        queryKey: ["users", search],
        queryFn: async () => {
            const { data } = await api.get("/users", {
                params: { search },
            });
            return data.items;
        },
    });
}
