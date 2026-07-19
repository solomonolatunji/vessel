import { beforeEach, describe, expect, it, vi } from 'vitest';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';
import { serviceVariablesService } from './service-variables';

vi.mock('#/lib/apiClient', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

vi.mock('#/lib/error', () => ({
  handleApiError: vi.fn((_err) => new Error('API Error')),
}));

describe('serviceVariablesService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('list calls GET /services/:serviceId/variables', async () => {
    const mockResponse = { data: [] };
    vi.mocked(apiClient.get).mockResolvedValueOnce(mockResponse);

    const result = await serviceVariablesService.list('1');

    expect(apiClient.get).toHaveBeenCalledWith('/services/1/variables');
    expect(result).toEqual(mockResponse);
  });

  it('create calls POST /services/:serviceId/variables', async () => {
    const payload = { key: 'TEST', value: '123' };
    const mockResponse = { data: payload };
    vi.mocked(apiClient.post).mockResolvedValueOnce(mockResponse);

    const result = await serviceVariablesService.create('1', payload);

    expect(apiClient.post).toHaveBeenCalledWith('/services/1/variables', payload);
    expect(result).toEqual(mockResponse);
  });

  it('update calls PUT /services/:serviceId/variables/:id', async () => {
    const payload = { key: 'TEST', value: '456' };
    const mockResponse = { data: payload };
    vi.mocked(apiClient.put).mockResolvedValueOnce(mockResponse);

    const result = await serviceVariablesService.update('1', 'var1', payload);

    expect(apiClient.put).toHaveBeenCalledWith('/services/1/variables/var1', payload);
    expect(result).toEqual(mockResponse);
  });

  it('delete calls DELETE /services/:serviceId/variables/:id', async () => {
    vi.mocked(apiClient.delete).mockResolvedValueOnce(undefined);

    await serviceVariablesService.delete('1', 'var1');

    expect(apiClient.delete).toHaveBeenCalledWith('/services/1/variables/var1');
  });

  it('handles errors correctly', async () => {
    vi.mocked(apiClient.get).mockRejectedValueOnce(new Error('Network error'));

    await expect(serviceVariablesService.list('1')).rejects.toThrow('API Error');
    expect(handleApiError).toHaveBeenCalled();
  });
});
