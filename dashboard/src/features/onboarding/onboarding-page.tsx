import { zodResolver } from '@hookform/resolvers/zod';
import { useNavigate } from '@tanstack/react-router';
import { ArrowRight, Building2 } from 'lucide-react';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useListWorkspaces, useUpdateWorkspace } from '#/hooks/useWorkspaces';
import type { Workspace } from '#/interfaces/workspace';
import { authStore } from '#/stores/authStore';

const onboardingSchema = z.object({
  name: z.string().min(2, 'Workspace name must be at least 2 characters'),
});

type OnboardingSchema = z.infer<typeof onboardingSchema>;

export const OnboardingPage = () => {
  const navigate = useNavigate();
  const user = authStore.state.user;

  const { data: workspacesResponse, isLoading } = useListWorkspaces();
  const { mutateAsync: updateWorkspace, isPending } = useUpdateWorkspace();

  const workspaces: Workspace[] =
    (workspacesResponse as { data?: Workspace[] } | undefined)?.data ?? [];
  const defaultWorkspace = workspaces[0];

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<OnboardingSchema>({
    resolver: zodResolver(onboardingSchema),
    defaultValues: {
      name: '',
    },
  });

  useEffect(() => {
    if (user?.name) {
      setValue('name', `${user.name}'s Workspace`);
    } else {
      setValue('name', 'Personal Workspace');
    }
  }, [user, setValue]);

  const onSubmit = async (data: OnboardingSchema) => {
    if (defaultWorkspace) {
      try {
        await updateWorkspace({
          id: defaultWorkspace.id,
          payload: { name: data.name },
        });
        navigate({ to: '/' });
      } catch {}
    } else {
      navigate({ to: '/' });
    }
  };

  const handleSkip = () => {
    navigate({ to: '/' });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center space-y-2">
          <div className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 mb-4">
            <Building2 className="h-6 w-6 text-primary" />
          </div>
          <h1 className="text-3xl font-bold tracking-tight">Welcome to Vessl!</h1>
          <p className="text-muted-foreground text-sm">Let's get your first workspace set up.</p>
        </div>

        <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name" className="text-sm font-medium">
                Workspace Name
              </Label>
              <Input
                id="name"
                placeholder="My Awesome Workspace"
                className="h-12 bg-background text-base"
                {...register('name')}
              />
              {errors.name && <p className="text-[13px] text-destructive">{errors.name.message}</p>}
            </div>

            <div className="flex flex-col space-y-3 pt-2">
              <Button
                type="submit"
                className="h-12 w-full text-base font-medium"
                disabled={isPending}
              >
                {isPending ? 'Saving...' : 'Continue to Dashboard'}
                {!isPending && <ArrowRight className="ml-2 h-4 w-4" />}
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="h-12 w-full text-muted-foreground"
                onClick={handleSkip}
                disabled={isPending}
              >
                Skip for now
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};
