import { Cloud, Database, Eye, EyeOff, Info, MoreVertical, Plus, Trash, X } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogTitle, DialogTrigger } from '#/components/ui/dialog';
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

  const { data: s3Destinations, isLoading } = useListS3Destinations('global');
  const createS3Dest = useCreateS3Destination();
  const deleteS3Dest = useDeleteS3Destination();

  const handleProviderChange = (value: string) => {
    setProvider(value);
    if (value === 'r2') {
      setEndpoint('https://<account_id>.r2.cloudflarestorage.com');
      setRegion('auto');
    } else if (value === 's3') {
      setEndpoint('https://s3.<region>.amazonaws.com');
      setRegion('us-east-1');
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
      await deleteS3Dest.mutateAsync({ id });
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
            <Button className="h-11 shrink-0 gap-2 bg-primary font-bold text-primary-foreground text-xs uppercase tracking-wider">
              <Plus className="h-4 w-4" />
              New Destination
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl [&>button]:hidden">
            <div className="flex items-center justify-between border-border/50 border-b px-6 py-4">
              <DialogTitle className="font-semibold text-lg">New S3 Storage</DialogTitle>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 rounded-full"
                onClick={() => setIsOpen(false)}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
            <div className="p-6">
              <p className="mb-6 text-muted-foreground text-sm">
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
                <div className="space-y-2">
                  <Label className="font-medium text-muted-foreground text-sm">Provider</Label>
                  <Select value={provider} onValueChange={handleProviderChange}>
                    <SelectTrigger className="h-11 bg-background/50">
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
                  <div className="space-y-2">
                    <Label className="font-medium text-muted-foreground text-sm">
                      Name <span className="text-yellow-500">*</span>
                    </Label>
                    <Input
                      required
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="h-11 bg-background/50 focus-visible:ring-yellow-500/50"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label className="font-medium text-muted-foreground text-sm">Description</Label>
                    <Input
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                      className="h-11 bg-background/50 focus-visible:ring-yellow-500/50"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label className="font-medium text-muted-foreground text-sm">
                    Endpoint <span className="text-yellow-500">*</span>
                  </Label>
                  <Input
                    required
                    value={endpoint}
                    onChange={(e) => setEndpoint(e.target.value)}
                    className="h-11 bg-background/50 font-mono text-sm focus-visible:ring-yellow-500/50"
                  />
                </div>

                <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label className="font-medium text-muted-foreground text-sm">
                      Bucket <span className="text-yellow-500">*</span>
                    </Label>
                    <Input
                      required
                      value={bucket}
                      onChange={(e) => setBucket(e.target.value)}
                      className="h-11 bg-background/50 font-mono focus-visible:ring-yellow-500/50"
                    />
                  </div>
                  <div className="space-y-2">
                    <div className="flex items-center gap-1.5">
                      <Label className="font-medium text-muted-foreground text-sm">
                        Region <span className="text-yellow-500">*</span>
                      </Label>
                      <Info className="h-3.5 w-3.5 text-yellow-500" />
                    </div>
                    <Input
                      required
                      value={region}
                      onChange={(e) => setRegion(e.target.value)}
                      className="h-11 bg-background/50 font-mono focus-visible:ring-yellow-500/50"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label className="font-medium text-muted-foreground text-sm">
                      Access Key <span className="text-yellow-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        required
                        type="password"
                        value={accessKeyId}
                        onChange={(e) => setAccessKeyId(e.target.value)}
                        className="h-11 bg-background/50 pr-10 font-mono focus-visible:ring-yellow-500/50"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute top-0 right-0 h-11 w-11 text-muted-foreground"
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <Label className="font-medium text-muted-foreground text-sm">
                      Secret Key <span className="text-yellow-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        required
                        type={showSecret ? 'text' : 'password'}
                        value={secretAccessKey}
                        onChange={(e) => setSecretAccessKey(e.target.value)}
                        className="h-11 bg-background/50 pr-10 font-mono focus-visible:ring-yellow-500/50"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => setShowSecret(!showSecret)}
                        className="absolute top-0 right-0 h-11 w-11 text-muted-foreground"
                      >
                        {showSecret ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </Button>
                    </div>
                  </div>
                </div>

                <div className="pt-2">
                  <Button
                    type="submit"
                    disabled={createS3Dest.isPending}
                    className="h-11 w-full border border-border bg-card font-medium text-foreground text-sm transition-colors hover:bg-muted"
                  >
                    {createS3Dest.isPending ? 'Validating...' : 'Validate Connection & Continue'}
                  </Button>
                </div>
              </form>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {list.length === 0 ? (
          <div className="col-span-full flex flex-col items-center justify-center rounded-2xl border border-border/50 border-dashed bg-card/20 py-16 text-center">
            <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 text-primary">
              <Cloud className="h-6 w-6" />
            </div>
            <h3 className="font-medium text-lg">No destinations</h3>
            <p className="mt-1 max-w-sm text-muted-foreground text-sm">
              Add an S3 destination to start backing up your databases remotely.
            </p>
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
                  <span className="max-w-[150px] truncate font-medium" title={dest.bucket}>
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
