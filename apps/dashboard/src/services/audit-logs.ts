import type { AuditLog } from '#/interfaces/audit';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const auditLogsService = {
  list: async (): Promise<BaseResponse<AuditLog[]>> => {
    try {
      return await apiClient.get<BaseResponse<AuditLog[]>>(`/audit-logs`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
