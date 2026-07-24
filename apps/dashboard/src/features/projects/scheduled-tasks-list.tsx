import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import {
  useDeleteScheduledTask,
  useListScheduledTasks,
  useTriggerScheduledTask,
} from '#/hooks/useScheduledTasks';

export function ScheduledTasksList() {
  const { data: tasksRes, isLoading } = useListScheduledTasks('');
  const triggerTask = useTriggerScheduledTask();
  const deleteTask = useDeleteScheduledTask();

  if (isLoading) {
    return <div>Loading scheduled tasks...</div>;
  }

  const tasks = tasksRes?.data || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Scheduled Tasks</CardTitle>
      </CardHeader>
      <CardContent>
        {tasks.length === 0 ? (
          <p className="text-gray-500 text-sm">No scheduled tasks found for this service.</p>
        ) : (
          <ul className="space-y-4">
            {tasks.map((task) => (
              <li key={task.id} className="flex flex-col gap-2 rounded border p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">{task.name}</h3>
                    <p className="text-gray-500 text-sm">
                      {task.schedule} | {task.status}
                    </p>
                    <code className="rounded bg-gray-100 px-1 py-0.5 text-xs">{task.command}</code>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => triggerTask.mutate({ id: task.id })}
                      disabled={triggerTask.isPending}
                    >
                      Trigger
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => deleteTask.mutate({ id: task.id })}
                      disabled={deleteTask.isPending}
                    >
                      Delete
                    </Button>
                  </div>
                </div>
                {task.lastRunAt && (
                  <div className="mt-2 text-gray-500 text-xs">
                    Last Run: {new Date(task.lastRunAt).toLocaleString()}
                  </div>
                )}
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
