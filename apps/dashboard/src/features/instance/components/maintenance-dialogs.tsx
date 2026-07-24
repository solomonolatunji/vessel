import { Check, RefreshCw, Trash2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';

interface MaintenanceDialogsProps {
  confirmCleanup: boolean;
  setConfirmCleanup: (val: boolean) => void;
  cleaning: boolean;
  handleCleanup: () => void;
  confirmRestart: boolean;
  setConfirmRestart: (val: boolean) => void;
  restarting: boolean;
  handleRestart: () => void;
}

export function MaintenanceDialogs({
  confirmCleanup,
  setConfirmCleanup,
  cleaning,
  handleCleanup,
  confirmRestart,
  setConfirmRestart,
  restarting,
  handleRestart,
}: MaintenanceDialogsProps) {
  return (
    <>
      <Dialog open={confirmCleanup} onOpenChange={setConfirmCleanup}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-100 [&>button]:hidden">
          <div className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <Trash2 className="h-5 w-5 text-primary" />
                  Run Docker Cleanup
                </DialogTitle>
                <DialogDescription>Removes unused images and volumes</DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setConfirmCleanup(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={handleCleanup}
              disabled={cleaning}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <Check className="h-3.5 w-3.5" />
              {cleaning ? 'Running...' : 'Run Cleanup'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={confirmRestart} onOpenChange={setConfirmRestart}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-100 [&>button]:hidden">
          <div className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                  <RefreshCw className="h-5 w-5" />
                  Restart Daemon
                </DialogTitle>
                <DialogDescription>All services will be briefly unavailable</DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setConfirmRestart(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={handleRestart}
              disabled={restarting}
              variant="destructive"
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <RefreshCw className="h-3.5 w-3.5" />
              {restarting ? 'Restarting...' : 'Restart'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
