import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { settingsService } from '#/services/settings';

export const useGetSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getSettings'].filter(Boolean),
    queryFn: () => settingsService.getSettings(),
  });
};

export const useGetPublicSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getPublicSettings'],
    queryFn: () => settingsService.getPublicSettings(),
  });
};

export const useGetSetupStatus = () => {
  return useQuery({
    queryKey: ['settings', 'getSetupStatus'],
    queryFn: () => settingsService.getSetupStatus(),
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

export const useGetAISettings = () => {
  return useQuery({
    queryKey: ['settings', 'getAISettings'].filter(Boolean),
    queryFn: () => settingsService.getAISettings(),
  });
};

export const useUpdateAISettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) => settingsService.updateAISettings(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getAISettings'] });
    },
  });
};

export const useGetNotificationSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getNotificationSettings'].filter(Boolean),
    queryFn: () => settingsService.getNotificationSettings(),
  });
};

export const useUpdateNotificationSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) =>
      settingsService.updateNotificationSettings(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getNotificationSettings'] });
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
    mutationFn: ({ provider, payload }: { provider: string; payload: Record<string, unknown> }) =>
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
    mutationFn: (payload: Record<string, unknown>) =>
      settingsService.exchangeGithubManifest(payload),
  });
};

export const useGetUpdateStatus = () => {
  return useQuery({
    queryKey: ['settings', 'getUpdateStatus'].filter(Boolean),
    queryFn: () => settingsService.getUpdateStatus(),
  });
};

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

export const useSaveNotificationChannel = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) =>
      settingsService.saveNotificationChannel(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings', 'getNotifications'] });
    },
  });
};

export const useTestNotification = () => {
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) => settingsService.testNotification(payload),
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
