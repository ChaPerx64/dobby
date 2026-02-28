import { useEffect, useRef } from 'react';
import type { Transaction } from '@/types/api';
import { formatCurrency } from '@/lib/format';

interface TransactionListProps {
  transactions: Transaction[];
  initialBalance: number;
}

export function TransactionList({ transactions, initialBalance }: TransactionListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'auto' });
  }, [transactions]);

  if (transactions.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center text-muted-foreground">
        No transactions found for this period.
      </div>
    );
  }

  // Sort by date ascending (oldest first)
  const sortedTransactions = [...transactions].sort((a, b) => {
    return new Date(a.date || 0).getTime() - new Date(b.date || 0).getTime();
  });

  // Calculate running balance
  let currentBalance = initialBalance;
  const transactionsWithBalance = sortedTransactions.map(tx => {
    currentBalance += tx.amount;
    return { ...tx, runningBalance: currentBalance };
  });

  return (
    <div className="flex-1 overflow-auto -mx-6 px-6">
      <table className="w-full text-sm text-left border-collapse">
        <thead className="sticky top-0 bg-background/95 backdrop-blur z-10 border-b border-border">
          <tr>
            <th className="py-3 pr-4 font-medium text-muted-foreground">Date</th>
            <th className="px-4 py-3 font-medium text-muted-foreground">Description</th>
            <th className="px-4 py-3 font-medium text-muted-foreground">Category</th>
            <th className="py-3 px-4 font-medium text-muted-foreground text-right">Amount</th>
            <th className="py-3 pl-4 font-medium text-muted-foreground text-right">Balance</th>
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
              <td className="py-3 pl-4 text-right font-medium">
                {formatCurrency(tx.runningBalance)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div ref={bottomRef} className="h-1" />
    </div>
  );
}
