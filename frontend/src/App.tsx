import { useState, useEffect } from 'react';
import { useCalculator } from './hooks/useCalculator';
import { Display } from './components/Display';
import { Keypad } from './components/Keypad';
import { HistoryPanel } from './components/HistoryPanel';
import { Header } from './components/Header';
import { checkHealth } from './services/api';
import { History, Keyboard } from 'lucide-react';

export function App() {
  const {
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
  } = useCalculator();

  const [showHistory, setShowHistory] = useState<boolean>(false);
  const [apiStatus, setApiStatus] = useState<'online' | 'offline' | 'checking'>('checking');

  // Verify backend connectivity
  useEffect(() => {
    let isMounted = true;
    const verifyApi = async () => {
      try {
        await checkHealth();
        if (isMounted) setApiStatus('online');
      } catch {
        if (isMounted) setApiStatus('offline');
      }
    };

    verifyApi();
    const interval = setInterval(verifyApi, 10000);
    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, []);

  // Physical keyboard support
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Avoid capturing shortcuts like Ctrl+R or F5
      if (e.ctrlKey || e.altKey || e.metaKey) return;

      const key = e.key;

      if (/^[0-9]$/.test(key)) {
        e.preventDefault();
        inputDigit(key);
      } else if (key === '.' || key === ',') {
        e.preventDefault();
        inputDecimal();
      } else if (key === '+') {
        e.preventDefault();
        setOperation('add');
      } else if (key === '-') {
        e.preventDefault();
        setOperation('subtract');
      } else if (key === '*') {
        e.preventDefault();
        setOperation('multiply');
      } else if (key === '/') {
        e.preventDefault();
        setOperation('divide');
      } else if (key === '^') {
        e.preventDefault();
        setOperation('power');
      } else if (key === '%') {
        e.preventDefault();
        executeUnary('percentage');
      } else if (key === 'Enter' || key === '=') {
        e.preventDefault();
        executeEquals();
      } else if (key === 'Backspace') {
        e.preventDefault();
        deleteDigit();
      } else if (key === 'Escape') {
        e.preventDefault();
        clear();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [
    inputDigit,
    inputDecimal,
    setOperation,
    executeEquals,
    executeUnary,
    deleteDigit,
    clear,
  ]);

  return (
    <div className="calculator-app-container">
      <Header apiStatus={apiStatus} />

      <main className="calculator-card">
        <Display
          value={display}
          equation={equation}
          loading={loading}
          error={error}
        />

        <Keypad
          onDigit={inputDigit}
          onDecimal={inputDecimal}
          onOperation={setOperation}
          onUnary={executeUnary}
          onEquals={executeEquals}
          onClear={clear}
          onDelete={deleteDigit}
          onToggleSign={toggleSign}
          activeOp={pendingOp}
          disabled={loading}
        />

        <div className="calc-toolbar">
          <button
            type="button"
            className="history-toggle-btn"
            onClick={() => setShowHistory((prev) => !prev)}
            data-testid="toggle-history-btn"
          >
            <History size={14} />
            <span>{showHistory ? 'Hide History' : `History (${history.length})`}</span>
          </button>

          <span className="keyboard-hint" title="You can use your numpad or keyboard">
            <Keyboard size={12} style={{ display: 'inline', marginRight: 4, verticalAlign: 'middle' }} />
            Keyboard enabled
          </span>
        </div>

        {showHistory && (
          <HistoryPanel
            history={history}
            onSelect={loadHistoryItem}
            onClear={clearHistory}
          />
        )}
      </main>
    </div>
  );
}

export default App;
