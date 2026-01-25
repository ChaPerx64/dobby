/**
 * Format para (1/100 RSD) to RSD with apostrophe thousands separator
 * @param para - Amount in para (e.g., 12000000 = 120'000.00 RSD)
 * @returns Formatted string (e.g., "120'000.00")
 */
export function formatCurrency(para: number): string {
  const rsd = para / 100;
  const formatted = rsd.toFixed(2);
  const [integer, decimal] = formatted.split('.');
  // Add apostrophe thousands separator
  const withSeparator = integer.replace(/\B(?=(\d{3})+(?!\d))/g, "'");
  return `${withSeparator}.${decimal}`;
}

/**
 * Format date range for display
 * @param start - Start date (YYYY-MM-DD)
 * @param end - End date (YYYY-MM-DD)
 * @returns Formatted string (e.g., "February, 5th - March, 5th")
 */
export function formatDateRange(start: string, end: string): string {
  const startDate = new Date(start);
  const endDate = new Date(end);
  const options: Intl.DateTimeFormatOptions = { month: 'long', day: 'numeric' };
  const startStr = startDate.toLocaleDateString('en-US', options);
  const endStr = endDate.toLocaleDateString('en-US', options);
  return `${startStr.replace(',', ', ')} - ${endStr.replace(',', ', ')}`;
}

/**
 * Get month name from date string
 * @param dateStr - Date string (YYYY-MM-DD)
 * @returns Month name (e.g., "February")
 */
export function getMonthName(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'long' });
}

/**
 * Convert currency amount to para
 * @param amount - Amount in currency units (e.g., 25.50)
 * @returns Amount in para (e.g., 2550)
 */
export function parseMoney(amount: number): number {
  return Math.round(amount * 100);
}
