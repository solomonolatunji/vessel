import { format } from 'date-fns';
import { Calendar, Clock, FolderOpen, Key, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { useListTokens } from '#/hooks/useProfile';
import { ApiKeyCreateDialog } from './components/api-key-create-dialog';
import { ApiKeyDeleteDialog } from './components/api-key-delete-dialog';
import { ApiKeyNewDialog } from './components/api-key-new-dialog';

export function ApiKeysList() {
  const { data: tokensResponse, isLoading } = useListTokens();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isNewKeyOpen, setIsNewKeyOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [newKeyPlain, setNewKeyPlain] = useState('');

  const handleCreateSuccess = (plainKey: string) => {
    setNewKeyPlain(plainKey);
    setIsCreateOpen(false);
    setIsNewKeyOpen(true);
  };

  const tokens = tokensResponse?.data || [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Key className="h-6 w-6" />
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
          <Button onClick={() => setIsCreateOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            CREATE API KEY
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
          <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
              <Key className="h-5 w-5 text-primary" />
            </div>
            <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">No API keys</h3>
            <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
              Create an API key to access Codedock programmatically.
            </p>
            <Button onClick={() => setIsCreateOpen(true)} className="mt-6 gap-2">
              <Plus className="h-4 w-4" />
              CREATE API KEY
            </Button>
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

      <ApiKeyCreateDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        onSuccess={handleCreateSuccess}
      />
      <ApiKeyNewDialog
        open={isNewKeyOpen}
        onOpenChange={setIsNewKeyOpen}
        newKeyPlain={newKeyPlain}
        onClose={() => {
          setNewKeyPlain('');
          setIsNewKeyOpen(false);
        }}
      />
      <ApiKeyDeleteDialog deleteId={deleteId} onClose={() => setDeleteId(null)} />
    </div>
  );
}
