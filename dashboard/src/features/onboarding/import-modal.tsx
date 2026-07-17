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
      <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-187.5 [&>button]:hidden">
        <div className="flex flex-col p-8 pb-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
                <Database className="h-6 w-6 text-primary" />
              </div>
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-2xl tracking-tight">
                  Import Vessl
                </DialogTitle>
                <DialogDescription className="mt-1 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  RESTORE AN ENCRYPTED BUNDLE FROM ANOTHER SERVER.
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

        <div className="grid grid-cols-1 gap-6 p-8 md:grid-cols-2 md:gap-8">
          <div className="space-y-3">
            <Label
              htmlFor="migrationBundle"
              className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]"
            >
              MIGRATION BUNDLE
            </Label>
            <div className="relative flex h-12 w-full items-center gap-3 rounded-xl border border-border/50 bg-background/80 px-4 transition-all duration-300 focus-within:border-primary/50 focus-within:ring-2 focus-within:ring-primary/20 hover:bg-background">
              <CloudUpload className="h-5 w-5 text-muted-foreground" />
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
          <div className="space-y-3">
            <Label
              htmlFor="passphrase"
              className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]"
            >
              PASSPHRASE
            </Label>
            <Input
              id="passphrase"
              type="password"
              value={passphrase}
              onChange={(e) => setPassphrase(e.target.value)}
              className="h-12 rounded-xl border-border/50 bg-background/80 px-4 text-sm transition-all duration-300 focus:border-primary/50 focus:ring-2 focus:ring-primary/20"
            />
          </div>
        </div>

        <div className="h-px w-full bg-border/50" />

        <div className="flex items-center justify-end gap-6 p-8 pt-6">
          <Button
            variant="ghost"
            onClick={() => {
              setFileName('');
              setPassphrase('');
              setSelectedFile(null);
              onOpenChange(false);
            }}
            className="flex h-11 items-center gap-2 rounded-xl px-6 font-semibold text-muted-foreground text-xs uppercase tracking-widest hover:text-foreground"
          >
            CANCEL
          </Button>
          <Button
            variant="outline"
            onClick={handleImport}
            disabled={isPending || !selectedFile || !passphrase}
            className="flex h-11 items-center gap-2 rounded-xl border-primary/20 bg-primary/10 px-6 font-semibold text-primary text-xs uppercase tracking-widest hover:bg-primary/20 hover:text-primary"
          >
            <Check className="h-4 w-4" /> {isPending ? 'IMPORTING...' : 'IMPORT'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
