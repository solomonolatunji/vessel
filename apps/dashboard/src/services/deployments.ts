import type {
  ExplainDeploymentResponse,
  GetDeploymentLogsResponse,
  GetDiagnosticsResponse,
  GetServiceMetricsResponse,
  ListDeploymentsResponse,
  ListPRPreviewsResponse,
  RollbackDeploymentResponse,
  TriggerDeploymentRequest,
  TriggerDeploymentResponse,
} from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const deploymentsService = {
  listByService: async (serviceId: string): Promise<ListDeploymentsResponse> => {
    try {
      return await apiClient.get<ListDeploymentsResponse>(`/services/${serviceId}/deployments`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listPRPreviews: async (serviceId: string): Promise<ListPRPreviewsResponse> => {
    try {
      return await apiClient.get<ListPRPreviewsResponse>(`/services/${serviceId}/previews`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  trigger: async (
    serviceId: string,
    payload?: TriggerDeploymentRequest
  ): Promise<TriggerDeploymentResponse> => {
    try {
      return await apiClient.post<TriggerDeploymentResponse>(
        `/services/${serviceId}/deploy`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  rollback: async (deploymentId: string): Promise<RollbackDeploymentResponse> => {
    try {
      return await apiClient.post<RollbackDeploymentResponse>(
        `/deployments/${deploymentId}/rollback`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getLogs: async (deploymentId: string): Promise<GetDeploymentLogsResponse> => {
    try {
      return await apiClient.get<GetDeploymentLogsResponse>(`/deployments/${deploymentId}/logs`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getMetrics: async (serviceId: string): Promise<GetServiceMetricsResponse> => {
    try {
      return await apiClient.get<GetServiceMetricsResponse>(`/services/${serviceId}/metrics`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  diagnostics: async (deploymentId: string): Promise<GetDiagnosticsResponse> => {
    try {
      return await apiClient.post<GetDiagnosticsResponse>(
        `/deployments/${deploymentId}/diagnostics`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  explainFailure: async (deploymentId: string): Promise<ExplainDeploymentResponse> => {
    try {
      return await apiClient.get<ExplainDeploymentResponse>(`/deployments/${deploymentId}/explain`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
