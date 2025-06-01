"use client";
import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { get, post } from "../lib/api";

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
    const [startDate, setStartDate] = useState("");
    const [endDate, setEndDate] = useState("");
    const [deadline, setDeadline] = useState("");
    const [deadlineTime, setDeadlineTime] = useState("00:00");
    const [createError, setCreateError] = useState("");
    const [createLoading, setCreateLoading] = useState(false);
    const [userRoles, setUserRoles] = useState<string[]>([]);
    const [userLoaded, setUserLoaded] = useState(false);
    const router = useRouter();

    // ユーザ情報取得
    useEffect(() => {
        async function fetchSession() {
            try {
                const res = await get(`/session`);
                if (res && res.ok) {
                    const data = await res.json();
                    setUserRoles(data.user?.roles || []);
                } else {
                    setUserRoles([]);
                }
            } catch {
                setUserRoles([]);
            } finally {
                setUserLoaded(true);
            }
        }
        fetchSession();
    }, []);

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();
        setCreateError("");
        setCreateLoading(true);
        try {
            // Format deadline as "yyyy-mm-dd HH:MM:SS"
            const formattedDeadline = `${deadline} ${deadlineTime}:00`;

            const res = await post(`/requests`, {
                start_date: startDate,
                end_date: endDate,
                deadline: formattedDeadline,
            });

            if (res && res.ok) {
                setStartDate("");
                setEndDate("");
                setDeadline("");
                setDeadlineTime("00:00");
                await fetchRequests();
            } else if (res) {
                const data = await res.json();
                setCreateError(data.error || "作成に失敗しました");
            }
        } catch (e) {
            setCreateError("通信エラーが発生しました");
        } finally {
            setCreateLoading(false);
        }
    };

    async function fetchRequests() {
        setLoading(true);
        setError("");
        try {
            const res = await get(`/requests`);
            if (res && res.ok) {
                const data = await res.json();
                setRequests(data);
            } else if (res) {
                const data = await res.json();
                setError(data.error || "取得に失敗しました");
            }
        } catch (e) {
            setError("通信エラーが発生しました");
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchRequests();
    }, []);

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-2xl p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">リクエスト一覧</h1>
                {/* manager権限のみ追加フォームを表示 */}
                {userLoaded && userRoles.includes("manager") && (
                    <form className="flex flex-col sm:flex-row gap-2 mb-6 items-end" onSubmit={handleCreate}>
                        <div className="flex flex-col">
                            <label className="text-sm">開始日</label>
                            <input type="date" className="border rounded px-2 py-1" value={startDate} onChange={e => setStartDate(e.target.value)} required disabled={createLoading} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm">終了日</label>
                            <input type="date" className="border rounded px-2 py-1" value={endDate} onChange={e => setEndDate(e.target.value)} required disabled={createLoading} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm">提出期限（日付）</label>
                            <input type="date" className="border rounded px-2 py-1" value={deadline} onChange={e => setDeadline(e.target.value)} required disabled={createLoading} />
                        </div>
                        <div className="flex flex-col">
                            <label className="text-sm">提出期限（時間）</label>
                            <input type="time" className="border rounded px-2 py-1" value={deadlineTime} onChange={e => setDeadlineTime(e.target.value)} required disabled={createLoading} />
                        </div>
                        <button type="submit" className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition disabled:opacity-50" disabled={createLoading}>追加</button>
                    </form>
                )}
                {createError && userRoles.includes("manager") && <div className="text-red-600 text-center mb-2">{createError}</div>}
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
                                <tr
                                    key={req.id}
                                    className="hover:bg-blue-50 cursor-pointer"
                                    onClick={() => router.push(`/requests/${req.id}`)}
                                >
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
