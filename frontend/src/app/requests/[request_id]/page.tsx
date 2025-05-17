"use client";
import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";

type Entry = {
    id: number;
    user: {
        id: number;
        name: string;
    };
    date: string;
    hour: number;
};

type RequestEntries = {
    id: number;
    entries: Entry[];
};

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export default function RequestDetailPage() {
    const params = useParams();
    const requestId = params?.request_id;
    const [entries, setEntries] = useState<RequestEntries | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    async function fetchEntries() {
        if (!requestId) return;
        setLoading(true);
        setError("");
        try {
            const res = await fetch(`${API_BASE_URL}/requests/${requestId}/entries`, {
                credentials: "include",
            });
            if (!res.ok) {
                const data = await res.json();
                setError(data.error || "取得に失敗しました");
            } else {
                const data = await res.json();
                setEntries(data);
            }
        } catch (e) {
            setError("通信エラーが発生しました");
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchEntries();
    }, [requestId]);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-2xl p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">リクエスト詳細・エントリー一覧</h1>
                <div className="mb-4 text-center">
                    <a
                        href={`/requests/${requestId}/submit`}
                        className="inline-block bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                    >
                        エントリー提出ページへ
                    </a>
                </div>
                {loading ? (
                    <div className="text-center">読み込み中...</div>
                ) : error ? (
                    <div className="text-red-600 text-center">{error}</div>
                ) : (
                    <table className="w-full border mt-4">
                        <thead>
                            <tr className="bg-gray-100">
                                <th className="border px-2 py-1">エントリーID</th>
                                <th className="border px-2 py-1">ユーザー名</th>
                                <th className="border px-2 py-1">日付</th>
                                <th className="border px-2 py-1">時刻</th>
                            </tr>
                        </thead>
                        <tbody>
                            {entries && Array.isArray(entries.entries) && entries.entries.length > 0 ? (
                                entries.entries.map((entry: Entry) => (
                                    <tr key={entry.id} className="hover:bg-gray-50">
                                        <td className="border px-2 py-1 text-center">{entry.id}</td>
                                        <td className="border px-2 py-1">{entry.user.name}</td>
                                        <td className="border px-2 py-1">{entry.date}</td>
                                        <td className="border px-2 py-1">{entry.hour}</td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td className="border px-2 py-1 text-center" colSpan={4}>エントリーがありません</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
}
