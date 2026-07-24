import { format } from 'date-fns';
import { Loader2, Network, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useCreateLogDrain, useDeleteLogDrain, useListLogDrains } from '#/hooks/useApps';

export function LogDrainsCard({ serviceId }: { serviceId: string }) {
  const { data: drains, isLoading } = useListLogDrains(serviceId);
  const createDrain = useCreateLogDrain();
  const deleteDrain = useDeleteLogDrain();

  const [isOpen, setIsOpen] = useState(false);
  const [drainType, setDrainType] = useState('webhook');
  const [endpoint, setEndpoint] = useState('');
  const [authToken, setAuthToken] = useState('');

  const handleCreate = async () => {
    if (!endpoint) {
      toast.error('Endpoint URL is required');
      return;
    }

    try {
      await createDrain.mutateAsync({
        appId: serviceId,
        payload: {
          drainType,
          endpointUrl: endpoint,
          authToken: authToken || undefined,
        },
      });
      toast.success('Log drain created successfully');
      setIsOpen(false);
      setEndpoint('');
      setAuthToken('');
      setDrainType('webhook');
    } catch (_error) {
      toast.error('Failed to create log drain');
    }
  };

  const handleDelete = async (drainId: string) => {
    if (!confirm('Are you sure you want to delete this log drain?')) return;
    try {
      await deleteDrain.mutateAsync({ appId: serviceId, drainId });
      toast.success('Log drain deleted successfully');
    } catch (_error) {
      toast.error('Failed to delete log drain');
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="text-xl">Log Drains</CardTitle>
          <CardDescription>Forward container logs to external monitoring services.</CardDescription>
        </div>
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Add Drain
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add Log Drain</DialogTitle>
              <DialogDescription>Configure an external service to receive logs.</DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label>Destination Type</Label>
                <Select value={drainType} onValueChange={setDrainType}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="webhook">Custom Webhook</SelectItem>
                    <SelectItem value="axiom">Axiom</SelectItem>
                    <SelectItem value="new_relic">New Relic</SelectItem>
                    <SelectItem value="datadog">Datadog</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <Label>Endpoint URL</Label>
                <Input
                  placeholder="https://example.com/logs"
                  value={endpoint}
                  onChange={(e) => setEndpoint(e.target.value)}
                />
              </div>
              <div className="grid gap-2">
                <Label>Auth Token (Optional)</Label>
                <Input
                  type="password"
                  placeholder="Bearer token or API key"
                  value={authToken}
                  onChange={(e) => setAuthToken(e.target.value)}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreate} disabled={createDrain.isPending}>
                {createDrain.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Create Drain
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center p-8">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : drains?.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-lg border border-dashed p-8 text-center">
            <Network className="mb-4 h-8 w-8 text-muted-foreground" />
            <h3 className="font-semibold">No log drains configured</h3>
            <p className="mt-1 text-muted-foreground text-sm">
              Add a log drain to stream logs to external services.
            </p>
          </div>
        ) : (
          <div className="rounded-md border">
            <div className="grid grid-cols-[1fr_2fr_1fr_auto] items-center gap-4 border-b bg-muted/50 p-4 font-medium text-muted-foreground text-sm">
              <div>Type</div>
              <div>Endpoint</div>
              <div>Added</div>
              <div className="w-8"></div>
            </div>
            <div className="divide-y">
              {drains?.map((drain: any) => (
                <div
                  key={drain.id}
                  className="grid grid-cols-[1fr_2fr_1fr_auto] items-center gap-4 p-4 text-sm"
                >
                  <div className="font-medium capitalize">{drain.drainType}</div>
                  <div className="truncate text-muted-foreground">{drain.endpointUrl}</div>
                  <div className="text-muted-foreground">
                    {format(new Date(drain.createdAt), 'MMM d, yyyy')}
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
                    onClick={() => handleDelete(drain.id)}
                    disabled={deleteDrain.isPending}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
