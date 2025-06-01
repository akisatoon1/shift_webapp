// Fetch APIをラップして、401エラーを自動的に処理するためのユーティリティ

/**
 * 共通のAPI呼び出し関数
 * 401エラー時に自動的にログインページにリダイレクトします
 */
export const fetchWithAuth = async (url: string, options: RequestInit = {}) => {
    // クッキーを含めるためのデフォルトオプション
    const defaultOptions: RequestInit = {
        credentials: 'include',
        ...options,
    };

    try {
        const response = await fetch(url, defaultOptions);

        // 401エラーが発生したらログインページにリダイレクト
        if (response.status === 401) {
            console.error('Authentication failed, redirecting to login page');
            window.location.href = '/login';
            return null;
        }

        return response;
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
};

/**
 * GETリクエスト用のショートハンド関数
 */
export const get = async (url: string, options: RequestInit = {}) => {
    return fetchWithAuth(url, {
        method: 'GET',
        ...options,
    });
};

/**
 * POSTリクエスト用のショートハンド関数
 */
export const post = async (url: string, data: any, options: RequestInit = {}) => {
    return fetchWithAuth(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {}),
        },
        body: JSON.stringify(data),
        ...options,
    });
};

/**
 * DELETEリクエスト用のショートハンド関数
 */
export const del = async (url: string, options: RequestInit = {}) => {
    return fetchWithAuth(url, {
        method: 'DELETE',
        ...options,
    });
};
