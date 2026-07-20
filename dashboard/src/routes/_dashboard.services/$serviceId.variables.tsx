import { zodResolver } from '@hookform/resolvers/zod';
import { createFileRoute } from '@tanstack/react-router';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Switch } from '#/components/ui/switch';
import { useCreateVariable, useDeleteVariable, useListVariables } from '#/hooks/useApps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/variables')({
  component: VariablesTab,
});

const variableSchema = z.object({
  key: z.string().min(1, 'Key is required'),
  value: z.string().min(1, 'Value is required'),
  isSecret: z.boolean().default(false),
});

type VariableFormValues = z.infer<typeof variableSchema>;

function VariablesTab() {
  const { serviceId } = Route.useParams();
  const { data: variablesData, isLoading } = useListVariables(serviceId);
  const { mutateAsync: createVar, isPending: isCreating } = useCreateVariable();
  const { mutateAsync: deleteVar } = useDeleteVariable();

  const [visibleValues, setVisibleValues] = useState<Record<string, boolean>>({});

  const form = useForm<VariableFormValues>({
    resolver: zodResolver(variableSchema),
    defaultValues: {
      key: '',
      value: '',
      isSecret: false,
    },
  });

  const onSubmit = async (data: VariableFormValues) => {
    try {
      await createVar({ appId: serviceId, payload: data });
      toast.success('Variable created successfully');
      form.reset();
    } catch (error: any) {
      toast.error(error?.message || 'Failed to create variable');
    }
  };

  const handleDelete = async (varId: string) => {
    if (!confirm('Are you sure you want to delete this variable?')) return;
    try {
      await deleteVar({ appId: serviceId, varId });
      toast.success('Variable deleted successfully');
    } catch (error: any) {
      toast.error(error?.message || 'Failed to delete variable');
    }
  };

  const toggleVisibility = (varId: string) => {
    setVisibleValues((prev) => ({
      ...prev,
      [varId]: !prev[varId],
    }));
  };

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const variables = variablesData?.data || [];

  return (
    <div className="space-y-6">
      <div>
        <h2 className="font-medium text-lg">Environment Variables</h2>
        <p className="text-muted-foreground text-sm">
          Manage environment variables and secrets for your service.
        </p>
      </div>

      <div className="rounded-lg border bg-card p-6">
        <h3 className="mb-4 font-medium text-sm">Add New Variable</h3>
        <form onSubmit={form.handleSubmit(onSubmit)} className="flex items-start gap-4">
          <div className="flex-1 space-y-2">
            <Input
              placeholder="API_KEY"
              {...form.register('key')}
              className={form.formState.errors.key ? 'border-destructive' : ''}
            />
            {form.formState.errors.key && (
              <p className="text-destructive text-xs">{form.formState.errors.key.message}</p>
            )}
          </div>
          <div className="flex-1 space-y-2">
            <Input
              placeholder="your-secret-value"
              {...form.register('value')}
              className={form.formState.errors.value ? 'border-destructive' : ''}
            />
            {form.formState.errors.value && (
              <p className="text-destructive text-xs">{form.formState.errors.value.message}</p>
            )}
          </div>
          <div className="flex items-center space-x-2 pt-2">
            <Switch
              checked={form.watch('isSecret')}
              onCheckedChange={(checked) => form.setValue('isSecret', checked)}
              id="isSecret"
            />
            <label htmlFor="isSecret" className="font-medium text-sm">
              Secret
            </label>
          </div>
          <Button type="submit" disabled={isCreating} className="pt-2">
            {isCreating ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Plus className="mr-2 h-4 w-4" />
            )}
            Add
          </Button>
        </form>
      </div>

      <div className="rounded-lg border bg-card">
        <div className="grid grid-cols-[1fr_1fr_auto_auto] items-center gap-4 border-b bg-muted/50 px-6 py-3 font-medium text-sm">
          <div>Key</div>
          <div>Value</div>
          <div>Type</div>
          <div className="w-[100px] text-right">Actions</div>
        </div>
        <div className="divide-y">
          {variables.length === 0 ? (
            <div className="p-6 text-center text-muted-foreground text-sm">
              No environment variables found.
            </div>
          ) : (
            variables.map((v) => (
              <div
                key={v.id}
                className="grid grid-cols-[1fr_1fr_auto_auto] items-center gap-4 px-6 py-4"
              >
                <div className="font-medium font-mono text-sm">{v.key}</div>
                <div className="flex items-center gap-2">
                  <div className="font-mono text-muted-foreground text-sm">
                    {v.isSecret && !visibleValues[v.id] ? '••••••••••••••••' : v.value}
                  </div>
                  {v.isSecret && (
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6"
                      onClick={() => toggleVisibility(v.id)}
                    >
                      {visibleValues[v.id] ? (
                        <EyeOff className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </Button>
                  )}
                </div>
                <div>
                  <span className="rounded-full bg-secondary px-2 py-1 font-medium text-xs">
                    {v.isSecret ? 'Secret' : 'Plain'}
                  </span>
                </div>
                <div className="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
                    onClick={() => handleDelete(v.id)}
                  >
                    <Trash className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
