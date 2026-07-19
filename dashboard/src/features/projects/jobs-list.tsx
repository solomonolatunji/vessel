import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { useDeleteJob, useListJobs, useTriggerJob } from '#/hooks/useJobs';

export function JobsList() {
  const { data: jobsRes, isLoading } = useListJobs();
  const triggerJob = useTriggerJob();
  const deleteJob = useDeleteJob();

  if (isLoading) {
    return <div>Loading jobs...</div>;
  }

  const jobs = jobsRes?.data || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Cron Jobs</CardTitle>
      </CardHeader>
      <CardContent>
        {jobs.length === 0 ? (
          <p className="text-gray-500 text-sm">No jobs found for this project.</p>
        ) : (
          <ul className="space-y-4">
            {jobs.map((job) => (
              <li key={job.id} className="flex flex-col gap-2 rounded border p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">{job.name}</h3>
                    <p className="text-gray-500 text-sm">
                      {job.schedule} | {job.status}
                    </p>
                    <code className="rounded bg-gray-100 px-1 py-0.5 text-xs">{job.command}</code>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => triggerJob.mutate({ id: job.id })}
                      disabled={triggerJob.isPending}
                    >
                      Trigger
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => deleteJob.mutate({ id: job.id })}
                      disabled={deleteJob.isPending}
                    >
                      Delete
                    </Button>
                  </div>
                </div>
                {job.lastRunAt && (
                  <div className="mt-2 text-gray-500 text-xs">
                    Last Run: {new Date(job.lastRunAt).toLocaleString()}
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
