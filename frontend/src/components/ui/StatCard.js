import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { motion } from "framer-motion";
export function StatCard({ title, value }) {
    return (_jsxs(motion.article, { initial: { opacity: 0, y: 6 }, animate: { opacity: 1, y: 0 }, className: "rounded-lg border border-slate-200 bg-white p-4 shadow-sm dark:border-slate-700 dark:bg-slate-900", children: [_jsx("h3", { className: "text-xs uppercase text-slate-500", children: title }), _jsx("p", { className: "mt-2 text-2xl font-semibold", children: value })] }));
}
