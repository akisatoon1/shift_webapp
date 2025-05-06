"use client";
import React, { useEffect, useState } from "react";

type Request = {
    id: number;
    creator: {
        id: number;
        name: string;
    };
    start_date: string;
    end_date: string;
    deadline: string;
    created_at: string;
};

export default function RequestsPage() {
    const [requests, setRequests] = useState<Request[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        const fetchRequests = async () => {
            setLoading(true);
            setError("");
            try {
                const res = await fetch("http://localhost:3000/api/requests", {
                    credentials: "include",
                });
                if (!res.ok) {
                    const data = await res.json();
                    setError(data.error || "取得に失敗しました");
                } else {
                    const data = await res.json();
                    setRequests(data);
                }
            } catch (e) {
                setError("通信エラーが発生しました");
            } finally {
                setLoading(false);
            }
        };
        fetchRequests();
    }, []);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-2xl p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">リクエスト一覧</h1>
                {loading ? (
                    <div className="text-center">読み込み中...</div>
                ) : error ? (
                    <div className="text-red-600 text-center">{error}</div>
                ) : (
                    <table className="w-full border mt-4">
                        <thead>
                            <tr className="bg-gray-100">
                                <th className="border px-2 py-1">ID</th>
                                <th className="border px-2 py-1">作成者</th>
                                <th className="border px-2 py-1">開始日</th>
                                <th className="border px-2 py-1">終了日</th>
                                <th className="border px-2 py-1">提出期限</th>
                                <th className="border px-2 py-1">作成日</th>
                            </tr>
                        </thead>
                        <tbody>
                            {requests.map((req) => (
                                <tr key={req.id} className="hover:bg-gray-50">
                                    <td className="border px-2 py-1 text-center">{req.id}</td>
                                    <td className="border px-2 py-1">{req.creator.name}</td>
                                    <td className="border px-2 py-1">{req.start_date}</td>
                                    <td className="border px-2 py-1">{req.end_date}</td>
                                    <td className="border px-2 py-1">{req.deadline}</td>
                                    <td className="border px-2 py-1">{req.created_at}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
}
