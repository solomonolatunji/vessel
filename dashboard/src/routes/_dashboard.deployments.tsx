import { createFileRoute } from '@tanstack/react-router';
import { Loader2, Rocket } from 'lucide-react';
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
import { useListByProject } from '#/hooks/useApps';
import { useListByService } from '#/hooks/useDeployments';
import { useListProjects } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/deployments')({
  component: DeploymentsPage,
});

function DeploymentsPage() {
  const [selectedProjectId, setSelectedProjectId] = useState<string>('');
  const [selectedServiceId, setSelectedServiceId] = useState<string>('');

  const { data: projectsResponse, isLoading: isLoadingProjects } = useListProjects();
  const projects = projectsResponse?.data?.records || [];

  const { data: appsResponse, isLoading: isLoadingApps } = useListByProject(selectedProjectId);
  const apps = appsResponse?.data || [];

  const { data: deploymentsResponse, isLoading: isLoadingDeployments } =
    useListByService(selectedServiceId);
  const deploymentsPaginated = deploymentsResponse?.data;
  const deployments = deploymentsPaginated?.records || [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Rocket className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Deployments</h1>
            <p className="text-muted-foreground text-sm">
              View your application deployments across projects.
            </p>
          </div>
        </div>
      </div>

      <Card>
        <CardHeader className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
          <CardTitle>Deployment History</CardTitle>
          <div className="flex flex-col gap-2 sm:flex-row">
            <Select
              value={selectedProjectId}
              onValueChange={(val) => {
                setSelectedProjectId(val);
                setSelectedServiceId('');
              }}
            >
              <SelectTrigger className="w-[200px]">
                <SelectValue placeholder="Select Project" />
              </SelectTrigger>
              <SelectContent>
                {projects.map(
                  (project: any /* biome-ignore lint/suspicious/noExplicitAny: any */) => (
                    <SelectItem key={project.id} value={project.id}>
                      {project.name}
                    </SelectItem>
                  )
                )}
              </SelectContent>
            </Select>

            <Select
              value={selectedServiceId}
              onValueChange={setSelectedServiceId}
              disabled={!selectedProjectId || apps.length === 0}
            >
              <SelectTrigger className="w-[200px]">
                <SelectValue placeholder="Select App" />
              </SelectTrigger>
              <SelectContent>
                {apps.map((app: any /* biome-ignore lint/suspicious/noExplicitAny: any */) => (
                  <SelectItem key={app.id} value={app.id}>
                    {app.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {isLoadingProjects || isLoadingApps || isLoadingDeployments ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : !selectedServiceId ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <Rocket className="mb-4 h-8 w-8 opacity-20" />
              <p>Select a project and app to view deployments.</p>
            </div>
          ) : deployments.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <Rocket className="mb-4 h-8 w-8 opacity-20" />
              <p>No deployments found for this app.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Status</TableHead>
                  <TableHead>Branch</TableHead>
                  <TableHead>Commit</TableHead>
                  <TableHead>Trigger</TableHead>
                  <TableHead>Created</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {deployments.map(
                  (deployment: any /* biome-ignore lint/suspicious/noExplicitAny: any */) => (
                    <TableRow key={deployment.id}>
                      <TableCell className="font-medium">{deployment.status}</TableCell>
                      <TableCell>{deployment.branch || '-'}</TableCell>
                      <TableCell className="max-w-[100px] truncate font-mono text-xs">
                        {deployment.commitHash ? deployment.commitHash.substring(0, 7) : '-'}
                      </TableCell>
                      <TableCell>{deployment.trigger || '-'}</TableCell>
                      <TableCell>{new Date(deployment.createdAt).toLocaleString()}</TableCell>
                    </TableRow>
                  )
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
