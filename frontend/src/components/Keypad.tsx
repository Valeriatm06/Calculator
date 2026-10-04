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
    <div className="keypad-grid" data-testid="calculator-keypad">
      {/* Row 1: Clear, Delete, Advanced (Sqrt, Power) */}
      <button
        type="button"
        className="key-btn action-clear"
        onClick={onClear}
        data-testid="key-clear"
        title="Clear All (Esc)"
      >
        AC
      </button>
      <button
        type="button"
        className="key-btn action-func"
        onClick={onDelete}
        data-testid="key-delete"
        title="Backspace"
      >
        <Delete size={20} />
      </button>
      <button
        type="button"
        className="key-btn action-func"
        onClick={() => onUnary('sqrt')}
        disabled={disabled}
        data-testid="key-sqrt"
        title="Square Root"
      >
        √
      </button>
      <button
        type="button"
        className={`key-btn action-func ${activeOp === 'power' ? 'active-op' : ''}`}
        onClick={() => onOperation('power')}
        disabled={disabled}
        data-testid="key-power"
        title="Exponentiation (^)"
      >
        xʸ
      </button>

      {/* Row 2: Percentage, +/- sign, Division */}
      <button
        type="button"
        className="key-btn action-func"
        onClick={() => onUnary('percentage')}
        disabled={disabled}
        data-testid="key-percentage"
        title="Percentage (%)"
      >
        %
      </button>
      <button
        type="button"
        className="key-btn action-func"
        onClick={onToggleSign}
        disabled={disabled}
        data-testid="key-plusminus"
        title="Toggle Sign (+/-)"
      >
        ±
      </button>
      <button
        type="button"
        className={`key-btn action-operator key-span-2 ${activeOp === 'divide' ? 'active-op' : ''}`}
        onClick={() => onOperation('divide')}
        disabled={disabled}
        data-testid="key-divide"
        title="Divide (/)"
      >
        ÷
      </button>

      {/* Row 3: 7, 8, 9, Multiply */}
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('7')}
        data-testid="key-7"
      >
        7
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('8')}
        data-testid="key-8"
      >
        8
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('9')}
        data-testid="key-9"
      >
        9
      </button>
      <button
        type="button"
        className={`key-btn action-operator ${activeOp === 'multiply' ? 'active-op' : ''}`}
        onClick={() => onOperation('multiply')}
        disabled={disabled}
        data-testid="key-multiply"
        title="Multiply (*)"
      >
        ×
      </button>

      {/* Row 4: 4, 5, 6, Subtract */}
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('4')}
        data-testid="key-4"
      >
        4
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('5')}
        data-testid="key-5"
      >
        5
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('6')}
        data-testid="key-6"
      >
        6
      </button>
      <button
        type="button"
        className={`key-btn action-operator ${activeOp === 'subtract' ? 'active-op' : ''}`}
        onClick={() => onOperation('subtract')}
        disabled={disabled}
        data-testid="key-subtract"
        title="Subtract (-)"
      >
        −
      </button>

      {/* Row 5: 1, 2, 3, Add */}
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('1')}
        data-testid="key-1"
      >
        1
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('2')}
        data-testid="key-2"
      >
        2
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={() => onDigit('3')}
        data-testid="key-3"
      >
        3
      </button>
      <button
        type="button"
        className={`key-btn action-operator ${activeOp === 'add' ? 'active-op' : ''}`}
        onClick={() => onOperation('add')}
        disabled={disabled}
        data-testid="key-add"
        title="Add (+)"
      >
        +
      </button>

      {/* Row 6: 0, ., Equals */}
      <button
        type="button"
        className="key-btn key-span-2"
        onClick={() => onDigit('0')}
        data-testid="key-0"
      >
        0
      </button>
      <button
        type="button"
        className="key-btn"
        onClick={onDecimal}
        data-testid="key-decimal"
      >
        .
      </button>
      <button
        type="button"
        className="key-btn action-equals"
        onClick={onEquals}
        disabled={disabled}
        data-testid="key-equals"
        title="Calculate (Enter or =)"
      >
        =
      </button>
    </div>
  );
};
