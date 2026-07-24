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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings'] });
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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getAISettings'] });
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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getNotificationSettings'] });
    },
  });
};

export const useGetNotifications = () => {
  return useQuery({
    queryKey: ['settings', 'getNotifications'].filter(Boolean),
    queryFn: () => settingsService.getNotifications(),
  });
};

export const useGetGitApps = () => {
  return useQuery({
    queryKey: ['settings', 'getGitApps', 'github'],
    queryFn: () => settingsService.getGitApps(),
  });
};

export const useSaveGitApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) => settingsService.saveGitApp(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', 'github'] });
    },
  });
};

export const useDeleteGitApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => settingsService.deleteGitApp(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', 'github'] });
    },
  });
};

export const useExchangeGithubManifest = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) =>
      settingsService.exchangeGithubManifest(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', 'github'] });
    },
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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getNotifications'] });
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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getNotifications'] });
    },
  });
};
