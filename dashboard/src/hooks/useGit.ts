import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { gitService } from '#/services/git';

export const useGetStatus = () => {
  return useQuery({
    queryKey: ['git', 'getStatus'].filter(Boolean),
    queryFn: () => gitService.getStatus(),
  });
};

export const useConnect = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof gitService.connect>[0] }) =>
      gitService.connect(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['git'] });
    },
  });
};

export const useDisconnect = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { provider: string }) => gitService.disconnect(payload.provider),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['git'] });
    },
  });
};

export const useListRepos = (provider: string) => {
  return useQuery({
    queryKey: ['git', 'listRepos', provider].filter(Boolean),
    queryFn: () => gitService.listRepos(provider),
  });
};
