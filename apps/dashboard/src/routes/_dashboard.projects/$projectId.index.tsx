import { createFileRoute, Link, useNavigate } from '@tanstack/react-router';
import { Activity, Box, Database, Folder, Plus, Server } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { ServiceIcon } from '#/components/ui/service-icon';
import { useGetCanvasSummary, useGetEnvironmentCanvas } from '#/hooks/useCanvas';
import { useGetProject } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/projects/$projectId/')({
  component: ProjectOverviewComponent,
});

const GithubIcon = ({ className }: { className?: string }) => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
    role="img"
    aria-label="GitHub"
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.2c3-.3 6-1.5 6-6.5 0-1.4-.5-2.5-1.5-3.4.1-.3.6-1.6-.1-3.3 0 0-1.2-.4-3.8 1.4a12.8 12.8 0 0 0-7 0C3.9 1.5 2.7 1.9 2.7 1.9c-.7 1.7-.2 3 .1 3.3-1 1-1.5 2-1.5 3.4 0 5 3 6.2 6 6.5-.4.4-.7 1-.8 2.2-.8.4-2.8.9-4-1.1 0 0-.7-1.3-2-1.4 0 0-1.3-.1-.1 1.2 0 0 1.2 1.8 3 2.5 1.5.5 3.3.4 3.3.4z" />
  </svg>
);

const getAppIcon = (runtimeMode: string) => {
  const mode = (runtimeMode || '').toLowerCase();
  if (mode.includes('docker') || mode.includes('container'))
    return <Box className="h-6 w-6 text-blue-500" />;
  if (mode.includes('node') || mode.includes('js'))
    return <Box className="h-6 w-6 text-yellow-500" />;
  if (mode.includes('go')) return <Box className="h-6 w-6 text-cyan-500" />;
  if (mode.includes('python')) return <Box className="h-6 w-6 text-blue-400" />;
  if (mode.includes('github')) return <GithubIcon className="h-6 w-6 text-foreground" />;
  return <Server className="h-6 w-6 text-primary" />;
};

const getDbIcon = (engine: string) => {
  const eng = (engine || '').toLowerCase();
  if (eng.includes('postgres')) return <Database className="h-6 w-6 text-blue-500" />;
  if (eng.includes('mysql') || eng.includes('mariadb'))
    return <Database className="h-6 w-6 text-blue-400" />;
  if (eng.includes('redis')) return <Database className="h-6 w-6 text-red-500" />;
  if (eng.includes('mongo')) return <Database className="h-6 w-6 text-green-500" />;
  return <Database className="h-6 w-6 text-primary" />;
};

function StatusBadge({ status }: { status: string }) {
  const s = (status || '').toLowerCase();
  let color = 'bg-gray-500/10 text-gray-500 border-gray-500/20';
  let dot = 'bg-gray-500';

  if (s === 'running' || s === 'online' || s === 'healthy') {
    color = 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
    dot = 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.4)]';
  } else if (s === 'failed' || s === 'error' || s === 'stopped') {
    color = 'bg-red-500/10 text-red-500 border-red-500/20';
    dot = 'bg-red-500';
  } else if (s === 'deploying' || s === 'pending' || s === 'building') {
    color = 'bg-amber-500/10 text-amber-500 border-amber-500/20';
    dot = 'bg-amber-500 animate-pulse';
  }

  return (
    <div
      className={`flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 font-medium text-xs ${color}`}
    >
      <div className={`h-1.5 w-1.5 rounded-full ${dot}`} />
      <span className="capitalize">{status || 'Unknown'}</span>
    </div>
  );
}

