import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { apiClient } from '#/lib/apiClient';

interface RuntimeModeCardProps {
  serviceId: string;
  initialData?: {
    runtimeMode?: 'web' | 'worker';
    internalPort?: number;
  };
}

export function RuntimeModeCard({ serviceId, initialData }: RuntimeModeCardProps) {
  const queryClient = useQueryClient();
  const [runtimeMode, setRuntimeMode] = useState<'web' | 'worker'>(
    initialData?.runtimeMode || 'web'
  );
  const [internalPort, setInternalPort] = useState(initialData?.internalPort || 3000);

  const updateMutation = useMutation({
    mutationFn: async (data: any) => {
      return apiClient.put(`/services/${serviceId}/runtime`, data);
    },
    onSuccess: () => {
      toast.success('Runtime mode updated');
      queryClient.invalidateQueries({ queryKey: ['service', serviceId] });
    },
  });

  const handleSave = () => {
    updateMutation.mutate({
      runtimeMode,
      internalPort: runtimeMode === 'web' ? internalPort : undefined,
    });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Runtime Mode</CardTitle>
        <CardDescription>Choose how this service behaves inside the cluster.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div
            className={`cursor-pointer rounded-lg border p-4 transition-colors ${
              runtimeMode === 'web'
                ? 'border-zinc-900 bg-zinc-50 dark:border-zinc-100 dark:bg-zinc-900'
                : 'border-zinc-200 hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700'
            }`}
            onClick={() => setRuntimeMode('web')}
            onKeyDown={(e) => {
              if (e.key === 'Enter') setRuntimeMode('web');
            }}
            tabIndex={0}
            role="button"
          >
            <div className="mb-2 flex items-center justify-between">
              <h3 className="font-semibold">Web Service</h3>
              {runtimeMode === 'web' && <Badge variant="default">Active</Badge>}
            </div>
            <p className="text-sm text-zinc-500">
              Exposes a public HTTP route and binds to an internal port. Monitored via /healthz.
            </p>
          </div>

          <div
            className={`cursor-pointer rounded-lg border p-4 transition-colors ${
              runtimeMode === 'worker'
                ? 'border-zinc-900 bg-zinc-50 dark:border-zinc-100 dark:bg-zinc-900'
                : 'border-zinc-200 hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700'
            }`}
            onClick={() => setRuntimeMode('worker')}
            onKeyDown={(e) => {
              if (e.key === 'Enter') setRuntimeMode('worker');
            }}
            tabIndex={0}
            role="button"
          >
            <div className="mb-2 flex items-center justify-between">
              <h3 className="font-semibold">Background Worker</h3>
              {runtimeMode === 'worker' && <Badge variant="default">Active</Badge>}
            </div>
            <p className="text-sm text-zinc-500">
              No exposed public route. Runs as a background process and monitored via uptime checks.
            </p>
          </div>
        </div>

        {runtimeMode === 'web' && (
          <div className="space-y-2 rounded-md bg-zinc-50 p-4 dark:bg-zinc-900/50">
            <Label htmlFor="internalPort">Internal Port</Label>
            <Input
              id="internalPort"
              type="number"
              value={internalPort}
              onChange={(e) => setInternalPort(parseInt(e.target.value, 10))}
              className="max-w-[200px]"
            />
            <p className="text-xs text-zinc-500">
              The port your application server listens on (e.g., 3000, 8080).
            </p>
          </div>
        )}

        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={updateMutation.isPending}>
            {updateMutation.isPending ? 'Saving...' : 'Save Mode'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
