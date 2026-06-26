import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api";

interface AuditLog {
  id: number;
  actorUsername: string;
  operation: string;
  objectType: string;
  objectId: string;
  result: string;
  createdAt: string;
}

export function AuditPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["audit"],
    queryFn: async () => {
      const { data } = await api.get<{ items: AuditLog[] }>("/audit/logs");
      return data.items;
    },
    refetchInterval: 10000,
  });

  return (
    <section>
      <h2 className="mb-4 text-2xl font-semibold">Auditoria</h2>
      {isLoading ? (
        <p className="text-sm text-slate-500">Carregando logs...</p>
      ) : (
        <div className="space-y-2">
          {(data ?? []).map((log) => (
            <article key={log.id} className="rounded border border-slate-200 p-3 dark:border-slate-700">
              <p className="font-medium">{log.actorUsername} - {log.operation} {log.objectType}</p>
              <p className="text-sm text-slate-500">{log.objectId} | {log.result} | {new Date(log.createdAt).toLocaleString()}</p>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
