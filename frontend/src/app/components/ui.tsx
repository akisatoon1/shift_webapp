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
        outline: 'bg-white text-blue-700 border-blue-400 border'
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

interface CardProps {
    children: React.ReactNode;
    className?: string;
}

export const Card: React.FC<CardProps> = ({ children, className = '' }) => {
    return (
        <div className={`p-8 bg-white rounded shadow-md ${className}`}>
            {children}
        </div>
    );
};

interface PageContainerProps {
    children: React.ReactNode;
    className?: string;
    maxWidth?: string;
}

export const PageContainer: React.FC<PageContainerProps> = ({
    children,
    className = '',
    maxWidth = 'max-w-2xl'
}) => {
    return (
        <div className={`flex flex-col items-center justify-center min-h-screen bg-gray-50 ${className}`}>
            <div className={`w-full ${maxWidth} p-8 bg-white rounded shadow-md`}>
                {children}
            </div>
        </div>
    );
};

interface PageTitleProps {
    children: React.ReactNode;
}

export const PageTitle: React.FC<PageTitleProps> = ({ children }) => {
    return <h1 className="text-2xl font-bold mb-6 text-center">{children}</h1>;
};

interface FormFieldProps {
    label: string;
    children: React.ReactNode;
}

export const FormField: React.FC<FormFieldProps> = ({ label, children }) => {
    return (
        <label className="flex flex-col gap-1">
            <span className="text-sm">{label}</span>
            {children}
        </label>
    );
};

interface LoadingIndicatorProps {
    message?: string;
}

export const LoadingIndicator: React.FC<LoadingIndicatorProps> = ({ message = '読み込み中...' }) => {
    return <div className="text-center">{message}</div>;
};

interface ErrorMessageProps {
    message: string;
}

export const ErrorMessage: React.FC<ErrorMessageProps> = ({ message }) => {
    return <div className="text-red-600 text-center">{message}</div>;
};

interface SuccessMessageProps {
    message: string;
}

export const SuccessMessage: React.FC<SuccessMessageProps> = ({ message }) => {
    return <div className="text-green-600 text-center text-sm mb-2">{message}</div>;
};

interface TabButtonProps {
    active: boolean;
    onClick: () => void;
    children: React.ReactNode;
}

export const TabButton: React.FC<TabButtonProps> = ({ active, onClick, children }) => {
    return (
        <button
            className={`inline-block px-4 py-2 rounded font-semibold border ml-2 
        ${active ? 'bg-blue-600 text-white border-blue-700' : 'bg-white text-blue-700 border-blue-400'}`}
            onClick={onClick}
        >
            {children}
        </button>
    );
};
