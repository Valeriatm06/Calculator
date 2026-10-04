import React from 'react';
import { AlertCircle, Loader2 } from 'lucide-react';

interface DisplayProps {
  value: string;
  equation: string;
  loading: boolean;
  error: string | null;
}

export const Display: React.FC<DisplayProps> = ({
  value,
  equation,
  loading,
  error,
}) => {
  return (
    <div className="calculator-display" data-testid="calculator-display">
      {/* Upper line: Equation / Expression (as seen in the screenshot) */}
      <div className="display-expression" data-testid="display-equation">
        {equation || '\u00A0'}
      </div>

      {/* Main line: Current Value / Result */}
      <div className="display-main">
        {loading && (
          <Loader2
            className="spinner"
            size={18}
            data-testid="loading-spinner"
          />
        )}
        <div
          className={`display-value ${loading ? 'loading' : ''}`}
          data-testid="display-value"
        >
          {value}
        </div>
      </div>

      {error && (
        <div className="display-error-banner" data-testid="display-error">
          <AlertCircle size={14} />
          <span>{error}</span>
        </div>
      )}
    </div>
  );
};
