import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";
export function AppShell() {
    return (_jsxs("div", { className: "flex min-h-screen", children: [_jsx(Sidebar, {}), _jsxs("div", { className: "flex flex-1 flex-col", children: [_jsx(Topbar, {}), _jsx("main", { className: "flex-1 p-6", children: _jsx(Outlet, {}) })] })] }));
}
