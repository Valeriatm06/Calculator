import React from 'react';

interface HeaderProps {
  apiStatus: 'online' | 'offline' | 'checking';
}

export const Header: React.FC<HeaderProps> = ({ apiStatus }) => {
  return (
    <header className="app-header">
      <div className="brand-section">
        <div className="brand-icon">∑</div>
        <div>
          <span className="brand-title">FullStack Calc</span>
          <span className="brand-badge" style={{ marginLeft: '8px' }}>Go + React</span>
        </div>
      </div>

      <div
        className={`status-badge ${apiStatus === 'offline' ? 'error' : ''}`}
        title={`API Status: ${apiStatus}`}
        data-testid="api-status-badge"
      >
        <span
          className={`status-dot ${apiStatus === 'offline' ? 'error' : ''}`}
        />
        <span>{apiStatus === 'online' ? 'API Online' : apiStatus === 'offline' ? 'API Offline' : 'Connecting...'}</span>
      </div>
    </header>
  );
};
