import { useListByService } from '#/hooks/useDeployments';
import { Loader2 } from 'lucide-react';

export function ServiceDeployments({ app }: { app: any }) {
  const { data, isLoading } = useListByService(app.id);

  if (isLoading) {
    return <Loader2 className="animate-spin w-6 h-6 text-gray-500" />;
  }

  const deployments = data?.data?.records || [];

  return (
    <div className="rounded border p-4 shadow-sm space-y-4">
      <h2 className="font-semibold text-lg">Deployments for {app.name}</h2>
      {deployments.length === 0 ? (
        <p className="text-sm text-gray-500">No deployments found.</p>
      ) : (
        <div className="space-y-2">
          {deployments.map((dep: any) => (
            <div key={dep.id} className="p-2 border rounded text-sm">
              <span className="font-medium">{dep.status}</span> - {dep.commitMessage || 'Manual deployment'}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
