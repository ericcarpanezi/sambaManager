import { useDashboard } from "../hooks/useDashboard";
import { StatCard } from "../components/ui/StatCard";

export function DashboardPage() {
  const { data, isLoading } = useDashboard();

  if (isLoading || !data) {
    return <p className="text-sm text-slate-500">Carregando dashboard...</p>;
  }

  return (
    <section>
      <h2 className="mb-4 text-2xl font-semibold">Dashboard</h2>
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard title="Usuários" value={data.usersTotal} />
        <StatCard title="Computadores" value={data.computersTotal} />
        <StatCard title="Grupos" value={data.groupsTotal} />
        <StatCard title="OUs" value={data.ousTotal} />
        <StatCard title="Usuários Bloqueados" value={data.lockedUsers} />
        <StatCard title="Contas Desabilitadas" value={data.disabledAccounts} />
        <StatCard title="Computadores Inativos" value={data.inactiveComputers} />
        <StatCard title="Eventos Recentes" value={data.recentEvents} />
      </div>
    </section>
  );
}
