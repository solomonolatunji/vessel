import type { BaseResponse } from '#/interfaces/base';
import type { GetHistoricalMetricsResponse } from '#/interfaces/metrics';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const metricsService = {
  getHistoricalMetrics: async (
    serviceId: string,
    range?: string
  ): Promise<BaseResponse<GetHistoricalMetricsResponse>> => {
    try {
      const query = new URLSearchParams();
      if (range) query.set('range', range);

      const queryString = query.toString() ? `?${query.toString()}` : '';
      return await apiClient.get<BaseResponse<GetHistoricalMetricsResponse>>(
        `/services/${serviceId}/metrics/historical${queryString}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
