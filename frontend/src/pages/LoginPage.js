import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../lib/api";
import { useAuthStore } from "../store/auth";
export function LoginPage() {
    const navigate = useNavigate();
    const setAuth = useAuthStore((s) => s.setAuth);
    const [username, setUsername] = useState("operador.demo");
    const [password, setPassword] = useState("demo");
    const [error, setError] = useState(null);
    async function onSubmit(event) {
        event.preventDefault();
        setError(null);
        try {
            const { data } = await api.post("/auth/login", { username, password });
            setAuth({
                accessToken: data.tokens.accessToken,
                refreshToken: data.tokens.refreshToken,
                username: data.user.username,
                permissions: data.permissions,
            });
            navigate("/");
        }
        catch {
            setError("Falha ao autenticar. Verifique usuário e senha.");
        }
    }
    return (_jsx("div", { className: "flex min-h-screen items-center justify-center bg-slate-100 dark:bg-slate-950", children: _jsxs("form", { onSubmit: onSubmit, className: "w-full max-w-sm rounded-lg bg-white p-6 shadow dark:bg-slate-900", children: [_jsx("h1", { className: "mb-2 text-xl font-semibold", children: "AG Directory Manager" }), _jsx("p", { className: "mb-4 text-sm text-slate-500", children: "Acesso operacional ao Samba AD" }), _jsx("label", { className: "mb-2 block text-sm", children: "Usu\u00E1rio" }), _jsx("input", { className: "mb-4 w-full rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-800", value: username, onChange: (e) => setUsername(e.target.value) }), _jsx("label", { className: "mb-2 block text-sm", children: "Senha" }), _jsx("input", { type: "password", className: "mb-4 w-full rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-800", value: password, onChange: (e) => setPassword(e.target.value) }), error && _jsx("p", { className: "mb-4 text-sm text-red-500", children: error }), _jsx("button", { className: "w-full rounded bg-sky-600 px-4 py-2 text-white hover:bg-sky-700", children: "Entrar" })] }) }));
}
