import { describe, it, expect, vi, beforeEach } from 'vitest';
import { calculate, checkHealth } from './api';

describe('API Service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('performs calculate request successfully', async () => {
    const mockResponse = {
      operation: 'add',
      a: 10,
      b: 5,
      result: 15,
      formatted: '10 + 5 = 15',
    };

    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => mockResponse,
    });

    const res = await calculate({ operation: 'add', a: 10, b: 5 });
    expect(res).toEqual(mockResponse);
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/calculate'),
      expect.objectContaining({
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
        body: JSON.stringify({ operation: 'add', a: 10, b: 5 }),
      })
    );
  });

  it('handles backend domain error (e.g. division by zero)', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({
        success: false,
        error: 'division by zero is undefined',
        code: 'DIVISION_BY_ZERO',
      }),
    });

    await expect(
      calculate({ operation: 'divide', a: 10, b: 0 })
    ).rejects.toThrow('division by zero is undefined');
  });

  it('handles network connection error', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('Failed to fetch'));

    await expect(
      calculate({ operation: 'add', a: 2, b: 2 })
    ).rejects.toThrow('Cannot reach backend server');
  });

  it('verifies checkHealth endpoint', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'healthy', service: 'calculator-api' }),
    });

    const res = await checkHealth();
    expect(res.status).toBe('healthy');
  });
});
