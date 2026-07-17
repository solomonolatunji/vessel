import { Database, Info } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Checkbox } from '#/components/ui/checkbox';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  useCreateS3Destination,
  useDeleteS3Destination,
  useListS3Destinations,
} from '#/hooks/useBackups';

export function S3DestinationsList() {
  const [accountId, setAccountId] = useState('');
  const [bucket, setBucket] = useState('');
  const [accessKeyId, setAccessKeyId] = useState('');
  const [secretAccessKey, setSecretAccessKey] = useState('');
  const [createOrVerify, setCreateOrVerify] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  const { data: s3Destinations, isLoading } = useListS3Destinations('global');
  const createS3Dest = useCreateS3Destination();
  const deleteS3Dest = useDeleteS3Destination();

  // If there's an existing configuration, populate it
  const existingDest = s3Destinations?.data?.[0];

  useEffect(() => {
    if (existingDest) {
      // For R2, the account ID is usually the subdomain of the endpoint.
      // But we can just store the endpoint directly or reconstruct it.
      // E.g., https://<account_id>.r2.cloudflarestorage.com
      const match = existingDest.endpoint.match(/https?:\/\/([^.]+)\.r2/);
      if (match) {
        setAccountId(match[1]);
      } else {
        setAccountId(existingDest.endpoint); // fallback
      }
      setBucket(existingDest.bucket);
      setAccessKeyId(existingDest.accessKeyId);
      // Secret access key is usually not returned or shouldn't be populated for security
    }
  }, [existingDest]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      const endpoint = accountId.includes('http')
        ? accountId
        : `https://${accountId}.r2.cloudflarestorage.com`;

      await createS3Dest.mutateAsync({
        payload: {
          projectId: 'global', // assuming global for instance-level settings
          name: 'default',
          endpoint,
          bucket,
          region: 'auto',
          accessKeyId,
          secretAccessKey,
        },
      });
      toast.success('R2 connection saved successfully');
    } catch (_error) {
      toast.error('Failed to save R2 connection');
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!existingDest) return;
    try {
      await deleteS3Dest.mutateAsync({ id: existingDest.id });
      toast.success('R2 connection deleted successfully');
      setAccountId('');
      setBucket('');
      setAccessKeyId('');
      setSecretAccessKey('');
    } catch (_error) {
      toast.error('Failed to delete R2 connection');
    }
  };

  if (isLoading) {
    return <div className="p-6">Loading configuration...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-6 pb-2 md:flex-row md:items-start">
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
              STORAGE & BACKUPS
            </p>
            <h1 className="font-bold text-3xl tracking-tight">Connect R2</h1>
          </div>
          <p className="text-muted-foreground text-sm leading-relaxed">
            Manage your R2 connection. Store R2 credentials in Vessl for database backups.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6">
        <div className="rounded-xl border border-blue-500/20 bg-blue-500/10 p-4">
          <div className="flex items-start gap-3">
            <div className="mt-0.5 rounded-full bg-blue-500/20 p-1">
              <Info className="h-4 w-4 text-blue-500" />
            </div>
            <div>
              <h3 className="font-medium text-blue-500 text-sm">Secure your credentials</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                These credentials will be used to automatically push database backups to your R2
                bucket. Ensure the API token has sufficient permissions to write to the bucket.
              </p>
            </div>
          </div>
        </div>

        <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
              <Database className="h-5 w-5" />
            </div>
            <div>
              <h2 className="font-semibold text-lg">R2 credentials</h2>
              <p className="text-muted-foreground text-sm">Enter your connection details</p>
            </div>
          </div>

          <form onSubmit={handleSave} className="space-y-4">
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label
                  htmlFor="accountId"
                  className="font-bold text-muted-foreground text-xs uppercase tracking-wider"
                >
                  ACCOUNT ID
                </Label>
                <Input
                  id="accountId"
                  type="text"
                  value={accountId}
                  onChange={(e) => setAccountId(e.target.value)}
                  placeholder="00d8..."
                  className="h-11 bg-background/50 font-mono"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label
                  htmlFor="bucket"
                  className="font-bold text-muted-foreground text-xs uppercase tracking-wider"
                >
                  BUCKET NAME
                </Label>
                <Input
                  id="bucket"
                  type="text"
                  value={bucket}
                  onChange={(e) => setBucket(e.target.value)}
                  placeholder="vessl-backups"
                  className="h-11 bg-background/50 font-mono"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label
                  htmlFor="accessKeyId"
                  className="font-bold text-muted-foreground text-xs uppercase tracking-wider"
                >
                  ACCESS KEY ID
                </Label>
                <Input
                  id="accessKeyId"
                  type="text"
                  value={accessKeyId}
                  onChange={(e) => setAccessKeyId(e.target.value)}
                  placeholder="Enter access key"
                  className="h-11 bg-background/50 font-mono"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label
                  htmlFor="secretAccessKey"
                  className="font-bold text-muted-foreground text-xs uppercase tracking-wider"
                >
                  SECRET ACCESS KEY
                </Label>
                <Input
                  id="secretAccessKey"
                  type="password"
                  value={secretAccessKey}
                  onChange={(e) => setSecretAccessKey(e.target.value)}
                  placeholder="Enter secret key"
                  className="h-11 bg-background/50 font-mono"
                  required
                />
              </div>
            </div>

            <div className="flex items-center space-x-2 pt-2 pb-2">
              <Checkbox
                id="verify"
                checked={createOrVerify}
                onCheckedChange={(checked) => setCreateOrVerify(checked as boolean)}
              />
              <Label htmlFor="verify" className="font-medium text-sm leading-none">
                Verify credentials and create bucket if missing
              </Label>
            </div>

            <div className="flex justify-end gap-3 pt-2">
              {existingDest && (
                <Button
                  type="button"
                  variant="destructive"
                  onClick={handleDelete}
                  disabled={deleteS3Dest.isPending}
                  className="h-11 font-bold text-xs uppercase tracking-wider"
                >
                  DELETE
                </Button>
              )}
              <Button
                type="submit"
                disabled={isSaving || createS3Dest.isPending}
                className="h-11 bg-primary font-bold text-primary-foreground text-xs uppercase tracking-wider"
              >
                {isSaving ? 'SAVING...' : 'SAVE R2 CONFIGURATION'}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
