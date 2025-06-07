"use client";
import React, { useState, useEffect } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import dayjs from "dayjs";
import { post, get } from "../../../lib/api";

// TODO: シフト提出更新

// 既存の提出データの型定義
type SubmissionEntry = {
    id: number;
    date: string;
    hour: number;
};

type Submission = {
    id: number;
    user: {
        id: number;
        name: string;
    };
    entries: SubmissionEntry[];
};

type SubmissionResponse = {
    submission: Submission | null;
};

export default function EntrySubmitPage() {
    const params = useParams();
    const router = useRouter();
    const searchParams = useSearchParams();
    const requestId = params?.request_id;
    const startDate = searchParams.get("start_date");
    const endDate = searchParams.get("end_date");
    const [submitError, setSubmitError] = useState("");
    const [submitLoading, setSubmitLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    // 提出済みデータ
    const [existingSubmission, setExistingSubmission] = useState<Submission | null>(null);
    const [loadingSubmission, setLoadingSubmission] = useState(true);
    // カレンダー用: 今月の日付リストを生成
    const today = dayjs();
    const year = today.year();
    const month = today.month(); // 0-indexed
    const daysInMonth = today.daysInMonth();
    // 0-23時のリスト
    const hourList = Array.from({ length: 24 }, (_, i) => i);
    // 選択状態: { [date: string]: Set<number> }
    const [selected, setSelected] = useState<{ [date: string]: Set<number> }>({});

    // start_date, end_date から日付リストを生成
    const dateList = React.useMemo(() => {
        if (!startDate || !endDate) return [];
        const start = dayjs(startDate);
        const end = dayjs(endDate);
        const dates = [];
        let d = start;
        while (d.isBefore(end) || d.isSame(end, 'day')) {
            dates.push(d.format("YYYY-MM-DD"));
            d = d.add(1, "day");
        }
        return dates;
    }, [startDate, endDate]);

    // start_date, end_dateがない場合のエラー表示
    const paramError = !startDate || !endDate;

    // 既存の提出データを取得
    useEffect(() => {
        if (!requestId) return;

        async function fetchSubmission() {
            setLoadingSubmission(true);
            try {
                const res = await get(`/requests/${requestId}/submissions/mine`);
                if (res && res.ok) {
                    const data: SubmissionResponse = await res.json();
                    setExistingSubmission(data.submission);

                    // 提出データが存在する場合、そのデータで選択状態を初期化
                    if (data.submission) {
                        const newSelected: { [date: string]: Set<number> } = {};
                        data.submission.entries.forEach(entry => {
                            if (!newSelected[entry.date]) {
                                newSelected[entry.date] = new Set<number>();
                            }
                            newSelected[entry.date].add(entry.hour);
                        });
                        setSelected(newSelected);
                    }
                }
            } catch (e) {
                console.error("提出データの取得に失敗しました:", e);
            } finally {
                setLoadingSubmission(false);
            }
        }

        fetchSubmission();
    }, [requestId]);

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
        try {
            const res = await post(`/requests/${requestId}/submissions`, submitEntries);
            if (res && res.ok) {
                setSuccess(true);
                setSelected({});
                setTimeout(() => {
                    router.push(`/requests/${requestId}`);
                }, 1200);
            } else if (res) {
                const data = await res.json();
                setSubmitError(data.error || "提出に失敗しました");
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
                        className="inline-block bg-gray-400 text-white px-4 py-2 rounded hover:bg-gray-500 transition"
                    >
                        エントリー一覧ページへ戻る
                    </a>
                </div>
                {paramError ? (
                    <div className="text-red-600 text-center text-sm mb-2">日付範囲が指定されていません（start_date, end_dateパラメータ必須）</div>
                ) : loadingSubmission ? (
                    <div className="text-center py-4">読み込み中...</div>
                ) : (
                    <form onSubmit={handleSubmit}>
                        <div className="overflow-x-auto">
                            {existingSubmission && (
                                <div className="bg-blue-50 p-4 mb-4 rounded border border-blue-200">
                                    <p className="font-semibold text-blue-700 text-center">
                                        すでに提出済みです（提出ID: {existingSubmission.id}）
                                    </p>
                                </div>
                            )}
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
                                                        className={`border px-1 py-1 text-center select-none transition w-8 h-8 align-middle ${isSelected
                                                            ? 'bg-blue-400 text-white'
                                                            : 'bg-white hover:bg-blue-100'
                                                            } ${existingSubmission
                                                                ? 'cursor-default'
                                                                : 'cursor-pointer'
                                                            }`}
                                                        style={{ minWidth: 32, maxWidth: 32, width: 32, minHeight: 32, height: 32 }}
                                                        onClick={() => !existingSubmission && !submitLoading && handleCellClick(date, hour)}
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
                        {!existingSubmission ? (
                            <button
                                type="submit"
                                className="bg-blue-600 text-white rounded px-4 py-2 font-semibold hover:bg-blue-700 transition disabled:opacity-50 w-full"
                                disabled={submitLoading}
                            >
                                {submitLoading ? "提出中..." : "まとめて提出"}
                            </button>
                        ) : (
                            <div className="bg-gray-300 text-gray-600 rounded px-4 py-2 font-semibold w-full text-center">
                                提出済み
                            </div>
                        )}
                    </form>
                )}
            </div>
        </div>
    );
}
