// useSession.ts
// セッション管理用のカスタムフック
import { useState, useEffect } from 'react';
import { get } from '@/app/lib/api';
import { User, SessionData } from '@/app/types';

interface UseSessionResult {
    user: User | null;
    roles: string[];
    isLoading: boolean;
    error: string | null;
    refetch: () => Promise<void>;
}

export function useSession(): UseSessionResult {
    const [user, setUser] = useState<User | null>(null);
    const [roles, setRoles] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    const fetchSession = async () => {
        setIsLoading(true);
        setError(null);

        try {
            const response = await get(`/session`);

            if (response && response.ok) {
                const data: SessionData = await response.json();
                setUser(data.user);
                setRoles(data.user?.roles || []);
            } else {
                setUser(null);
                setRoles([]);
                setError('セッション情報の取得に失敗しました');
            }
        } catch (err) {
            setUser(null);
            setRoles([]);
            setError('セッション情報の取得中にエラーが発生しました');
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchSession();
    }, []);

    return {
        user,
        roles,
        isLoading,
        error,
        refetch: fetchSession
    };
}
