import { createFileRoute } from '@tanstack/react-router';
import { Calendar, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useListJobs } from '#/hooks/useJobs';
import { useListProjects } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/jobs')({
  component: JobsPage,
});

function JobsPage() {
  const [selectedProjectId, setSelectedProjectId] = useState<string>('');

  const { data: projectsResponse, isLoading: isLoadingProjects } = useListProjects();
  const projects = projectsResponse?.data?.records || [];

  const { data: jobsResponse, isLoading: isLoadingJobs } = useListJobs();
  const jobs = jobsResponse?.data || [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Calendar className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Cron Jobs</h1>
            <p className="text-muted-foreground text-sm">Manage and monitor scheduled tasks.</p>
          </div>
        </div>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <CardTitle>All Jobs</CardTitle>
          <div className="w-[200px]">
            <Select value={selectedProjectId} onValueChange={setSelectedProjectId}>
              <SelectTrigger>
                <SelectValue placeholder="All Projects" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All Projects</SelectItem>
                {projects.map(
                  (project: any /* biome-ignore lint/suspicious/noExplicitAny: any */) => (
                    <SelectItem key={project.id} value={project.id}>
                      {project.name}
                    </SelectItem>
                  )
                )}
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {isLoadingProjects || isLoadingJobs ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : jobs.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <Calendar className="mb-4 h-8 w-8 opacity-20" />
              <p>No jobs found.</p>
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
                </TableRow>
              </TableHeader>
              <TableBody>
                {jobs.map((job) => (
                  <TableRow key={job.id}>
                    <TableCell className="font-medium">{job.name}</TableCell>
                    <TableCell>{job.schedule}</TableCell>
                    <TableCell className="font-mono text-xs">{job.command}</TableCell>
                    <TableCell>{job.status}</TableCell>
                    <TableCell>
                      {job.lastRunAt ? new Date(job.lastRunAt).toLocaleString() : 'Never'}
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
