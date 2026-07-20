import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateServiceVarRequest,
  UpdateServiceVarRequest,
  Variable,
  EnvExampleVariableSuggestion,
} from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const serviceVarsService = {
  list: async (serviceId: string): Promise<BaseResponse<Variable[]>> => {
    try {
      return await apiClient.get<BaseResponse<Variable[]>>(`/services/${serviceId}/variables`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getEnvSuggestions: async (serviceId: string): Promise<BaseResponse<EnvExampleVariableSuggestion[]>> => {
    try {
      return await apiClient.get<BaseResponse<EnvExampleVariableSuggestion[]>>(`/services/${serviceId}/env-suggestions`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (
    serviceId: string,
    payload: CreateServiceVarRequest
  ): Promise<BaseResponse<Variable>> => {
    try {
      return await apiClient.post<BaseResponse<Variable>>(
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
    payload: UpdateServiceVarRequest
  ): Promise<BaseResponse<Variable>> => {
    try {
      return await apiClient.put<BaseResponse<Variable>>(
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
