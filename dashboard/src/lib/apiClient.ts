import { env } from '#/env';
import { authStore } from '#/stores/authStore';

const API_BASE_URL = env.VITE_API_URL;

export class ApiError extends Error {
  public status: number;
  public data: unknown;

  constructor(status: number, message: string, data?: unknown) {
    super(message);
    this.status = status;
    this.data = data;
    this.name = 'ApiError';
  }
}

export const apiClient = {
  /**
   * Executes an HTTP request to the Vessl API.
   * If VITE_IS_CLOUD is enabled and the endpoint is not a cloud-native route,
   * the request is automatically proxied through the active server tunnel.
   */
  async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const isCloud = env.VITE_IS_CLOUD;
    let rewrittenEndpoint = endpoint;

    if (isCloud) {
      const isCloudNative =
        endpoint.startsWith('/auth/') ||
        endpoint.startsWith('/system/') ||
        endpoint.startsWith('/billing/') ||
        endpoint.startsWith('/teams/') ||
        endpoint.startsWith('/users/') ||
        endpoint.startsWith('/servers/') ||
        endpoint.startsWith('/cloud/'); // fallback just in case

      if (!isCloudNative) {
        const activeServerId = localStorage.getItem('vessl_active_server_id');
        if (activeServerId) {
          rewrittenEndpoint = `/servers/${activeServerId}/proxy/api${endpoint}`;
        }
      }
    }

    const url = `${API_BASE_URL}${rewrittenEndpoint}`;

    const headers = new Headers(options.headers || {});
    if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
      headers.set('Content-Type', 'application/json');
    }

    const token = authStore.state.token;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (response.status === 204) {
      return {} as T;
    }

    const isJson = response.headers.get('content-type')?.includes('application/json');
    const data = isJson ? await response.json() : await response.text();

    if (!response.ok) {
      throw new ApiError(
        response.status,
        data?.error || response.statusText || 'An error occurred',
        data
      );
    }

    return data as T;
  },

  get<T>(endpoint: string, options?: RequestInit) {
    return this.fetch<T>(endpoint, { ...options, method: 'GET' });
  },

  post<T>(endpoint: string, body?: unknown, options?: RequestInit) {
    return this.fetch<T>(endpoint, {
      ...options,
      method: 'POST',
      body: body instanceof FormData ? body : JSON.stringify(body),
    });
  },

  put<T>(endpoint: string, body?: unknown, options?: RequestInit) {
    return this.fetch<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: body instanceof FormData ? body : JSON.stringify(body),
    });
  },

  delete<T>(endpoint: string, options?: RequestInit) {
    return this.fetch<T>(endpoint, { ...options, method: 'DELETE' });
  },
};
