import { useMutation } from '@tanstack/react-query';
import { Copy, Globe, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';

import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Switch } from '#/components/ui/switch';
import { apiClient } from '#/lib/apiClient';

interface DatabaseNetworkingProps {
  database: {
    id: string;
    isPublic: boolean;
    publicEndpoint?: string;
  };
  onUpdate: () => void;
}

export function DatabaseNetworking({ database, onUpdate }: DatabaseNetworkingProps) {
  const togglePublicAccess = useMutation({
    mutationFn: async (isPublic: boolean) => {
      return apiClient.put(`/databases/${database.id}`, { isPublic });
    },
    onSuccess: () => {
      toast.success('Networking settings updated');
      onUpdate();
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to update networking settings');
    },
  });

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard');
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Globe className="h-5 w-5" />
          Public Networking
        </CardTitle>
        <CardDescription>
          Expose your database to the public internet using a secure TCP endpoint.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="flex items-center justify-between rounded-lg border p-4">
          <div className="space-y-0.5">
            <h3 className="font-medium text-sm">Public Access</h3>
            <p className="text-muted-foreground text-sm">
              Allow external connections via ExternalDNS
            </p>
          </div>
          <Switch
            checked={database.isPublic}
            onCheckedChange={(checked) => togglePublicAccess.mutate(checked)}
            disabled={togglePublicAccess.isPending}
          />
        </div>

        {database.isPublic && database.publicEndpoint && (
          <div className="space-y-4 rounded-lg border bg-muted/50 p-4">
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="gap-1">
                <ShieldCheck className="h-3 w-3" />
                Let's Encrypt TCP SNI enabled
              </Badge>
            </div>

            <div className="space-y-2">
              <label className="font-medium text-sm">Public TCP Endpoint</label>
              <div className="flex gap-2">
                <Input readOnly value={database.publicEndpoint} className="font-mono" />
                <Button variant="secondary" onClick={() => handleCopy(database.publicEndpoint!)}>
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
