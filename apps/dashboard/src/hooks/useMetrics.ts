import { useQuery } from '@tanstack/react-query';
import { metricsService } from '#/services/metrics';

export const useHistoricalMetrics = (serviceId: string, range?: string) => {
  return useQuery({
    queryKey: ['historicalMetrics', serviceId, range],
    queryFn: () => metricsService.getHistoricalMetrics(serviceId, range),
    enabled: !!serviceId,
  });
};
