import { useQuery } from '@tanstack/react-query';
import { deploymentsService } from '#/services';

export const useListPRPreviews = (serviceId: string) => {
  return useQuery({
    queryKey: ['pr-previews', serviceId],
    queryFn: () => deploymentsService.listPRPreviews(serviceId),
    enabled: !!serviceId,
  });
};
