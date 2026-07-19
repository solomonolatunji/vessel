import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { domainsService } from '#/services/domains';

export const useListByProject = (projectId: string) => {
  return useQuery({
    queryKey: ['domains', 'listByProject', projectId].filter(Boolean),
    queryFn: () => domainsService.listByProject(projectId),
  });
};

export const useCreate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      projectId: string;
      payload: { domainName: string; redirectTo?: string; pathPrefix?: string };
    }) => domainsService.create(payload.projectId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['domains'] });
    },
  });
};

export const useDelete = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => domainsService.delete(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['domains'] });
    },
  });
};
