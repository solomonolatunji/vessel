import { useQuery } from '@tanstack/react-query';
import { logsService } from '#/services/logs';

export const useHistoricalLogs = (serviceId: string, range?: string, limit?: number) => {
  return useQuery({
    queryKey: ['historicalLogs', serviceId, range, limit],
    queryFn: () => logsService.getHistoricalLogs(serviceId, range, limit),
    enabled: !!serviceId,
  });
};
