import { Check, CloudUpload, Database } from 'lucide-react';
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
import { useImportSystem } from '#/hooks/useSystem';

export const ImportModal = ({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) => {
  const [fileName, setFileName] = useState<string>('');
  const [passphrase, setPassphrase] = useState<string>('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const { mutateAsync: importSystem, isPending } = useImportSystem();

  const handleImport = async () => {
    if (!selectedFile || !passphrase) return;
    try {
      const formData = new FormData();
      formData.append('bundle', selectedFile);
      formData.append('passphrase', passphrase);
      await importSystem(formData);
      toast.success('Successfully imported Vessl bundle!');
      onOpenChange(false);
      setFileName('');
      setPassphrase('');
      setSelectedFile(null);
    } catch (_error) {
      toast.error('Failed to import bundle.');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-[600px] [&>button]:hidden">
        <div className="px-5 pt-5 pb-4">
          <div className="flex items-start justify-between">
            <div className="flex flex-col">
              <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                <Database className="h-5 w-5 text-primary" />
                Import Vessl
              </DialogTitle>
              <DialogDescription>Restore an encrypted bundle</DialogDescription>
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

        <div className="grid grid-cols-1 gap-5 p-5 md:grid-cols-2">
          <div className="space-y-2.5">
            <Label
              htmlFor="migrationBundle"
              className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]"
            >
              MIGRATION BUNDLE
            </Label>
            <div className="relative flex h-10 w-full items-center gap-3 rounded-lg border border-border/50 bg-background/80 px-3 transition-all duration-300 focus-within:border-primary/50 focus-within:ring-1 focus-within:ring-primary/20 hover:bg-background">
              <CloudUpload className="h-4 w-4 text-muted-foreground" />
              <span className="truncate font-mono text-foreground/90 text-sm">
                {fileName || 'Choose .vessl file'}
              </span>
              <Input
                id="migrationBundle"
                type="file"
                accept=".vessl"
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    setFileName(file.name);
                    setSelectedFile(file);
                  } else {
                    setFileName('');
                    setSelectedFile(null);
                  }
                }}
                className="absolute inset-0 h-full w-full cursor-pointer opacity-0"
              />
            </div>
          </div>
          <div className="space-y-2.5">
            <Label
              htmlFor="passphrase"
              className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]"
            >
              PASSPHRASE
            </Label>
            <Input
              id="passphrase"
              type="password"
              value={passphrase}
              onChange={(e) => setPassphrase(e.target.value)}
              className="h-10 rounded-lg border-border/50 bg-background/80 px-3 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-1 focus:ring-primary/20"
            />
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 p-5 pt-0">
          <Button
            variant="ghost"
            onClick={() => {
              setFileName('');
              setPassphrase('');
              setSelectedFile(null);
              onOpenChange(false);
            }}
            className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            Cancel
          </Button>
          <Button
            onClick={handleImport}
            disabled={isPending || !selectedFile || !passphrase}
            className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
          >
            <Check className="h-3.5 w-3.5" /> {isPending ? 'Importing...' : 'Import'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
