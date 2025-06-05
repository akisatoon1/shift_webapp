// useRequests.ts
// リクエスト操作用のカスタムフック
import { useState } from 'react';
import { get, post } from '@/app/lib/api';
import { Request, RequestDetail } from '@/app/types';

interface UseRequestsResult {
    requests: Request[];
    isLoading: boolean;
    error: string | null;
    fetchRequests: () => Promise<void>;
    createRequest: (data: CreateRequestData) => Promise<boolean>;
    isCreating: boolean;
    createError: string | null;
}

interface CreateRequestData {
    start_date: string;
    end_date: string;
    deadline: string;
}

export function useRequests(): UseRequestsResult {
    const [requests, setRequests] = useState<Request[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [isCreating, setIsCreating] = useState<boolean>(false);
    const [createError, setCreateError] = useState<string | null>(null);

    const fetchRequests = async () => {
        setIsLoading(true);
        setError(null);

        try {
            const res = await get(`/requests`);
            if (res && res.ok) {
                const data = await res.json();
                setRequests(data);
            } else if (res) {
                const data = await res.json();
                setError(data.error || "取得に失敗しました");
            }
        } catch (err) {
            setError("通信エラーが発生しました");
        } finally {
            setIsLoading(false);
        }
    };

    const createRequest = async (data: CreateRequestData): Promise<boolean> => {
        setIsCreating(true);
        setCreateError(null);

        try {
            const res = await post(`/requests`, data);

            if (res && res.ok) {
                await fetchRequests();
                return true;
            } else if (res) {
                const errorData = await res.json();
                setCreateError(errorData.error || "作成に失敗しました");
            }
            return false;
        } catch (err) {
            setCreateError("通信エラーが発生しました");
            return false;
        } finally {
            setIsCreating(false);
        }
    };

    return {
        requests,
        isLoading,
        error,
        fetchRequests,
        createRequest,
        isCreating,
        createError
    };
}

interface UseRequestDetailResult {
    requestDetail: RequestDetail | null;
    isLoading: boolean;
    error: string | null;
    fetchRequestDetail: (id: string) => Promise<void>;
}

export function useRequestDetail(): UseRequestDetailResult {
    const [requestDetail, setRequestDetail] = useState<RequestDetail | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const fetchRequestDetail = async (id: string) => {
        if (!id) return;

        setIsLoading(true);
        setError(null);

        try {
            const res = await get(`/requests/${id}`);
            if (res && res.ok) {
                const data = await res.json();
                setRequestDetail(data);
            } else if (res) {
                const data = await res.json();
                setError(data.error || "取得に失敗しました");
            }
        } catch (err) {
            setError("通信エラーが発生しました");
        } finally {
            setIsLoading(false);
        }
    };

    return {
        requestDetail,
        isLoading,
        error,
        fetchRequestDetail
    };
}
