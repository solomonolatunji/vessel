/**
 * A lightweight fetch wrapper for interacting with the Codedock Go Daemon API.
 * Designed to be used seamlessly with TanStack Query.
 */

import { toast } from 'sonner';
import { env } from '#/env';
import { useAuthStore } from '#/stores/authStore';

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

    const token = useAuthStore.getState().token;
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
      const authState = useAuthStore.getState();
      if (authState.refreshToken) {
        try {
          const refreshRes = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refreshToken: authState.refreshToken }),
          });
          if (refreshRes.ok) {
            const data = await refreshRes.json();
            if (data.data?.token && data.data?.user) {
              authState.setAuth(data.data.token, data.data.refreshToken || null, data.data.user);
              // Retry the original request
              headers.set('Authorization', `Bearer ${data.data.token}`);
              const retryResponse = await fetch(url, {
                ...options,
                headers,
                credentials: 'include',
              });
              if (retryResponse.ok) {
                const isJson = retryResponse.headers
                  .get('content-type')
                  ?.includes('application/json');
                return (isJson ? await retryResponse.json() : await retryResponse.text()) as T;
              }
            }
          }
        } catch (e) {
          console.error('Failed to refresh token', e);
        }
      }

      useAuthStore.getState().logout();
      if (
        !window.location.pathname.startsWith('/signin') &&
        !window.location.pathname.startsWith('/signup') &&
        !window.location.pathname.startsWith('/forgot-password') &&
        !window.location.pathname.startsWith('/reset-password') &&
        !window.location.pathname.startsWith('/onboarding')
      ) {
        window.location.href = '/signin';
      }
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
    const token = useAuthStore.getState().token;
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
      const authState = useAuthStore.getState();
      if (authState.refreshToken) {
        try {
          const refreshRes = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refreshToken: authState.refreshToken }),
          });
          if (refreshRes.ok) {
            const data = await refreshRes.json();
            if (data.data?.token && data.data?.user) {
              authState.setAuth(data.data.token, data.data.refreshToken || null, data.data.user);
              headers.set('Authorization', `Bearer ${data.data.token}`);
              const retryResponse = await fetch(url, {
                ...options,
                method: 'GET',
                headers,
                credentials: 'include',
              });
              if (retryResponse.ok) {
                return retryResponse.blob();
              }
            }
          }
        } catch (e) {
          console.error('Failed to refresh token', e);
        }
      }

      useAuthStore.getState().logout();
      if (
        !window.location.pathname.startsWith('/signin') &&
        !window.location.pathname.startsWith('/signup') &&
        !window.location.pathname.startsWith('/forgot-password') &&
        !window.location.pathname.startsWith('/reset-password') &&
        !window.location.pathname.startsWith('/setup')
      ) {
        window.location.href = '/signin';
      }
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

  async postBlob(endpoint: string, body?: unknown, options?: RequestInit): Promise<Blob> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = new Headers(options?.headers || {});
    if (!headers.has('Content-Type') && !(body instanceof FormData)) {
      headers.set('Content-Type', 'application/json');
    }
    const token = useAuthStore.getState().token;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
    const csrfToken = getCookie('csrf_token');
    if (csrfToken) {
      headers.set('X-CSRF-Token', csrfToken);
    }
    const response = await fetch(url, {
      ...options,
      method: 'POST',
      body: body instanceof FormData ? body : JSON.stringify(body),
      headers,
      credentials: 'include',
    });

    if (response.status === 401) {
      const authState = useAuthStore.getState();
      if (authState.refreshToken) {
        try {
          const refreshRes = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refreshToken: authState.refreshToken }),
          });
          if (refreshRes.ok) {
            const data = await refreshRes.json();
            if (data.data?.token && data.data?.user) {
              authState.setAuth(data.data.token, data.data.refreshToken || null, data.data.user);
              headers.set('Authorization', `Bearer ${data.data.token}`);
              const retryResponse = await fetch(url, {
                ...options,
                method: 'POST',
                body: body instanceof FormData ? body : JSON.stringify(body),
                headers,
                credentials: 'include',
              });
              if (retryResponse.ok) {
                return retryResponse.blob();
              }
            }
          }
        } catch (e) {
          console.error('Failed to refresh token', e);
        }
      }

      useAuthStore.getState().logout();
      if (
        !window.location.pathname.startsWith('/signin') &&
        !window.location.pathname.startsWith('/signup') &&
        !window.location.pathname.startsWith('/forgot-password') &&
        !window.location.pathname.startsWith('/reset-password') &&
        !window.location.pathname.startsWith('/setup')
      ) {
        window.location.href = '/signin';
      }
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
