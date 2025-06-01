// Fetch APIをラップして、401エラーを自動的に処理するためのユーティリティ

// API ベースURLの設定
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

/**
 * 共通のAPI呼び出し関数
 * 401エラー時に自動的にログインページにリダイレクトします
 */
export const fetchWithAuth = async (endpoint: string, options: RequestInit = {}) => {
    // API ベースURLとエンドポイントを結合
    const url = `${API_BASE_URL}${endpoint}`;
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
export const get = async (endpoint: string, options: RequestInit = {}) => {
    return fetchWithAuth(endpoint, {
        method: 'GET',
        ...options,
    });
};

/**
 * POSTリクエスト用のショートハンド関数
 */
export const post = async (endpoint: string, data: any, options: RequestInit = {}) => {
    return fetchWithAuth(endpoint, {
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
export const del = async (endpoint: string, options: RequestInit = {}) => {
    return fetchWithAuth(endpoint, {
        method: 'DELETE',
        ...options,
    });
};
