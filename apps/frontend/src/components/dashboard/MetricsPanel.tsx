import { Item, ItemContent, ItemTitle, ItemDescription, ItemGroup } from '@/components/ui/item';
import { formatCurrency } from '@/lib/format';
import { CreateAllocationModal } from './CreateAllocationModal';

interface MetricsPanelProps {
  allocated: number;
  spent: number;
  remaining: number;
  projectedBalance: number;
  envelopes: Array<{ id: string; name: string }>;
  defaultEnvelopeId?: string;
  onAllocationCreated: () => void;
}

export function MetricsPanel({
  allocated,
  spent,
  remaining,
  projectedBalance,
  envelopes,
  defaultEnvelopeId,
  onAllocationCreated,
}: MetricsPanelProps) {
  const metrics = [
    { label: 'Allocated', value: allocated, highlight: false },
    { label: 'Spent', value: spent, highlight: false },
    { label: 'Remaining', value: remaining, highlight: false },
    {
      label: 'Projected ending balance',
      value: projectedBalance,
      highlight: projectedBalance < 0,
    },
  ];

  return (
    <div className="w-full md:w-80 bg-background p-4 md:p-6 border-b border-border md:border-b-0 md:border-r">
      <ItemGroup className="grid grid-cols-2 md:flex md:flex-col gap-4">
        {metrics.map((metric, index) => (
          <Item key={index} variant="muted">
            <ItemContent>
              <ItemTitle className="text-sm font-medium text-muted-foreground">
                {metric.label}
              </ItemTitle>
              <ItemDescription
                className={`text-2xl font-bold ${
                  metric.highlight ? 'text-destructive' : 'text-foreground'
                }`}
              >
                {formatCurrency(metric.value)}
              </ItemDescription>
            </ItemContent>
          </Item>
        ))}
      </ItemGroup>
      <div className="mt-4 md:mt-6">
        <CreateAllocationModal
          envelopes={envelopes}
          defaultEnvelopeId={defaultEnvelopeId}
          onAllocationCreated={onAllocationCreated}
          buttonVariant="outline"
          buttonClassName="w-full justify-center gap-2"
        />
      </div>
    </div>
  );
}
