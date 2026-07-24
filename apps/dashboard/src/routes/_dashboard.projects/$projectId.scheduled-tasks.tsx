import { createFileRoute } from '@tanstack/react-router';
import { ScheduledTasksList } from '#/features/projects/scheduled-tasks-list';

export const Route = createFileRoute('/_dashboard/projects/$projectId/scheduled-tasks')({
  component: ScheduledTasksRouteComponent,
});

function ScheduledTasksRouteComponent() {
  Route.useParams();

  return (
    <div className="flex flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Scheduled Tasks</h1>
      <ScheduledTasksList />
    </div>
  );
}
