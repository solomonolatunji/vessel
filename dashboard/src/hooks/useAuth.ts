import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { toast } from 'sonner';
import { authService } from '#/services/auth';
import { authActions } from '#/stores/authStore';

export const useLogin = () => {
  const queryClient = useQueryClient();
  const router = useRouter();
  return useMutation({
    mutationFn: (payload: { credentials: Parameters<typeof authService.login>[0] }) =>
      authService.login(payload.credentials),
    onSuccess: async (data) => {
      if (data?.token && data?.user) {
        authActions.setAuth(data.token, data.user);
      }
      await queryClient.invalidateQueries({ queryKey: ['auth'] });
      await router.invalidate();

      toast.success('Logged in successfully');
      router.navigate({ to: '/' });
    },
  });
};

export const useRegister = () => {
  const queryClient = useQueryClient();
  const router = useRouter();
  return useMutation({
    mutationFn: (payload: { details: Parameters<typeof authService.register>[0] }) =>
      authService.register(payload.details),
    onSuccess: async (data) => {
      if (data?.token && data?.user) {
        authActions.setAuth(data.token, data.user);
      }
      await queryClient.invalidateQueries({ queryKey: ['auth'] });
      await router.invalidate();

      toast.success('Account created successfully');
      router.navigate({ to: '/' });
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
      await router.invalidate();
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
    mutationFn: (payload: { payload: Parameters<typeof authService.verify2FA>[0] }) =>
      authService.verify2FA(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useDisable2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof authService.disable2FA>[0] }) =>
      authService.disable2FA(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};
