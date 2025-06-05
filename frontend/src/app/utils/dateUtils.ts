// dateUtils.ts
// 日付関連のユーティリティ関数
import dayjs from 'dayjs';

/**
 * 開始日と終了日の間のすべての日付を含む配列を返す
 */
export function getDateRange(startDate: string, endDate: string): string[] {
    if (!startDate || !endDate) return [];

    const dateList: string[] = [];
    const start = dayjs(startDate);
    const end = dayjs(endDate);

    let current = start;
    while (current.isBefore(end) || current.isSame(end, 'day')) {
        dateList.push(current.format('YYYY-MM-DD'));
        current = current.add(1, 'day');
    }

    return dateList;
}

/**
 * 日時フォーマット：日付と時間を結合
 */
export function formatDatetime(date: string, time: string): string {
    return `${date} ${time}:00`;
}

/**
 * 人が読みやすい形式に日付をフォーマット
 */
export function formatDateForDisplay(dateStr: string): string {
    return dayjs(dateStr).format('YYYY年MM月DD日');
}

/**
 * 時間の配列を生成（0-23）
 */
export function getHourRange(): number[] {
    return Array.from({ length: 24 }, (_, i) => i);
}
