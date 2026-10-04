import React from 'react';
import { Delete } from 'lucide-react';
import type { OperationType } from '../types/calculator';

interface KeypadProps {
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperation: (op: OperationType) => void;
  onUnary: (op: 'sqrt' | 'percentage') => void;
  onEquals: () => void;
  onClear: () => void;
  onDelete: () => void;
  onToggleSign: () => void;
  activeOp?: OperationType | null;
  disabled?: boolean;
}

export const Keypad: React.FC<KeypadProps> = ({
  onDigit,
  onDecimal,
  onOperation,
  onUnary,
  onEquals,
  onClear,
  onDelete,
  onToggleSign,
  activeOp = null,
  disabled = false,
}) => {
  return (
    <div className="keypad-container" data-testid="calculator-keypad">
      {/* Mini Function Bar (as seen in Google/Android Calculator) */}
      <div className="mini-function-bar">
        <button
          type="button"
          className="mini-func-btn"
          onClick={() => onUnary('sqrt')}
          disabled={disabled}
          data-testid="key-sqrt"
          title="Square Root"
        >
          √
        </button>
        <button
          type="button"
          className="mini-func-btn"
          onClick={() => onOperation('power')}
          disabled={disabled}
          data-testid="key-power"
          title="Exponentiation (^)"
        >
          ^
        </button>
        <button
          type="button"
          className="mini-func-btn"
          onClick={() => onUnary('percentage')}
          disabled={disabled}
          data-testid="key-percentage"
          title="Percentage (%)"
        >
          %
        </button>
        <button
          type="button"
          className="mini-func-btn"
          onClick={onToggleSign}
          disabled={disabled}
          data-testid="key-plusminus"
          title="Toggle Sign (±)"
        >
          ±
        </button>
      </div>

      {/* Main Circular Keypad Grid (4 columns x 5 rows) */}
      <div className="circle-keypad-grid">
        {/* Row 1: AC (soft green), ± (soft blue), ^ (soft blue), ÷ (soft blue) */}
        <button
          type="button"
          className="circle-key key-ac"
          onClick={onClear}
          data-testid="key-clear"
          title="Clear All (Esc)"
        >
          AC
        </button>
        <button
          type="button"
          className={`circle-key key-op ${activeOp === 'power' ? 'active-op' : ''}`}
          onClick={() => onOperation('power')}
          disabled={disabled}
          title="Exponentiation (^)"
        >
          ( )
        </button>
        <button
          type="button"
          className={`circle-key key-op ${activeOp === 'power' ? 'active-op' : ''}`}
          onClick={() => onOperation('power')}
          disabled={disabled}
          title="Exponentiation (^)"
        >
          ^
        </button>
        <button
          type="button"
          className={`circle-key key-op ${activeOp === 'divide' ? 'active-op' : ''}`}
          onClick={() => onOperation('divide')}
          disabled={disabled}
          data-testid="key-divide"
          title="Divide (÷)"
        >
          ÷
        </button>

        {/* Row 2: 7, 8, 9, × */}
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('7')}
          data-testid="key-7"
        >
          7
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('8')}
          data-testid="key-8"
        >
          8
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('9')}
          data-testid="key-9"
        >
          9
        </button>
        <button
          type="button"
          className={`circle-key key-op ${activeOp === 'multiply' ? 'active-op' : ''}`}
          onClick={() => onOperation('multiply')}
          disabled={disabled}
          data-testid="key-multiply"
          title="Multiply (×)"
        >
          ×
        </button>

        {/* Row 3: 4, 5, 6, − */}
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('4')}
          data-testid="key-4"
        >
          4
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('5')}
          data-testid="key-5"
        >
          5
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('6')}
          data-testid="key-6"
        >
          6
        </button>
        <button
          type="button"
          className={`circle-key key-op ${activeOp === 'subtract' ? 'active-op' : ''}`}
          onClick={() => onOperation('subtract')}
          disabled={disabled}
          data-testid="key-subtract"
          title="Subtract (−)"
        >
          −
        </button>

        {/* Row 4: 1, 2, 3, + */}
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('1')}
          data-testid="key-1"
        >
          1
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('2')}
          data-testid="key-2"
        >
          2
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('3')}
          data-testid="key-3"
        >
          3
        </button>
        <button
          type="button"
          className={`circle-key key-op ${activeOp === 'add' ? 'active-op' : ''}`}
          onClick={() => onOperation('add')}
          disabled={disabled}
          data-testid="key-add"
          title="Add (+)"
        >
          +
        </button>

        {/* Row 5: 0, ., ⌫, = */}
        <button
          type="button"
          className="circle-key key-num"
          onClick={() => onDigit('0')}
          data-testid="key-0"
        >
          0
        </button>
        <button
          type="button"
          className="circle-key key-num"
          onClick={onDecimal}
          data-testid="key-decimal"
        >
          .
        </button>
        <button
          type="button"
          className="circle-key key-num key-backspace"
          onClick={onDelete}
          data-testid="key-delete"
          title="Backspace"
        >
          <Delete size={22} />
        </button>
        <button
          type="button"
          className="circle-key key-equals"
          onClick={onEquals}
          disabled={disabled}
          data-testid="key-equals"
          title="Calculate (=)"
        >
          =
        </button>
      </div>
    </div>
  );
};
