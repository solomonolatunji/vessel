import { format } from 'date-fns';
import { Database, Download, Play, RefreshCw, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogTrigger } from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreate, useDelete, useList, useListRecords, useTrigger } from '#/hooks/useBackups';

export function BackupsList() {
  const { data: configsData, isLoading } = useList('global');
  const configs = configsData?.data || [];

  const createBackup = useCreate();
  const deleteBackup = useDelete();
  const triggerBackup = useTrigger();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [name, setName] = useState('');
  const [schedule, setSchedule] = useState('0 0 * * *');
  const [retentionDays, setRetentionDays] = useState('7');

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createBackup.mutateAsync({
        payload: {
          projectId: 'global',
          name,
          schedule,
          retentionDays: parseInt(retentionDays, 10),
        },
      });
      toast.success('Backup configuration created');
      setIsCreateOpen(false);
      setName('');
    } catch {
      toast.error('Failed to create backup configuration');
    }
  };

  const handleTrigger = async (id: string) => {
    try {
      await triggerBackup.mutateAsync({ id });
      toast.success('Backup triggered successfully');
    } catch {
      toast.error('Failed to trigger backup');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteBackup.mutateAsync({ id });
      toast.success('Backup configuration deleted');
    } catch {
      toast.error('Failed to delete backup configuration');
    }
  };

  if (isLoading) {
    return <div className="p-6 text-muted-foreground">Loading backups...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-6 pb-2 md:flex-row md:items-start">
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              STORAGE & BACKUPS
            </p>
            <h1 className="font-bold text-3xl tracking-tight">System Backups</h1>
          </div>
          <p className="max-w-2xl text-muted-foreground text-sm leading-relaxed">
            Manage your scheduled system backups, view history, and trigger manual backups.
          </p>
        </div>

        <div className="flex shrink-0 flex-col items-end gap-4">
          <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
            <DialogTrigger asChild>
              <Button className="h-11 px-8 font-bold text-xs uppercase tracking-wider">
                Create Backup Config
              </Button>
            </DialogTrigger>
            <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
              <form onSubmit={handleCreate} className="flex flex-col gap-6 p-6">
                <div>
                  <h2 className="font-bold text-lg">New Backup Configuration</h2>
                  <p className="text-muted-foreground text-sm">Schedule a recurring backup.</p>
                </div>

                <div className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="name">Name</Label>
                    <Input
                      id="name"
                      placeholder="e.g. Daily DB Backup"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="schedule">Cron Schedule</Label>
                    <Input
                      id="schedule"
                      placeholder="0 0 * * *"
                      value={schedule}
                      onChange={(e) => setSchedule(e.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="retentionDays">Retention Days</Label>
                    <Input
                      id="retentionDays"
                      type="number"
                      min="1"
                      value={retentionDays}
                      onChange={(e) => setRetentionDays(e.target.value)}
                      required
                    />
                  </div>
                </div>

                <div className="flex justify-end gap-3 pt-4">
                  <Button type="button" variant="ghost" onClick={() => setIsCreateOpen(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={createBackup.isPending}>
                    {createBackup.isPending ? 'Creating...' : 'Create'}
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <div className="grid gap-4">
        {configs.length === 0 ? (
          <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
            No backup configurations found.
          </div>
        ) : (
          configs.map((config) => (
            <div
              key={config.id}
              className="group relative flex flex-col justify-between gap-4 overflow-hidden rounded-xl border border-border/50 bg-card/50 p-6 transition-all hover:border-border hover:bg-card md:flex-row md:items-center"
            >
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
                  <Database className="h-5 w-5" />
                </div>
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold">{config.name}</h3>
                    <Badge variant="secondary" className="text-[10px]">
                      {config.status}
                    </Badge>
                  </div>
                  <p className="text-muted-foreground text-sm">
                    Schedule:{' '}
                    <code className="rounded bg-muted px-1.5 py-0.5 text-xs">
                      {config.schedule}
                    </code>{' '}
                    • Retains {config.retentionDays} days
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-2">
                <RecordsDialog configId={config.id} configName={config.name} />

                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleTrigger(config.id)}
                  disabled={triggerBackup.isPending}
                >
                  <Play className="mr-2 h-4 w-4" />
                  Run Now
                </Button>
                <Button
                  variant="destructive"
                  size="icon"
                  className="h-9 w-9 opacity-0 transition-opacity group-hover:opacity-100"
                  onClick={() => handleDelete(config.id)}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

function RecordsDialog({ configId, configName }: { configId: string; configName: string }) {
  const [isOpen, setIsOpen] = useState(false);
  const { data: recordsData, isLoading } = useListRecords(configId);
  const records = recordsData?.data || [];

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <RefreshCw className="mr-2 h-4 w-4" />
          History
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-3xl gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl [&>button]:hidden">
        <div className="flex flex-col gap-6 p-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="font-bold text-lg">{configName} History</h2>
              <p className="text-muted-foreground text-sm">
                Recent backup executions and artifacts.
              </p>
            </div>
            <Button variant="ghost" size="icon" onClick={() => setIsOpen(false)}>
              <span className="sr-only">Close</span>
              <svg
                width="15"
                height="15"
                viewBox="0 0 15 15"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                className="h-4 w-4"
              >
                <path
                  d="M11.7816 4.03157C12.0062 3.80702 12.0062 3.44295 11.7816 3.2184C11.5571 2.99385 11.193 2.99385 10.9685 3.2184L7.50005 6.68682L4.03164 3.2184C3.80708 2.99385 3.44301 2.99385 3.21846 3.2184C2.99391 3.44295 2.99391 3.80702 3.21846 4.03157L6.68688 7.49999L3.21846 10.9684C2.99391 11.193 2.99391 11.557 3.21846 11.7816C3.44301 12.0061 3.80708 12.0061 4.03164 11.7816L7.50005 8.31316L10.9685 11.7816C11.193 12.0061 11.5571 12.0061 11.7816 11.7816C12.0062 11.557 12.0062 11.193 11.7816 10.9684L8.31322 7.49999L11.7816 4.03157Z"
                  fill="currentColor"
                  fillRule="evenodd"
                  clipRule="evenodd"
                ></path>
              </svg>
            </Button>
          </div>

          <div className="max-h-[60vh] overflow-y-auto pr-2">
            {isLoading ? (
              <div className="py-8 text-center text-muted-foreground">Loading records...</div>
            ) : records.length === 0 ? (
              <div className="py-8 text-center text-muted-foreground">No backup records yet.</div>
            ) : (
              <div className="space-y-3">
                {records.map((record) => (
                  <div
                    key={record.id}
                    className="flex items-center justify-between rounded-lg border border-border/50 bg-background/50 p-4"
                  >
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-sm">
                          {record.startedAt
                            ? format(new Date(record.startedAt), 'MMM d, yyyy HH:mm:ss')
                            : 'Unknown time'}
                        </span>
                        <Badge
                          variant={
                            record.status === 'completed'
                              ? 'default'
                              : record.status === 'failed'
                                ? 'destructive'
                                : 'secondary'
                          }
                          className="text-[10px]"
                        >
                          {record.status}
                        </Badge>
                      </div>
                      <p className="text-muted-foreground text-xs">
                        Size: {(record.fileSizeBytes / 1024 / 1024).toFixed(2)} MB
                      </p>
                    </div>
                    {record.status === 'completed' && record.s3Url && (
                      <Button variant="outline" size="sm" asChild>
                        <a href={record.s3Url} target="_blank" rel="noreferrer">
                          <Download className="mr-2 h-4 w-4" />
                          Download
                        </a>
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
