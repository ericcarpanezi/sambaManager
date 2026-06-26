export function SecurityPage() {
  return (
    <section className="space-y-4">
      <h2 className="text-2xl font-semibold">Segurança</h2>
      <ul className="list-disc space-y-1 pl-5 text-sm text-slate-600 dark:text-slate-300">
        <li>JWT com Access + Refresh Token</li>
        <li>CSRF em operações mutáveis</li>
        <li>Rate Limiting por IP</li>
        <li>Headers de CSP, XSS e hardening</li>
        <li>Auditoria imutável por operação</li>
      </ul>
    </section>
  );
}
