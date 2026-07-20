import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreate, useDelete, useListByService } from '#/hooks/useDomains';

export function ServiceDomains({ serviceId }: { serviceId: string }) {
  const { data: domainsRes, isLoading } = useListByService(serviceId);
  const createDomain = useCreate();
  const deleteDomain = useDelete();

  const [newDomain, setNewDomain] = useState('');

  const handleCreate = () => {
    if (!newDomain.trim()) return;
    createDomain.mutate(
      { serviceId, payload: { domainName: newDomain } },
      { onSuccess: () => setNewDomain('') }
    );
  };

  const domains = domainsRes?.data || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Domains</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        {isLoading ? (
          <div>Loading domains...</div>
        ) : domains.length === 0 ? (
          <p className="text-gray-500 text-sm">No domains configured.</p>
        ) : (
          <ul className="space-y-2">
            {domains.map((domain) => (
              <li key={domain.id} className="flex items-center justify-between rounded border p-2">
                <div>
                  <div className="font-semibold">{domain.domainName}</div>
                  <div className="text-gray-500 text-xs">Status: {domain.sslCertStatus}</div>
                </div>
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={() => deleteDomain.mutate({ id: domain.id })}
                  disabled={deleteDomain.isPending}
                >
                  Remove
                </Button>
              </li>
            ))}
          </ul>
        )}

        <div className="mt-4 flex flex-col gap-2 border-t pt-4">
          <Label htmlFor="new-domain">Add New Domain</Label>
          <div className="flex gap-2">
            <Input
              id="new-domain"
              placeholder="example.com"
              value={newDomain}
              onChange={(e) => setNewDomain(e.target.value)}
            />
            <Button onClick={handleCreate} disabled={createDomain.isPending}>
              Add
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
