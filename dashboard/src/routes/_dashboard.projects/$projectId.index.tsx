import { createFileRoute } from '@tanstack/react-router';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { useGetCanvasSummary, useGetEnvironmentCanvas } from '#/hooks/useCanvas';
import { useGetProject } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/projects/$projectId/')({
  component: ProjectOverviewComponent,
});

function ProjectOverviewComponent() {
  const { projectId } = Route.useParams();
  const { data: projectRes, isLoading: projectLoading } = useGetProject(projectId);
  const { data: summaryRes, isLoading: summaryLoading } = useGetCanvasSummary(projectId);

  const envId = summaryRes?.data?.defaultEnvironment?.id;
  const { data: envRes, isLoading: envLoading } = useGetEnvironmentCanvas(envId || '');

  if (projectLoading || summaryLoading || (envId && envLoading)) {
    return <div className="p-4">Loading project...</div>;
  }

  const project = projectRes?.data;
  const envData = envRes?.data;

  return (
    <div className="flex flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">{project?.name || 'Project Overview'}</h1>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Apps / Services</CardTitle>
          </CardHeader>
          <CardContent>
            {envData?.apps && envData.apps.length > 0 ? (
              <ul className="space-y-2">
                {envData.apps.map((app) => (
                  <li key={app.id} className="flex items-center justify-between rounded border p-2">
                    <div>
                      <div className="font-semibold">{app.name}</div>
                      <div className="text-gray-500 text-sm">
                        {app.runtimeMode} | {app.status}
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-gray-500 text-sm">No apps found.</p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Databases</CardTitle>
          </CardHeader>
          <CardContent>
            {envData?.databases && envData.databases.length > 0 ? (
              <ul className="space-y-2">
                {envData.databases.map((db) => (
                  <li key={db.id} className="flex items-center justify-between rounded border p-2">
                    <div>
                      <div className="font-semibold">{db.name}</div>
                      <div className="text-gray-500 text-sm">
                        {db.engine} | {db.status}
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-gray-500 text-sm">No databases found.</p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
