/**
 * A lightweight fetch wrapper for interacting with the Vessl Go Daemon API.
 * Designed to be used seamlessly with TanStack Query.
 */

import { toast } from 'sonner';
import { env } from '#/env';
import { authActions, authStore } from '#/stores/authStore';

const API_BASE_URL = env.VITE_API_URL;

function getCookie(name: string) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift();
  return null;
}

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
  async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    const headers = new Headers(options.headers || {});
    if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
      headers.set('Content-Type', 'application/json');
    }

    const token = authStore.state.token;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }

    const csrfToken = getCookie('csrf_token');
    if (csrfToken) {
      headers.set('X-CSRF-Token', csrfToken);
    }

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: 'include',
    });

    if (response.status === 401) {
      authActions.logout();
      window.location.href = '/login';
      throw new ApiError(401, 'Session expired. Please log in again.');
    }

    if (response.status === 204) {
      return {} as T;
    }

    const isJson = response.headers.get('content-type')?.includes('application/json');
    const data = isJson ? await response.json() : await response.text();

    if (!response.ok) {
      const errorMessage =
        data?.message || data?.error || response.statusText || 'An error occurred';
      toast.error(errorMessage);
      throw new ApiError(response.status, errorMessage, data);
    }

    return data as T;
  },

  get<T>(endpoint: string, options?: RequestInit) {
    return this.fetch<T>(endpoint, { ...options, method: 'GET' });
  },

  async getBlob(endpoint: string, options?: RequestInit): Promise<Blob> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = new Headers(options?.headers || {});
    const token = authStore.state.token;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
    const csrfToken = getCookie('csrf_token');
    if (csrfToken) {
      headers.set('X-CSRF-Token', csrfToken);
    }
    const response = await fetch(url, {
      ...options,
      method: 'GET',
      headers,
      credentials: 'include',
    });

    if (response.status === 401) {
      authActions.logout();
      window.location.href = '/login';
      throw new ApiError(401, 'Session expired. Please log in again.');
    }

    if (!response.ok) {
      const isJson = response.headers.get('content-type')?.includes('application/json');
      const data = isJson ? await response.json() : await response.text();
      const errorMessage =
        data?.message || data?.error || response.statusText || 'An error occurred';
      toast.error(errorMessage);
      throw new ApiError(response.status, errorMessage, data);
    }
    return response.blob();
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
