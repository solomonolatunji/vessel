import { useQuery } from '@tanstack/react-query';
import { Database } from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import type { GetDatabasesResponse } from '#/interfaces/database';
import { apiClient } from '#/lib/apiClient';

export function DatabaseList() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['databases'],
    queryFn: () => apiClient.get<GetDatabasesResponse>('/databases'),
  });

  if (isLoading) {
    return <div className="animate-pulse p-4 text-muted-foreground">Loading databases...</div>;
  }

  if (error) {
    return <div className="p-4 text-red-500">Failed to load databases</div>;
  }

  const databases = data?.data || [];

  if (databases.length === 0) {
    return (
      <Card className="flex flex-col items-center justify-center p-12 text-center">
        <div className="mb-4 rounded-full bg-secondary p-4">
          <Database className="h-8 w-8 text-muted-foreground" />
        </div>
        <CardTitle className="mb-2">No databases found</CardTitle>
        <CardDescription>
          You haven't provisioned any databases yet. Spin up a new instance to get started.
        </CardDescription>
      </Card>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
      {databases.map((db) => (
        <Card key={db.id} className="transition-colors hover:border-primary/50">
          <CardHeader className="flex flex-row items-start justify-between pb-3">
            <div>
              <CardTitle className="flex items-center gap-2 text-base">
                <Database className="h-4 w-4 text-muted-foreground" />
                {db.name}
              </CardTitle>
              <CardDescription className="mt-1 flex items-center gap-2">
                {db.engine}{' '}
                <Badge variant="outline" className="text-xs">
                  {db.version}
                </Badge>
              </CardDescription>
            </div>
            <Badge
              variant={
                db.status === 'running'
                  ? 'default'
                  : db.status === 'error'
                    ? 'destructive'
                    : 'secondary'
              }
            >
              {db.status}
            </Badge>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Internal DNS:</span>
              <span className="font-mono">{db.internalDns || 'Pending...'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Port:</span>
              <span className="font-mono">{db.port}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Volume:</span>
              <span className="max-w-30 truncate font-mono text-xs" title={db.volumePath}>
                {db.volumePath}
              </span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
