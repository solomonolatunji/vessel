import { Database, Eye, EyeOff, Info, MoreVertical, Plus, Trash } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import {
  useCreateS3Destination,
  useDeleteS3Destination,
  useListS3Destinations,
} from '#/hooks/useBackups';

export function S3DestinationsList() {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [provider, setProvider] = useState('r2');
  const [endpoint, setEndpoint] = useState('');
  const [bucket, setBucket] = useState('');
  const [region, setRegion] = useState('us-east-1');
  const [accessKeyId, setAccessKeyId] = useState('');
  const [secretAccessKey, setSecretAccessKey] = useState('');
  const [showSecret, setShowSecret] = useState(false);

  const { data: s3Destinations, isLoading } = useListS3Destinations();
  const createS3Dest = useCreateS3Destination();
  const deleteS3Dest = useDeleteS3Destination();

  const handleProviderChange = (value: string) => {
    setProvider(value);
    if (value === 'r2') {
      setEndpoint('https://<account_id>.r2.cloudflarestorage.com');
      setRegion('auto');
    } else if (value === 's3') {
      setRegion('us-east-1');
      setEndpoint('https://s3.us-east-1.amazonaws.com');
    } else {
      setEndpoint('');
      setRegion('');
    }
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createS3Dest.mutateAsync({
        payload: {
          projectId: 'global',
          name,
          description,
          provider,
          endpoint,
          bucket,
          region,
          accessKeyId,
          secretAccessKey,
        },
      });
      toast.success('S3 destination added successfully');
      setIsOpen(false);
      setName('');
      setDescription('');
      setProvider('r2');
      setEndpoint('');
      setBucket('');
      setRegion('us-east-1');
      setAccessKeyId('');
      setSecretAccessKey('');
    } catch (_error) {
      toast.error('Failed to add S3 destination');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteS3Dest.mutateAsync({ id, projectId: 'global' });
      toast.success('S3 destination deleted successfully');
    } catch (_error) {
      toast.error('Failed to delete S3 destination');
    }
  };

  if (isLoading) {
    return <div className="p-6">Loading configuration...</div>;
  }

  const list = s3Destinations?.data || [];

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Database className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">S3 Destinations</h1>
            <p className="text-muted-foreground text-sm">
              Manage your S3 compatible storage connections. Store credentials in Vessl to use them
              as database backup targets.
            </p>
          </div>
        </div>
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              NEW DESTINATION
            </Button>
          </DialogTrigger>
          <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-150 [&>button]:hidden">
            <div className="px-5 pt-5 pb-4">
              <div className="flex items-start justify-between">
                <div className="flex flex-col">
                  <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                    <Database className="h-5 w-5 text-primary" />
                    New S3 Storage
                  </DialogTitle>
                  <DialogDescription>Connect compatible S3 storage</DialogDescription>
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
            <div className="px-5 pt-4 pb-5">
              <p className="mb-5 text-muted-foreground text-sm">
                For more details, please visit the{' '}
                <a
                  href="https://docs.vessl.dev"
                  className="text-primary underline"
                  target="_blank"
                  rel="noreferrer"
                >
                  Vessl Docs
                </a>
                .
              </p>
              <form onSubmit={handleSave} className="space-y-5">
                <div className="space-y-2.5">
                  <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                    PROVIDER
                  </Label>
                  <Select value={provider} onValueChange={handleProviderChange}>
                    <SelectTrigger className="h-10 bg-background/50 text-sm">
                      <SelectValue placeholder="Select provider" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="r2">Cloudflare R2</SelectItem>
                      <SelectItem value="s3">AWS S3</SelectItem>
                      <SelectItem value="custom">Custom S3 Compatible</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
                  <div className="space-y-2.5">
                    <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                      NAME <span className="text-yellow-500">*</span>
                    </Label>
                    <Input
                      required
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="h-10 bg-background/50 text-sm focus-visible:ring-yellow-500/50"
                    />
                  </div>
                  <div className="space-y-2.5">
                    <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                      DESCRIPTION
                    </Label>
                    <Input
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                      className="h-10 bg-background/50 text-sm focus-visible:ring-yellow-500/50"
                    />
                  </div>
                </div>

                <div className="space-y-2.5">
                  <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                    ENDPOINT <span className="text-yellow-500">*</span>
                  </Label>
                  <Input
                    required
                    value={endpoint}
                    onChange={(e) => setEndpoint(e.target.value)}
                    className="h-10 bg-background/50 font-mono text-sm focus-visible:ring-yellow-500/50"
                  />
                </div>

                <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
                  <div className="space-y-2.5">
                    <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                      BUCKET <span className="text-yellow-500">*</span>
                    </Label>
                    <Input
                      required
                      value={bucket}
                      onChange={(e) => setBucket(e.target.value)}
                      className="h-10 bg-background/50 font-mono text-sm focus-visible:ring-yellow-500/50"
                    />
                  </div>
                  <div className="space-y-2.5">
                    <div className="flex items-center gap-1.5">
                      <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                        REGION <span className="text-yellow-500">*</span>
                      </Label>
                      <Info className="h-3.5 w-3.5 text-yellow-500" />
                    </div>
                    <Input
                      required
                      value={region}
                      onChange={(e) => setRegion(e.target.value)}
                      className="h-10 bg-background/50 font-mono text-sm focus-visible:ring-yellow-500/50"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
                  <div className="space-y-2.5">
                    <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                      ACCESS KEY <span className="text-yellow-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        required
                        type="password"
                        value={accessKeyId}
                        onChange={(e) => setAccessKeyId(e.target.value)}
                        className="h-10 bg-background/50 pr-10 font-mono text-sm focus-visible:ring-yellow-500/50"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute top-0 right-0 h-10 w-10 text-muted-foreground"
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                  <div className="space-y-2.5">
                    <Label className="font-mono font-semibold text-[10px] text-muted-foreground uppercase tracking-[0.2em]">
                      SECRET KEY <span className="text-yellow-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        required
                        type={showSecret ? 'text' : 'password'}
                        value={secretAccessKey}
                        onChange={(e) => setSecretAccessKey(e.target.value)}
                        className="h-10 bg-background/50 pr-10 font-mono text-sm focus-visible:ring-yellow-500/50"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => setShowSecret(!showSecret)}
                        className="absolute top-0 right-0 h-10 w-10 text-muted-foreground"
                      >
                        {showSecret ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </Button>
                    </div>
                  </div>
                </div>

                <div className="flex items-center justify-end gap-3 pt-2">
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => setIsOpen(false)}
                    className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    disabled={createS3Dest.isPending}
                    className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
                  >
                    <Plus className="h-3.5 w-3.5" />
                    {createS3Dest.isPending ? 'Validating...' : 'Validate Connection'}
                  </Button>
                </div>
              </form>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {list.length === 0 ? (
          <div className="col-span-full flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-primary/20 bg-primary/10">
              <Database className="h-5 w-5 text-primary" />
            </div>
            <h3 className="mt-4 font-bold text-foreground text-lg tracking-tight">
              No destinations
            </h3>
            <p className="mt-1 max-w-sm text-center text-muted-foreground text-sm">
              Add an S3 destination to start backing up your databases remotely.
            </p>
            <Button className="mt-6 gap-2" onClick={() => setIsOpen(true)}>
              <Plus className="h-4 w-4" />
              NEW DESTINATION
            </Button>
          </div>
        ) : (
          list.map((dest) => (
            <div
              key={dest.id}
              className="group relative flex flex-col rounded-2xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-border"
            >
              <div className="mb-4 flex items-start justify-between gap-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
                  <Database className="h-5 w-5" />
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 opacity-0 transition-opacity group-hover:opacity-100"
                    >
                      <MoreVertical className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    align="end"
                    className="w-40 border-border/50 bg-card/95 backdrop-blur-xl"
                  >
                    <DropdownMenuItem
                      className="text-destructive focus:bg-destructive/10 focus:text-destructive"
                      onClick={() => handleDelete(dest.id)}
                    >
                      <Trash className="mr-2 h-4 w-4" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
              <div className="space-y-1">
                <h3 className="font-semibold text-lg">{dest.name}</h3>
                {dest.description && (
                  <p className="text-muted-foreground text-sm">{dest.description}</p>
                )}
              </div>
              <div className="mt-4 flex flex-col gap-2 border-border/50 border-t pt-4 text-sm">
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Provider</span>
                  <span className="font-medium capitalize">{dest.provider}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Bucket</span>
                  <span className="max-w-37.5 truncate font-medium" title={dest.bucket}>
                    {dest.bucket}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Region</span>
                  <span className="font-medium uppercase">{dest.region}</span>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
