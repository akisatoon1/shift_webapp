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

export default function RequestDetailPage() {
    const params = useParams();
    const requestId = params?.request_id;
    const [entries, setEntries] = useState<RequestEntries[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        if (!requestId) return;
        const fetchEntries = async () => {
            setLoading(true);
            setError("");
            try {
                const res = await fetch(`http://localhost:3000/api/requests/${requestId}/entries`, {
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
        };
        fetchEntries();
    }, [requestId]);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-2xl p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">リクエスト詳細・エントリー一覧</h1>
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
                            {entries.map((entryGroup) =>
                                entryGroup.entries.map((entry) => (
                                    <tr key={entry.id} className="hover:bg-gray-50">
                                        <td className="border px-2 py-1 text-center">{entry.id}</td>
                                        <td className="border px-2 py-1">{entry.user.name}</td>
                                        <td className="border px-2 py-1">{entry.date}</td>
                                        <td className="border px-2 py-1">{entry.hour}</td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
}
