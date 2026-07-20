import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import { Calendar, Loader2, Play, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useGetApp } from '#/hooks/useApps';
import { apiClient } from '#/lib/apiClient';

export const Route = createFileRoute('/_dashboard/services/$serviceId/scheduled-tasks')({
  component: ServiceJobsRoute,
});

function ServiceJobsRoute() {
  const { serviceId } = Route.useParams();
  const queryClient = useQueryClient();
  const { data: appData } = useGetApp(serviceId);
  const app = appData?.data;

  const [isCreating, setIsCreating] = useState(false);
  const [name, setName] = useState('');
  const [schedule, setSchedule] = useState('0 0 * * *');
  const [command, setCommand] = useState('');

  const { data: jobsResponse, isLoading: isLoadingJobs } = useQuery({
    queryKey: ['jobs', 'service', serviceId],
    queryFn: async () => {
      const res: any = await apiClient.get(`/jobs?serviceId=${serviceId}`);
      return res.data;
    },
  });

  const jobs = jobsResponse || [];

  const createMutation = useMutation({
    mutationFn: async (data: any) => {
      return apiClient.post('/jobs', data);
    },
    onSuccess: () => {
      toast.success('Job created successfully');
      queryClient.invalidateQueries({ queryKey: ['jobs', 'service', serviceId] });
      setIsCreating(false);
      setName('');
      setSchedule('0 0 * * *');
      setCommand('');
    },
    onError: () => {
      toast.error('Failed to create job');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (jobId: string) => {
      return apiClient.delete(`/jobs/${jobId}`);
    },
    onSuccess: () => {
      toast.success('Job deleted successfully');
      queryClient.invalidateQueries({ queryKey: ['jobs', 'service', serviceId] });
    },
    onError: () => {
      toast.error('Failed to delete job');
    },
  });

  const runMutation = useMutation({
    mutationFn: async (jobId: string) => {
      return apiClient.post(`/jobs/${jobId}/trigger`);
    },
    onSuccess: () => {
      toast.success('Job triggered successfully');
      queryClient.invalidateQueries({ queryKey: ['jobs', 'service', serviceId] });
    },
    onError: () => {
      toast.error('Failed to trigger job');
    },
  });

  const handleCreate = () => {
    if (!name || !schedule || !command) {
      toast.error('Please fill all fields');
      return;
    }
    createMutation.mutate({
      projectId: app?.projectId,
      serviceId: serviceId,
      name,
      schedule,
      command,
    });
  };

  if (!app) return null;

  return (
    <div className="mx-auto max-w-5xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-bold text-2xl">Scheduled Tasks</h2>
          <p className="text-muted-foreground">Run cron jobs inside this service's container.</p>
        </div>
        <Button onClick={() => setIsCreating(!isCreating)} className="gap-2">
          <Plus className="h-4 w-4" />
          Add Job
        </Button>
      </div>

      {isCreating && (
        <Card className="border-primary/20 bg-primary/5">
          <CardHeader>
            <CardTitle>Create New Cron Job</CardTitle>
            <CardDescription>Jobs run inside the active container of this service.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label>Job Name</Label>
                <Input
                  placeholder="e.g., Daily DB Backup"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Cron Expression</Label>
                <Input
                  placeholder="0 0 * * *"
                  value={schedule}
                  onChange={(e) => setSchedule(e.target.value)}
                  className="font-mono"
                />
              </div>
              <div className="space-y-2 md:col-span-2">
                <Label>Command to Run</Label>
                <Input
                  placeholder="e.g., php artisan schedule:run or npm run db:prune"
                  value={command}
                  onChange={(e) => setCommand(e.target.value)}
                  className="font-mono"
                />
              </div>
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button variant="ghost" onClick={() => setIsCreating(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreate} disabled={createMutation.isPending}>
                {createMutation.isPending ? 'Creating...' : 'Create Job'}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardContent className="p-0">
          {isLoadingJobs ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : jobs.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <Calendar className="mb-4 h-8 w-8 opacity-20" />
              <p>No scheduled tasks configured.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Schedule</TableHead>
                  <TableHead>Command</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Last Run</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {jobs.map((job: any) => (
                  <TableRow key={job.id}>
                    <TableCell className="font-medium">{job.name}</TableCell>
                    <TableCell className="font-mono text-xs">{job.schedule}</TableCell>
                    <TableCell className="font-mono text-xs">{job.command}</TableCell>
                    <TableCell>{job.status}</TableCell>
                    <TableCell className="text-muted-foreground text-xs">
                      {job.lastRunAt ? new Date(job.lastRunAt).toLocaleString() : 'Never'}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          variant="outline"
                          size="icon"
                          title="Run now"
                          onClick={() => runMutation.mutate(job.id)}
                          disabled={runMutation.isPending}
                        >
                          <Play className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="destructive"
                          size="icon"
                          title="Delete"
                          onClick={() => deleteMutation.mutate(job.id)}
                          disabled={deleteMutation.isPending}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
