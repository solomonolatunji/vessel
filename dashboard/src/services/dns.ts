import type { BaseResponse } from '#/interfaces/base';
import type { CreateDNSRecordRequest, DNSRecord, UpdateDNSRecordRequest } from '#/interfaces/dns';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const dnsService = {
  list: async (): Promise<BaseResponse<DNSRecord[]>> => {
    try {
      return await apiClient.get<BaseResponse<DNSRecord[]>>(`/dns`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (payload: CreateDNSRecordRequest): Promise<BaseResponse<DNSRecord>> => {
    try {
      return await apiClient.post<BaseResponse<DNSRecord>>(`/dns`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  update: async (id: string, payload: UpdateDNSRecordRequest): Promise<BaseResponse<DNSRecord>> => {
    try {
      return await apiClient.put<BaseResponse<DNSRecord>>(`/dns/${id}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/dns/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
