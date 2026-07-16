import type { BaseResponse } from '#/interfaces/base';
import type { GetHistoricalLogsResponse } from '#/interfaces/logs';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const logsService = {
  getHistoricalLogs: async (
    serviceId: string,
    range?: string,
    limit?: number
  ): Promise<BaseResponse<GetHistoricalLogsResponse>> => {
    try {
      const query = new URLSearchParams();
      if (range) query.set('range', range);
      if (limit) query.set('limit', limit.toString());

      const queryString = query.toString() ? `?${query.toString()}` : '';
      return await apiClient.get<BaseResponse<GetHistoricalLogsResponse>>(
        `/services/${serviceId}/logs/historical${queryString}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
