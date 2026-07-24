import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

import { Button } from '#/components/ui/button';
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
import { apiClient } from '#/lib/apiClient';

interface RowEditorModalProps {
  isOpen: boolean;
  onClose: () => void;
  databaseId: string;
  tableName: string;
  columns: string[];
  initialData?: Record<string, unknown> | null;
}

export function RowEditorModal({
  isOpen,
  onClose,
  databaseId,
  tableName,
  columns,
  initialData,
}: RowEditorModalProps) {
  const [formData, setFormData] = useState<Record<string, string>>({});
  const queryClient = useQueryClient();
  const isEditing = !!initialData;

  useEffect(() => {
    if (isOpen) {
      const initial = {} as Record<string, string>;
      for (const col of columns) {
        initial[col] = initialData ? String(initialData[col] ?? '') : '';
      }
      setFormData(initial);
    }
  }, [isOpen, initialData, columns]);

  const saveRow = useMutation({
    mutationFn: async (data: Record<string, string>) => {
      if (isEditing && initialData) {
        return apiClient.put(`/databases/${databaseId}/data/${tableName}`, {
          keys: initialData,
          data: data,
        });
      }
      return apiClient.post(`/databases/${databaseId}/data/${tableName}`, data);
    },
    onSuccess: () => {
      toast.success(isEditing ? 'Row updated successfully' : 'Row added successfully');
      queryClient.invalidateQueries({ queryKey: ['table-data', databaseId, tableName] });
      onClose();
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to save row');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    saveRow.mutate(formData);
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-106.25">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>{isEditing ? 'Edit Row' : 'Add New Row'}</DialogTitle>
            <DialogDescription>
              {isEditing
                ? `Make changes to the row in ${tableName}.`
                : `Insert a new record into ${tableName}.`}
            </DialogDescription>
          </DialogHeader>
          <div className="grid max-h-[60vh] gap-4 overflow-y-auto py-4">
            {columns.map((col) => (
              <div key={col} className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor={col} className="text-right">
                  {col}
                </Label>
                <Input
                  id={col}
                  value={formData[col] || ''}
                  onChange={(e) => setFormData({ ...formData, [col]: e.target.value })}
                  className="col-span-3 font-mono text-sm"
                />
              </div>
            ))}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={saveRow.isPending}>
              {saveRow.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save changes
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
