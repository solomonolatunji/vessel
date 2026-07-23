import { ArrowRightLeft, Download, Info, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Row, Section } from '#/components/ui/section';
import { useExportSystem } from '#/hooks/useSystem';

export function MigrationSettings() {
  const [passphrase, setPassphrase] = useState('');
  const [confirmPassphrase, setConfirmPassphrase] = useState('');
  const exportSystem = useExportSystem();

  const handleExport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (passphrase !== confirmPassphrase) {
      toast.error('Passphrases do not match');
      return;
    }
    if (passphrase.length < 8) {
      toast.error('Passphrase must be at least 8 characters long');
      return;
    }
    try {
      const blob = await exportSystem.mutateAsync(passphrase);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `codedock-export-${new Date().toISOString().split('T')[0]}.codedock`;
      a.click();
      window.URL.revokeObjectURL(url);

      toast.success('Instance bundle downloaded successfully');
      setPassphrase('');
      setConfirmPassphrase('');
    } catch {
      toast.error('Failed to export instance bundle');
    }
  };

  const passphrasesMismatch = confirmPassphrase.length > 0 && passphrase !== confirmPassphrase;

  return (
    <form onSubmit={handleExport}>
      <div className="space-y-6 pb-12">
        <div className="mb-5 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
              <ArrowRightLeft className="h-6 w-6" />
            </div>
            <div>
              <h1 className="font-bold text-xl">Migration</h1>
              <p className="text-muted-foreground text-sm">
                Export an encrypted bundle of your instance for migration or backup.
              </p>
            </div>
          </div>
        </div>

        <div className="flex items-start gap-3 rounded-xl border border-blue-500/20 bg-blue-500/10 p-4">
          <div className="mt-0.5 rounded-full bg-blue-500/20 p-1">
            <Info className="h-4 w-4 text-blue-500" />
          </div>
          <div>
            <h3 className="font-medium text-blue-500 text-sm">Keep your passphrase safe</h3>
            <p className="mt-1 text-muted-foreground text-sm">
              The migration bundle is heavily encrypted. Without the exact passphrase, you will not
              be able to decrypt and import this data on a new instance.
            </p>
          </div>
        </div>

        <Section
          icon={<ArrowRightLeft className="h-4 w-4" />}
          title="Export Instance Bundle"
          action={
            <Button
              type="submit"
              size="sm"
              disabled={exportSystem.isPending || !passphrase || passphrasesMismatch}
            >
              {exportSystem.isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Download className="mr-2 h-4 w-4" />
              )}
              {exportSystem.isPending ? 'Exporting...' : 'Download Bundle'}
            </Button>
          }
        >
          <Row
            label="Migration Passphrase"
            description="Used to encrypt the exported bundle. Store it somewhere safe."
          >
            <Input
              id="passphrase"
              type="password"
              value={passphrase}
              onChange={(e) => setPassphrase(e.target.value)}
              placeholder="Enter a secure passphrase"
              className="h-10 w-full font-mono"
              required
            />
          </Row>
          <Row
            label="Confirm Passphrase"
            description={
              passphrasesMismatch
                ? 'Passphrases do not match.'
                : 'Re-enter your passphrase to confirm.'
            }
          >
            <Input
              id="confirmPassphrase"
              type="password"
              value={confirmPassphrase}
              onChange={(e) => setConfirmPassphrase(e.target.value)}
              placeholder="Confirm your passphrase"
              className={`h-10 w-full font-mono ${passphrasesMismatch ? 'border-destructive focus-visible:ring-destructive/20' : ''}`}
              required
            />
          </Row>
        </Section>
      </div>
    </form>
  );
}
