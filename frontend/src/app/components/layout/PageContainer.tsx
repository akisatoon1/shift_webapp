// PageContainer.tsx
// ページ全体をラップするコンテナコンポーネント
import React from 'react';

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
