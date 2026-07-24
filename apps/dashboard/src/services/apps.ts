import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateAppResponse,
  CreateAppServiceRequest,
  CreateServiceVarRequest,
  GetAppResponse,
  ListAppsResponse,
  ListVariablesResponse,
  UpdateAppResponse,
  UpdateAppServiceRequest,
  UpdateServiceVarRequest,
  Variable,
} from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const appsService = {
  listByProject: async (projectId: string): Promise<ListAppsResponse> => {
    try {
      return await apiClient.get<ListAppsResponse>(`/projects/${projectId}/apps`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listByEnvironment: async (environmentId: string): Promise<ListAppsResponse> => {
    try {
      return await apiClient.get<ListAppsResponse>(`/environments/${environmentId}/apps`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getApp: async (appId: string): Promise<GetAppResponse> => {
    try {
      return await apiClient.get<GetAppResponse>(`/apps/${appId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createApp: async (
    environmentId: string,
    payload: CreateAppServiceRequest
  ): Promise<CreateAppResponse> => {
    try {
      return await apiClient.post<CreateAppResponse>(
        `/environments/${environmentId}/apps`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateApp: async (
    appId: string,
    payload: UpdateAppServiceRequest
  ): Promise<UpdateAppResponse> => {
    try {
      return await apiClient.put<UpdateAppResponse>(`/apps/${appId}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteApp: async (appId: string): Promise<void> => {
    try {
      await apiClient.delete(`/apps/${appId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  stopApp: async (appId: string): Promise<void> => {
    try {
      await apiClient.post(`/apps/${appId}/stop`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  redeployApp: async (appId: string): Promise<void> => {
    try {
      await apiClient.post(`/apps/${appId}/redeploy`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  restartApp: async (appId: string): Promise<void> => {
    try {
      await apiClient.post(`/apps/${appId}/restart`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listVariables: async (appId: string): Promise<ListVariablesResponse> => {
    try {
      return await apiClient.get<ListVariablesResponse>(`/services/${appId}/variables`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createVariable: async (
    appId: string,
    payload: CreateServiceVarRequest
  ): Promise<BaseResponse<Variable>> => {
    try {
      return await apiClient.post<BaseResponse<Variable>>(`/services/${appId}/variables`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateVariable: async (
    appId: string,
    varId: string,
    payload: UpdateServiceVarRequest
  ): Promise<BaseResponse<Variable>> => {
    try {
      return await apiClient.put<BaseResponse<Variable>>(
        `/services/${appId}/variables/${varId}`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteVariable: async (appId: string, varId: string): Promise<void> => {
    try {
      await apiClient.delete(`/services/${appId}/variables/${varId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listLogDrains: async (appId: string): Promise<any[]> => {
    try {
      const response = await apiClient.get<any>(`/apps/${appId}/log-drains`);
      return response.data?.data || [];
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createLogDrain: async (appId: string, data: any): Promise<any> => {
    try {
      const response = await apiClient.post<any>(`/apps/${appId}/log-drains`, data);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteLogDrain: async (appId: string, drainId: string): Promise<void> => {
    try {
      await apiClient.delete(`/apps/${appId}/log-drains/${drainId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
