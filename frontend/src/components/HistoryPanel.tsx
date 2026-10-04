import React from 'react';
import type { HistoryItem } from '../types/calculator';
import { Trash2 } from 'lucide-react';

interface HistoryPanelProps {
  history: HistoryItem[];
  onSelect: (item: HistoryItem) => void;
  onClear: () => void;
}

export const HistoryPanel: React.FC<HistoryPanelProps> = ({
  history,
  onSelect,
  onClear,
}) => {
  return (
    <div className="history-panel" data-testid="history-panel">
      <div className="history-panel-header">
        <span>CALCULATION HISTORY ({history.length})</span>
        {history.length > 0 && (
          <button
            type="button"
            className="history-clear-btn"
            onClick={onClear}
            data-testid="clear-history-btn"
            title="Clear History"
          >
            <Trash2 size={13} style={{ display: 'inline', marginRight: 4 }} />
            Clear
          </button>
        )}
      </div>

      {history.length === 0 ? (
        <div className="history-empty">No calculations yet</div>
      ) : (
        history.map((item) => (
          <div
            key={item.id}
            className="history-item"
            onClick={() => onSelect(item)}
            title="Click to load result"
            data-testid={`history-item-${item.id}`}
          >
            <span className="history-expr">{item.expression}</span>
            <span className="history-res">{item.result}</span>
          </div>
        ))
      )}
    </div>
  );
};
