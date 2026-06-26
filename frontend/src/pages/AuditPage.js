import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api";
export function AuditPage() {
    const { data, isLoading } = useQuery({
        queryKey: ["audit"],
        queryFn: async () => {
            const { data } = await api.get("/audit/logs");
            return data.items;
        },
        refetchInterval: 10000,
    });
    return (_jsxs("section", { children: [_jsx("h2", { className: "mb-4 text-2xl font-semibold", children: "Auditoria" }), isLoading ? (_jsx("p", { className: "text-sm text-slate-500", children: "Carregando logs..." })) : (_jsx("div", { className: "space-y-2", children: (data ?? []).map((log) => (_jsxs("article", { className: "rounded border border-slate-200 p-3 dark:border-slate-700", children: [_jsxs("p", { className: "font-medium", children: [log.actorUsername, " - ", log.operation, " ", log.objectType] }), _jsxs("p", { className: "text-sm text-slate-500", children: [log.objectId, " | ", log.result, " | ", new Date(log.createdAt).toLocaleString()] })] }, log.id))) }))] }));
}
