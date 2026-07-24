import { createFileRoute } from '@tanstack/react-router';
import { ProjectMembers } from '#/features/projects/project-members';
import { ProjectTokens } from '#/features/projects/project-tokens';

export const Route = createFileRoute('/_dashboard/projects/$projectId/settings')({
  component: SettingsRouteComponent,
});

function SettingsRouteComponent() {
  const { projectId } = Route.useParams();

  return (
    <div className="mx-auto flex w-full max-w-5xl flex-col gap-6 p-4">
      <h1 className="font-bold text-2xl">Project Settings</h1>

      <div className="grid grid-cols-1 gap-8">
        <section className="rounded-lg border bg-white p-6 shadow-sm">
          <ProjectMembers projectId={projectId} />
        </section>

        <section className="rounded-lg border bg-white p-6 shadow-sm">
          <ProjectTokens projectId={projectId} />
        </section>
      </div>
    </div>
  );
}
