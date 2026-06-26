import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useDashboard } from "../hooks/useDashboard";
import { StatCard } from "../components/ui/StatCard";
export function DashboardPage() {
    const { data, isLoading } = useDashboard();
    if (isLoading || !data) {
        return _jsx("p", { className: "text-sm text-slate-500", children: "Carregando dashboard..." });
    }
    return (_jsxs("section", { children: [_jsx("h2", { className: "mb-4 text-2xl font-semibold", children: "Dashboard" }), _jsxs("div", { className: "grid gap-4 md:grid-cols-2 xl:grid-cols-4", children: [_jsx(StatCard, { title: "Usu\u00E1rios", value: data.usersTotal }), _jsx(StatCard, { title: "Computadores", value: data.computersTotal }), _jsx(StatCard, { title: "Grupos", value: data.groupsTotal }), _jsx(StatCard, { title: "OUs", value: data.ousTotal }), _jsx(StatCard, { title: "Usu\u00E1rios Bloqueados", value: data.lockedUsers }), _jsx(StatCard, { title: "Contas Desabilitadas", value: data.disabledAccounts }), _jsx(StatCard, { title: "Computadores Inativos", value: data.inactiveComputers }), _jsx(StatCard, { title: "Eventos Recentes", value: data.recentEvents })] })] }));
}
