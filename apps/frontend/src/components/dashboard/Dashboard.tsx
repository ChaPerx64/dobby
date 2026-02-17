import { useState, useEffect } from 'react';
import { Sidebar } from './Sidebar';
import { MetricsPanel } from './MetricsPanel';
import { SpendingChart } from './SpendingChart';
import { apiClient } from '@/api/client';
import type { Period } from '@/types/api';
import type { CategoryItem } from '@/types/dashboard';

export function Dashboard() {
  const [selectedCategory, setSelectedCategory] = useState('total');
  const [period, setPeriod] = useState<Period | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleEnvelopeCreated = (newEnvelope: { id: string; name: string }) => {
    if (!period) return;

    // Check if it already exists in the summaries to avoid duplicates
    if (period.envelopeSummaries.some(s => s.envelopeId === newEnvelope.id)) return;

    const newSummary = {
      envelopeId: newEnvelope.id,
      envelopeName: newEnvelope.name,
      amount: 0,
      spent: 0,
      remaining: 0,
    };

    setPeriod({
      ...period,
      envelopeSummaries: [...period.envelopeSummaries, newSummary],
    });
  };

  useEffect(() => {
    async function loadData() {
      try {
        setLoading(true);
        const { data: periodData, error: periodError } = await apiClient.getCurrentPeriod();
        if (periodError) throw new Error(periodError.message || 'Failed to fetch period');
        if (!periodData) throw new Error('No period data received');
        
        setPeriod(periodData);
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

  const allocationCategories: CategoryItem[] = period.envelopeSummaries.map(s => ({
    id: s.envelopeId,
    name: s.envelopeName,
    allocated: s.amount,
    spent: s.spent,
    remaining: s.remaining,
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
        onEnvelopeCreated={handleEnvelopeCreated}
      />
      <MetricsPanel
        allocated={currentCategory.allocated}
        spent={currentCategory.spent}
        remaining={currentCategory.remaining}
        projectedBalance={period.projectedEndingBalance ?? 0}
      />
      <SpendingChart data={[]} />
    </div>
  );
}
