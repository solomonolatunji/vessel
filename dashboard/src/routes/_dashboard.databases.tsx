import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/databases')({
  component: DatabasesPage,
});

function DatabasesPage() {
  return (
    <div className="p-6">
      <h1 className="mb-4 font-semibold text-2xl">Databases</h1>
      <p className="text-muted-foreground">Databases content goes here.</p>
    </div>
  );
}
