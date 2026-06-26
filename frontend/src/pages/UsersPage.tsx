import { useMemo, useState } from "react";
import { useUsers } from "../hooks/useUsers";

export function UsersPage() {
  const [search, setSearch] = useState("");
  const { data, isLoading } = useUsers(search);

  const rows = useMemo(() => data ?? [], [data]);

  return (
    <section>
      <h2 className="mb-4 text-2xl font-semibold">Usuários</h2>
      <input
        className="mb-4 w-full max-w-sm rounded border border-slate-300 px-3 py-2 dark:border-slate-700 dark:bg-slate-900"
        placeholder="Pesquisar usuário..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />

      {isLoading ? (
        <p className="text-sm text-slate-500">Carregando usuários...</p>
      ) : (
        <div className="overflow-hidden rounded-lg border border-slate-200 dark:border-slate-700">
          <table className="min-w-full text-sm">
            <thead className="bg-slate-100 dark:bg-slate-800">
              <tr>
                <th className="px-3 py-2 text-left">Nome</th>
                <th className="px-3 py-2 text-left">Login</th>
                <th className="px-3 py-2 text-left">Departamento</th>
                <th className="px-3 py-2 text-left">OU</th>
                <th className="px-3 py-2 text-left">Status</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((user) => (
                <tr key={user.id} className="border-t border-slate-100 dark:border-slate-800">
                  <td className="px-3 py-2">{user.displayName}</td>
                  <td className="px-3 py-2">{user.username}</td>
                  <td className="px-3 py-2">{user.department}</td>
                  <td className="px-3 py-2">{user.ou}</td>
                  <td className="px-3 py-2">{user.enabled ? "Habilitado" : "Desabilitado"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
