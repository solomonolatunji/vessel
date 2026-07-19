import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateDNSRecordRequest, UpdateDNSRecordRequest } from '#/interfaces/dns';
import { dnsService } from '#/services/dns';

export const useListDNS = () => {
  return useQuery({
    queryKey: ['dns'],
    queryFn: () => dnsService.list(),
  });
};

export const useCreateDNS = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateDNSRecordRequest) => dnsService.create(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['dns'] });
    },
  });
};

export const useUpdateDNS = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateDNSRecordRequest }) =>
      dnsService.update(id, payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['dns'] });
    },
  });
};

export const useDeleteDNS = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => dnsService.delete(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['dns'] });
    },
  });
};
