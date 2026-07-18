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

      authActions.setAuth(data.token, data.user);

      queryClient.clear();

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
      await router.navigate({ to: '/' });

      toast.success('Account created! Welcome to Vessl.');
    },
    onError: (error: Error) => {
      toast.error(error?.message || 'Registration failed. Please try again.');
    },
  });
};

export const useSetup = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: (details: Parameters<typeof authService.setup>[0]) => authService.setup(details),
    onSuccess: async () => {
      queryClient.clear();
      toast.success('Setup complete! Please sign in to continue.');
      await router.navigate({ to: '/signin' });
    },
    onError: (error: Error) => {
      toast.error(error?.message || 'Setup failed. Please try again.');
    },
  });
};

export const useForgotPassword = () => {
  return useMutation({
    mutationFn: (email: string) => authService.forgotPassword(email),
    onSuccess: async (data) => {
      toast.success(
        data?.message ||
          'If an account with that email exists, a password reset link has been sent.'
      );
    },
    onError: (error: Error) => {
      toast.error(error?.message || 'Failed to request password reset. Please try again.');
    },
  });
};

export const useResetPassword = () => {
  const router = useRouter();
  return useMutation({
    mutationFn: ({ token, newPassword }: { token: string; newPassword: string }) =>
      authService.resetPassword(token, newPassword),
    onSuccess: async () => {
      toast.success('Password successfully reset. You can now sign in.');
      await router.navigate({ to: '/signin' });
    },
    onError: (error: Error) => {
      toast.error(error?.message || 'Failed to reset password. The link might be expired.');
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
      await router.navigate({ to: '/signin' });
    },
    onError: () => {
      authActions.logout();
      queryClient.clear();
      router.navigate({ to: '/signin' });
    },
  });
};

export const useSetup2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.setup2FA(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useVerify2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof authService.verify2FA>[0]) =>
      authService.verify2FA(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useDisable2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof authService.disable2FA>[0]) =>
      authService.disable2FA(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};
