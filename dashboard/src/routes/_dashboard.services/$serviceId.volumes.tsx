import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import { HardDrive, Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useGetApp } from '#/hooks/useApps';
import { apiClient } from '#/lib/apiClient';

export const Route = createFileRoute('/_dashboard/services/$serviceId/volumes')({
  component: ServiceVolumesRoute,
});

function ServiceVolumesRoute() {
  const { serviceId } = Route.useParams();
  const queryClient = useQueryClient();
  const { data: appData } = useGetApp(serviceId);
  const app = appData?.data;

  const [isCreating, setIsCreating] = useState(false);
  const [hostPath, setHostPath] = useState('');
  const [containerPath, setContainerPath] = useState('');

  const { data: volumesResponse, isLoading: isLoadingVolumes } = useQuery({
    queryKey: ['volumes', 'service', serviceId],
    queryFn: async () => {
      const res: any = await apiClient.get(`/apps/${serviceId}/volumes`);
      return res.data;
    },
  });

  const volumes = volumesResponse?.data || [];

  const createMutation = useMutation({
    mutationFn: async (data: any) => {
      return apiClient.post(`/apps/${serviceId}/volumes`, data);
    },
    onSuccess: () => {
      toast.success('Volume mapped successfully');
      queryClient.invalidateQueries({ queryKey: ['volumes', 'service', serviceId] });
      setIsCreating(false);
      setHostPath('');
      setContainerPath('');
    },
    onError: () => {
      toast.error('Failed to map volume');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (volumeId: string) => {
      return apiClient.delete(`/apps/${serviceId}/volumes/${volumeId}`);
    },
    onSuccess: () => {
      toast.success('Volume unmapped successfully');
      queryClient.invalidateQueries({ queryKey: ['volumes', 'service', serviceId] });
    },
    onError: () => {
      toast.error('Failed to unmap volume');
    },
  });

  const handleCreate = () => {
    if (!hostPath || !containerPath) {
      toast.error('Please fill all fields');
      return;
    }
    createMutation.mutate({
      hostPath,
      containerPath,
    });
  };

  if (!app) return null;

  return (
    <div className="mx-auto max-w-5xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-bold text-2xl">Persistent Storage</h2>
          <p className="text-muted-foreground">
            Mount persistent Docker volumes so data survives container restarts.
          </p>
        </div>
        <Button onClick={() => setIsCreating(!isCreating)} className="gap-2">
          <Plus className="h-4 w-4" />
          Add Volume
        </Button>
      </div>

      {isCreating && (
        <Card className="border-primary/20 bg-primary/5">
          <CardHeader>
            <CardTitle>Map New Volume</CardTitle>
            <CardDescription>
              Specify the host path (or named volume) and the path inside the container.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label>Host Path / Volume Name</Label>
                <Input
                  placeholder="e.g., /opt/vessl/data or my-app-data"
                  value={hostPath}
                  onChange={(e) => setHostPath(e.target.value)}
                  className="font-mono"
                />
              </div>
              <div className="space-y-2">
                <Label>Container Path</Label>
                <Input
                  placeholder="e.g., /var/www/html/storage"
                  value={containerPath}
                  onChange={(e) => setContainerPath(e.target.value)}
                  className="font-mono"
                />
              </div>
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button variant="ghost" onClick={() => setIsCreating(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreate} disabled={createMutation.isPending}>
                {createMutation.isPending ? 'Mapping...' : 'Map Volume'}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardContent className="p-0">
          {isLoadingVolumes ? (
            <div className="flex justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : volumes.length === 0 ? (
            <div className="flex flex-col items-center justify-center p-12 text-center text-muted-foreground">
              <HardDrive className="mb-4 h-8 w-8 opacity-20" />
              <p>No volumes mapped for this service.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Host Path</TableHead>
                  <TableHead>Container Path</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {volumes.map((vol: any) => (
                  <TableRow key={vol.id}>
                    <TableCell className="font-mono text-sm">{vol.hostPath}</TableCell>
                    <TableCell className="font-mono text-sm">{vol.containerPath}</TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="destructive"
                        size="icon"
                        title="Delete"
                        onClick={() => deleteMutation.mutate(vol.id)}
                        disabled={deleteMutation.isPending}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
