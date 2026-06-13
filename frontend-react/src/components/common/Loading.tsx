/**
 * 加载组件
 * 显示加载动画或骨架屏
 */

import React from 'react';

export interface LoadingProps {
    size?: 'small' | 'medium' | 'large';
    text?: string;
    variant?: 'spinner' | 'skeleton' | 'dots';
    className?: string;
}

const Loading: React.FC<LoadingProps> = ({
    size = 'medium',
    text,
    variant = 'spinner',
    className = ''
}) => {
    const classes = [`loading`, `loading-${variant}`, `loading-${size}`, className].join(' ');

    if (variant === 'skeleton') {
        return (
            <div className={classes}>
                <div className="skeleton skeleton-text" />
                <div className="skeleton skeleton-text" />
                <div className="skeleton skeleton-text skeleton-text-short" />
            </div>
        );
    }

    if (variant === 'dots') {
        return (
            <div className={classes}>
                <div className="loading-dots">
                    <span className="dot" />
                    <span className="dot" />
                    <span className="dot" />
                </div>
                {text && <p className="loading-text">{text}</p>}
            </div>
        );
    }

    // 默认 spinner
    return (
        <div className={classes}>
            <div className="spinner" />
            {text && <p className="loading-text">{text}</p>}
        </div>
    );
};

export default Loading;