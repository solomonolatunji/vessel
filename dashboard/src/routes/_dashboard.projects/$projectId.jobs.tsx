import { createFileRoute } from '@tanstack/react-router';
import { JobsList } from '#/features/projects/jobs-list';

export const Route = createFileRoute('/_dashboard/projects/$projectId/jobs')({
  component: JobsRouteComponent,
});

function JobsRouteComponent() {
  const { projectId } = Route.useParams();

  return (
    <div className="flex flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Jobs</h1>
      <JobsList projectId={projectId} />
    </div>
  );
}
