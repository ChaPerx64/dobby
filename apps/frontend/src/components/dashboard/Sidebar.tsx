import { useAuth } from 'react-oidc-context';
import { Separator } from '@/components/ui/separator';
import type { CategoryItem } from '@/types/dashboard';
import { getMonthName, formatDateRange } from '@/lib/format';
import { LogOut } from 'lucide-react';
import { CreateEnvelopeModal } from './CreateEnvelopeModal';
import { CreateAllocationModal } from './CreateAllocationModal';
import type { components } from '@/api/types';

interface SidebarProps {
  period: { startDate: string; endDate: string };
  categories: CategoryItem[];
  selectedCategory: string;
  onSelectCategory: (categoryId: string) => void;
  onEnvelopeCreated: (envelope: components["schemas"]["Envelope"]) => void;
  onAllocationCreated: () => void;
}

export function Sidebar({
  period,
  categories,
  selectedCategory,
  onSelectCategory,
  onEnvelopeCreated,
  onAllocationCreated,
}: SidebarProps) {
  const auth = useAuth();
  const monthName = getMonthName(period.startDate);
  const dateRange = formatDateRange(period.startDate, period.endDate);

  const handleLogout = () => {
    auth.signoutRedirect();
  };

  // Filter out the 'Total' category for the allocation modal
  const envelopes = categories
    .filter((cat) => cat.id !== 'total')
    .map((cat) => ({ id: cat.id, name: cat.name }));

  const defaultEnvelopeId = selectedCategory !== 'total' ? selectedCategory : undefined;

  return (
    <div className="w-64 bg-card border-r border-border p-6 flex flex-col">
      <div className="mb-4">
        <h1 className="text-2xl font-bold text-foreground">{monthName}</h1>
        <p className="text-sm text-muted-foreground">{dateRange}</p>
      </div>

      <Separator className="mb-4" />

      <nav className="flex-1 overflow-y-auto">
        <div className="mb-2 space-y-1">
          <CreateEnvelopeModal onEnvelopeCreated={onEnvelopeCreated} />
          <CreateAllocationModal
            envelopes={envelopes}
            onAllocationCreated={onAllocationCreated}
            defaultEnvelopeId={defaultEnvelopeId}
          />
        </div>
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

      <div className="mt-auto pt-4">
        <Separator className="mb-4" />
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-2 px-3 py-2 rounded-md text-sm text-destructive hover:bg-destructive/10 transition-colors"
        >
          <LogOut size={16} />
          Logout
        </button>
      </div>
    </div>
  );
}
