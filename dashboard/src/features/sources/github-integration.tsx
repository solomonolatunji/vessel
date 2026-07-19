import { Check, Edit, ExternalLink, Plus, Trash } from 'lucide-react';
import { useEffect, useState } from 'react';
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
import { Textarea } from '#/components/ui/textarea';
import {
  useDeleteGitApp,
  useExchangeGithubManifest,
  useGetGitApps,
  useSaveGitApp,
} from '#/hooks/useSettings';
import type { GithubApp } from '#/interfaces/settings';

const GithubIcon = ({ className }: { className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
    <path d="M9 18c-4.51 2-5-2-7-2" />
  </svg>
);

export function GithubIntegration() {
  const { data, isLoading } = useGetGitApps();
  const saveMutation = useSaveGitApp();
  const deleteMutation = useDeleteGitApp();
  const exchangeMutation = useExchangeGithubManifest();

  const apps = (data?.data as GithubApp[]) || [];

  const [isEditing, setIsEditing] = useState(false);
  const [editingApp, setEditingApp] = useState<GithubApp | null>(null);
  const [deletingApp, setDeletingApp] = useState<string | null>(null);

  // Form state
  const [accessToken, setAccessToken] = useState('');
  const [webhookSecret, setWebhookSecret] = useState('');
  const [appId, setAppId] = useState('');
  const [clientId, setClientId] = useState('');
  const [appSlug, setAppSlug] = useState('');
  const [privateKey, setPrivateKey] = useState('');

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    if (code && !exchangeMutation.isPending && !exchangeMutation.isSuccess) {
      window.history.replaceState({}, document.title, window.location.pathname);
      exchangeMutation.mutate(
        { code },
        {
          onSuccess: () => {
            toast.success('GitHub App connected successfully!');
            setIsEditing(false);
            setEditingApp(null);
          },
          onError: (err) => {
            toast.error(err.message || 'Failed to connect GitHub App');
          },
        }
      );
    }
  }, [exchangeMutation.isPending, exchangeMutation.isSuccess, exchangeMutation.mutate]);

  useEffect(() => {
    if (editingApp) {
      setAppId(editingApp.appId || '');
      setClientId(editingApp.clientId || '');
      setWebhookSecret(editingApp.webhookSecret ? '********' : '');
      setAccessToken(editingApp.clientSecret ? '********' : '');
      setAppSlug(editingApp.name || '');
      setPrivateKey(editingApp.privateKey ? '********' : '');
    } else {
      setAppId('');
      setClientId('');
      setWebhookSecret('');
      setAccessToken('');
      setAppSlug('');
      setPrivateKey('');
    }
  }, [editingApp]);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();

    const payload = {
      ...(editingApp?.id ? { id: editingApp.id } : {}),
      appId,
      clientId,
      name: appSlug,
      ...(accessToken !== '********' ? { clientSecret: accessToken } : {}),
      ...(webhookSecret !== '********' ? { webhookSecret } : {}),
      ...(privateKey !== '********' ? { privateKey } : {}),
    };

    saveMutation.mutate(payload, {
      onSuccess: () => {
        setIsEditing(false);
        setEditingApp(null);
        toast.success('GitHub settings saved successfully');
      },
      onError: (err: Error) => {
        toast.error(err.message || 'Failed to save GitHub settings');
      },
    });
  };

  const confirmDelete = () => {
    if (!deletingApp) return;
    deleteMutation.mutate(deletingApp, {
      onSuccess: () => {
        if (editingApp?.id === deletingApp) {
          setIsEditing(false);
          setEditingApp(null);
        }
        toast.success('GitHub connection removed');
        setDeletingApp(null);
      },
      onError: (err: Error) => {
        toast.error(err.message || 'Failed to remove GitHub connection');
      },
    });
  };

  const manifestStr =
    typeof window !== 'undefined'
      ? JSON.stringify({
          name: `vessl-${Math.random().toString(36).substring(7)}`,
          url: window.location.origin,
          hook_attributes: {
            url: `${window.location.origin}/api/webhooks/github/services/generic`,
          },
          redirect_url: `${window.location.origin}/dashboard/sources`,
          public: false,
          default_permissions: {
            contents: 'read',
            metadata: 'read',
            pull_requests: 'read',
            emails: 'read',
          },
          default_events: ['push', 'pull_request'],
        })
      : '{}';

  if (isLoading) {
    return <div className="h-64 animate-pulse rounded-xl bg-card" />;
  }

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <GithubIcon className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Connected GitHub Apps</h1>
            <p className="text-muted-foreground text-sm">
              Connect GitHub Apps to automatically deploy pushed commits.
            </p>
          </div>
        </div>
        <Button
          className="gap-2"
          onClick={() => {
            setEditingApp(null);
            setIsEditing(true);
          }}
        >
          <Plus className="h-4 w-4" />
          ADD GITHUB APP
        </Button>
      </div>

      {apps.length > 0 ? (
        <div className="grid grid-cols-1 gap-4">
          {apps.map((app) => (
            <div key={app.id} className="rounded-xl border border-border bg-card p-6">
              <div className="flex items-start justify-between">
                <div className="flex items-start gap-4">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
                    <GithubIcon className="h-6 w-6" />
                  </div>
                  <div className="space-y-1">
                    <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                      GITHUB INTEGRATION
                    </p>
                    <div className="flex items-center gap-3">
                      <h2 className="font-bold text-xl tracking-tight">
                        {app.name || 'GitHub App'}
                      </h2>
                      <div className="rounded border border-primary/30 bg-primary/10 px-2 py-0.5 font-semibold text-[10px] text-primary uppercase tracking-widest">
                        CONNECTED
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-10 w-10 border-border/50 bg-transparent hover:bg-card"
                    onClick={() => {
                      setEditingApp(app);
                      setIsEditing(true);
                    }}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-10 w-10 border-border/50 bg-transparent hover:bg-destructive/10 hover:text-destructive"
                    onClick={() => setDeletingApp(app.id)}
                    disabled={deleteMutation.isPending}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <div className="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2">
                <div className="rounded-lg border border-border/50 bg-background/50 p-4">
                  <p className="font-medium text-[10px] text-muted-foreground uppercase tracking-widest">
                    APP SLUG
                  </p>
                  <p className="mt-2 font-mono text-sm">{app.name || 'Not set'}</p>
                </div>
                <div className="rounded-lg border border-border/50 bg-background/50 p-4">
                  <p className="font-medium text-[10px] text-muted-foreground uppercase tracking-widest">
                    APP ID
                  </p>
                  <p className="mt-2 truncate font-mono text-sm">{app.appId || 'Not set'}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
            <GithubIcon className="h-5 w-5 text-primary" />
          </div>
          <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">
            No GitHub Apps connected
          </h3>
          <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
            Connect a GitHub App to deploy repositories and receive webhooks.
          </p>
          <Button
            className="mt-6 gap-2"
            onClick={() => {
              setEditingApp(null);
              setIsEditing(true);
            }}
          >
            <Plus className="h-4 w-4" />
            CONNECT GITHUB APP
          </Button>
        </div>
      )}

      <Dialog open={isEditing} onOpenChange={setIsEditing}>
        <DialogContent className="max-w-200 gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <GithubIcon className="h-5 w-5 text-primary" />
                  {editingApp ? 'Edit GitHub App' : 'Connect GitHub App'}
                </DialogTitle>
                <DialogDescription>Configure Github Integration</DialogDescription>
              </div>
              <div className="flex items-center gap-3">
                {!editingApp && (
                  <Button
                    asChild
                    variant="ghost"
                    className="h-9 gap-2 font-mono font-semibold text-[11px] text-primary uppercase tracking-wider hover:text-primary"
                  >
                    <a href="https://github.com/settings/apps/new" target="_blank" rel="noreferrer">
                      <ExternalLink className="h-3.5 w-3.5" />
                      Create App
                    </a>
                  </Button>
                )}
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
          </div>

          <div className="h-px w-full bg-border/50" />

          <div className="px-5 pt-4 pb-5">
            {!editingApp && (
              <div className="mt-4 rounded-lg border border-primary/20 bg-primary/5 p-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="font-bold text-[10px] text-primary uppercase tracking-widest">
                      ONE-CLICK CONNECT
                    </p>
                    <p className="text-sm">
                      Create the GitHub App and fill in every credential automatically.
                    </p>
                  </div>
                  <form
                    action="https://github.com/settings/apps/new"
                    method="post"
                    target="_blank"
                    rel="noopener"
                  >
                    <input type="hidden" name="manifest" value={manifestStr} />
                    <Button
                      type="submit"
                      className="h-10 gap-2 bg-primary/20 font-bold text-primary text-xs uppercase tracking-wider hover:bg-primary/30"
                    >
                      <GithubIcon className="h-4 w-4" />
                      CONNECT WITH GITHUB
                    </Button>
                  </form>
                </div>
              </div>
            )}

            {!editingApp && (
              <div className="relative mt-12 mb-8 flex w-full justify-center border-border/50 border-t">
                <span className="absolute -top-3 bg-card px-4 font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                  OR ENTER MANUALLY
                </span>
              </div>
            )}

            <form onSubmit={handleSave} className="space-y-6">
              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_CLIENT_SECRET
                  </Label>
                  <Input
                    type="password"
                    placeholder="GitHub personal access token"
                    value={accessToken}
                    onChange={(e) => setAccessToken(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_WEBHOOK_SECRET
                  </Label>
                  <Input
                    type="password"
                    placeholder="Webhook secret"
                    value={webhookSecret}
                    onChange={(e) => setWebhookSecret(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_APP_ID
                  </Label>
                  <Input
                    placeholder="e.g. 4322334"
                    value={appId}
                    onChange={(e) => setAppId(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
                <div className="space-y-2">
                  <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                    GITHUB_CLIENT_ID
                  </Label>
                  <Input
                    placeholder="Iv1.xxxxxxxxxxxx"
                    value={clientId}
                    onChange={(e) => setClientId(e.target.value)}
                    className="h-11 bg-background/50 font-mono"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                  GITHUB_APP_SLUG
                </Label>
                <Input
                  placeholder="my-vessl-app"
                  value={appSlug}
                  onChange={(e) => setAppSlug(e.target.value)}
                  className="h-11 bg-background/50 font-mono"
                />
              </div>

              <div className="space-y-2">
                <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                  GITHUB_PRIVATE_KEY
                </Label>
                <Textarea
                  placeholder="-----BEGIN RSA PRIVATE KEY-----"
                  value={privateKey}
                  onChange={(e) => setPrivateKey(e.target.value)}
                  className="min-h-40 bg-background/50 font-mono"
                />
              </div>

              <div className="mt-8 flex justify-end gap-3 pt-6">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setIsEditing(false)}
                  className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={saveMutation.isPending}
                  className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
                >
                  <Check className="h-3.5 w-3.5" />
                  {saveMutation.isPending ? 'Saving...' : 'Save Settings'}
                </Button>
              </div>
            </form>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={!!deletingApp} onOpenChange={(open) => !open && setDeletingApp(null)}>
        <DialogContent className="max-w-100 gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl [&>button]:hidden">
          <div className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                  <Trash className="h-5 w-5" />
                  Remove GitHub App
                </DialogTitle>
                <DialogDescription>This will break existing deployments</DialogDescription>
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
              onClick={() => setDeletingApp(null)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={(e) => {
                e.preventDefault();
                confirmDelete();
              }}
              disabled={deleteMutation.isPending}
              variant="destructive"
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              <Trash className="h-3.5 w-3.5" />
              {deleteMutation.isPending ? 'Removing...' : 'Remove App'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
