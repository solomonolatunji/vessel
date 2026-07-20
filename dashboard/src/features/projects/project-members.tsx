import { Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import { useAddMember, useListMembers, useRemoveMember } from '#/hooks/useProjectSettings';

export function ProjectMembers({ projectId }: { projectId: string }) {
  const { data: members, isLoading } = useListMembers(projectId);
  const { mutateAsync: addMember, isPending: isAdding } = useAddMember();
  const { mutateAsync: removeMember, isPending: isRemoving } = useRemoveMember();

  const [email, setEmail] = useState('');
  const [permission, setPermission] = useState<any>('member');

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;

    try {
      await addMember({
        projectId,
        payload: { email, permission },
      });
      setEmail('');
      setPermission('member');
      toast.success('Member added successfully');
    } catch (err: any) {
      toast.error(err.message || 'Failed to add member');
    }
  };

  const handleRemove = async (id: string) => {
    try {
      await removeMember({ projectId, memberId: id });
      toast.success('Member removed successfully');
    } catch (err: any) {
      toast.error(err.message || 'Failed to remove member');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2">
        <h2 className="font-medium text-lg">Project Members</h2>
        <p className="text-gray-500 text-sm">
          Manage who has access to this project. Members will only be able to view and manage
          resources within this project.
        </p>
      </div>

      <form onSubmit={handleAdd} className="flex items-end gap-4 rounded-lg border bg-gray-50 p-4">
        <div className="flex-1 space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="colleague@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="w-48 space-y-2">
          <Label htmlFor="permission">Role</Label>
          <Select value={permission} onValueChange={setPermission}>
            <SelectTrigger id="permission">
              <SelectValue placeholder="Select role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="owner">Owner</SelectItem>
              <SelectItem value="admin">Admin</SelectItem>
              <SelectItem value="member">Member</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Button type="submit" disabled={isAdding}>
          {isAdding ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            <Plus className="mr-2 h-4 w-4" />
          )}
          Add Member
        </Button>
      </form>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>Joined At</TableHead>
              <TableHead className="w-25"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center">
                  <Loader2 className="mx-auto h-6 w-6 animate-spin text-gray-400" />
                </TableCell>
              </TableRow>
            ) : members?.data?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center text-gray-500">
                  No members found.
                </TableCell>
              </TableRow>
            ) : (
              members?.data?.map((member: any) => (
                <TableRow key={member.id}>
                  <TableCell>{member.email}</TableCell>
                  <TableCell className="capitalize">{member.permission}</TableCell>
                  <TableCell>{new Date(member.created_at).toLocaleDateString()}</TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-red-600 hover:bg-red-50 hover:text-red-700"
                      onClick={() => handleRemove(member.user_id)}
                      disabled={isRemoving}
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
