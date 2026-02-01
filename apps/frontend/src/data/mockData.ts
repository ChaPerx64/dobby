import type { Period } from '@/types/api';
import type { CategoryItem, ChartDataPoint } from '@/types/dashboard';

/**
 * Mock data for the finance dashboard
 * Period: Feb 5 - Mar 5
 * Budget: 120'000 RSD (12000000 para)
 * Spent: 23'000 RSD (2300000 para)
 * Remaining: 97'000 RSD (9700000 para)
 */

export const mockPeriod: Period = {
  id: '1',
  startDate: '2026-02-05',
  endDate: '2026-03-05',
  isActive: true,
  totalBudget: 12000000, // 120'000.00 RSD
  totalRemaining: 9700000, // 97'000.00 RSD
  totalSpent: 2300000, // 23'000.00 RSD
  projectedEndingBalance: -300000, // -3'000.00 RSD
};

export const mockCategories: CategoryItem[] = [
  {
    id: 'total',
    name: 'Total',
    allocated: 12000000, // 120'000.00
    spent: 2300000, // 23'000.00
    remaining: 9700000, // 97'000.00
  },
  {
    id: 'groceries',
    name: 'Groceries',
    allocated: 5000000, // 50'000.00
    spent: 1200000, // 12'000.00
    remaining: 3800000, // 38'000.00
  },
  {
    id: 'chaian',
    name: "Chaian's pocket money",
    allocated: 3500000, // 35'000.00
    spent: 800000, // 8'000.00
    remaining: 2700000, // 27'000.00
  },
  {
    id: 'sophia',
    name: "Sophia's pocket money",
    allocated: 3500000, // 35'000.00
    spent: 300000, // 3'000.00
    remaining: 3200000, // 32'000.00
  },
];

export const mockChartData: ChartDataPoint[] = [
  { date: '2026-02-05', remaining: 12000000, spent: 0 },
  { date: '2026-02-08', remaining: 11500000, spent: 500000 },
  { date: '2026-02-11', remaining: 11000000, spent: 1000000 },
  { date: '2026-02-14', remaining: 10500000, spent: 1500000 },
  { date: '2026-02-17', remaining: 10000000, spent: 2000000 },
  { date: '2026-02-20', remaining: 9800000, spent: 2200000 },
  { date: '2026-02-23', remaining: 9700000, spent: 2300000 },
  { date: '2026-02-26', remaining: 9700000, spent: 2300000 },
  { date: '2026-03-01', remaining: 9700000, spent: 2300000 },
  { date: '2026-03-05', remaining: 9700000, spent: 2300000 },
];

export const mockProjectedBalance = mockPeriod.projectedEndingBalance!; // -300000 (-3'000.00 RSD)
