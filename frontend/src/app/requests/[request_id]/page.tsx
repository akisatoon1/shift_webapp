"use client";
import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";

// TODO: エントリー表示のUI

type Entry = {
    id: number;
    user: {
        id: number;
        name: string;
    };
    date: string;
    hour: number;
};

type RequestDetail = {
    id: number;
    creator: {
        id: number;
        name: string;
    };
    start_date: string;
    end_date: string;
    deadline: string;
    created_at: string;
    entries: Entry[];
};

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export default function RequestDetailPage() {
    const params = useParams();
    const requestId = params?.request_id;
    const [requestDetail, setRequestDetail] = useState<RequestDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    async function fetchRequestDetail() {
        if (!requestId) return;
        setLoading(true);
        setError("");
        try {
            const res = await fetch(`${API_BASE_URL}/requests/${requestId}`, {
                credentials: "include",
            });
            if (!res.ok) {
                const data = await res.json();
                setError(data.error || "取得に失敗しました");
            } else {
                const data = await res.json();
                setRequestDetail(data);
            }
        } catch (e) {
            setError("通信エラーが発生しました");
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchRequestDetail();
    }, [requestId]);

    // 日付範囲をrequestDetailから取得
    let startDate = "";
    let endDate = "";
    if (requestDetail) {
        startDate = requestDetail.start_date;
        endDate = requestDetail.end_date;
    }

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-2xl p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">リクエスト詳細・エントリー一覧</h1>
                {requestDetail && (
                    <div className="mb-6 text-sm text-gray-700">
                        <div><span className="font-semibold">作成者:</span> {requestDetail.creator.name}</div>
                        <div><span className="font-semibold">開始日:</span> {requestDetail.start_date}</div>
                        <div><span className="font-semibold">終了日:</span> {requestDetail.end_date}</div>
                        <div><span className="font-semibold">締切:</span> {requestDetail.deadline}</div>
                        <div><span className="font-semibold">作成日時:</span> {requestDetail.created_at}</div>
                    </div>
                )}
                <div className="mb-4 text-center">
                    <a
                        href="/requests"
                        className="inline-block bg-gray-400 text-white px-4 py-2 rounded hover:bg-gray-500 transition mr-2"
                    >
                        リクエスト一覧へ戻る
                    </a>
                    <a
                        href={`/requests/${requestId}/submit${startDate && endDate ? `?start_date=${startDate}&end_date=${endDate}` : ''}`}
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
                            {requestDetail && Array.isArray(requestDetail.entries) && requestDetail.entries.length > 0 ? (
                                requestDetail.entries.map((entry: Entry) => (
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
