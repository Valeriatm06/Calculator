import React from 'react';

interface HeaderProps {
  apiStatus: 'online' | 'offline' | 'checking';
}

export const Header: React.FC<HeaderProps> = ({ apiStatus }) => {
  return (
    <header className="app-header">
      <div className="brand-section">
        <div className="brand-icon">∑</div>
        <span className="brand-title">Calculator</span>
      </div>

      <a
        href="http://localhost:8080"
        target="_blank"
        rel="noreferrer"
        className={`status-badge ${apiStatus === 'offline' ? 'error' : ''}`}
        title="Click to open API documentation"
        data-testid="api-status-badge"
        style={{ textDecoration: 'none', cursor: 'pointer' }}
      >
        <span
          className={`status-dot ${apiStatus === 'offline' ? 'error' : ''}`}
        />
        <span>{apiStatus === 'online' ? 'API Online' : apiStatus === 'offline' ? 'API Offline' : 'Connecting...'}</span>
      </a>
    </header>
  );
};
