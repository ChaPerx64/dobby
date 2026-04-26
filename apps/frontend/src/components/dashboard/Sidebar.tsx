import { useState } from 'react';
import { useAuth } from 'react-oidc-context';
import { Separator } from '@/components/ui/separator';
import { Button } from '@/components/ui/button';
import type { CategoryItem } from '@/types/dashboard';
import { getMonthName, formatDateRange } from '@/lib/format';
import { LogOut, X, Star } from 'lucide-react';
import { CreateEnvelopeModal } from './CreateEnvelopeModal';
import { CreateSpendingModal } from './CreateSpendingModal';
import { apiClient } from '@/api/client';
import type { components } from '@/api/types';

interface SidebarProps {
  period: { startDate: string; endDate: string; id: string; defaultEnvelopeId?: string };
  categories: CategoryItem[];
  selectedCategory: string;
  onSelectCategory: (categoryId: string) => void;
  onEnvelopeCreated: (envelope: components["schemas"]["Envelope"]) => void;
  onAllocationCreated: () => void;
  onDefaultEnvelopeChanged: () => void;
  isOpen: boolean;
  onClose: () => void;
}

export function Sidebar({
  period,
  categories,
  selectedCategory,
  onSelectCategory,
  onEnvelopeCreated,
  onAllocationCreated,
  onDefaultEnvelopeChanged,
  isOpen,
  onClose,
}: SidebarProps) {
  const auth = useAuth();
  const [settingDefaultFor, setSettingDefaultFor] = useState<string | null>(null);
  const monthName = getMonthName(period.startDate);
  const dateRange = formatDateRange(period.startDate, period.endDate);

  const handleLogout = () => {
    auth.signoutRedirect();
  };

  const handleSetDefault = async (envelopeId: string) => {
    const isCurrentDefault = period.defaultEnvelopeId === envelopeId;
    setSettingDefaultFor(envelopeId);
    try {
      await apiClient.updatePeriod(period.id, {
        defaultEnvelopeId: isCurrentDefault ? null : envelopeId,
      });
      onDefaultEnvelopeChanged();
    } finally {
      setSettingDefaultFor(null);
    }
  };

  const envelopes = categories
    .filter((cat) => cat.id !== 'total')
    .map((cat) => ({ id: cat.id, name: cat.name }));

  const defaultEnvelopeId = selectedCategory !== 'total' ? selectedCategory : undefined;

  return (
    <div className={`fixed inset-y-0 left-0 z-50 w-64 bg-card border-r border-border p-6 flex flex-col transform transition-transform duration-200 ease-in-out md:relative md:translate-x-0 ${isOpen ? 'translate-x-0' : '-translate-x-full'}`}>
      <div className="mb-4 flex justify-between items-start">
        <div>
          <h1 className="text-2xl font-bold text-foreground">{monthName}</h1>
          <p className="text-sm text-muted-foreground">{dateRange}</p>
        </div>
        <Button variant="ghost" size="icon" className="md:hidden" onClick={onClose}>
          <X size={20} />
        </Button>
      </div>

      <Separator className="mb-4" />

      <nav className="flex-1 overflow-y-auto">
        <div className="mb-2 space-y-1">
          <CreateEnvelopeModal onEnvelopeCreated={onEnvelopeCreated} />
          <CreateSpendingModal
            envelopes={envelopes}
            onSpendingCreated={onAllocationCreated}
            defaultEnvelopeId={defaultEnvelopeId}
            buttonVariant="ghost"
            buttonClassName="w-full justify-start gap-2 px-3 py-2 h-auto font-normal text-muted-foreground hover:text-foreground"
          />
        </div>
        <ul className="space-y-1">
          {categories.map((category) => (
            <li key={category.id} className="flex items-center gap-1">
              <button
                onClick={() => {
                  onSelectCategory(category.id);
                  onClose();
                }}
                className={`flex-1 text-left px-3 py-2 rounded-md text-sm transition-colors ${
                  selectedCategory === category.id
                    ? 'bg-accent text-accent-foreground font-medium'
                    : 'text-foreground hover:bg-accent/50'
                }`}
              >
                {category.name}
              </button>
              {category.id !== 'total' && (
                <button
                  onClick={() => handleSetDefault(category.id)}
                  disabled={settingDefaultFor === category.id}
                  title={period.defaultEnvelopeId === category.id ? 'Remove default' : 'Set as default envelope'}
                  className={`p-1.5 rounded-md transition-colors flex-shrink-0 ${
                    period.defaultEnvelopeId === category.id
                      ? 'text-yellow-500 hover:text-yellow-600'
                      : 'text-muted-foreground hover:text-foreground'
                  } disabled:opacity-40`}
                >
                  <Star
                    size={14}
                    fill={period.defaultEnvelopeId === category.id ? 'currentColor' : 'none'}
                  />
                </button>
              )}
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
