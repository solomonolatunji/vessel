import { format } from 'date-fns';
import { Calendar, Check, Clock, Copy, FolderOpen, Key, Plus, Trash2, X } from 'lucide-react';
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
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreateToken, useDeleteToken, useListTokens } from '#/hooks/useProfile';

export function ApiKeysList() {
  const { data: tokensResponse, isLoading } = useListTokens();
  const createToken = useCreateToken();
  const deleteToken = useDeleteToken();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isNewKeyOpen, setIsNewKeyOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [newKeyPlain, setNewKeyPlain] = useState('');
  const [copied, setCopied] = useState(false);

  const [name, setName] = useState('');
  const [accessLevel, setAccessLevel] = useState<'read' | 'read_write'>('read');
  const [projectScope, setProjectScope] = useState<'all' | 'specific'>('all');
  const [expirationDays, setExpirationDays] = useState<number | null>(30);

  const handleCopy = () => {
    navigator.clipboard.writeText(newKeyPlain);
    setCopied(true);
    toast.success('API key copied to clipboard');
    setTimeout(() => setCopied(false), 2000);
  };

  const handleCreate = () => {
    if (!name.trim()) {
      toast.error('Please enter a name for the API key');
      return;
    }

    let expiresAt: string | undefined;
    if (expirationDays !== null) {
      const date = new Date();
      date.setDate(date.getDate() + expirationDays);
      expiresAt = date.toISOString();
    }

    createToken.mutate(
      {
        payload: {
          name,
          accessLevel,
          projectScope,
          allowedProjects: [],
          expiresAt,
        },
      },
      {
        onSuccess: (data) => {
          setNewKeyPlain(data.plain);
          setIsCreateOpen(false);
          setIsNewKeyOpen(true);
          setName('');
          setAccessLevel('read');
          setProjectScope('all');
          setExpirationDays(30);
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to create API key');
        },
      }
    );
  };

  const handleDelete = () => {
    if (!deleteId) return;
    deleteToken.mutate(
      { id: deleteId },
      {
        onSuccess: () => {
          toast.success('API key deleted');
          setDeleteId(null);
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to delete API key');
        },
      }
    );
  };

  const tokens = tokensResponse?.data || [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Key className="h-4.5 w-4.5" />
          </div>
          <div>
            <h1 className="font-bold text-xl">API keys</h1>
            <p className="text-muted-foreground text-sm">
              Manage API keys for external integrations, CLI access, and programmatic control.
            </p>
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-4">
          <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
            {tokens.length} {tokens.length === 1 ? 'KEY' : 'KEYS'}
          </p>
          <Button
            onClick={() => setIsCreateOpen(true)}
            className="h-11 rounded-xl border-primary/20 bg-primary/10 px-6 font-semibold text-primary text-xs uppercase tracking-widest hover:bg-primary/20 hover:text-primary"
          >
            <Plus className="mr-2 h-4 w-4" /> CREATE API KEY
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6">
        {isLoading ? (
          <div className="rounded-2xl border border-border/50 bg-card/40 p-12 text-center">
            <span className="font-medium text-muted-foreground text-sm uppercase tracking-widest">
              LOADING...
            </span>
          </div>
        ) : tokens.length === 0 ? (
          <div className="flex items-center gap-3 rounded-2xl border border-border/50 bg-card/40 p-6">
            <Key className="h-4 w-4 text-muted-foreground" />
            <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              NO API KEYS
            </span>
          </div>
        ) : (
          tokens.map((token) => (
            <div key={token.id} className="rounded-2xl border border-border/50 bg-card/40 p-6">
              <div className="mb-6 flex items-start justify-between">
                <div>
                  <div className="flex items-center gap-3">
                    <h3 className="font-bold text-foreground text-xl">{token.name}</h3>
                    <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-bold text-[10px] text-primary uppercase tracking-widest">
                      ACTIVE
                    </div>
                  </div>
                  <p className="mt-2 font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
                    {token.prefix}
                  </p>
                </div>
                <Button
                  variant="outline"
                  onClick={() => setDeleteId(token.id)}
                  className="h-10 w-10 border-border/50 bg-transparent p-0 text-muted-foreground transition-colors hover:border-destructive/30 hover:bg-destructive/10 hover:text-destructive"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>

              <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
                <div className="rounded-xl border border-border/50 bg-background/30 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <Key className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      ACCESS
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">
                    {token.accessLevel === 'read_write' ? 'Read and Write' : 'Read'}
                  </div>
                </div>

                <div className="rounded-xl border border-border/50 bg-background/30 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <FolderOpen className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      PROJECTS
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">
                    {token.projectScope === 'specific' ? 'Specific projects' : 'All projects'}
                  </div>
                </div>

                <div className="rounded-xl border border-border/50 bg-background/30 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <Calendar className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      EXPIRATION
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">
                    {token.expiresAt
                      ? `Expires ${format(new Date(token.expiresAt), 'MMM d, hh:mm a')}`
                      : 'No expiration'}
                  </div>
                </div>

                <div className="rounded-xl border border-border/50 bg-background/30 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <Clock className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      LAST USED
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">Never</div>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-2xl [&>button]:hidden">
          <div className="flex flex-col p-8 pb-6">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                  <Key className="h-5 w-5 text-primary" />
                </div>
                <div className="flex flex-col">
                  <DialogTitle className="font-bold text-2xl text-foreground tracking-tight">
                    Create API key
                  </DialogTitle>
                  <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                    BEARER TOKEN ACCESS
                  </DialogDescription>
                </div>
              </div>
              <DialogClose asChild>
                <Button className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground">
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="space-y-8 p-8">
            <div className="space-y-3">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                NAME
              </Label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Production deploys"
                className="h-12 rounded-xl border-border/50 bg-background/50 px-4"
              />
            </div>

            <div className="space-y-3">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                ACCESS
              </Label>
              <div className="grid grid-cols-2 gap-4">
                <button
                  type="button"
                  onClick={() => setAccessLevel('read')}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    accessLevel === 'read'
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <Key className="h-4 w-4" /> READ
                </button>
                <button
                  type="button"
                  onClick={() => setAccessLevel('read_write')}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    accessLevel === 'read_write'
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <Key className="h-4 w-4" /> READ AND WRITE
                </button>
              </div>
            </div>

            <div className="space-y-3">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                PROJECTS
              </Label>
              <div className="grid grid-cols-2 gap-4">
                <button
                  type="button"
                  onClick={() => setProjectScope('all')}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    projectScope === 'all'
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <FolderOpen className="h-4 w-4" /> ALL PROJECTS
                </button>
                <button
                  type="button"
                  onClick={() => setProjectScope('specific')}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    projectScope === 'specific'
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <FolderOpen className="h-4 w-4" /> SPECIFIC PROJECTS
                </button>
              </div>
            </div>

            <div className="space-y-3">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                EXPIRATION
              </Label>
              <div className="grid grid-cols-2 gap-4">
                <button
                  type="button"
                  onClick={() => setExpirationDays(7)}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    expirationDays === 7
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <Clock className="h-4 w-4" /> 7 DAYS
                </button>
                <button
                  type="button"
                  onClick={() => setExpirationDays(30)}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    expirationDays === 30
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <Clock className="h-4 w-4" /> 30 DAYS
                </button>
                <button
                  type="button"
                  onClick={() => setExpirationDays(90)}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    expirationDays === 90
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <Clock className="h-4 w-4" /> 90 DAYS
                </button>
                <button
                  type="button"
                  onClick={() => setExpirationDays(null)}
                  className={`flex h-12 items-center justify-center gap-2 rounded-xl border transition-colors ${
                    expirationDays === null
                      ? 'border-primary/50 bg-primary/10 text-primary'
                      : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
                  }`}
                >
                  <Clock className="h-4 w-4" /> NO EXPIRATION
                </button>
              </div>
            </div>

            <div className="pt-2">
              <Button
                onClick={handleCreate}
                disabled={createToken.isPending}
                className="h-12 rounded-xl border-primary/20 bg-primary/10 px-8 font-semibold text-primary text-xs uppercase tracking-widest hover:bg-primary/20 hover:text-primary"
              >
                <Check className="mr-2 h-4 w-4" />{' '}
                {createToken.isPending ? 'CREATING...' : 'CREATE KEY'}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={isNewKeyOpen} onOpenChange={setIsNewKeyOpen}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-2xl [&>button]:hidden">
          <div className="flex flex-col p-8 pb-6">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                  <Key className="h-5 w-5 text-primary" />
                </div>
                <div className="flex flex-col">
                  <DialogTitle className="font-bold text-2xl text-foreground tracking-tight">
                    New API key
                  </DialogTitle>
                  <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                    SHOWN ONLY ONCE
                  </DialogDescription>
                </div>
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

          <div className="p-8">
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
                  {copied ? (
                    <Check className="h-4 w-4 text-primary" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
          </div>
        </DialogContent>
      </Dialog>
      <Dialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="flex flex-col p-8 pb-6">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-destructive/20 bg-destructive/10">
                  <Trash2 className="h-5 w-5 text-destructive" />
                </div>
                <div className="flex flex-col">
                  <DialogTitle className="font-bold text-2xl text-foreground tracking-tight">
                    Delete API Key
                  </DialogTitle>
                </div>
              </div>
            </div>
            <p className="mt-4 text-muted-foreground text-sm">
              Are you sure you want to delete this API key? Any applications or scripts using it
              will immediately lose access. This action cannot be undone.
            </p>
          </div>
          <div className="flex justify-end gap-3 p-8 pt-6">
            <Button
              variant="ghost"
              onClick={() => setDeleteId(null)}
              className="h-11 px-8 font-bold text-muted-foreground text-xs uppercase tracking-wider hover:bg-muted"
            >
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete} disabled={deleteToken.isPending}>
              <Trash2 className="mr-2 h-4 w-4" />
              {deleteToken.isPending ? 'Deleting...' : 'Delete Key'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
