import { useQuery } from '@tanstack/react-query';
import { auditLogsService } from '#/services/audit-logs';

export const useAuditLogs = () => {
  return useQuery({
    queryKey: ['auditLogs'],
    queryFn: () => auditLogsService.list(),
  });
};
