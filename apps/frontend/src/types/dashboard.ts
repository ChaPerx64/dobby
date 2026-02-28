/**
 * Frontend-only type definitions for the dashboard
 * These are not part of the OpenAPI schema
 */

/**
 * Represents a category in the sidebar with aggregated financial data
 * Can be either the "Total" (from Period) or individual allocations
 */
export interface CategoryItem {
  id: string;
  name: string;
  allocated: number; // in para (currency cents)
  spent: number; // in para (currency cents)
  remaining: number; // in para (currency cents)
}

/**
 * Data point for the spending chart
 */
export interface ChartDataPoint {
  date: string;
  remaining: number; // in para (currency cents)
}
