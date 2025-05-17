"use client";
import React, { useState } from "react";
import { useParams, useRouter } from "next/navigation";

// TODO: エントリー提出をカレンダー形式で行う

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export default function EntrySubmitPage() {
    const params = useParams();
    const router = useRouter();
    const requestId = params?.request_id;
    const [entryDate, setEntryDate] = useState("");
    const [entryHour, setEntryHour] = useState<number | "">("");
    const [submitError, setSubmitError] = useState("");
    const [submitLoading, setSubmitLoading] = useState(false);
    const [success, setSuccess] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitError("");
        setSubmitLoading(true);
        setSuccess(false);
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
                setSubmitError(data.error || "提出に失敗しました");
            } else {
                setSuccess(true);
                setEntryDate("");
                setEntryHour("");
                setTimeout(() => {
                    router.push(`/requests/${requestId}`);
                }, 1200);
            }
        } catch (e) {
            setSubmitError("通信エラーが発生しました");
        } finally {
            setSubmitLoading(false);
        }
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-md p-8 bg-white rounded shadow-md">
                <h1 className="text-2xl font-bold mb-6 text-center">エントリー提出フォーム</h1>
                <div className="mb-4 text-center">
                    <a
                        href={`/requests/${requestId}`}
                        className="inline-block bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                    >
                        エントリー一覧ページへ戻る
                    </a>
                </div>
                <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                    <div className="flex flex-col">
                        <label className="text-sm">日付</label>
                        <input
                            type="date"
                            className="border rounded px-2 py-1"
                            value={entryDate}
                            onChange={e => setEntryDate(e.target.value)}
                            required
                            disabled={submitLoading}
                        />
                    </div>
                    <div className="flex flex-col">
                        <label className="text-sm">時刻</label>
                        <input
                            type="number"
                            className="border rounded px-2 py-1"
                            value={entryHour}
                            onChange={e => setEntryHour(e.target.value === '' ? '' : Number(e.target.value))}
                            min={0}
                            max={23}
                            required
                            disabled={submitLoading}
                        />
                    </div>
                    {submitError && (
                        <div className="text-red-600 text-center text-sm">{submitError}</div>
                    )}
                    {success && (
                        <div className="text-green-600 text-center text-sm">提出が完了しました。リダイレクトします...</div>
                    )}
                    <button
                        type="submit"
                        className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition disabled:opacity-50"
                        disabled={submitLoading}
                    >
                        {submitLoading ? "提出中..." : "提出"}
                    </button>
                </form>
            </div>
        </div>
    );
}
