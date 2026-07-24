import { Loader2 } from 'lucide-react';
import { useListByService } from '#/hooks/useDeployments';
import { DeploymentFailureAi } from './deployment-failure-ai';

export function ServiceDeployments({ app }: { app: any }) {
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
          {deployments.map((dep: any) => (
            <div key={dep.id} className="rounded border p-3 text-sm">
              <div className="mb-2 flex items-center justify-between">
                <div>
                  <span
                    className={`rounded-full px-2 py-0.5 font-medium text-xs ${
                      dep.status === 'READY'
                        ? 'bg-green-100 text-green-700'
                        : dep.status === 'FAILED'
                          ? 'bg-red-100 text-red-700'
                          : 'bg-blue-100 text-blue-700'
                    }`}
                  >
                    {dep.status}
                  </span>
                  <span className="ml-2 text-gray-600">
                    {dep.commitMessage || 'Manual deployment'}
                  </span>
                </div>
                <span className="text-gray-400 text-xs">
                  {new Date(dep.createdAt).toLocaleString()}
                </span>
              </div>

              {dep.status === 'FAILED' && <DeploymentFailureAi deploymentId={dep.id} />}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
