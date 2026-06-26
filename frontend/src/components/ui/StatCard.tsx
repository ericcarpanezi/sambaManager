import { motion } from "framer-motion";

interface StatCardProps {
  title: string;
  value: number;
}

export function StatCard({ title, value }: StatCardProps) {
  return (
    <motion.article
      initial={{ opacity: 0, y: 6 }}
      animate={{ opacity: 1, y: 0 }}
      className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm dark:border-slate-700 dark:bg-slate-900"
    >
      <h3 className="text-xs uppercase text-slate-500">{title}</h3>
      <p className="mt-2 text-2xl font-semibold">{value}</p>
    </motion.article>
  );
}
