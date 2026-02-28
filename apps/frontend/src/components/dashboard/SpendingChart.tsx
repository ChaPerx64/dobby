import type { ChartConfig } from '@/components/ui/chart';
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import type { ChartDataPoint } from '@/types/dashboard';
import { formatCurrency } from '@/lib/format';

interface SpendingChartProps {
  data: ChartDataPoint[];
}

const chartConfig = {
  remaining: {
    label: 'Remaining Balance',
    color: 'hsl(var(--chart-1))',
  },
} satisfies ChartConfig;

export function SpendingChart({ data }: SpendingChartProps) {
  // Convert para to RSD for chart display
  const chartData = data.map((point) => ({
    date: new Date(point.date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    }),
    remaining: point.remaining / 100,
  }));

  return (
    <div className="flex-1 bg-background p-6">
      <h2 className="text-xl font-bold text-foreground mb-4">
        Balance Over Time
      </h2>
      <ChartContainer config={chartConfig} className="h-[500px] w-full">
        <AreaChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="date"
            tickLine={false}
            axisLine={false}
            tickMargin={8}
          />
          <YAxis
            tickLine={false}
            axisLine={false}
            tickMargin={8}
            tickFormatter={(value) => formatCurrency(value * 100)}
          />
          <ChartTooltip
            content={
              <ChartTooltipContent
                labelFormatter={(value) => `Date: ${value}`}
                formatter={(value) =>
                  formatCurrency((value as number) * 100)
                }
              />
            }
          />
          <Area
            type="monotone"
            dataKey="remaining"
            stroke={chartConfig.remaining.color}
            fill={chartConfig.remaining.color}
            fillOpacity={0.6}
          />
        </AreaChart>
      </ChartContainer>
    </div>
  );
}
