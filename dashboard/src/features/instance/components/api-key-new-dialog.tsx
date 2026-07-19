import { Check, Copy, Key, X } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';

interface ApiKeyNewDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  newKeyPlain: string;
}

export function ApiKeyNewDialog({ open, onOpenChange, newKeyPlain }: ApiKeyNewDialogProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(newKeyPlain);
    setCopied(true);
    toast.success('API key copied to clipboard');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-2xl [&>button]:hidden">
        <div className="px-5 pt-5 pb-4">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-xl tracking-tight">
                <Key className="h-5 w-5 text-primary" />
                New API key
              </DialogTitle>
              <DialogDescription className="mt-1.5 flex items-center gap-1.5 font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                <Key className="h-3 w-3" />
                Shown only once
              </DialogDescription>
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

        <div className="h-px w-full bg-border/50" />

        <div className="p-6">
          <div className="rounded-xl border border-primary/20 bg-primary/5 p-6">
            <div className="mb-4 flex items-start justify-between">
              <div className="flex items-center gap-2">
                <Key className="h-4 w-4 text-primary" />
                <span className="font-bold text-[10px] text-primary uppercase tracking-[0.15em]">
                  NEW API KEY
                </span>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
                >
                  <X className="h-4 w-4" />
                </Button>
              </DialogClose>
            </div>

            <p className="mb-4 text-foreground/90 text-sm">
              This full key is shown only once. Store it somewhere secure before closing this
              dialog.
            </p>

            <div className="flex items-center justify-between rounded-xl border border-border/50 bg-background/50 p-1 pl-4">
              <span className="break-all py-2 pr-4 font-mono text-foreground/90 text-sm">
                {newKeyPlain}
              </span>
              <Button
                variant="ghost"
                size="icon"
                onClick={handleCopy}
                className="h-10 w-10 shrink-0 text-muted-foreground hover:bg-background/80 hover:text-foreground"
              >
                {copied ? <Check className="h-4 w-4 text-primary" /> : <Copy className="h-4 w-4" />}
              </Button>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
