import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../lib/api";
export function SettingsPage() {
    const queryClient = useQueryClient();
    const [message, setMessage] = useState(null);
    const { data } = useQuery({
        queryKey: ["ldap-settings"],
        queryFn: async () => {
            const { data } = await api.get("/settings/ldap");
            return data;
        },
    });
    const save = useMutation({
        mutationFn: async (payload) => {
            await api.put("/settings/ldap", payload);
        },
        onSuccess: () => {
            setMessage("Configurações salvas com sucesso.");
            queryClient.invalidateQueries({ queryKey: ["ldap-settings"] });
        },
    });
    function onSubmit(event) {
        event.preventDefault();
        const form = new FormData(event.currentTarget);
        save.mutate({
            serverUrl: String(form.get("serverUrl") || ""),
            baseDn: String(form.get("baseDn") || ""),
            bindDn: String(form.get("bindDn") || ""),
            bindSecret: String(form.get("bindSecret") || ""),
            startTls: form.get("startTls") === "on",
            skipVerify: form.get("skipVerify") === "on",
        });
    }
    return (_jsxs("section", { children: [_jsx("h2", { className: "mb-4 text-2xl font-semibold", children: "Configura\u00E7\u00F5es LDAP" }), _jsxs("form", { onSubmit: onSubmit, className: "grid max-w-xl gap-3", children: [_jsx("input", { name: "serverUrl", defaultValue: data?.serverUrl, placeholder: "Servidor LDAP", className: "rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" }), _jsx("input", { name: "baseDn", defaultValue: data?.baseDn, placeholder: "Base DN", className: "rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" }), _jsx("input", { name: "bindDn", defaultValue: data?.bindDn, placeholder: "Bind DN", className: "rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" }), _jsx("input", { name: "bindSecret", type: "password", defaultValue: data?.bindSecret, placeholder: "Senha de servi\u00E7o", className: "rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" }), _jsxs("label", { className: "text-sm", children: [_jsx("input", { name: "startTls", type: "checkbox", defaultChecked: data?.startTls, className: "mr-2" }), "StartTLS"] }), _jsxs("label", { className: "text-sm", children: [_jsx("input", { name: "skipVerify", type: "checkbox", defaultChecked: data?.skipVerify, className: "mr-2" }), "Ignorar valida\u00E7\u00E3o de certificado"] }), _jsx("button", { disabled: save.isPending, className: "rounded bg-sky-600 px-4 py-2 text-white hover:bg-sky-700 disabled:opacity-60", children: save.isPending ? "Salvando..." : "Salvar" })] }), message && _jsx("p", { className: "mt-3 text-sm text-emerald-500", children: message })] }));
}
