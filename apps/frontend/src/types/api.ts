// Re-export OpenAPI types for easier imports
import type { components } from '@/api/types';

export type User = components['schemas']['User'];
export type Period = components['schemas']['PeriodSummary'];
export type PeriodListItem = components['schemas']['PeriodListItem'];
export type Envelope = components['schemas']['Envelope'];
export type EnvelopeSummary = components['schemas']['EnvelopeSummary'];
export type Transaction = components['schemas']['Transaction'];
export type CreateTransaction = components['schemas']['CreateTransaction'];
export type ApiError = components['schemas']['Error'];
