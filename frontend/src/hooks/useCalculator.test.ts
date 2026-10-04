import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useCalculator } from './useCalculator';
import * as api from '../services/api';

describe('useCalculator Hook', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
  });

  it('initializes with default values', () => {
    const { result } = renderHook(() => useCalculator());
    expect(result.current.display).toBe('0');
    expect(result.current.equation).toBe('');
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(result.current.history).toEqual([]);
  });

  it('inputs digits correctly without leading zero issues', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('5');
    });
    expect(result.current.display).toBe('5');

    act(() => {
      result.current.inputDigit('8');
    });
    expect(result.current.display).toBe('58');
  });

  it('handles decimal input correctly and prevents duplicate dots', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('3');
      result.current.inputDecimal();
      result.current.inputDigit('1');
      result.current.inputDecimal(); // should be ignored
      result.current.inputDigit('4');
    });

    expect(result.current.display).toBe('3.14');
  });

  it('toggles sign correctly', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('9');
      result.current.toggleSign();
    });
    expect(result.current.display).toBe('-9');

    act(() => {
      result.current.toggleSign();
    });
    expect(result.current.display).toBe('9');
  });

  it('deletes digits correctly', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('1');
      result.current.inputDigit('2');
      result.current.inputDigit('3');
    });
    expect(result.current.display).toBe('123');

    act(() => {
      result.current.deleteDigit();
    });
    expect(result.current.display).toBe('12');
  });

  it('executes addition successfully via API and updates history', async () => {
    vi.spyOn(api, 'calculate').mockResolvedValueOnce({
      operation: 'add',
      a: 10,
      b: 25,
      result: 35,
      formatted: '10 + 25 = 35',
    });

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('1');
      result.current.inputDigit('0');
    });

    await act(async () => {
      result.current.setOperation('add');
    });

    act(() => {
      result.current.inputDigit('2');
      result.current.inputDigit('5');
    });

    await act(async () => {
      await result.current.executeEquals();
    });

    expect(result.current.display).toBe('35');
    expect(result.current.equation).toBe('10 + 25 =');
    expect(result.current.history.length).toBe(1);
    expect(result.current.history[0].result).toBe(35);
  });

  it('handles backend division by zero error', async () => {
    vi.spyOn(api, 'calculate').mockRejectedValueOnce(
      new api.CalculatorApiError('division by zero is undefined', 'DIVISION_BY_ZERO', 400)
    );

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('8');
    });

    await act(async () => {
      result.current.setOperation('divide');
    });

    act(() => {
      result.current.inputDigit('0');
    });

    await act(async () => {
      await result.current.executeEquals();
    });

    expect(result.current.error).toBe('division by zero is undefined');
  });

  it('executes unary square root operation', async () => {
    vi.spyOn(api, 'calculate').mockResolvedValueOnce({
      operation: 'sqrt',
      a: 64,
      result: 8,
      formatted: '√64 = 8',
    });

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('6');
      result.current.inputDigit('4');
    });

    await act(async () => {
      await result.current.executeUnary('sqrt');
    });

    expect(result.current.display).toBe('8');
    expect(result.current.equation).toBe('√(64) =');
    expect(result.current.history.length).toBe(1);
  });

  it('clears all state', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('9');
      result.current.setOperation('add');
      result.current.clear();
    });

    expect(result.current.display).toBe('0');
    expect(result.current.equation).toBe('');
    expect(result.current.error).toBeNull();
  });
});
