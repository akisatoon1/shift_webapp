// Button.tsx
// 共通のボタンコンポーネント
import React from 'react';

interface ButtonProps {
    children: React.ReactNode;
    onClick?: () => void;
    type?: 'button' | 'submit' | 'reset';
    variant?: 'primary' | 'secondary' | 'outline';
    disabled?: boolean;
    className?: string;
    fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
    children,
    onClick,
    type = 'button',
    variant = 'primary',
    disabled = false,
    className = '',
    fullWidth = false,
}) => {
    // バリアントに基づいてスタイルを設定
    const variantClasses = {
        primary: 'bg-blue-600 text-white hover:bg-blue-700',
        secondary: 'bg-gray-400 text-white hover:bg-gray-500',
        outline: 'bg-white text-blue-700 border-blue-400 border hover:bg-blue-50'
    };

    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled}
            className={`
                rounded px-4 py-2 font-semibold transition
                ${variantClasses[variant]}
                ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
                ${fullWidth ? 'w-full' : ''}
                ${className}
            `}
        >
            {children}
        </button>
    );
};
