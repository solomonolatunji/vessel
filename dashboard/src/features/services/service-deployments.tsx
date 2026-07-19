import { Loader2 } from 'lucide-react';
import { useListByService } from '#/hooks/useDeployments';

export function ServiceDeployments({
  app,
}: {
  app: any /* biome-ignore lint/suspicious/noExplicitAny: any */;
}) {
  const { data, isLoading } = useListByService(app.id);

  if (isLoading) {
    return <Loader2 className="h-6 w-6 animate-spin text-gray-500" />;
  }

  const deployments = data?.data?.records || [];

  return (
    <div className="space-y-4 rounded border p-4 shadow-sm">
      <h2 className="font-semibold text-lg">Deployments for {app.name}</h2>
      {deployments.length === 0 ? (
        <p className="text-gray-500 text-sm">No deployments found.</p>
      ) : (
        <div className="space-y-2">
          {deployments.map((dep: any /* biome-ignore lint/suspicious/noExplicitAny: any */) => (
            <div key={dep.id} className="rounded border p-2 text-sm">
              <span className="font-medium">{dep.status}</span> -{' '}
              {dep.commitMessage || 'Manual deployment'}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
