import type { BaseResponse } from '#/interfaces/base';
import type { DomainConfig } from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const domainsService = {
  listByService: async (serviceId: string): Promise<BaseResponse<DomainConfig[]>> => {
    try {
      return await apiClient.get<BaseResponse<DomainConfig[]>>(`/services/${serviceId}/domains`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (
    serviceId: string,
    payload: { domainName: string; redirectTo?: string; pathPrefix?: string }
  ): Promise<BaseResponse<DomainConfig>> => {
    try {
      return await apiClient.post<BaseResponse<DomainConfig>>(
        `/services/${serviceId}/domains`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/domains/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
