import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const serviceVariablesService = {
  list: async (serviceId: string): Promise<BaseResponse<unknown[]>> => {
    try {
      return await apiClient.get<BaseResponse<unknown[]>>(`/services/${serviceId}/variables`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (serviceId: string, payload: unknown): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>(
        `/services/${serviceId}/variables`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  update: async (
    serviceId: string,
    id: string,
    payload: unknown
  ): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.put<BaseResponse<unknown>>(
        `/services/${serviceId}/variables/${id}`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (serviceId: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/services/${serviceId}/variables/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
