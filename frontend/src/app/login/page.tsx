"use client";
import React, { useState } from "react";

export default function LoginPage() {
    const [loginId, setLoginId] = useState("");
    const [password, setPassword] = useState("");

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-sm p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">ログイン</h1>
                <form className="flex flex-col gap-4">
                    <label className="flex flex-col gap-1">
                        <span className="text-sm">ログインID</span>
                        <input
                            type="text"
                            className="border rounded px-3 py-2"
                            value={loginId}
                            onChange={e => setLoginId(e.target.value)}
                            autoComplete="username"
                            required
                        />
                    </label>
                    <label className="flex flex-col gap-1">
                        <span className="text-sm">パスワード</span>
                        <input
                            type="password"
                            className="border rounded px-3 py-2"
                            value={password}
                            onChange={e => setPassword(e.target.value)}
                            autoComplete="current-password"
                            required
                        />
                    </label>
                    <button
                        type="submit"
                        className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition"
                    >
                        ログイン
                    </button>
                </form>
            </div>
        </div>
    );
}
