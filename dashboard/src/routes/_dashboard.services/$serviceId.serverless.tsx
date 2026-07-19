import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { useGetApp } from '#/hooks/useApps';
import { useGetCode } from '#/hooks/useServerless';

export const Route = createFileRoute('/_dashboard/services/$serviceId/serverless')({
  component: ServiceServerlessRoute,
});

function ServiceServerlessRoute() {
  const { serviceId } = Route.useParams();
  const { data: appData, isLoading: appLoading } = useGetApp(serviceId);
  const { data: codeData, isLoading: codeLoading } = useGetCode(serviceId);

  if (appLoading || codeLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const app = appData?.data;
  const codeInfo = codeData?.data;

  if (!app) {
    return <div>Service not found.</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="font-bold text-2xl">Serverless Functions</h1>
      <div className="rounded border p-4 shadow-sm">
        <h2 className="font-semibold text-lg">Code</h2>
        <pre className="mt-2 rounded bg-muted p-4">
          {codeInfo?.codeContent || '// No code found'}
        </pre>
      </div>
    </div>
  );
}
