import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { ImportDatabaseResponse } from '#/interfaces/database';
import { apiClient } from '#/lib/apiClient';

const schema = z.object({
  sourceUrl: z.string().min(1, 'Source URL is required').url('Must be a valid URL'),
});

type FormData = z.infer<typeof schema>;

interface Props {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  databaseId: string;
}

export function DatabaseImportModal({ isOpen, onOpenChange, databaseId }: Props) {
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      sourceUrl: '',
    },
  });

  const importMutation = useMutation({
    mutationFn: (data: FormData) => {
      return apiClient.post<ImportDatabaseResponse>(`/databases/${databaseId}/import`, data);
    },
    onSuccess: () => {
      toast.success('Import process started successfully');
      queryClient.invalidateQueries({ queryKey: ['databases', databaseId] });
      onOpenChange(false);
    },
  });

  const onSubmit = (data: FormData) => {
    importMutation.mutate(data);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Import Data</DialogTitle>
          <DialogDescription>
            Import data from an external database using a connection string (e.g. postgres:// or
            redis://).
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="mt-4 space-y-4">
          <div className="space-y-2">
            <Label>Source Connection URL</Label>
            <Input {...register('sourceUrl')} placeholder="postgres://user:pass@host:5432/db" />
            {errors.sourceUrl && <p className="text-red-500 text-sm">{errors.sourceUrl.message}</p>}
          </div>

          <div className="flex justify-end space-x-2 pt-4">
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={importMutation.isPending}>
              {importMutation.isPending ? 'Starting Import...' : 'Import Data'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
