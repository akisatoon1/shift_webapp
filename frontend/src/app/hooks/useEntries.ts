// useEntries.ts
// エントリー提出用のカスタムフック
import { useState } from 'react';
import { post } from '@/app/lib/api';

interface EntryData {
    date: string;
    hour: number;
}

interface UseEntriesResult {
    submitEntries: (requestId: string, entries: EntryData[]) => Promise<boolean>;
    isSubmitting: boolean;
    error: string | null;
    success: boolean;
}

export function useEntries(): UseEntriesResult {
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<boolean>(false);

    const submitEntries = async (requestId: string, entries: EntryData[]): Promise<boolean> => {
        setIsSubmitting(true);
        setError(null);
        setSuccess(false);

        try {
            const res = await post(`/requests/${requestId}/submissions`, entries);

            if (res && res.ok) {
                setSuccess(true);
                return true;
            } else if (res) {
                const data = await res.json();
                setError(data.error || "提出に失敗しました");
            }
            return false;
        } catch (err) {
            setError("通信エラーが発生しました");
            return false;
        } finally {
            setIsSubmitting(false);
        }
    };

    return {
        submitEntries,
        isSubmitting,
        error,
        success
    };
}