function ProjectOverviewComponent() {
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const { data: projectRes, isLoading: projectLoading } = useGetProject(projectId);
  const { data: summaryRes, isLoading: summaryLoading } = useGetCanvasSummary(projectId);

  const envId = summaryRes?.data?.defaultEnvironment?.id;
  const { data: envRes, isLoading: envLoading } = useGetEnvironmentCanvas(envId || '');

  if (projectLoading || summaryLoading || (envId && envLoading)) {
    return (
      <div className="flex h-full min-h-100 items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <Activity className="h-8 w-8 animate-pulse text-primary" />
          <p className="text-muted-foreground">Loading project workspace...</p>
        </div>
      </div>
    );
  }

  const project = projectRes?.data;
  const envData = envRes?.data;
  const apps = envData?.apps || [];
  const dbs = envData?.databases || [];

  const totalResources = apps.length + dbs.length;

  return (
    <div className="flex flex-col gap-8 p-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="font-bold text-3xl tracking-tight">
            {project?.name || 'Project Overview'}
          </h1>
          <div className="mt-2 flex items-center gap-3 text-muted-foreground text-sm">
            <div className="flex items-center gap-1.5">
              <Folder className="h-4 w-4" />
              <span>{project?.description || 'No description provided'}</span>
            </div>
            {summaryRes?.data?.defaultEnvironment && (
              <>
                <span>&bull;</span>
                <span className="flex items-center gap-1.5">
                  <div className="h-2 w-2 rounded-full bg-emerald-500" />
                  {summaryRes.data.defaultEnvironment.name} Environment
                </span>
              </>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="lg"
            className="shadow-sm"
            onClick={() => navigate({ to: '/projects/$projectId/settings', params: { projectId } })}
          >
            Settings
          </Button>
          <Button
            size="lg"
            className="shadow-sm"
            onClick={() => navigate({ to: '/projects/$projectId/new', params: { projectId } })}
          >
            <Plus className="mr-2 h-5 w-5" />
            New Resource
          </Button>
        </div>
      </div>

      <div className="space-y-4">
        <div className="flex items-center justify-between border-b pb-2">
          <h2 className="font-semibold text-lg">Resources ({totalResources})</h2>
        </div>

        {totalResources === 0 ? (
          <div className="flex h-64 flex-col items-center justify-center rounded-2xl border border-dashed bg-card/30">
            <Box className="mb-4 h-12 w-12 text-muted-foreground/50" />
            <h3 className="font-semibold text-lg">No resources yet</h3>
            <p className="mt-1 mb-6 max-w-md text-center text-muted-foreground text-sm">
              This project is empty. Deploy a new application, database, or service to get started.
            </p>
            <Button
              onClick={() => navigate({ to: '/projects/$projectId/new', params: { projectId } })}
            >
              <Plus className="mr-2 h-4 w-4" />
              Deploy First Resource
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {apps.map((app: any) => (
              <Link
                key={app.id}
                to="/services/$serviceId"
                params={{ serviceId: app.id }}
                className="group flex flex-col rounded-xl border bg-card p-5 transition-all hover:border-primary/50 hover:shadow-md"
              >
                <div className="flex items-start justify-between">
                  <div className="flex h-12 w-12 items-center justify-center overflow-hidden rounded-lg border bg-background/50 shadow-sm transition-colors group-hover:border-primary/20 group-hover:bg-primary/5">
                    {app.icon && app.icon !== 'git' ? (
                      <ServiceIcon icon={app.icon} className="h-8 w-8" />
                    ) : (
                      getAppIcon(app.runtimeMode || app.sourceType)
                    )}
                  </div>
                  <StatusBadge status={app.status || app.deploymentStatus} />
                </div>

                <div className="mt-4 flex-1">
                  <h3 className="font-semibold text-foreground transition-colors group-hover:text-primary">
                    {app.name}
                  </h3>
                  <p className="mt-1 line-clamp-1 text-muted-foreground text-xs uppercase tracking-wider">
                    {app.runtimeMode || 'Application'}
                  </p>
                </div>

                <div className="mt-5 border-t pt-4">
                  <div className="flex items-center justify-between text-muted-foreground text-xs">
                    <span className="flex items-center gap-1.5">
                      <Server className="h-3.5 w-3.5" />
                      {app.internalPort ? `Port ${app.internalPort}` : 'No Port'}
                    </span>
                  </div>
                </div>
              </Link>
            ))}

            {dbs.map((db: any) => (
              <div
                key={db.id}
                className="group flex cursor-default flex-col rounded-xl border bg-card p-5 transition-all hover:border-primary/50 hover:shadow-md"
              >
                <div className="flex items-start justify-between">
                  <div className="flex h-12 w-12 items-center justify-center rounded-lg border bg-background/50 shadow-sm transition-colors group-hover:border-primary/20 group-hover:bg-primary/5">
                    {getDbIcon(db.engine)}
                  </div>
                  <StatusBadge status={db.status} />
                </div>

                <div className="mt-4 flex-1">
                  <h3 className="font-semibold text-foreground transition-colors group-hover:text-primary">
                    {db.name}
                  </h3>
                  <p className="mt-1 line-clamp-1 text-muted-foreground text-xs uppercase tracking-wider">
                    {db.engine || 'Database'}
                  </p>
                </div>

                <div className="mt-5 border-t pt-4">
                  <div className="flex items-center justify-between text-muted-foreground text-xs">
                    <span className="flex items-center gap-1.5">
                      <Database className="h-3.5 w-3.5" />
                      {db.internalPort ? `Port ${db.internalPort}` : 'Default Port'}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
