import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { oauthService } from '#/services/oauth';

export const oauthKeys = {
  all: ['oauth'] as const,
  enabled: ['oauth', 'enabled'] as const,
};

export const useListOAuthProviders = () => {
  return useQuery({
    queryKey: oauthKeys.all,
    queryFn: () => oauthService.listProviders(),
  });
};

export const useEnabledOAuthProviders = () => {
  return useQuery({
    queryKey: oauthKeys.enabled,
    queryFn: () => oauthService.getEnabledProviders(),
  });
};

export const useSaveOAuthProvider = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof oauthService.saveProvider>[0] }) =>
      oauthService.saveProvider(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: oauthKeys.all });
    },
  });
};
