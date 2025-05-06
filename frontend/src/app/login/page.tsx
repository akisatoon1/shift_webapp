"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";

export default function LoginPage() {
    const [loginId, setLoginId] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const router = useRouter();
    const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setLoading(true);
        try {
            const res = await fetch(`${API_BASE_URL}/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({
                    login_id: loginId,
                    password: password,
                }),
            });
            if (res.ok) {
                router.push("/requests");
            } else {
                const data = await res.json();
                setError(data.error || "ログインに失敗しました");
            }
        } catch (err) {
            setError("通信エラーが発生しました");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-sm p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">ログイン</h1>
                <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                    <label className="flex flex-col gap-1">
                        <span className="text-sm">ログインID</span>
                        <input
                            type="text"
                            className="border rounded px-3 py-2"
                            value={loginId}
                            onChange={e => setLoginId(e.target.value)}
                            autoComplete="username"
                            required
                            disabled={loading}
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
                            disabled={loading}
                        />
                    </label>
                    {error && (
                        <div className="text-red-600 text-sm text-center">{error}</div>
                    )}
                    <button
                        type="submit"
                        className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition disabled:opacity-50"
                        disabled={loading}
                    >
                        {loading ? "ログイン中..." : "ログイン"}
                    </button>
                </form>
            </div>
        </div>
    );
}
