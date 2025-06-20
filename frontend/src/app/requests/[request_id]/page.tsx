"use client";
import React, { useEffect, useState } from "react";
import { useParams, useSearchParams } from "next/navigation";
import dayjs from "dayjs";
import { get } from "../../lib/api";

type Entry = {
    id: number;
    user: {
        id: number;
        name: string;
    };
    date: string;
    hour: number;
};

type Submission = {
    submitter: {
        id: number;
        name: string;
    };
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
    submissions: Submission[];
    entries: Entry[];
};

export default function RequestDetailPage() {
    const params = useParams();
    const searchParams = useSearchParams();
    const requestId = params?.request_id;
    const [requestDetail, setRequestDetail] = useState<RequestDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [userRoles, setUserRoles] = useState<string[]>([]);
    const [userLoaded, setUserLoaded] = useState(false);

    async function fetchRequestDetail() {
        if (!requestId) return;
        setLoading(true);
        setError("");
        try {
            const res = await get(`/requests/${requestId}`);
            if (res && res.ok) {
                const data = await res.json();
                setRequestDetail(data);
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

    useEffect(() => {
        fetchRequestDetail();
    }, [requestId]);

    // 日付範囲をrequestDetailから取得
    let startDate = "";
    let endDate = "";
    let dateList: string[] = [];
    if (requestDetail) {
        startDate = requestDetail.start_date;
        endDate = requestDetail.end_date;
        // start_date, end_date から日付リストを生成
        const start = dayjs(startDate);
        const end = dayjs(endDate);
        let d = start;
        while (d.isBefore(end) || d.isSame(end, 'day')) {
            dateList.push(d.format("YYYY-MM-DD"));
            d = d.add(1, "day");
        }
    }
    
    // ユーザーごとにentriesをグループ化
    const userEntriesMap: { [userId: number]: { user: Entry["user"]; entries: Entry[] } } = {};
    if (requestDetail && Array.isArray(requestDetail.entries)) {
        requestDetail.entries.forEach((entry) => {
            if (!userEntriesMap[entry.user.id]) {
                userEntriesMap[entry.user.id] = { user: entry.user, entries: [] };
            }
            userEntriesMap[entry.user.id].entries.push(entry);
        });
    }

    // 提出者（submitters）の一覧を作成
    const submitterMap: { [submitterId: number]: { submitter: Submission["submitter"]; entries: Entry[] } } = {};
    if (requestDetail) {
        // すべての提出者を追加
        if (Array.isArray(requestDetail.submissions)) {
            requestDetail.submissions.forEach((submission) => {
                if (!submitterMap[submission.submitter.id]) {
                    submitterMap[submission.submitter.id] = {
                        submitter: submission.submitter,
                        entries: []
                    };
                }
            });
        }
        
        // 各提出者のエントリーを追加
        if (Array.isArray(requestDetail.entries)) {
            requestDetail.entries.forEach((entry) => {
                if (submitterMap[entry.user.id]) {
                    submitterMap[entry.user.id].entries.push(entry);
                } else {
                    // エントリーがあるのに提出者リストにない場合（通常はないはず）
                    submitterMap[entry.user.id] = {
                        submitter: { id: entry.user.id, name: entry.user.name },
                        entries: [entry]
                    };
                }
            });
        }
    }
    
    // 0-23時のリスト
    const hourList = Array.from({ length: 24 }, (_, i) => i);

    // 表示モード: "user" or "date" をパラメータから取得
    const paramView = searchParams.get("view");
    const [viewMode, setViewMode] = useState<'user' | 'date'>(paramView === 'date' ? 'date' : 'user');

    // ページ遷移せずにURLパラメータのみ変更
    const handleChangeView = (mode: 'user' | 'date') => {
        const url = new URL(window.location.href);
        url.searchParams.set('view', mode);
        window.history.replaceState(null, '', url.toString());
        setViewMode(mode);
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
            <div className="w-full max-w-[80vw] p-8 bg-white rounded shadow-md">
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
                    {userLoaded && userRoles.includes("employee") && (
                        <a
                            href={`/requests/${requestId}/submit${startDate && endDate ? `?start_date=${startDate}&end_date=${endDate}` : ''}`}
                            className="inline-block bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition mr-2"
                        >
                            エントリー提出ページへ
                        </a>
                    )}
                    {/* 表示切り替えボタン */}
                    <button
                        className={`inline-block px-4 py-2 rounded font-semibold border ml-2 ${viewMode === 'user' ? 'bg-blue-600 text-white border-blue-700' : 'bg-white text-blue-700 border-blue-400'}`}
                        onClick={() => handleChangeView('user')}
                    >
                        ユーザーごと
                    </button>
                    <button
                        className={`inline-block px-4 py-2 rounded font-semibold border ml-2 ${viewMode === 'date' ? 'bg-blue-600 text-white border-blue-700' : 'bg-white text-blue-700 border-blue-400'}`}
                        onClick={() => handleChangeView('date')}
                    >
                        日付ごと
                    </button>
                </div>
                {loading ? (
                    <div className="text-center">読み込み中...</div>
                ) : error ? (
                    <div className="text-red-600 text-center">{error}</div>
                ) : (
                    <div>
                        {viewMode === 'user' ? (
                            // --- 提出者ごと表示 ---
                            Object.values(submitterMap).length === 0 ? (
                                <div className="text-center text-gray-500">提出者がありません</div>
                            ) : (
                                Object.values(submitterMap).map(({ submitter, entries }) => (
                                    <div key={submitter.id} className="mb-8">
                                        <div className="font-bold mb-2 text-lg">{submitter.name} さんのエントリー</div>
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
                                                                const hasEntry = entries.some(e => e.date === date && e.hour === hour);
                                                                return (
                                                                    <td
                                                                        key={hour}
                                                                        className={`border px-1 py-1 text-center select-none w-8 h-8 align-middle ${hasEntry ? 'bg-blue-400 text-white font-bold' : 'bg-white'}`}
                                                                        style={{ minWidth: 32, maxWidth: 32, width: 32, minHeight: 32, height: 32 }}
                                                                    >
                                                                        <span style={{ display: 'inline-block', width: '1em', textAlign: 'center' }}>
                                                                            {hasEntry ? '●' : '\u00A0'}
                                                                        </span>
                                                                    </td>
                                                                );
                                                            })}
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>
                                ))
                            )
                        ) : (
                            // --- 日付ごとに表を分ける表示 ---
                            dateList.length === 0 ? (
                                <div className="text-center text-gray-500">エントリーがありません</div>
                            ) : (
                                dateList.map(date => {
                                    // 提出者一覧
                                    const submitters = Object.values(submitterMap).map(({ submitter }) => submitter);
                                    return (
                                        <div key={date} className="mb-8">
                                            <div className="font-bold mb-2 text-lg">{date} のエントリー</div>
                                            <div className="overflow-x-auto">
                                                <table className="border mb-4 min-w-max">
                                                    <thead>
                                                        <tr className="bg-gray-100">
                                                            <th className="border px-2 py-1 w-28">提出者</th>
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
                                                        {submitters.map(submitter => {
                                                            // その提出者のその日付のエントリー
                                                            const entries = requestDetail?.entries.filter(e => e.user.id === submitter.id && e.date === date) || [];
                                                            return (
                                                                <tr key={submitter.id}>
                                                                    <td className="border px-2 py-1 text-center text-xs">{submitter.name}</td>
                                                                    {hourList.map(hour => {
                                                                        const hasEntry = entries.some(e => e.hour === hour);
                                                                        return (
                                                                            <td
                                                                                key={hour}
                                                                                className={`border px-1 py-1 text-center select-none w-8 h-8 align-middle ${hasEntry ? 'bg-blue-400 text-white font-bold' : 'bg-white'}`}
                                                                                style={{ minWidth: 32, maxWidth: 32, width: 32, minHeight: 32, height: 32 }}
                                                                            >
                                                                                <span style={{ display: 'inline-block', width: '1em', textAlign: 'center' }}>
                                                                                    {hasEntry ? '●' : '\u00A0'}
                                                                                </span>
                                                                            </td>
                                                                        );
                                                                    })}
                                                                </tr>
                                                            );
                                                        })}
                                                    </tbody>
                                                </table>
                                            </div>
                                        </div>
                                    );
                                })
                            )
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
