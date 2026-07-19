import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { useGetApp } from '#/hooks/useApps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/domains')({
  component: ServiceDomainsRoute,
});

function ServiceDomainsRoute() {
  const { serviceId } = Route.useParams();
  const { data: appData, isLoading } = useGetApp(serviceId);

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;

  if (!app) {
    return <div>Service not found.</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="font-bold text-2xl">Domains</h1>
      <p>Domain management for {app.name} is coming soon.</p>
    </div>
  );
}
