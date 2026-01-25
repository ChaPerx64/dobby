import { Separator } from '@/components/ui/separator';
import type { CategoryItem } from '@/types/dashboard';
import { getMonthName, formatDateRange } from '@/lib/format';

interface SidebarProps {
  period: { startDate: string; endDate: string };
  categories: CategoryItem[];
  selectedCategory: string;
  onSelectCategory: (categoryId: string) => void;
}

export function Sidebar({
  period,
  categories,
  selectedCategory,
  onSelectCategory,
}: SidebarProps) {
  const monthName = getMonthName(period.startDate);
  const dateRange = formatDateRange(period.startDate, period.endDate);

  return (
    <div className="w-64 bg-card border-r border-border p-6 flex flex-col">
      <div className="mb-4">
        <h1 className="text-2xl font-bold text-foreground">{monthName}</h1>
        <p className="text-sm text-muted-foreground">{dateRange}</p>
      </div>

      <Separator className="mb-4" />

      <nav className="flex-1">
        <ul className="space-y-1">
          {categories.map((category) => (
            <li key={category.id}>
              <button
                onClick={() => onSelectCategory(category.id)}
                className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                  selectedCategory === category.id
                    ? 'bg-accent text-accent-foreground font-medium'
                    : 'text-foreground hover:bg-accent/50'
                }`}
              >
                {category.name}
              </button>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  );
}
