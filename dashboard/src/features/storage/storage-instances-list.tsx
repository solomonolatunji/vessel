import { formatDistanceToNow } from 'date-fns';
import { Database, HardDrive, Play, Square, Trash } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { useDelete, useList, useStart, useStop } from '#/hooks/useStorage';
import type { Storage } from '#/interfaces/storage';
import { CreateStorageModal } from './create-storage-modal';

export const StorageInstancesList = () => {
  const { data: storageData, isLoading } = useList();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const startMutation = useStart();
  const stopMutation = useStop();
  const deleteMutation = useDelete();

  const handleStart = async (id: string) => {
    try {
      await startMutation.mutateAsync({ id });
      toast.success('Storage container started');
    } catch {
      toast.error('Failed to start container');
    }
  };

  const handleStop = async (id: string) => {
    try {
      await stopMutation.mutateAsync({ id });
      toast.success('Storage container stopped');
    } catch {
      toast.error('Failed to stop container');
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this storage container?')) return;
    try {
      await deleteMutation.mutateAsync({ id });
      toast.success('Storage container deleted');
    } catch {
      toast.error('Failed to delete container');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-6 pb-2 md:flex-row md:items-start">
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              STORAGE
            </p>
            <h1 className="font-bold text-3xl tracking-tight">Storage Instances</h1>
          </div>
          <p className="max-w-2xl text-muted-foreground text-sm leading-relaxed">
            Manage your self-hosted S3-compatible storage containers (MinIO).
          </p>
        </div>
        <div className="flex shrink-0 items-center gap-3">
          <Button onClick={() => setIsCreateModalOpen(true)}>Create Instance</Button>
        </div>
      </div>

      <div className="rounded-xl border border-border/50 bg-card/40">
        <div className="border-border/50 border-b p-4">
          <h2 className="font-semibold text-lg">Instances</h2>
        </div>
        <div className="p-4">
          {isLoading ? (
            <div className="space-y-4">
              {[...Array(3)].map((_, i) => (
                <Skeleton key={i} className="h-20 w-full rounded-xl" />
              ))}
            </div>
          ) : !storageData?.data?.length ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
                <HardDrive className="h-6 w-6" />
              </div>
              <h3 className="font-semibold text-lg">No storage instances</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                Create a MinIO storage container to store object data.
              </p>
            </div>
          ) : (
            <div className="grid gap-4">
              {storageData.data.map((instance: Storage) => (
                <div
                  key={instance.id}
                  className="flex flex-col gap-4 rounded-xl border border-border/50 bg-background/50 p-4 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="flex items-start gap-4">
                    <div className="mt-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <Database className="h-5 w-5" />
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <h4 className="font-semibold">{instance.name}</h4>
                        <Badge variant={instance.status === 'running' ? 'default' : 'secondary'}>
                          {instance.status}
                        </Badge>
                      </div>
                      <div className="mt-1 flex flex-col gap-1 text-muted-foreground text-sm">
                        <span>Bucket: {instance.bucketName}</span>
                        <span>Created {formatDistanceToNow(new Date(instance.createdAt))} ago</span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {instance.status === 'stopped' ? (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStart(instance.id)}
                        disabled={startMutation.isPending}
                      >
                        <Play className="mr-2 h-4 w-4" />
                        Start
                      </Button>
                    ) : (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleStop(instance.id)}
                        disabled={stopMutation.isPending}
                      >
                        <Square className="mr-2 h-4 w-4" />
                        Stop
                      </Button>
                    )}
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDelete(instance.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {isCreateModalOpen && <CreateStorageModal onClose={() => setIsCreateModalOpen(false)} />}
    </div>
  );
};
