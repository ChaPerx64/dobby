import { useState, useEffect } from 'react';
import { Sidebar } from './Sidebar';
import { MetricsPanel } from './MetricsPanel';
import { SpendingChart } from './SpendingChart';
import { apiClient } from '@/api/client';
import type { Period, Allocation } from '@/types/api';
import type { CategoryItem } from '@/types/dashboard';
import { mockChartData } from '@/data/mockData';

export function Dashboard() {
  const [selectedCategory, setSelectedCategory] = useState('total');
  const [period, setPeriod] = useState<Period | null>(null);
  const [allocations, setAllocations] = useState<Allocation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadData() {
      try {
        setLoading(true);
        const { data: periodData, error: periodError } = await apiClient.getCurrentPeriod();
        if (periodError) throw new Error(periodError.message || 'Failed to fetch period');
        if (!periodData) throw new Error('No period data received');
        
        setPeriod(periodData);

        const { data: allocData, error: allocError } = await apiClient.listAllocations(periodData.id);
        if (allocError) throw new Error(allocError.message || 'Failed to fetch allocations');
        
        setAllocations(allocData || []);
      } catch (err: any) {
        console.error(err);
        setError(err.message || 'An error occurred');
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  if (loading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <p className="text-muted-foreground">Loading dashboard...</p>
      </div>
    );
  }

  if (error || !period) {
    return (
      <div className="h-screen flex items-center justify-center">
        <p className="text-destructive">Error: {error}</p>
      </div>
    );
  }

  // Transform data for view
  const totalCategory: CategoryItem = {
    id: 'total',
    name: 'Total',
    allocated: period.totalBudget,
    spent: period.totalSpent,
    remaining: period.totalRemaining,
  };

  const allocationCategories: CategoryItem[] = allocations.map(a => ({
    id: a.envelopeId,
    name: a.envelope,
    allocated: a.amount,
    spent: a.spent,
    remaining: a.remaining,
  }));

  const categories = [totalCategory, ...allocationCategories];

  const currentCategory =
    categories.find((cat) => cat.id === selectedCategory) || categories[0];

  return (
    <div className="h-screen flex">
      <Sidebar
        period={period}
        categories={categories}
        selectedCategory={selectedCategory}
        onSelectCategory={setSelectedCategory}
      />
      <MetricsPanel
        allocated={currentCategory.allocated}
        spent={currentCategory.spent}
        remaining={currentCategory.remaining}
        projectedBalance={period.projectedEndingBalance ?? 0}
      />
      <SpendingChart data={mockChartData} />
    </div>
  );
}