/**
 * Shared money formatting helpers for the OA approval pages: display
 * formatting with the yuan sign and RMB uppercase words required on finance
 * archive printouts.
 */

const DIGIT_WORDS = [
  '零',
  '壹',
  '贰',
  '叁',
  '肆',
  '伍',
  '陆',
  '柒',
  '捌',
  '玖',
];

const INTEGER_UNITS = ['', '拾', '佰', '仟'];
const GROUP_UNITS = ['', '万', '亿', '兆'];

function integerPartToWords(value: string): string {
  // Split into 4-digit groups from the right. Each 万/亿 unit prints only
  // when its own group carries a nonzero digit; zeros between groups collapse
  // into one 零 connector.
  const groups: string[] = [];
  for (let end = value.length; end > 0; end -= 4) {
    groups.unshift(value.slice(Math.max(0, end - 4), end));
  }

  const digitWords = DIGIT_WORDS;
  const groupToWords = (group: string): { text: string; leadingZero: boolean } => {
    let groupResult = '';
    let zeroPending = false;
    let leadingZero = false;
    for (let i = 0; i < group.length; i += 1) {
      const digit = Number(group[i]);
      const unitIndex = (group.length - 1 - i) % 4;
      if (digit === 0) {
        zeroPending = true;
        continue;
      }
      if (zeroPending) {
        if (groupResult !== '') {
          groupResult += digitWords[0];
        } else {
          leadingZero = true;
        }
      }
      groupResult += digitWords[digit] + INTEGER_UNITS[unitIndex];
      zeroPending = false;
    }
    return { text: groupResult, leadingZero };
  };

  let result = '';
  let zeroCarry = false;
  let previousEndedWithZero = false;
  groups.forEach((group, groupIndex) => {
    const { text: groupResult, leadingZero } = groupToWords(group);
    if (groupResult === '') {
      zeroCarry = true;
      return;
    }
    if (result !== '' && (zeroCarry || previousEndedWithZero || leadingZero)) {
      result += digitWords[0];
    }
    result += groupResult + GROUP_UNITS[groups.length - 1 - groupIndex];
    zeroCarry = false;
    previousEndedWithZero = group[group.length - 1] === '0';
  });
  return result;
}

export function formatYuan(amount: number | string | null | undefined): string {
  const value = Number(amount ?? 0);
  const safe = Number.isFinite(value) ? value : 0;
  return `¥${safe.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
}

/** Convert one amount into formal RMB uppercase words, e.g. 壹仟贰佰叁拾肆元伍角陆分. */
export function formatRmbUppercase(
  amount: number | string | null | undefined,
): string {
  const value = Number(amount ?? 0);
  if (!Number.isFinite(value)) return '';
  const negative = value < 0;
  const abs = Math.abs(value);
  if (abs >= 1e16) return '';

  const [intPart = '0', decimalPart = ''] = abs.toFixed(2).split('.');
  const jiao = Number(decimalPart[0] ?? '0');
  const fen = Number(decimalPart[1] ?? '0');

  let result = '';
  if (abs >= 1) {
    result += integerPartToWords(intPart) + '元';
  } else {
    result += '';
  }

  if (jiao === 0 && fen === 0) {
    result += '整';
  } else {
    if (jiao > 0) {
      result += DIGIT_WORDS[jiao] + '角';
    } else if (fen > 0 && abs >= 1) {
      result += DIGIT_WORDS[0];
    }
    if (fen > 0) {
      result += DIGIT_WORDS[fen] + '分';
    }
    if (jiao === 0 && fen === 0) {
      result += '整';
    }
  }

  return (negative ? '负' : '') + result;
}
