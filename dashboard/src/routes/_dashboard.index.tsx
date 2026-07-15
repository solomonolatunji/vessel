import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_dashboard/")({
  component: DashboardPage,
});

function DashboardPage() {
  return (
    <div className="p-6">
      <h1 className="mb-4 font-semibold text-2xl">Dashboard</h1>
      <p className="text-muted-foreground">Dashboard content goes here.</p>
    </div>
  );
}
