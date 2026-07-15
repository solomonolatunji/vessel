import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/projects')({
  component: ProjectsPage,
});

function ProjectsPage() {
  return (
    <div className="p-6">
      <h1 className="mb-4 font-semibold text-2xl">Projects</h1>
      <p className="text-muted-foreground">Projects content goes here.</p>
    </div>
  );
}
