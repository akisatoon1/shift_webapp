// Table.tsx
// テーブルコンポーネント
import React from 'react';

interface TableProps {
    children: React.ReactNode;
    className?: string;
}

export const Table: React.FC<TableProps> = ({ children, className = '' }) => {
    return (
        <div className="overflow-x-auto w-full">
            <table className={`w-full border ${className}`}>
                {children}
            </table>
        </div>
    );
};

interface TableHeaderProps {
    children: React.ReactNode;
}

export const TableHeader: React.FC<TableHeaderProps> = ({ children }) => {
    return <thead className="bg-gray-100">{children}</thead>;
};

interface TableBodyProps {
    children: React.ReactNode;
}

export const TableBody: React.FC<TableBodyProps> = ({ children }) => {
    return <tbody>{children}</tbody>;
};

interface TableRowProps {
    children: React.ReactNode;
    onClick?: () => void;
    className?: string;
}

export const TableRow: React.FC<TableRowProps> = ({
    children,
    onClick,
    className = ''
}) => {
    return (
        <tr
            className={`${onClick ? 'hover:bg-blue-50 cursor-pointer' : ''} ${className}`}
            onClick={onClick}
        >
            {children}
        </tr>
    );
};

interface TableCellProps {
    children: React.ReactNode;
    header?: boolean;
    className?: string;
}

export const TableCell: React.FC<TableCellProps> = ({
    children,
    header = false,
    className = ''
}) => {
    const Tag = header ? 'th' : 'td';
    return (
        <Tag className={`border px-2 py-1 ${className}`}>
            {children}
        </Tag>
    );
};
