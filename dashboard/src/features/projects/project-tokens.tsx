import { Copy, Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
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
import { useCreateToken, useDeleteToken, useListTokens } from '#/hooks/useProjectSettings';

export function ProjectTokens({ projectId }: { projectId: string }) {
  const { data: tokens, isLoading } = useListTokens(projectId);
  const { mutateAsync: createToken, isPending: isCreating } = useCreateToken();
  const { mutateAsync: deleteToken, isPending: isDeleting } = useDeleteToken();

  const [name, setName] = useState('');
  const [newTokenValue, setNewTokenValue] = useState<string | null>(null);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name) return;

    try {
      const res = await createToken({
        projectId,
        payload: { name, environmentId: '', scopes: [] },
      });
      setName('');
      setNewTokenValue((res?.data as any)?.token || null);
      toast.success('Token created successfully');
    } catch (err: any) {
      toast.error(err.message || 'Failed to create token');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteToken({ projectId, id });
      toast.success('Token deleted successfully');
    } catch (err: any) {
      toast.error(err.message || 'Failed to delete token');
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard');
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2">
        <h2 className="font-medium text-lg">API Tokens</h2>
        <p className="text-gray-500 text-sm">
          Create API tokens to access the Vessl API for this specific project.
        </p>
      </div>

      <form
        onSubmit={handleCreate}
        className="flex items-end gap-4 rounded-lg border bg-gray-50 p-4"
      >
        <div className="flex-1 space-y-2">
          <Label htmlFor="tokenName">Token Name</Label>
          <Input
            id="tokenName"
            placeholder="e.g. CI/CD Pipeline"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <Button type="submit" disabled={isCreating}>
          {isCreating ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            <Plus className="mr-2 h-4 w-4" />
          )}
          Create Token
        </Button>
      </form>

      {newTokenValue && (
        <div className="space-y-2 rounded-lg border border-amber-200 bg-amber-50 p-4">
          <h3 className="font-medium text-amber-800">Save your new token</h3>
          <p className="text-amber-700 text-sm">
            This is the only time you will be able to see this token. Please save it somewhere safe.
          </p>
          <div className="mt-2 flex items-center gap-2">
            <code className="flex-1 break-all rounded bg-amber-100 px-3 py-2 text-amber-900">
              {newTokenValue}
            </code>
            <Button variant="outline" size="icon" onClick={() => copyToClipboard(newTokenValue)}>
              <Copy className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead className="w-25"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={3} className="h-24 text-center">
                  <Loader2 className="mx-auto h-6 w-6 animate-spin text-gray-400" />
                </TableCell>
              </TableRow>
            ) : tokens?.data?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} className="h-24 text-center text-gray-500">
                  No API tokens found.
                </TableCell>
              </TableRow>
            ) : (
              tokens?.data?.map((token: any) => (
                <TableRow key={token.id}>
                  <TableCell className="font-medium">{token.name}</TableCell>
                  <TableCell>{new Date(token.created_at).toLocaleDateString()}</TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-red-600 hover:bg-red-50 hover:text-red-700"
                      onClick={() => handleDelete(token.id)}
                      disabled={isDeleting}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
