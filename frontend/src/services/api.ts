import type { CalculateRequest, CalculateResponse, ApiErrorResponse } from '../types/calculator';

// Use environment variable if provided, fallback to relative path (handled by Vite proxy) or localhost
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export class CalculatorApiError extends Error {
  code?: string;
  statusCode?: number;

  constructor(message: string, code?: string, statusCode?: number) {
    super(message);
    this.name = 'CalculatorApiError';
    this.code = code;
    this.statusCode = statusCode;
  }
}

/**
 * Sends a calculation request to the Go backend API.
 */
export async function calculate(payload: CalculateRequest): Promise<CalculateResponse> {
  const url = `${BASE_URL}/calculate`;

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
      body: JSON.stringify(payload),
    });

    const data = await response.json();

    if (!response.ok) {
      const errorData = data as ApiErrorResponse;
      throw new CalculatorApiError(
        errorData.error || 'Calculation failed on server',
        errorData.code,
        response.status
      );
    }

    return data as CalculateResponse;
  } catch (err: unknown) {
    if (err instanceof CalculatorApiError) {
      throw err;
    }

    if (err instanceof Error) {
      if (err.message.includes('Failed to fetch') || err.message.includes('NetworkError')) {
        throw new CalculatorApiError(
          'Cannot reach backend server. Please verify Go service is running on port 8080.',
          'NETWORK_ERROR'
        );
      }
      throw new CalculatorApiError(err.message, 'CLIENT_ERROR');
    }

    throw new CalculatorApiError('An unknown error occurred', 'UNKNOWN_ERROR');
  }
}

/**
 * Checks backend health status.
 */
export async function checkHealth(): Promise<{ status: string; service: string }> {
  const response = await fetch(`${BASE_URL}/health`);
  if (!response.ok) {
    throw new Error('Backend is unavailable');
  }
  return response.json();
}
