import { useState } from 'react';
import { Sidebar } from './Sidebar';
import { MetricsPanel } from './MetricsPanel';
import { SpendingChart } from './SpendingChart';
import {
  mockPeriod,
  mockCategories,
  mockChartData,
  mockProjectedBalance,
} from '@/data/mockData';

export function Dashboard() {
  const [selectedCategory, setSelectedCategory] = useState('total');

  // Find the selected category data
  const currentCategory =
    mockCategories.find((cat) => cat.id === selectedCategory) ||
    mockCategories[0];

  return (
    <div className="h-screen flex">
      <Sidebar
        period={mockPeriod}
        categories={mockCategories}
        selectedCategory={selectedCategory}
        onSelectCategory={setSelectedCategory}
      />
      <MetricsPanel
        allocated={currentCategory.allocated}
        spent={currentCategory.spent}
        remaining={currentCategory.remaining}
        projectedBalance={mockProjectedBalance}
      />
      <SpendingChart data={mockChartData} />
    </div>
  );
}
