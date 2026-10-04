import { useState, useCallback, useEffect, useRef } from 'react';
import { calculate, CalculatorApiError } from '../services/api';
import type { HistoryItem, OperationType } from '../types/calculator';

interface UseCalculatorReturn {
  display: string;
  equation: string;
  loading: boolean;
  error: string | null;
  history: HistoryItem[];
  pendingOp: OperationType | null;
  inputDigit: (digit: string) => void;
  inputDecimal: () => void;
  setOperation: (op: OperationType) => void;
  executeEquals: () => Promise<void>;
  executeUnary: (op: 'sqrt' | 'percentage') => Promise<void>;
  toggleSign: () => void;
  deleteDigit: () => void;
  clear: () => void;
  clearHistory: () => void;
  loadHistoryItem: (item: HistoryItem) => void;
}

export function useCalculator(): UseCalculatorReturn {
  const [display, setDisplay] = useState<string>('0');
  const [equation, setEquation] = useState<string>('');
  const [pendingValue, setPendingValue] = useState<number | null>(null);
  const [pendingOp, setPendingOp] = useState<OperationType | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Synchronous ref for new-input tracking to avoid stale-closure issues
  const isNewInputRef = useRef<boolean>(true);

  const [history, setHistory] = useState<HistoryItem[]>(() => {
    try {
      const saved = localStorage.getItem('calculator_history');
      if (saved) {
        const parsed = JSON.parse(saved);
        return parsed.map((item: any) => ({
          ...item,
          timestamp: new Date(item.timestamp),
        }));
      }
    } catch {
      // Ignore local storage parse error
    }
    return [];
  });

  // Save history to localStorage
  useEffect(() => {
    try {
      localStorage.setItem('calculator_history', JSON.stringify(history));
    } catch {
      // Ignore storage errors
    }
  }, [history]);

  const clear = useCallback(() => {
    setDisplay('0');
    setEquation('');
    setPendingValue(null);
    setPendingOp(null);
    isNewInputRef.current = true;
    setError(null);
  }, []);

  const clearHistory = useCallback(() => {
    setHistory([]);
    try {
      localStorage.removeItem('calculator_history');
    } catch {
      // Ignore storage errors
    }
  }, []);

  const inputDigit = useCallback((digit: string) => {
    setError(null);
    if (isNewInputRef.current) {
      isNewInputRef.current = false;
      setDisplay(digit);
      return;
    }
    setDisplay((prev) => {
      if (prev === '0') {
        return digit;
      }
      if (prev.replace('-', '').replace('.', '').length >= 14) {
        return prev;
      }
      return prev + digit;
    });
  }, []);

  const inputDecimal = useCallback(() => {
    setError(null);
    if (isNewInputRef.current) {
      isNewInputRef.current = false;
      setDisplay('0.');
      return;
    }
    setDisplay((prev) => {
      if (!prev.includes('.')) {
        return prev + '.';
      }
      return prev;
    });
  }, []);

  const toggleSign = useCallback(() => {
    setError(null);
    setDisplay((prev) => {
      if (prev === '0') return '0';
      if (prev.startsWith('-')) {
        return prev.substring(1);
      }
      return '-' + prev;
    });
  }, []);

  const deleteDigit = useCallback(() => {
    setError(null);
    if (isNewInputRef.current) return;

    setDisplay((prev) => {
      if (prev.length === 1 || (prev.length === 2 && prev.startsWith('-'))) {
        isNewInputRef.current = true;
        return '0';
      }
      return prev.slice(0, -1);
    });
  }, []);

  const getSymbol = (op: OperationType): string => {
    switch (op) {
      case 'add': return '+';
      case 'subtract': return '−';
      case 'multiply': return '×';
      case 'divide': return '÷';
      case 'power': return '^';
      default: return op;
    }
  };

  const setOperation = useCallback(async (op: OperationType) => {
    setError(null);
    const currentValue = parseFloat(display);

    if (isNaN(currentValue)) {
      setError('Invalid number format');
      return;
    }

    // If operator changed before entering a new number, simply change the operator
    if (pendingValue !== null && pendingOp && isNewInputRef.current) {
      setPendingOp(op);
      setEquation(`${pendingValue} ${getSymbol(op)}`);
      return;
    }

    // Chaining operations (e.g. 5 + 5 + -> evaluates first and keeps next pending)
    if (pendingValue !== null && pendingOp && !isNewInputRef.current) {
      setLoading(true);
      try {
        const resp = await calculate({
          operation: pendingOp,
          a: pendingValue,
          b: currentValue,
        });

        const newHistoryItem: HistoryItem = {
          id: Math.random().toString(36).substring(2, 9),
          expression: resp.formatted,
          result: resp.result,
          timestamp: new Date(),
        };
        setHistory((prev) => [newHistoryItem, ...prev.slice(0, 19)]);

        setPendingValue(resp.result);
        setPendingOp(op);
        setEquation(`${resp.result} ${getSymbol(op)}`);
        setDisplay('0');
        isNewInputRef.current = true;
      } catch (err: unknown) {
        if (err instanceof CalculatorApiError) {
          setError(err.message);
        } else {
          setError('Calculation error');
        }
      } finally {
        setLoading(false);
      }
      return;
    }

    setPendingValue(currentValue);
    setPendingOp(op);
    setEquation(`${currentValue} ${getSymbol(op)}`);
    setDisplay('0');
    isNewInputRef.current = true;
  }, [display, pendingOp, pendingValue]);

  const executeEquals = useCallback(async () => {
    if (pendingValue === null || pendingOp === null) {
      return;
    }

    const currentValue = parseFloat(display);
    if (isNaN(currentValue)) {
      setError('Invalid number');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const resp = await calculate({
        operation: pendingOp,
        a: pendingValue,
        b: currentValue,
      });

      const fullExpression = `${pendingValue} ${getSymbol(pendingOp)} ${currentValue} =`;
      setEquation(fullExpression);
      setDisplay(String(resp.result));

      const newHistoryItem: HistoryItem = {
        id: Math.random().toString(36).substring(2, 9),
        expression: fullExpression,
        result: resp.result,
        timestamp: new Date(),
      };
      setHistory((prev) => [newHistoryItem, ...prev.slice(0, 19)]);

      setPendingValue(null);
      setPendingOp(null);
      isNewInputRef.current = true;
    } catch (err: unknown) {
      if (err instanceof CalculatorApiError) {
        setError(err.message);
      } else {
        setError('Error calculating result');
      }
    } finally {
      setLoading(false);
    }
  }, [display, pendingOp, pendingValue]);

  const executeUnary = useCallback(async (op: 'sqrt' | 'percentage') => {
    const currentValue = parseFloat(display);
    if (isNaN(currentValue)) {
      setError('Invalid number');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const resp = await calculate({
        operation: op,
        a: currentValue,
      });

      const formattedEq = op === 'sqrt' ? `√(${currentValue}) =` : `${currentValue}% =`;
      setEquation(formattedEq);
      setDisplay(String(resp.result));

      const newHistoryItem: HistoryItem = {
        id: Math.random().toString(36).substring(2, 9),
        expression: formattedEq,
        result: resp.result,
        timestamp: new Date(),
      };
      setHistory((prev) => [newHistoryItem, ...prev.slice(0, 19)]);

      isNewInputRef.current = true;
    } catch (err: unknown) {
      if (err instanceof CalculatorApiError) {
        setError(err.message);
      } else {
        setError('Error calculating result');
      }
    } finally {
      setLoading(false);
    }
  }, [display]);

  const loadHistoryItem = useCallback((item: HistoryItem) => {
    setDisplay(String(item.result));
    setEquation(item.expression);
    setPendingValue(null);
    setPendingOp(null);
    isNewInputRef.current = true;
    setError(null);
  }, []);

  return {
    display,
    equation,
    loading,
    error,
    history,
    pendingOp,
    inputDigit,
    inputDecimal,
    setOperation,
    executeEquals,
    executeUnary,
    toggleSign,
    deleteDigit,
    clear,
    clearHistory,
    loadHistoryItem,
  };
}
