import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { settingsService } from '#/services/settings';

export const useGetSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getSettings'].filter(Boolean),
    queryFn: () => settingsService.getSettings(),
  });
};

export const useUpdateSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof settingsService.updateSettings>[0] }) =>
      settingsService.updateSettings(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings'] });
    },
  });
};

export const useGetNotifications = () => {
  return useQuery({
    queryKey: ['settings', 'getNotifications'].filter(Boolean),
    queryFn: () => settingsService.getNotifications(),
  });
};

export const useGetGitApps = (provider: string) => {
  return useQuery({
    queryKey: ['settings', 'getGitApps', provider].filter(Boolean),
    queryFn: () => settingsService.getGitApps(provider),
  });
};

export const useSaveGitApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ provider, payload }: { provider: string; payload: unknown }) =>
      settingsService.saveGitApp(provider, payload),
    onSuccess: (_, { provider }) => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', provider] });
    },
  });
};

export const useDeleteGitApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ provider, id }: { provider: string; id: string }) =>
      settingsService.deleteGitApp(provider, id),
    onSuccess: (_, { provider }) => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', provider] });
    },
  });
};

export const useExchangeGithubManifest = () => {
  return useMutation({
    mutationFn: (payload: unknown) => settingsService.exchangeGithubManifest(payload),
  });
};

// --- UPDATES & SYSTEM ---

export const useCheckUpdate = () => {
  return useMutation({
    mutationFn: () => settingsService.checkUpdate(),
  });
};

export const useDeployUpdate = () => {
  return useMutation({
    mutationFn: () => settingsService.deployUpdate(),
  });
};

export const useActivateLicense = () => {
  return useMutation({
    mutationFn: (payload: unknown) => settingsService.activateLicense(payload),
  });
};

export const useSaveNotificationChannel = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: unknown) => settingsService.saveNotificationChannel(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getNotifications'] });
    },
  });
};

export const useTestNotification = () => {
  return useMutation({
    mutationFn: (payload: unknown) => settingsService.testNotification(payload),
  });
};

export const useDeleteNotificationChannel = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => settingsService.deleteNotificationChannel(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getNotifications'] });
    },
  });
};
