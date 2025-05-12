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
    const [entries, setEntries] = useState<RequestEntries[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [entryDate, setEntryDate] = useState("");
    const [entryHour, setEntryHour] = useState<number | "">("");
    const [createError, setCreateError] = useState("");
    const [createLoading, setCreateLoading] = useState(false);

    const handleCreateEntry = async (e: React.FormEvent) => {
        e.preventDefault();
        setCreateError("");
        setCreateLoading(true);
        try {
            const res = await fetch(`${API_BASE_URL}/requests/${requestId}/entries`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify([
                    {
                        date: entryDate,
                        hour: Number(entryHour),
                    },
                ]),
            });
            if (!res.ok) {
                const data = await res.json();
                setCreateError(data.error || "作成に失敗しました");
            } else {
                setEntryDate("");
                setEntryHour("");
                await fetchEntries();
            }
        } catch (e) {
            setCreateError("通信エラーが発生しました");
        } finally {
            setCreateLoading(false);
        }
    };

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
                setEntries(Array.isArray(data) ? data : []);
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
                <form className="flex flex-col sm:flex-row gap-2 mb-6 items-end" onSubmit={handleCreateEntry}>
                    <div className="flex flex-col">
                        <label className="text-sm">日付</label>
                        <input type="date" className="border rounded px-2 py-1" value={entryDate} onChange={e => setEntryDate(e.target.value)} required disabled={createLoading} />
                    </div>
                    <div className="flex flex-col">
                        <label className="text-sm">時刻</label>
                        <input type="number" className="border rounded px-2 py-1" value={entryHour} onChange={e => setEntryHour(e.target.value === '' ? '' : Number(e.target.value))} min={0} max={23} required disabled={createLoading} />
                    </div>
                    <button type="submit" className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition disabled:opacity-50" disabled={createLoading}>追加</button>
                </form>
                {createError && <div className="text-red-600 text-center mb-2">{createError}</div>}
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
