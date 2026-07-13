import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { toast } from 'sonner';
import { authService } from '#/services/auth';
import { authActions } from '#/stores/authStore';

export const useLogin = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: (credentials: Parameters<typeof authService.login>[0]) =>
      authService.login(credentials),
    onSuccess: async (data) => {
      if (!data?.token || !data?.user) {
        toast.error('Login failed: invalid response from server');
        return;
      }

      // 1. Persist auth state FIRST
      authActions.setAuth(data.token, data.user);

      // 2. Clear stale queries
      queryClient.clear();

      // 3. Navigate — invalidate happens automatically via beforeLoad
      await router.navigate({ to: '/' });

      toast.success('Welcome back!');
    },
    onError: (error: Error) => {
      toast.error(error?.message || 'Login failed. Please try again.');
    },
  });
};

export const useRegister = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: (details: Parameters<typeof authService.register>[0]) =>
      authService.register(details),
    onSuccess: async (data) => {
      if (!data?.token || !data?.user) {
        toast.error('Registration failed: invalid response from server');
        return;
      }

      authActions.setAuth(data.token, data.user);
      queryClient.clear();
      await router.navigate({ to: '/onboarding' });

      toast.success('Account created! Welcome to Vessl.');
    },
    onError: (error: Error) => {
      toast.error(error?.message || 'Registration failed. Please try again.');
    },
  });
};

export const useLogout = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: () => authService.logout(),
    onSuccess: async () => {
      authActions.logout();
      queryClient.clear();
      await router.navigate({ to: '/login' });
    },
    onError: () => {
      // Even if the server logout fails, clear local state
      authActions.logout();
      queryClient.clear();
      router.navigate({ to: '/login' });
    },
  });
};

export const useSetup2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.setup2FA(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useVerify2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof authService.verify2FA>[0]) =>
      authService.verify2FA(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useDisable2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof authService.disable2FA>[0]) =>
      authService.disable2FA(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};
