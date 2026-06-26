import { FormEvent, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../lib/api";

interface LDAPSettings {
  serverUrl: string;
  baseDn: string;
  bindDn: string;
  bindSecret: string;
  startTls: boolean;
  skipVerify: boolean;
}

export function SettingsPage() {
  const queryClient = useQueryClient();
  const [message, setMessage] = useState<string | null>(null);

  const { data } = useQuery({
    queryKey: ["ldap-settings"],
    queryFn: async () => {
      const { data } = await api.get<LDAPSettings>("/settings/ldap");
      return data;
    },
  });

  const save = useMutation({
    mutationFn: async (payload: LDAPSettings) => {
      await api.put("/settings/ldap", payload);
    },
    onSuccess: () => {
      setMessage("Configurações salvas com sucesso.");
      queryClient.invalidateQueries({ queryKey: ["ldap-settings"] });
    },
  });

  function onSubmit(event: FormEvent<HTMLFormElement>) {
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

  return (
    <section>
      <h2 className="mb-4 text-2xl font-semibold">Configurações LDAP</h2>
      <form onSubmit={onSubmit} className="grid max-w-xl gap-3">
        <input name="serverUrl" defaultValue={data?.serverUrl} placeholder="Servidor LDAP" className="rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" />
        <input name="baseDn" defaultValue={data?.baseDn} placeholder="Base DN" className="rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" />
        <input name="bindDn" defaultValue={data?.bindDn} placeholder="Bind DN" className="rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" />
        <input name="bindSecret" type="password" defaultValue={data?.bindSecret} placeholder="Senha de serviço" className="rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900" />

        <label className="text-sm"><input name="startTls" type="checkbox" defaultChecked={data?.startTls} className="mr-2" />StartTLS</label>
        <label className="text-sm"><input name="skipVerify" type="checkbox" defaultChecked={data?.skipVerify} className="mr-2" />Ignorar validação de certificado</label>

        <button disabled={save.isPending} className="rounded bg-sky-600 px-4 py-2 text-white hover:bg-sky-700 disabled:opacity-60">
          {save.isPending ? "Salvando..." : "Salvar"}
        </button>
      </form>
      {message && <p className="mt-3 text-sm text-emerald-500">{message}</p>}
    </section>
  );
}
