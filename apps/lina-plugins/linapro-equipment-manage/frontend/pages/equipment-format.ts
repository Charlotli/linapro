/** Shared money formatting helpers for equipment pages. */

/** Format one amount as `¥1,234.56` with grouped thousands. */
export function formatYuan(amount: number | string | null | undefined): string {
  const value = Number(amount ?? 0);
  const safe = Number.isFinite(value) ? value : 0;
  return `¥${safe.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
}
