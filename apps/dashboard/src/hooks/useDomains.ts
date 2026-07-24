import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { domainsService } from '#/services/domains';

export const useListByService = (serviceId: string) => {
  return useQuery({
    queryKey: ['domains', 'listByService', serviceId].filter(Boolean),
    queryFn: () => domainsService.listByService(serviceId),
  });
};

export const useCreate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      payload: { domainName: string; redirectTo?: string; pathPrefix?: string };
    }) => domainsService.create(payload.serviceId, payload.payload),
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
