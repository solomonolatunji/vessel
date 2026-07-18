import { CheckCircle, Download, RefreshCw } from 'lucide-react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { useCheckUpdate, useDeployUpdate, useGetUpdateStatus } from '#/hooks/useSettings';

export const UpdatesPage = () => {
  const { data, isLoading, refetch } = useGetUpdateStatus();
  const { mutateAsync: checkUpdate, isPending: checking } = useCheckUpdate();
  const { mutateAsync: deployUpdate, isPending: deploying } = useDeployUpdate();

  const info = data?.data;
  const hasUpdate = info?.hasUpdate ?? false;

  const handleCheck = async () => {
    try {
      await checkUpdate();
      await refetch();
      toast.success('Update check complete');
    } catch {
      toast.error('Failed to check for updates');
    }
  };

  const handleDeploy = async () => {
    try {
      await deployUpdate();
      toast.success('Update started — daemon will restart shortly');
    } catch {
      toast.error('Failed to deploy update');
    }
  };

  return (
    <div className="space-y-6">
      {/* Top Section */}
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <RefreshCw className="h-4.5 w-4.5" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Vessl release channel</h1>
            <p className="text-muted-foreground text-sm">
              Compare this install with the official Vessl releases, review changelogs, and update
              the daemon when you are ready.
            </p>
          </div>
        </div>

        <div className="flex shrink-0 flex-col items-end gap-4">
          {!isLoading && (
            <Badge
              variant="outline"
              className={`rounded-md border border-primary/50 bg-primary/10 px-3 py-1 font-bold text-[10px] text-primary uppercase tracking-[0.15em]`}
            >
              {hasUpdate ? 'UPDATE AVAILABLE' : 'UP TO DATE'}
            </Badge>
          )}
          <Button
            variant="outline"
            onClick={handleCheck}
            disabled={checking || deploying}
            className="flex h-11 items-center gap-2 rounded-xl border-border/50 bg-background/50 px-6 font-semibold text-foreground text-xs uppercase tracking-widest hover:bg-background"
          >
            <RefreshCw className={`h-4 w-4 ${checking ? 'animate-spin' : ''}`} />
            {checking ? 'CHECKING...' : 'CHECK UPDATES'}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div className="flex flex-col justify-center space-y-2 rounded-2xl border border-border/50 bg-card/40 p-6">
          <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            REPOSITORY
          </p>
          <p className="font-mono text-sm">vesslhq/vessl</p>
        </div>
        <div className="flex flex-col justify-center space-y-2 rounded-2xl border border-border/50 bg-card/40 p-6">
          <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            INSTALLED
          </p>
          <p className="font-mono text-sm">
            {isLoading ? <Skeleton className="h-5 w-20" /> : info?.currentVersion || 'unknown'}
          </p>
        </div>
        <div className="flex flex-col justify-center space-y-2 rounded-2xl border border-border/50 bg-card/40 p-6">
          <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            GITHUB LATEST
          </p>
          <p className="font-mono text-sm">
            {isLoading ? <Skeleton className="h-5 w-20" /> : info?.latestVersion || 'unknown'}
          </p>
        </div>
      </div>

      {/* Main Status Banner */}
      {!isLoading && (
        <div className="flex flex-col items-center justify-center space-y-6 rounded-2xl border border-border/50 bg-card/40 px-6 py-12 text-center">
          <div
            className={`flex h-16 w-16 items-center justify-center rounded-2xl border border-primary/20 bg-primary/10 text-primary`}
          >
            {hasUpdate ? <Download className="h-8 w-8" /> : <CheckCircle className="h-8 w-8" />}
          </div>
          <div className="space-y-2">
            <h2 className="font-bold text-2xl tracking-tight">
              {hasUpdate ? 'Update available' : 'Vessl is up to date'}
            </h2>
            <p className="text-muted-foreground text-sm">
              {hasUpdate
                ? `Version ${info?.latestVersion} is ready to be installed.`
                : 'Installed version matches latest release.'}
            </p>
          </div>

          {hasUpdate && (
            <Button
              onClick={handleDeploy}
              disabled={deploying}
              className="mt-4 flex h-12 items-center gap-2 rounded-xl border-primary/20 bg-primary/10 px-8 font-bold text-primary text-xs uppercase tracking-widest transition-all hover:bg-primary/20 hover:text-primary"
            >
              <Download className="h-4 w-4" />
              {deploying ? 'UPDATING...' : 'INSTALL UPDATE'}
            </Button>
          )}
        </div>
      )}

      {info?.releaseNotes && hasUpdate && (
        <div className="space-y-4 rounded-2xl border border-border/50 bg-card/40 p-8">
          <h3 className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            Release Notes
          </h3>
          <pre className="overflow-x-auto whitespace-pre-wrap rounded-xl border border-border/50 bg-background/50 p-6 font-mono text-muted-foreground text-xs leading-relaxed">
            {info.releaseNotes}
          </pre>
        </div>
      )}
    </div>
  );
};
