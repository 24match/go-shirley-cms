/**
 * 错误边界组件
 * 捕获子组件树中的 JavaScript 错误并显示降级 UI
 */

import React, { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
    children: ReactNode;
    fallback?: ReactNode;
}

interface State {
    hasError: boolean;
    error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('ErrorBoundary caught an error:', error, errorInfo);
    }

    private handleRetry = () => {
        this.setState({ hasError: false, error: undefined });
        window.location.reload();
    };

    public render() {
        if (this.state.hasError) {
            if (this.props.fallback) {
                return this.props.fallback;
            }

            return (
                <div style={styles.container}>
                    <div style={styles.content}>
                        <h1 style={styles.title}>Something went wrong</h1>
                        <p style={styles.message}>
                            We're sorry, but something unexpected happened.
                        </p>
                        {this.state.error && (
                            <details style={styles.details}>
                                <summary style={styles.summary}>Error Details</summary>
                                <pre style={styles.pre}>{this.state.error.toString()}</pre>
                            </details>
                        )}
                        <button 
                            style={styles.button}
                            onClick={this.handleRetry}
                        >
                            Reload Page
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}

const styles: Record<string, React.CSSProperties> = {
    container: {
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#f5f5f5',
        padding: '20px'
    },
    content: {
        maxWidth: '500px',
        textAlign: 'center',
        backgroundColor: '#fff',
        padding: '40px',
        borderRadius: '8px',
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
    },
    title: {
        fontSize: '24px',
        color: '#333',
        marginBottom: '16px'
    },
    message: {
        fontSize: '16px',
        color: '#666',
        marginBottom: '24px'
    },
    details: {
        textAlign: 'left',
        marginBottom: '24px',
        backgroundColor: '#f8f8f8',
        padding: '16px',
        borderRadius: '4px'
    },
    summary: {
        cursor: 'pointer',
        color: '#666',
        fontSize: '14px'
    },
    pre: {
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-word',
        fontSize: '12px',
        color: '#999',
        marginTop: '8px'
    },
    button: {
        padding: '12px 24px',
        fontSize: '16px',
        backgroundColor: '#06a499',
        color: '#fff',
        border: 'none',
        borderRadius: '4px',
        cursor: 'pointer',
        transition: 'background-color 0.3s'
    }
};

export default ErrorBoundary;