"use client";
import React, { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import dayjs from "dayjs";

// TODO: api設計を変える(entries周り)
// TODO: 日付を正しく列挙(提出ページのパラメータで渡す)
// TODO: 更新や削除はどうする？

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export default function EntrySubmitPage() {
    const params = useParams();
    const router = useRouter();
    const requestId = params?.request_id;
    const [submitError, setSubmitError] = useState("");
    const [submitLoading, setSubmitLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    // カレンダー用: 今月の日付リストを生成
    const today = dayjs();
    const year = today.year();
    const month = today.month(); // 0-indexed
    const daysInMonth = today.daysInMonth();
    // 0-23時のリスト
    const hourList = Array.from({ length: 24 }, (_, i) => i);
    // 選択状態: { [date: string]: Set<number> }
    const [selected, setSelected] = useState<{ [date: string]: Set<number> }>({});

    const dateList = Array.from({ length: daysInMonth }, (_, i) =>
        dayjs(new Date(year, month, i + 1)).format("YYYY-MM-DD")
    );

    const handleCellClick = (date: string, hour: number) => {
        setSelected(prev => {
            const prevSet = prev[date] ? new Set<number>(prev[date]) : new Set<number>();
            if (prevSet.has(hour)) {
                prevSet.delete(hour);
            } else {
                prevSet.add(hour);
            }
            return { ...prev, [date]: prevSet };
        });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitError("");
        setSubmitLoading(true);
        setSuccess(false);
        // 選択された日付・時刻のペアを抽出
        const submitEntries = Object.entries(selected)
            .flatMap(([date, hours]) =>
                Array.from(hours).map(hour => ({ date, hour }))
            );
        if (submitEntries.length === 0) {
            setSubmitError("1つ以上のセルを選択してください");
            setSubmitLoading(false);
            return;
        }
        try {
            const res = await fetch(`${API_BASE_URL}/requests/${requestId}/entries`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(submitEntries),
            });
            if (!res.ok) {
                const data = await res.json();
                setSubmitError(data.error || "提出に失敗しました");
            } else {
                setSuccess(true);
                setSelected({});
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
            <div className="w-full max-w-[80vw] p-4 sm:p-8 bg-white rounded shadow-md overflow-x-auto mx-auto">
                <h1 className="text-2xl font-bold mb-6 text-center">エントリー提出フォーム</h1>
                <div className="mb-4 text-center">
                    <a
                        href={`/requests/${requestId}`}
                        className="inline-block bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                    >
                        エントリー一覧ページへ戻る
                    </a>
                </div>
                <form onSubmit={handleSubmit}>
                    <div className="overflow-x-auto">
                        <table className="border mb-4 min-w-max">
                            <thead>
                                <tr className="bg-gray-100">
                                    <th className="border px-2 py-1 w-28">日付</th>
                                    {hourList.map(hour => (
                                        <th
                                            key={hour}
                                            className="border px-1 py-1 text-xs w-16"
                                            style={{ minWidth: 48, maxWidth: 64, width: 48 }}
                                        >
                                            {hour}~{hour + 1}
                                        </th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody>
                                {dateList.map(date => (
                                    <tr key={date}>
                                        <td className="border px-2 py-1 text-center text-xs">{date}</td>
                                        {hourList.map(hour => {
                                            const isSelected = selected[date]?.has(hour);
                                            return (
                                                <td
                                                    key={hour}
                                                    className={`border px-1 py-1 text-center cursor-pointer select-none transition w-8 h-8 align-middle ${isSelected ? 'bg-blue-400 text-white' : 'bg-white hover:bg-blue-100'}`}
                                                    style={{ minWidth: 32, maxWidth: 32, width: 32, minHeight: 32, height: 32 }}
                                                    onClick={() => !submitLoading && handleCellClick(date, hour)}
                                                >
                                                    <span style={{ display: 'inline-block', width: '1em', textAlign: 'center' }}>
                                                        {isSelected ? '●' : '\u00A0'}
                                                    </span>
                                                </td>
                                            );
                                        })}
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                    {submitError && (
                        <div className="text-red-600 text-center text-sm mb-2">{submitError}</div>
                    )}
                    {success && (
                        <div className="text-green-600 text-center text-sm mb-2">提出が完了しました。リダイレクトします...</div>
                    )}
                    <button
                        type="submit"
                        className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition disabled:opacity-50 w-full"
                        disabled={submitLoading}
                    >
                        {submitLoading ? "提出中..." : "まとめて提出"}
                    </button>
                </form>
            </div>
        </div>
    );
}
