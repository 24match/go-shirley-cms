/**
 * 按钮组件
 * 支持多种变体和尺寸
 */

import React, { ButtonHTMLAttributes, AnchorHTMLAttributes } from 'react';

export type ButtonVariant = 'primary' | 'outline' | 'text';
export type ButtonSize = 'small' | 'medium' | 'large';

export interface ButtonProps {
    variant?: ButtonVariant;
    size?: ButtonSize;
    href?: string;
    children: React.ReactNode;
    className?: string;
    disabled?: boolean;
    loading?: boolean;
    onClick?: (e: React.MouseEvent<HTMLButtonElement | HTMLAnchorElement>) => void;
}

const Button: React.FC<ButtonProps> = ({
    variant = 'primary',
    size = 'medium',
    href,
    children,
    className = '',
    disabled = false,
    loading = false,
    onClick,
    ...props
}) => {
    const classes = [
        'btn',
        `btn-${variant}`,
        `btn-${size}`,
        className,
        loading ? 'btn-loading' : '',
        disabled ? 'btn-disabled' : ''
    ].filter(Boolean).join(' ');

    const handleClick = (e: React.MouseEvent<HTMLButtonElement | HTMLAnchorElement>) => {
        if (disabled || loading) {
            e.preventDefault();
            return;
        }
        onClick?.(e);
    };

    if (href) {
        return (
            <a
                href={href}
                className={classes}
                onClick={handleClick}
                {...(props as AnchorHTMLAttributes<HTMLAnchorElement>)}
            >
                {loading && <span className="btn-spinner" />}
                {children}
            </a>
        );
    }

    return (
        <button
            className={classes}
            onClick={handleClick}
            disabled={disabled || loading}
            {...(props as ButtonHTMLAttributes<HTMLButtonElement>)}
        >
            {loading && <span className="btn-spinner" />}
            {children}
        </button>
    );
};

export default Button;