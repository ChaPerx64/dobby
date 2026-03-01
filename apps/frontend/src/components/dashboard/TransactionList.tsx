import { useEffect, useRef, useState } from 'react';
import { Pencil, Trash2 } from 'lucide-react';
import type { Transaction } from '@/types/api';
import { formatCurrency } from '@/lib/format';
import { Button } from '@/components/ui/button';
import { apiClient } from '@/api/client';
import { EditTransactionModal } from './EditTransactionModal';
import { CreateSpendingModal } from './CreateSpendingModal';

interface TransactionListProps {
  transactions: Transaction[];
  initialBalance: number;
  envelopes: Array<{ id: string; name: string }>;
  defaultEnvelopeId?: string;
  onTransactionChange: () => void;
}

export function TransactionList({ transactions, initialBalance, envelopes, defaultEnvelopeId, onTransactionChange }: TransactionListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const [editingTx, setEditingTx] = useState<Transaction | null>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'auto' });
  }, [transactions]);

  const handleDelete = async (txId: string) => {
    if (!window.confirm('Are you sure you want to delete this transaction?')) return;
    
    try {
      const { error } = await apiClient.deleteTransaction(txId);
      if (error) {
        alert('Failed to delete transaction: ' + error.message);
        return;
      }
      onTransactionChange();
    } catch (err) {
      console.error(err);
      alert('An unexpected error occurred');
    }
  };

  // Sort by date ascending (oldest first)
  const sortedTransactions = [...transactions].sort((a, b) => {
    return new Date(a.date || 0).getTime() - new Date(b.date || 0).getTime();
  });

  // Calculate running balance
  const transactionsWithBalance = sortedTransactions.reduce((acc, tx) => {
    const lastBalance = acc.length > 0 ? acc[acc.length - 1].runningBalance : initialBalance;
    const runningBalance = lastBalance + tx.amount;
    acc.push({ ...tx, runningBalance });
    return acc;
  }, [] as (Transaction & { runningBalance: number })[]);

  const content = transactions.length === 0 ? (
    <div className="flex-1 flex items-center justify-center text-muted-foreground py-12">
      No transactions found for this period.
    </div>
  ) : (
    <>
      <div className="w-full overflow-x-auto">
        <table className="w-full min-w-[600px] text-sm text-left border-collapse">
          <thead className="sticky top-0 bg-background/95 backdrop-blur z-10 border-b border-border">
          <tr>
            <th className="py-3 pr-4 font-medium text-muted-foreground">Date</th>
            <th className="px-4 py-3 font-medium text-muted-foreground">Description</th>
            <th className="px-4 py-3 font-medium text-muted-foreground">Category</th>
            <th className="py-3 px-4 font-medium text-muted-foreground text-right">Amount</th>
            <th className="py-3 px-4 font-medium text-muted-foreground text-right">Balance</th>
            <th className="py-3 pl-4 font-medium text-muted-foreground text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {transactionsWithBalance.map((tx) => (
            <tr key={tx.id} className="hover:bg-muted/50 transition-colors">
              <td className="py-3 pr-4">
                {tx.date ? new Date(tx.date).toLocaleDateString() : '-'}
              </td>
              <td className="px-4 py-3">{tx.description || '-'}</td>
              <td className="px-4 py-3 text-muted-foreground">{tx.category || '-'}</td>
              <td className={`py-3 px-4 text-right font-medium ${tx.amount < 0 ? 'text-foreground' : 'text-green-600'}`}>
                {tx.amount > 0 ? '+' : ''}{formatCurrency(tx.amount)}
              </td>
              <td className="py-3 px-4 text-right font-medium">
                {formatCurrency(tx.runningBalance)}
              </td>
              <td className="py-3 pl-4 text-right">
                <div className="flex justify-end gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-muted-foreground hover:text-foreground"
                    onClick={() => setEditingTx(tx)}
                  >
                    <Pencil size={14} />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-muted-foreground hover:text-destructive"
                    onClick={() => handleDelete(tx.id)}
                  >
                    <Trash2 size={14} />
                  </Button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
        </table>
      </div>
      <div ref={bottomRef} className="h-1" />
    </>
  );

  return (
    <div className="flex-1 overflow-y-auto overflow-x-hidden -mx-6 px-6">
      <div className="flex justify-end mb-4 pt-2">
        <CreateSpendingModal 
          envelopes={envelopes}
          defaultEnvelopeId={defaultEnvelopeId}
          onSpendingCreated={onTransactionChange}
        />
      </div>
      {content}
      <EditTransactionModal
        transaction={editingTx}
        envelopes={envelopes}
        open={!!editingTx}
        onOpenChange={(open) => !open && setEditingTx(null)}
        onSuccess={onTransactionChange}
      />
    </div>
  );
}
