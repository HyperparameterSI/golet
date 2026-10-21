// Author: HyperparameterSI

/** The prerendered startup screen works before React finishes loading. */
export function LoadingScreen() {
  return (
    <main className="loading-screen" role="status" aria-live="polite" aria-label="Loading Golet">
      <div className="loading-content">
        <div className="loading-emblem" aria-hidden="true">
          <span className="loading-orbit" />
          <div className="loading-logo">
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
              <path d="m16 5 12 6-12 6L4 11 16 5Z" fill="currentColor" fillOpacity=".12" />
              <path
                d="m4 11 12-6 12 6-12 6-12-6Zm0 5 12 6 12-6M4 21l12 6 12-6"
                stroke="currentColor"
                strokeWidth="1.7"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
        </div>
        <h1 className="loading-title">
          Golet
          <span className="loading-title-dot" aria-hidden="true">
            .
          </span>
        </h1>
        <p className="loading-caption">正在准备你的工作空间</p>
        <div className="loading-track" aria-hidden="true">
          <span />
        </div>
      </div>
    </main>
  );
}
