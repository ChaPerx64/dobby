// Re-export OpenAPI types for easier imports
import type { components } from '@/api/types';

export type User = components['schemas']['User'];
export type Period = components['schemas']['Period'];
export type Envelope = components['schemas']['Envelope'];
export type Allocation = components['schemas']['Allocation'];
export type Transaction = components['schemas']['Transaction'];
export type CreateTransaction = components['schemas']['CreateTransaction'];
export type ApiError = components['schemas']['Error'];
