import { createFileRoute } from '@tanstack/react-router';
import { ProjectDomains } from '#/features/projects/project-domains';

export const Route = createFileRoute('/_dashboard/projects/$projectId/settings')({
  component: SettingsRouteComponent,
});

function SettingsRouteComponent() {
  const { projectId } = Route.useParams();

  return (
    <div className="flex flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Project Settings</h1>

      <div className="grid grid-cols-1 gap-4">
        <ProjectDomains projectId={projectId} />
      </div>
    </div>
  );
}
