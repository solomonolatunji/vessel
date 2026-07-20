import { createFileRoute, Link, Outlet, useLocation } from '@tanstack/react-router';
import {
  Activity,
  AlertTriangle,
  ArrowLeft,
  BarChart2,
  Calendar,
  Code,
  GitPullRequest,
  Globe,
  HardDrive,
  Loader2,
  Settings,
  Terminal,
  Variable,
  Webhook,
  Wrench,
} from 'lucide-react';
import { Button } from '#/components/ui/button';
import { ServiceIcon } from '#/components/ui/service-icon';
import { useGetApp } from '#/hooks/useApps';
import { useGetDatabase } from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/services/$serviceId')({
  component: ServiceLayoutRoute,
});

function ServiceLayoutRoute() {
  const { serviceId } = Route.useParams();
  const location = useLocation();

  const { data: appData, isLoading: appLoading } = useGetApp(serviceId);
  const { data: dbData, isLoading: dbLoading } = useGetDatabase(serviceId);

  if (appLoading || dbLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;
  const db = dbData?.data;

  // Tabs for Application Services
  const appTabs = [
    { name: 'Configuration', href: `/services/${serviceId}`, icon: Settings, exact: true },
    { name: 'Build Settings', href: `/services/${serviceId}/build`, icon: Wrench },
    { name: 'Deployments', href: `/services/${serviceId}/deployments`, icon: Activity },
    { name: 'Metrics', href: `/services/${serviceId}/metrics`, icon: BarChart2 },
    { name: 'Webhooks', href: `/services/${serviceId}/webhooks`, icon: Webhook },
    { name: 'Scheduled Tasks', href: `/services/${serviceId}/scheduled-tasks`, icon: Calendar },
    { name: 'Storage', href: `/services/${serviceId}/volumes`, icon: HardDrive },
    { name: 'Domains', href: `/services/${serviceId}/domains`, icon: Globe },
    { name: 'Variables', href: `/services/${serviceId}/variables`, icon: Variable },
    { name: 'Terminal', href: `/services/${serviceId}/terminal`, icon: Terminal },
    { name: 'Serverless Editor', href: `/services/${serviceId}/serverless`, icon: Code },
    { name: 'PR Previews', href: `/services/${serviceId}/previews`, icon: GitPullRequest },
    { name: 'Danger Zone', href: `/services/${serviceId}/danger`, icon: AlertTriangle },
  ];

  // Tabs for Databases
  const dbTabs = [
    { name: 'Overview', href: `/services/${serviceId}`, icon: Settings, exact: true },
    { name: 'Danger Zone', href: `/services/${serviceId}/danger`, icon: AlertTriangle },
  ];

  const tabs = app ? appTabs : dbTabs;
  const resourceName = app?.name || db?.name || 'Service';
  const projectId = app?.projectId || db?.projectId;

  return (
    <div className="flex min-h-screen flex-col">
      <div className="border-b bg-card">
        <div className="px-6 py-4">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" asChild className="h-8 w-8">
              <Link to="/projects/$projectId" params={{ projectId: projectId as string }}>
                <ArrowLeft className="h-4 w-4" />
              </Link>
            </Button>
            {app?.icon && app.icon !== 'git' && (
              <ServiceIcon icon={app.icon} className="h-10 w-10 rounded-lg border" />
            )}
            <div>
              <h1 className="font-bold text-xl">{resourceName}</h1>
              <p className="text-muted-foreground text-sm">
                {app ? 'Application Service' : 'Database'}
              </p>
            </div>
          </div>
        </div>

        <div className="flex space-x-1 overflow-x-auto px-6">
          {tabs.map((tab) => {
            const isActive = tab.exact
              ? location.pathname === tab.href || location.pathname === `${tab.href}/`
              : location.pathname.startsWith(tab.href);
            return (
              <Link
                key={tab.name}
                to={tab.href}
                className={`flex items-center gap-2 whitespace-nowrap border-b-2 px-4 py-3 font-medium text-sm transition-colors ${
                  isActive
                    ? 'border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:border-muted hover:text-foreground'
                }`}
              >
                <tab.icon className="h-4 w-4" />
                {tab.name}
              </Link>
            );
          })}
        </div>
      </div>

      <div className="flex-1 p-6">
        <div className="mx-auto max-w-6xl">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
