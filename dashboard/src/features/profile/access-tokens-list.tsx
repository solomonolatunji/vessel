import { Key, Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useCreateToken, useDeleteToken, useListTokens } from '#/hooks/useProfile';

export function AccessTokensList() {
  const { data: response, isLoading } = useListTokens();
  const tokens = response?.data || [];

  const createToken = useCreateToken();
  const deleteToken = useDeleteToken();

  const [createOpen, setCreateOpen] = useState(false);
  const [newTokenName, setNewTokenName] = useState('');
  const [createdToken, setCreatedToken] = useState<string | null>(null);

  const handleCreateToken = async () => {
    if (!newTokenName.trim()) return;

    try {
      const res = await createToken.mutateAsync({
        payload: {
          name: newTokenName,
          accessLevel: 'read_write',
          projectScope: 'all',
          allowedProjects: [],
        },
      });
      setCreatedToken((res as any)?.data?.plain || 'token-created-successfully');
    } catch (error) {
      console.error('Failed to create token', error);
    }
  };

  const handleCloseCreate = () => {
    setCreateOpen(false);
    setNewTokenName('');
    setCreatedToken(null);
  };

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this token?')) {
      await deleteToken.mutateAsync({ id });
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <div>
          <CardTitle>Personal Access Tokens</CardTitle>
          <CardDescription className="mt-1.5">
            Generate tokens for API access or CLI authentication.
          </CardDescription>
        </div>
        <Button onClick={() => setCreateOpen(true)} className="gap-2">
          <Plus className="h-4 w-4" />
          Generate New Token
        </Button>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center p-6">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : tokens.length === 0 ? (
          <div className="flex flex-col items-center justify-center p-6 text-center text-muted-foreground">
            <Key className="mb-4 h-8 w-8 opacity-20" />
            <p>No access tokens generated yet.</p>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Last Used</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {tokens.map((token: any) => (
                <TableRow key={token.id}>
                  <TableCell className="font-medium">{token.name}</TableCell>
                  <TableCell>{new Date(token.createdAt).toLocaleString()}</TableCell>
                  <TableCell>
                    {token.lastUsedAt ? new Date(token.lastUsedAt).toLocaleString() : 'Never'}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-destructive hover:bg-destructive/10 hover:text-destructive"
                      onClick={() => handleDelete(token.id)}
                      disabled={deleteToken.isPending}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>

      <Dialog open={createOpen} onOpenChange={(open) => !open && handleCloseCreate()}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Generate Access Token</DialogTitle>
            <DialogDescription>
              {createdToken
                ? 'Make sure to copy your personal access token now. You won’t be able to see it again!'
                : 'Enter a name for your new personal access token.'}
            </DialogDescription>
          </DialogHeader>

          {createdToken ? (
            <div className="space-y-4 py-4">
              <div className="break-all rounded-md bg-muted p-4 font-mono text-sm">
                {createdToken}
              </div>
            </div>
          ) : (
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="token-name">Token Name</Label>
                <Input
                  id="token-name"
                  placeholder="e.g. CLI access, CI/CD pipeline"
                  value={newTokenName}
                  onChange={(e) => setNewTokenName(e.target.value)}
                />
              </div>
            </div>
          )}

          <DialogFooter>
            {createdToken ? (
              <Button onClick={handleCloseCreate}>Done</Button>
            ) : (
              <>
                <Button variant="outline" onClick={handleCloseCreate}>
                  Cancel
                </Button>
                <Button
                  onClick={handleCreateToken}
                  disabled={!newTokenName.trim() || createToken.isPending}
                >
                  {createToken.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                  Generate Token
                </Button>
              </>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
}
