import { useState, useEffect } from 'react';
import { Sidebar } from './Sidebar';
import { MetricsPanel } from './MetricsPanel';
import { SpendingChart } from './SpendingChart';
import { TransactionList } from './TransactionList';
import { apiClient } from '@/api/client';
import type { Period, Transaction } from '@/types/api';
import type { CategoryItem, ChartDataPoint } from '@/types/dashboard';

export function Dashboard() {
  const [selectedCategory, setSelectedCategory] = useState('total');
  const [activeTab, setActiveTab] = useState<'balance' | 'transactions'>('balance');
  const [period, setPeriod] = useState<Period | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadData = async () => {
    try {
      const { data: periodData, error: periodError } = await apiClient.getCurrentPeriod();
      if (periodError) throw new Error(periodError.message || 'Failed to fetch period');
      if (!periodData) throw new Error('No period data received');
      
      setPeriod(periodData);

      const { data: txData, error: txError } = await apiClient.listTransactions(periodData.id);
      if (txError) throw new Error(txError.message || 'Failed to fetch transactions');
      if (txData) {
        setTransactions(txData);
      }
    } catch (err: unknown) {
      console.error(err);
      setError(err instanceof Error ? err.message : 'An error occurred');
    }
  };

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
    async function init() {
      setLoading(true);
      await loadData();
      setLoading(false);
    }
    init();
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

  const chartData: ChartDataPoint[] = [];
  let filteredTx: Transaction[] = [];
  let initialBalance = 0;

  if (period) {
    filteredTx = selectedCategory === 'total'
      ? transactions
      : transactions.filter(tx => tx.envelopeId === selectedCategory);

    const txByDay = new Map<string, { allocated: number; spent: number }>();
    let totalAllocatedInPeriod = 0;

    for (const tx of filteredTx) {
      const dayStr = tx.date.split('T')[0];
      const existing = txByDay.get(dayStr) || { allocated: 0, spent: 0 };
      if (tx.amount > 0) {
        existing.allocated += tx.amount;
        totalAllocatedInPeriod += tx.amount;
      } else {
        existing.spent += Math.abs(tx.amount);
      }
      txByDay.set(dayStr, existing);
    }

    initialBalance = currentCategory.allocated - totalAllocatedInPeriod;
    let cumulativeAllocated = initialBalance;
    let cumulativeSpent = 0;

    const startDate = new Date(period.startDate);
    const endDate = new Date(period.endDate);
    const today = new Date();
    const chartEndDate = today < endDate ? today : endDate;

    const utcCurrent = new Date(Date.UTC(startDate.getFullYear(), startDate.getMonth(), startDate.getDate()));
    const utcEnd = new Date(Date.UTC(chartEndDate.getFullYear(), chartEndDate.getMonth(), chartEndDate.getDate()));

    while (utcCurrent <= utcEnd) {
      const dateStr = utcCurrent.toISOString().split('T')[0];

      const dailyTx = txByDay.get(dateStr) || { allocated: 0, spent: 0 };
      cumulativeAllocated += dailyTx.allocated;
      cumulativeSpent += dailyTx.spent;

      chartData.push({
        date: dateStr,
        remaining: cumulativeAllocated - cumulativeSpent,
      });

      utcCurrent.setUTCDate(utcCurrent.getUTCDate() + 1);
    }
  }

  return (
    <div className="h-screen flex">
      <Sidebar
        period={period}
        categories={categories}
        selectedCategory={selectedCategory}
        onSelectCategory={setSelectedCategory}
        onEnvelopeCreated={handleEnvelopeCreated}
        onAllocationCreated={loadData}
      />
      <MetricsPanel
        allocated={currentCategory.allocated}
        spent={currentCategory.spent}
        remaining={currentCategory.remaining}
        projectedBalance={period.projectedEndingBalance ?? 0}
      />
      <div className="flex-1 flex flex-col bg-background">
        <div className="flex border-b border-border">
          <button
            className={`px-6 py-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'balance'
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
            }`}
            onClick={() => setActiveTab('balance')}
          >
            Balance Over Time
          </button>
          <button
            className={`px-6 py-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'transactions'
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground hover:border-border'
            }`}
            onClick={() => setActiveTab('transactions')}
          >
            Transactions
          </button>
        </div>
        <div className="flex-1 p-6 overflow-hidden flex flex-col">
          {activeTab === 'balance' ? (
            <SpendingChart data={chartData} />
          ) : (
            <TransactionList transactions={filteredTx} initialBalance={initialBalance} />
          )}
        </div>
      </div>
    </div>
  );
}
