import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Moon, Sun } from "lucide-react";
import { useThemeStore } from "../../store/theme";
export function Topbar() {
    const { theme, toggle } = useThemeStore();
    return (_jsxs("header", { className: "flex h-14 items-center justify-between border-b border-slate-800/20 bg-white px-6 dark:border-slate-700 dark:bg-slate-900", children: [_jsx("div", { className: "text-sm text-slate-600 dark:text-slate-400", children: "Painel Administrativo Samba AD" }), _jsx("button", { onClick: toggle, className: "rounded-md border border-slate-300 p-2 hover:bg-slate-100 dark:border-slate-700 dark:hover:bg-slate-800", "aria-label": "Alternar tema", children: theme === "dark" ? _jsx(Sun, { size: 16 }) : _jsx(Moon, { size: 16 }) })] }));
}
