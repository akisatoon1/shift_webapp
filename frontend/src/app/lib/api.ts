// Fetch APIをラップして、401エラーを自動的に処理するためのユーティリティ
import { APIResponse } from '../types';

// API ベースURLの設定
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

/**
 * API 呼び出し中にエラーが発生した場合のエラー型
 */
export class ApiError extends Error {
    status: number;
    data?: any;

    constructor(message: string, status: number, data?: any) {
        super(message);
        this.name = 'ApiError';
        this.status = status;
        this.data = data;
    }
}

/**
 * 共通のAPI呼び出し関数
 * 401エラー時に自動的にログインページにリダイレクトします
 */
export const fetchWithAuth = async <T>(endpoint: string, options: RequestInit = {}): Promise<Response> => {
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
            return null as any;
        }

        return response;
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
};

/**
 * レスポンスの処理
 */
async function handleResponse<T>(response: Response | null): Promise<T> {
    if (!response) {
        throw new ApiError('Response is null', 0);
    }

    // JSONデータの取得を試みる
    let data: any;
    try {
        data = await response.json();
    } catch (error) {
        // JSONでない場合はテキストを取得
        data = { message: await response.text() };
    }

    if (!response.ok) {
        throw new ApiError(
            data.error || `API error: ${response.status}`,
            response.status,
            data
        );
    }

    return data as T;
}

/**
 * GETリクエスト用のショートハンド関数
 */
export const get = async <T>(endpoint: string, options: RequestInit = {}): Promise<Response> => {
    return fetchWithAuth(endpoint, {
        method: 'GET',
        ...options,
    });
};

/**
 * POSTリクエスト用のショートハンド関数
 */
export const post = async <T>(endpoint: string, data: any, options: RequestInit = {}): Promise<Response> => {
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
export const del = async <T>(endpoint: string, options: RequestInit = {}): Promise<Response> => {
    return fetchWithAuth(endpoint, {
        method: 'DELETE',
        ...options,
    });
};

/**
 * PUTリクエスト用のショートハンド関数
 */
export const put = async <T>(endpoint: string, data: any, options: RequestInit = {}): Promise<Response> => {
    return fetchWithAuth(endpoint, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {}),
        },
        body: JSON.stringify(data),
        ...options,
    });
};
