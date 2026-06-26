import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { LayoutDashboard, Users, Shield, Settings, ScrollText } from "lucide-react";
import { NavLink } from "react-router-dom";
const links = [
    { to: "/", label: "Dashboard", icon: LayoutDashboard },
    { to: "/users", label: "Usuários", icon: Users },
    { to: "/audit", label: "Auditoria", icon: ScrollText },
    { to: "/security", label: "Segurança", icon: Shield },
    { to: "/settings", label: "Configurações", icon: Settings },
];
export function Sidebar() {
    return (_jsxs("aside", { className: "w-64 border-r border-slate-800/20 bg-slate-200/60 p-4 dark:border-slate-700 dark:bg-slate-900", children: [_jsx("h1", { className: "mb-6 text-lg font-semibold", children: "AG Directory Manager" }), _jsx("nav", { className: "flex flex-col gap-2", children: links.map(({ to, label, icon: Icon }) => (_jsxs(NavLink, { to: to, className: ({ isActive }) => `flex items-center gap-2 rounded-md px-3 py-2 text-sm transition ${isActive
                        ? "bg-sky-600 text-white"
                        : "text-slate-700 hover:bg-slate-300 dark:text-slate-200 dark:hover:bg-slate-800"}`, children: [_jsx(Icon, { size: 16 }), label] }, to))) })] }));
}
