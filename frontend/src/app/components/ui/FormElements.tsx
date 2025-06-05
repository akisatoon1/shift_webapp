// FormElements.tsx
// フォーム関連のコンポーネント
import React from 'react';

interface FormFieldProps {
    label: string;
    children: React.ReactNode;
    error?: string;
}

export const FormField: React.FC<FormFieldProps> = ({ label, children, error }) => {
    return (
        <label className="flex flex-col gap-1">
            <span className="text-sm">{label}</span>
            {children}
            {error && <span className="text-red-600 text-xs">{error}</span>}
        </label>
    );
};

interface TextInputProps {
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    type?: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    name?: string;
    id?: string;
    autoComplete?: string;
    className?: string;
}

export const TextInput: React.FC<TextInputProps> = ({
    value,
    onChange,
    type = 'text',
    placeholder,
    required = false,
    disabled = false,
    name,
    id,
    autoComplete,
    className = '',
}) => {
    return (
        <input
            type={type}
            value={value}
            onChange={onChange}
            placeholder={placeholder}
            required={required}
            disabled={disabled}
            name={name}
            id={id}
            autoComplete={autoComplete}
            className={`border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 ${className}`}
        />
    );
};

interface ErrorMessageProps {
    children: React.ReactNode;
}

export const ErrorMessage: React.FC<ErrorMessageProps> = ({ children }) => {
    if (!children) return null;
    return <div className="text-red-600 text-sm text-center">{children}</div>;
};

interface SuccessMessageProps {
    children: React.ReactNode;
}

export const SuccessMessage: React.FC<SuccessMessageProps> = ({ children }) => {
    if (!children) return null;
    return <div className="text-green-600 text-sm text-center">{children}</div>;
};
