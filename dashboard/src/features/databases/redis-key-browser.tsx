import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Clock, Database, Search, Trash } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { apiClient } from '#/lib/apiClient';

interface RedisKeyBrowserProps {
  databaseId: string;
}

interface RedisKey {
  key: string;
  type: string;
  value: string;
  ttl: number;
}

export function RedisKeyBrowser({ databaseId }: RedisKeyBrowserProps) {
  const [search, setSearch] = useState('*');
  const queryClient = useQueryClient();

  const { data: keys, isLoading } = useQuery({
    queryKey: ['redis-keys', databaseId, search],
    queryFn: () =>
      apiClient.get<RedisKey[]>(
        `/databases/${databaseId}/redis/keys?pattern=${encodeURIComponent(search)}`
      ),
  });

  const deleteKey = useMutation({
    mutationFn: (key: string) =>
      apiClient.delete(`/databases/${databaseId}/redis/keys/${encodeURIComponent(key)}`),
    onSuccess: () => {
      toast.success('Key deleted successfully');
      queryClient.invalidateQueries({ queryKey: ['redis-keys', databaseId] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to delete key');
    },
  });

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'string':
        return 'bg-blue-500/10 text-blue-500';
      case 'hash':
        return 'bg-purple-500/10 text-purple-500';
      case 'list':
        return 'bg-orange-500/10 text-orange-500';
      case 'set':
        return 'bg-green-500/10 text-green-500';
      default:
        return 'bg-gray-500/10 text-gray-500';
    }
  };

  return (
    <Card className="flex h-[calc(100vh-12rem)] flex-col">
      <CardHeader className="border-b py-4">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Database className="h-5 w-5" />
              Redis Data Browser
            </CardTitle>
            <CardDescription>View and manage keys, TTLs, and data types.</CardDescription>
          </div>
          <div className="relative w-72">
            <Search className="absolute top-2.5 left-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="text"
              placeholder="Search pattern (e.g., user:*)"
              className="pl-9"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Key</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Value Summary</TableHead>
              <TableHead>TTL</TableHead>
              <TableHead className="w-[100px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                  Loading keys...
                </TableCell>
              </TableRow>
            ) : !keys || keys.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                  No keys found matching "{search}"
                </TableCell>
              </TableRow>
            ) : (
              keys.map((k) => (
                <TableRow key={k.key}>
                  <TableCell className="font-mono text-sm">{k.key}</TableCell>
                  <TableCell>
                    <Badge variant="secondary" className={getTypeColor(k.type)}>
                      {k.type}
                    </Badge>
                  </TableCell>
                  <TableCell className="max-w-xs truncate font-mono text-muted-foreground text-xs">
                    {k.value}
                  </TableCell>
                  <TableCell>
                    {k.ttl === -1 ? (
                      <Badge variant="outline">No Expiry</Badge>
                    ) : (
                      <span className="flex items-center gap-1 text-muted-foreground text-sm">
                        <Clock className="h-3 w-3" />
                        {k.ttl}s
                      </span>
                    )}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        if (confirm(`Delete key ${k.key}?`)) {
                          deleteKey.mutate(k.key);
                        }
                      }}
                    >
                      <Trash className="h-4 w-4 text-destructive" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
