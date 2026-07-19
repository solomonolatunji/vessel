import { beforeEach, describe, expect, it, vi } from 'vitest';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';
import { projectsService } from './projects';

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

describe('projectsService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('listProjects calls GET /projects', async () => {
    const mockResponse = { data: [] };
    vi.mocked(apiClient.get).mockResolvedValueOnce(mockResponse);

    const result = await projectsService.listProjects();

    expect(apiClient.get).toHaveBeenCalledWith('/projects');
    expect(result).toEqual(mockResponse);
  });

  it('getProject calls GET /projects/:id', async () => {
    const mockResponse = { data: { id: '1' } };
    vi.mocked(apiClient.get).mockResolvedValueOnce(mockResponse);

    const result = await projectsService.getProject('1');

    expect(apiClient.get).toHaveBeenCalledWith('/projects/1');
    expect(result).toEqual(mockResponse);
  });

  it('createProject calls POST /projects', async () => {
    const payload = { name: 'test' };
    const mockResponse = { data: { id: '1', name: 'test' } };
    vi.mocked(apiClient.post).mockResolvedValueOnce(mockResponse);

    const result = await projectsService.createProject(payload);

    expect(apiClient.post).toHaveBeenCalledWith('/projects', payload);
    expect(result).toEqual(mockResponse);
  });

  it('deleteProject calls DELETE /projects/:id', async () => {
    vi.mocked(apiClient.delete).mockResolvedValueOnce(undefined);

    await projectsService.deleteProject('1');

    expect(apiClient.delete).toHaveBeenCalledWith('/projects/1');
  });

  it('getVars calls GET /projects/:id/env', async () => {
    const mockResponse = { data: { KEY: 'VALUE' } };
    vi.mocked(apiClient.get).mockResolvedValueOnce(mockResponse);

    const result = await projectsService.getVars('1');

    expect(apiClient.get).toHaveBeenCalledWith('/projects/1/env');
    expect(result).toEqual(mockResponse);
  });

  it('setVars calls PUT /projects/:id/env', async () => {
    const payload = { variables: { KEY: 'NEW_VALUE' } };
    const mockResponse = { data: null };
    vi.mocked(apiClient.put).mockResolvedValueOnce(mockResponse);

    const result = await projectsService.setVars('1', payload);

    expect(apiClient.put).toHaveBeenCalledWith('/projects/1/env', payload);
    expect(result).toEqual(mockResponse);
  });

  it('handles errors correctly', async () => {
    vi.mocked(apiClient.get).mockRejectedValueOnce(new Error('Network error'));

    await expect(projectsService.listProjects()).rejects.toThrow('API Error');
    expect(handleApiError).toHaveBeenCalled();
  });
});
