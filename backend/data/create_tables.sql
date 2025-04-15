-- SQLite3で外部キー制約を有効にする (アプリケーションで実行する必要あり)
-- PRAGMA foreign_keys = ON;

CREATE TABLE users (
    -- app実装ではidは変更できない。後に変更可能にするかも
    id TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- "yyyy-mm-dd"形式かチェックする
-- CHECK (strftime('%Y-%m-%d', date) IS NOT NULL AND date LIKE '____-__-__')

-- 1. シフト提出要請テーブル (shift_requests)
CREATE TABLE shift_requests (
    id TEXT PRIMARY KEY,
    manager_id TEXT NOT NULL,
    start_date TEXT NOT NULL CHECK (strftime('%Y-%m-%d', start_date) IS NOT NULL AND start_date LIKE '____-__-__'),
    end_date TEXT NOT NULL CHECK (strftime('%Y-%m-%d', end_date) IS NOT NULL AND end_date LIKE '____-__-__'),
    created_at TEXT DEFAULT (datetime('now', 'localtime')),

    -- adminテーブルがないので、usersで代用中。後に変更予定
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE NO ACTION
);

-- 2. シフト提出テーブル (shift_submissions)
CREATE TABLE shift_submissions (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL,
    staff_id TEXT NOT NULL,
    -- shift_requests(start_date) <= submission_date <= shift_requests(end_date)
    submission_date TEXT NOT NULL CHECK (strftime('%Y-%m-%d', submission_date) IS NOT NULL AND submission_date LIKE '____-__-__'),
    submitted_at TEXT DEFAULT (datetime('now', 'localtime')),

    FOREIGN KEY (request_id) REFERENCES shift_requests(id) ON DELETE CASCADE,
    FOREIGN KEY (staff_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (request_id, staff_id, submission_date)
);

-- 3. シフトエントリテーブル (shift_entries)
CREATE TABLE shift_entries (
    id TEXT PRIMARY KEY,
    submission_id TEXT NOT NULL,
    shift_hour INTEGER NOT NULL CHECK (9 <= shift_hour AND shift_hour <= 22),
    created_at TEXT DEFAULT (datetime('now', 'localtime')),

    FOREIGN KEY (submission_id) REFERENCES shift_submissions(id) ON DELETE CASCADE,
    UNIQUE (submission_id, shift_hour)
);

.tables

-- 改善ポイント
-- shift_requests(start_date) <= submission_date <= shift_requests(end_date)
-- update_at
-- user_idの変更ができるかどうか
-- adminテーブルの必要性を考える
