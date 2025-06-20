"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { post } from "../lib/api";
import { Button } from "../components/ui";

// TODO: ログイン失敗時の挙動おかしい

export default function LoginPage() {
    const [loginId, setLoginId] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setLoading(true);
        try {
            const res = await post(`/login`, {
                login_id: loginId,
                password: password,
            });

            if (res && res.ok) {
                router.push("/requests");
            } else if (res) {
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
                            className="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            value={loginId}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setLoginId(e.target.value)}
                            autoComplete="username"
                            required
                            disabled={loading}
                        />
                    </label>
                    <label className="flex flex-col gap-1">
                        <span className="text-sm">パスワード</span>
                        <input
                            type="password"
                            className="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            value={password}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                            autoComplete="current-password"
                            required
                            disabled={loading}
                        />
                    </label>
                    {error && (
                        <div className="text-red-600 text-sm text-center">{error}</div>
                    )}
                    <Button
                        type="submit"
                        disabled={loading}
                        fullWidth={true}
                    >
                        {loading ? "ログイン中..." : "ログイン"}
                    </Button>
                </form>
            </div>
        </div>
    );
}
