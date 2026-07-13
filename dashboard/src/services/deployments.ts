import type { BaseResponse, PaginatedData } from '#/interfaces/base';
import type { Deployment, ServiceMetric, TriggerDeploymentRequest } from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const deploymentsService = {
  listByService: async (serviceId: string): Promise<BaseResponse<PaginatedData<Deployment>>> => {
    try {
      return await apiClient.get<BaseResponse<PaginatedData<Deployment>>>(
        `/services/${serviceId}/deployments`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  trigger: async (
    serviceId: string,
    payload?: TriggerDeploymentRequest
  ): Promise<BaseResponse<Deployment>> => {
    try {
      return await apiClient.post<BaseResponse<Deployment>>(
        `/services/${serviceId}/deploy`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  rollback: async (deploymentId: string): Promise<BaseResponse<Deployment>> => {
    try {
      return await apiClient.post<BaseResponse<Deployment>>(
        `/deployments/${deploymentId}/rollback`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getLogs: async (deploymentId: string): Promise<BaseResponse<string>> => {
    try {
      return await apiClient.get<BaseResponse<string>>(`/deployments/${deploymentId}/logs`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getMetrics: async (serviceId: string): Promise<BaseResponse<ServiceMetric>> => {
    try {
      return await apiClient.get<BaseResponse<ServiceMetric>>(`/services/${serviceId}/metrics`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // POST /deployments/:id/diagnostics
  // TODO: Update the return type once we know the exact shape of the AI Diagnostics response
  diagnostics: async (deploymentId: string): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>(
        `/deployments/${deploymentId}/diagnostics`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
