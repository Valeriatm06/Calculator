import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react';
import App from './App';
import * as api from './services/api';

describe('Calculator App Integration', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    vi.spyOn(api, 'checkHealth').mockResolvedValue({ status: 'healthy', service: 'calculator-api' });
  });

  it('renders calculator initial UI elements', async () => {
    render(<App />);

    expect(screen.getByTestId('display-value')).toHaveTextContent('0');
    expect(screen.getByTestId('key-add')).toBeInTheDocument();
    expect(screen.getByText('Calculator')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByTestId('api-status-badge')).toHaveTextContent('API Online');
    });
  });

  it('performs full calculation flow via UI buttons', async () => {
    vi.spyOn(api, 'calculate').mockResolvedValueOnce({
      operation: 'add',
      a: 12,
      b: 34,
      result: 46,
      formatted: '12 + 34 = 46',
    });

    render(<App />);

    // Click 1, 2
    fireEvent.click(screen.getByTestId('key-1'));
    fireEvent.click(screen.getByTestId('key-2'));
    expect(screen.getByTestId('display-value')).toHaveTextContent('12');

    // Click +
    fireEvent.click(screen.getByTestId('key-add'));
    expect(screen.getByTestId('display-equation')).toHaveTextContent('12 +');

    // Click 3, 4
    fireEvent.click(screen.getByTestId('key-3'));
    fireEvent.click(screen.getByTestId('key-4'));
    expect(screen.getByTestId('display-value')).toHaveTextContent('34');

    // Click =
    fireEvent.click(screen.getByTestId('key-equals'));

    await waitFor(() => {
      expect(screen.getByTestId('display-value')).toHaveTextContent('46');
    });
  });

  it('supports keyboard input for calculation', async () => {
    vi.spyOn(api, 'calculate').mockResolvedValueOnce({
      operation: 'multiply',
      a: 9,
      b: 9,
      result: 81,
      formatted: '9 × 9 = 81',
    });

    render(<App />);

    // Type 9
    fireEvent.keyDown(window, { key: '9' });
    expect(screen.getByTestId('display-value')).toHaveTextContent('9');

    // Type *
    fireEvent.keyDown(window, { key: '*' });

    // Type 9
    fireEvent.keyDown(window, { key: '9' });

    // Type Enter
    fireEvent.keyDown(window, { key: 'Enter' });

    await waitFor(() => {
      expect(screen.getByTestId('display-value')).toHaveTextContent('81');
    });
  });

  it('toggles and displays calculation history', async () => {
    vi.spyOn(api, 'calculate').mockResolvedValueOnce({
      operation: 'multiply',
      a: 5,
      b: 5,
      result: 25,
      formatted: '5 × 5 = 25',
    });

    render(<App />);

    fireEvent.click(screen.getByTestId('key-5'));
    fireEvent.click(screen.getByTestId('key-multiply'));
    fireEvent.click(screen.getByTestId('key-5'));
    fireEvent.click(screen.getByTestId('key-equals'));

    await waitFor(() => {
      expect(screen.getByTestId('display-value')).toHaveTextContent('25');
    });

    // Toggle history panel
    const historyBtn = screen.getByTestId('toggle-history-btn');
    fireEvent.click(historyBtn);

    const historyPanel = screen.getByTestId('history-panel');
    expect(historyPanel).toBeInTheDocument();
    expect(within(historyPanel).getByText('5 × 5 =')).toBeInTheDocument();
    expect(within(historyPanel).getByText('25')).toBeInTheDocument();
  });

  it('correctly handles 55 + 55: clears main display on operator and replaces cleanly without 5555', async () => {
    vi.spyOn(api, 'calculate').mockResolvedValueOnce({
      operation: 'add',
      a: 55,
      b: 55,
      result: 110,
      formatted: '55 + 55 = 110',
    });

    render(<App />);

    // Type 55
    fireEvent.click(screen.getByTestId('key-5'));
    fireEvent.click(screen.getByTestId('key-5'));
    expect(screen.getByTestId('display-value')).toHaveTextContent('55');

    // Click + (55 moves to equation, main display clears/resets to 0, active operator glows)
    fireEvent.click(screen.getByTestId('key-add'));
    expect(screen.getByTestId('display-equation')).toHaveTextContent('55 +');
    expect(screen.getByTestId('display-value')).toHaveTextContent('0');
    expect(screen.getByTestId('key-add')).toHaveClass('active-op');

    // Type 55 again
    fireEvent.click(screen.getByTestId('key-5'));
    expect(screen.getByTestId('display-value')).toHaveTextContent('5'); // NOT 555
    fireEvent.click(screen.getByTestId('key-5'));
    expect(screen.getByTestId('display-value')).toHaveTextContent('55'); // NOT 5555

    // Click =
    fireEvent.click(screen.getByTestId('key-equals'));

    await waitFor(() => {
      expect(screen.getByTestId('display-value')).toHaveTextContent('110');
      expect(screen.getByTestId('display-equation')).toHaveTextContent('55 + 55 =');
    });
  });
});
