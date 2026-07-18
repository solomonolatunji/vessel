import { Download, Info, Package, Upload } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

export function MigrationSettings() {
  const [passphrase, setPassphrase] = useState('');
  const [confirmPassphrase, setConfirmPassphrase] = useState('');
  const [isExporting, setIsExporting] = useState(false);

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

    setIsExporting(true);
    try {
      // TODO: Implement actual export logic
      await new Promise((resolve) => setTimeout(resolve, 2000));
      toast.success('Instance bundle downloaded successfully');
      setPassphrase('');
      setConfirmPassphrase('');
    } catch (_error) {
      toast.error('Failed to export instance bundle');
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Upload className="h-4.5 w-4.5" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Export instance</h1>
            <p className="text-muted-foreground text-sm">
              Export and import your instance data. Create an encrypted bundle of your entire
              instance (database and configuration) for migration or backup purposes.
            </p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6">
        <div className="rounded-xl border border-blue-500/20 bg-blue-500/10 p-4">
          <div className="flex items-start gap-3">
            <div className="mt-0.5 rounded-full bg-blue-500/20 p-1">
              <Info className="h-4 w-4 text-blue-500" />
            </div>
            <div>
              <h3 className="font-medium text-blue-500 text-sm">Keep your passphrase safe</h3>
              <p className="mt-1 text-muted-foreground text-sm">
                The migration bundle is heavily encrypted. Without the exact passphrase, you will
                not be able to decrypt and import this data on a new instance.
              </p>
            </div>
          </div>
        </div>

        <div className="space-y-6 rounded-2xl border border-border/50 bg-card/40 p-6">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
              <Package className="h-5 w-5" />
            </div>
            <div>
              <h2 className="font-semibold text-lg">Export instance bundle</h2>
              <p className="text-muted-foreground text-sm">
                Download a secure bundle of your instance
              </p>
            </div>
          </div>

          <form onSubmit={handleExport} className="space-y-4">
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label
                  htmlFor="passphrase"
                  className="font-bold text-muted-foreground text-xs uppercase tracking-wider"
                >
                  MIGRATION PASSPHRASE
                </Label>
                <Input
                  id="passphrase"
                  type="password"
                  value={passphrase}
                  onChange={(e) => setPassphrase(e.target.value)}
                  placeholder="Enter a secure passphrase"
                  className="h-11 bg-background/50 font-mono"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label
                  htmlFor="confirmPassphrase"
                  className="font-bold text-muted-foreground text-xs uppercase tracking-wider"
                >
                  CONFIRM PASSPHRASE
                </Label>
                <Input
                  id="confirmPassphrase"
                  type="password"
                  value={confirmPassphrase}
                  onChange={(e) => setConfirmPassphrase(e.target.value)}
                  placeholder="Confirm your passphrase"
                  className="h-11 bg-background/50 font-mono"
                  required
                />
              </div>
            </div>

            <div className="flex justify-end pt-2">
              <Button
                type="submit"
                disabled={isExporting}
                className="h-11 bg-primary font-bold text-primary-foreground text-xs uppercase tracking-wider"
              >
                <Download className="mr-2 h-4 w-4" />
                {isExporting ? 'EXPORTING...' : 'DOWNLOAD BUNDLE'}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
