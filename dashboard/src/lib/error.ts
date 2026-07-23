import { ApiError } from '#/lib/apiClient';

export class CodedockError extends Error {
  public status: number;
  public data: unknown;

  constructor(message: string, status: number, data?: unknown) {
    super(message);
    this.name = 'CodedockError';
    this.status = status;
    this.data = data;
  }
}

export const handleApiError = (error: unknown): CodedockError => {
  if (error instanceof ApiError) {
    let message: string;

    switch (error.status) {
      case 400:
        message = error.message || 'Bad request';
        break;
      case 401:
        message = error.message || 'Unauthorized access';
        break;
      case 403:
        message = error.message || 'Forbidden action';
        break;
      case 404:
        message = error.message || 'Resource not found';
        break;
      case 422:
        message = error.message || 'Validation error';
        break;
      case 429:
        message = error.message || 'Too many requests. Please try again later';
        break;
      case 500:
        message = error.message || 'Internal server error';
        break;
      default:
        message = error.message || 'An unexpected error occurred';
    }

    return new CodedockError(message, error.status, error.data);
  }

  if (error instanceof Error) {
    return new CodedockError(error.message, 0);
  }

  return new CodedockError('An unexpected error occurred', 0);
};

export const extractErrorMessage = (
  error: unknown,
  fallback = 'An unexpected error occurred'
): string => {
  if (typeof error === 'string') return error;
  if (error instanceof Error) return error.message;
  if (
    error &&
    typeof error === 'object' &&
    'message' in error &&
    typeof (error as Record<string, unknown>).message === 'string'
  ) {
    return (error as Record<string, unknown>).message as string;
  }
  return fallback;
};
