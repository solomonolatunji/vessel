import { createFileRoute } from '@tanstack/react-router';
import { ComposeDeployForm } from '#/features/projects/compose-deploy-form';

export const Route = createFileRoute('/_dashboard/projects/$projectId/compose')({
  component: ComposeRouteComponent,
});

function ComposeRouteComponent() {
  const { projectId } = Route.useParams();

  return (
    <div className="flex flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Docker Compose</h1>
      <ComposeDeployForm projectId={projectId} />
    </div>
  );
}
