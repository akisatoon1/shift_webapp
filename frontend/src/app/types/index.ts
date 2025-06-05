// 共通型定義
export type User = {
    id: number;
    name: string;
    roles?: string[];
};

export type Request = {
    id: number;
    creator: User;
    start_date: string;
    end_date: string;
    deadline: string;
    created_at: string;
};

export type RequestDetail = Request & {
    submissions: Submission[];
    entries: Entry[];
};

export type Entry = {
    id: number;
    user: User;
    date: string;
    hour: number;
};

export type Submission = {
    submitter: User;
};

export type APIResponse<T = any> = {
    data?: T;
    error?: string;
};

export type SessionData = {
    user: User;
};
