import { Activity, AlertTriangle, HardDrive, RefreshCw, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { useCleanupSystem, useGetSystemStats, useRestartSystem } from '#/hooks/useSystem';
import { MaintenanceDialogs } from './components/maintenance-dialogs';

const ProgressBar = ({ value, colorClass }: { value: number; colorClass: string }) => (
  <div className="flex h-1.5 w-full overflow-hidden rounded-full bg-background">
    <div className={`h-full ${colorClass}`} style={{ width: `${value}%` }} />
  </div>
);

export const MaintenancePage = () => {
  const { data: statsData, refetch } = useGetSystemStats();
  const { mutateAsync: cleanup, isPending: cleaning } = useCleanupSystem();
  const { mutateAsync: restart, isPending: restarting } = useRestartSystem();
  const [confirmCleanup, setConfirmCleanup] = useState(false);
  const [confirmRestart, setConfirmRestart] = useState(false);

  const stats = statsData?.data;

  const handleCleanup = async () => {
    try {
      await cleanup();
      toast.success('Docker cleanup completed');
    } catch {
      toast.error('Cleanup failed');
    } finally {
      setConfirmCleanup(false);
    }
  };

  const handleRestart = async () => {
    try {
      await restart();
      toast.success('Restart initiated');
    } catch {
      toast.error('Restart failed');
    } finally {
      setConfirmRestart(false);
    }
  };

  const usedPercent = stats?.disk.percent ? Number(stats.disk.percent.toFixed(1)) : 0;
  const freeGb = stats?.disk.freeGb ? stats.disk.freeGb.toFixed(1) : '0';
  const usedGb = stats?.disk.usedGb ? stats.disk.usedGb.toFixed(1) : '0';
  const totalGb = stats?.disk.totalGb ? stats.disk.totalGb.toFixed(1) : '0';

  const dockerPercent =
    stats?.disk?.totalGb && stats?.docker?.reclaimableGb
      ? Number(((stats.docker.reclaimableGb / stats.disk.totalGb) * 100).toFixed(1))
      : 0;
  const reclaimableGb = stats?.docker?.reclaimableGb ? stats.docker.reclaimableGb.toFixed(2) : '0';

  const buildCacheReclaimableStr = stats?.docker?.buildCache?.reclaimable || '0';
  const buildCacheGb = parseFloat(buildCacheReclaimableStr)
    ? parseFloat(buildCacheReclaimableStr).toFixed(2)
    : '0';

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Activity className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Host health and cleanup</h1>
            <p className="text-muted-foreground text-sm">
              Manage unused resources, dangling images, and system caches to reclaim space.
            </p>
          </div>
        </div>

        <div className="flex shrink-0 flex-col items-end gap-4">
          {Number(reclaimableGb) > 3 ? (
            <Badge
              variant="outline"
              className="border-destructive/50 bg-destructive/10 px-3 py-1 font-bold text-[10px] text-destructive uppercase tracking-widest"
            >
              ATTENTION
            </Badge>
          ) : (
            <Badge
              variant="outline"
              className="border-primary/50 bg-primary/10 px-3 py-1 font-bold text-[10px] text-primary uppercase tracking-widest"
            >
              0 ISSUES
            </Badge>
          )}
          <Button
            variant="outline"
            onClick={() => refetch()}
            className="flex h-11 items-center gap-2 rounded-xl border-border/50 bg-background/50 px-6 font-semibold text-foreground text-xs uppercase tracking-widest hover:bg-background"
          >
            <RefreshCw className="h-4 w-4" /> REFRESH
          </Button>
        </div>
      </div>

      {stats?.docker?.reclaimableGb && stats.docker.reclaimableGb > 3 ? (
        <div className="flex w-full items-center gap-3 rounded-lg border border-destructive/30 bg-destructive/10 p-4 font-medium text-destructive text-sm">
          <AlertTriangle className="h-4 w-4" /> Docker has more than 3 GB reclaimable.
        </div>
      ) : null}

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div className="flex flex-col justify-between space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
          <div className="flex items-center justify-between">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              ROOT DISK FREE
            </p>
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              {usedPercent}% USED
            </p>
          </div>
          <h2 className="font-bold text-3xl">{freeGb} GB</h2>
          <div className="space-y-2">
            <ProgressBar value={usedPercent} colorClass="bg-primary" />
            <p className="text-muted-foreground text-xs">
              {usedGb} GB used of {totalGb} GB on /.
            </p>
          </div>
        </div>

        <div className="flex flex-col justify-between space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
          <div className="flex items-center justify-between">
            <p className="w-32 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              DOCKER CLEANUP CANDIDATES
            </p>
            <p className="w-20 text-right font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              {dockerPercent}% OF DISK
            </p>
          </div>
          <h2 className="font-bold text-3xl">{reclaimableGb} GB</h2>
          <div className="space-y-2">
            <ProgressBar value={dockerPercent} colorClass="bg-yellow-500" />
            <p className="text-muted-foreground text-xs">
              {buildCacheGb} GB is build cache. Safe cleanup can usually clear that.
            </p>
          </div>
        </div>

        <div className="flex flex-col justify-between space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
          <div className="flex items-center justify-between">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              BUILD ARTIFACTS
            </p>
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              0%
            </p>
          </div>
          <h2 className="font-bold text-3xl">0 GB</h2>
          <div className="space-y-2">
            <ProgressBar value={0} colorClass="bg-primary" />
            <p className="text-muted-foreground text-xs">No build artifact directory yet.</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div className="flex h-32 flex-col rounded-2xl border border-border/50 bg-card/40 p-6">
          <p className="mb-4 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            DISK TREND
          </p>
          <div className="flex flex-1 items-end gap-1">
            {Array.from({ length: 15 }).map((_, i) => (
              <div
                key={i}
                className="flex-1 rounded-sm bg-primary/40"
                style={{ height: `${40 + i * 2}%` }}
              />
            ))}
          </div>
        </div>
        <div className="flex h-32 flex-col rounded-2xl border border-border/50 bg-card/40 p-6">
          <p className="mb-4 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            DOCKER RECLAIMABLE TREND
          </p>
          <div className="flex flex-1 items-end gap-1">
            {Array.from({ length: 15 }).map((_, i) => (
              <div
                key={i}
                className="flex-1 rounded-sm bg-yellow-500/40"
                style={{ height: `${60 - i}%` }}
              />
            ))}
          </div>
        </div>
        <div className="flex h-32 flex-col rounded-2xl border border-border/50 bg-card/40 p-6">
          <p className="mb-4 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            BUILD ARTIFACT TREND
          </p>
          <div className="flex flex-1 items-end gap-1">
            {Array.from({ length: 15 }).map((_, i) => (
              <div key={i} className="h-1 flex-1 rounded-sm bg-primary/10" />
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div className="overflow-hidden rounded-2xl border border-border/50 bg-card/40 lg:col-span-2">
          <div className="flex items-center justify-between border-border/50 border-b p-6">
            <h3 className="font-bold text-xl">Docker storage</h3>
            <Badge
              variant="outline"
              className="border-primary/50 bg-primary/10 px-3 py-1 font-bold text-[10px] text-primary uppercase tracking-widest"
            >
              AVAILABLE
            </Badge>
          </div>

          <div className="divide-y divide-border/50">
            <div className="grid grid-cols-3 items-center p-6">
              <div>
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  IMAGES
                </p>
                <p className="mt-1 text-muted-foreground text-xs">
                  {stats?.docker?.images?.active || 0}/{stats?.docker?.images?.totalCount || 0}{' '}
                  active
                </p>
              </div>
              <div className="col-span-2 flex items-center gap-6">
                <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
                  <div className="h-full w-1/3 bg-muted-foreground" />
                </div>
                <div className="space-x-2 font-mono text-sm">
                  <span className="text-foreground">{stats?.docker?.images?.size || '0 B'}</span>
                  <span className="text-yellow-500">
                    {stats?.docker?.images?.reclaimable || '0 B'} candidate
                  </span>
                </div>
              </div>
            </div>
            <div className="grid grid-cols-3 items-center p-6">
              <div>
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  CONTAINERS
                </p>
                <p className="mt-1 text-muted-foreground text-xs">
                  {stats?.docker?.containers?.active || 0}/
                  {stats?.docker?.containers?.totalCount || 0} active
                </p>
              </div>
              <div className="col-span-2 flex items-center gap-6">
                <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
                  <div className="h-full w-1/12 bg-muted-foreground" />
                </div>
                <div className="space-x-2 font-mono text-sm">
                  <span className="text-foreground">
                    {stats?.docker?.containers?.size || '0 B'}
                  </span>
                  <span className="text-muted-foreground/50">
                    {stats?.docker?.containers?.reclaimable || '0 B'} candidate
                  </span>
                </div>
              </div>
            </div>
            <div className="grid grid-cols-3 items-center p-6">
              <div>
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  LOCAL VOLUMES
                </p>
                <p className="mt-1 text-muted-foreground text-xs">
                  {stats?.docker?.volumes?.active || 0}/{stats?.docker?.volumes?.totalCount || 0}{' '}
                  active
                </p>
              </div>
              <div className="col-span-2 flex items-center gap-6">
                <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
                  <div className="h-full w-1/6 bg-muted-foreground" />
                </div>
                <div className="space-x-2 font-mono text-sm">
                  <span className="text-foreground">{stats?.docker?.volumes?.size || '0 B'}</span>
                  <span className="text-muted-foreground/50">
                    {stats?.docker?.volumes?.reclaimable || '0 B'} candidate
                  </span>
                </div>
              </div>
            </div>
            <div className="grid grid-cols-3 items-center p-6">
              <div>
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  BUILD CACHE
                </p>
                <p className="mt-1 text-muted-foreground text-xs">
                  {stats?.docker?.buildCache?.active || 0}/
                  {stats?.docker?.buildCache?.totalCount || 0} active
                </p>
              </div>
              <div className="col-span-2 flex items-center gap-6">
                <div className="flex h-2 w-48 overflow-hidden rounded-full bg-background">
                  <div className="h-full w-[80%] bg-muted-foreground" />
                </div>
                <div className="space-x-2 font-mono text-sm">
                  <span className="text-foreground">
                    {stats?.docker?.buildCache?.size || '0 B'}
                  </span>
                  <span className="text-yellow-500">
                    {stats?.docker?.buildCache?.reclaimable || '0 B'} candidate
                  </span>
                </div>
              </div>
            </div>
            <div className="bg-background/30 p-6 text-muted-foreground text-xs">
              Docker can keep image layers listed as candidates after safe cleanup when running
              services still reference them.
            </div>
          </div>
        </div>

        <div className="flex flex-col space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6 lg:col-span-1">
          <div className="flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
              <HardDrive className="h-5 w-5" />
            </div>
            <div>
              <h3 className="font-bold text-xl">Cleanup</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                Current disk and Docker cleanup targets.
              </p>
            </div>
          </div>

          <div className="divide-y divide-border/50 rounded-xl border border-border/50 text-sm">
            <div className="flex justify-between bg-background/30 p-4">
              <span className="text-muted-foreground">Top Docker candidate</span>
              <span className="font-mono">
                {stats?.docker?.buildCache?.reclaimable || '0 B'} Build Cache
              </span>
            </div>
            <div className="flex justify-between p-4">
              <span className="text-muted-foreground">Vessl data</span>
              <span className="font-mono">1.81 MB</span>
            </div>
            <div className="flex justify-between p-4">
              <span className="text-muted-foreground">Backups</span>
              <span className="font-mono">{stats?.backups?.size || '0 B'}</span>
            </div>
            <div className="flex justify-between p-4">
              <span className="text-muted-foreground">Package cache</span>
              <span className="font-mono">{stats?.packageCache?.size || '0 B'}</span>
            </div>
            <div className="flex justify-between p-4">
              <span className="text-muted-foreground">System logs</span>
              <span className="font-mono">{stats?.systemLogs?.size || '0 B'}</span>
            </div>
          </div>

          <div className="mt-auto flex flex-col gap-3 pt-4">
            <Button
              variant="outline"
              onClick={() => setConfirmCleanup(true)}
              disabled={cleaning}
              className="h-12 border-primary/30 bg-primary/5 font-semibold text-primary text-xs uppercase tracking-widest hover:bg-primary/10 hover:text-primary"
            >
              <RefreshCw className="mr-2 h-4 w-4" /> {cleaning ? 'CLEANING...' : 'SAFE CLEANUP'}
            </Button>
            <Button
              variant="outline"
              onClick={() => setConfirmRestart(true)}
              disabled={restarting}
              className="h-12 border-destructive/30 bg-destructive/5 font-semibold text-destructive text-xs uppercase tracking-widest hover:bg-destructive/10 hover:text-destructive"
            >
              <Trash2 className="mr-2 h-4 w-4" /> RESTART DAEMON
            </Button>
          </div>
        </div>
      </div>

      <MaintenanceDialogs
        confirmCleanup={confirmCleanup}
        setConfirmCleanup={setConfirmCleanup}
        cleaning={cleaning}
        handleCleanup={handleCleanup}
        confirmRestart={confirmRestart}
        setConfirmRestart={setConfirmRestart}
        restarting={restarting}
        handleRestart={handleRestart}
      />
    </div>
  );
};
