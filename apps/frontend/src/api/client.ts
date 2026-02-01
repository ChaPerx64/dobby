import createClient from 'openapi-fetch';
import type { paths } from './types';

// Create base client
export const api = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
});

// Middleware for JWT authentication
export function setAuthToken(token: string) {
  api.use({
    async onRequest({ request }) {
      request.headers.set('Authorization', `Bearer ${token}`);
      return request;
    },
  });
}

// Clear auth token
export function clearAuthToken() {
  api.use({ onRequest: async ({ request }) => request });
}

// Type-safe API methods
export const apiClient = {
  // Users
  async getCurrentUser() {
    return api.GET('/me');
  },

  async listUsers() {
    return api.GET('/users');
  },

  // Envelopes
  async listEnvelopes() {
    return api.GET('/envelopes');
  },

  // Allocations
  async listAllocations(periodId: string, userId?: string) {
    return api.GET('/allocations', {
      params: {
        query: { periodId, userId },
      },
    });
  },

  // Periods
  async getCurrentPeriod() {
    return api.GET('/periods/current');
  },

  async listPeriods() {
    return api.GET('/periods');
  },

  async getPeriod(periodId: string) {
    return api.GET('/periods/{periodId}', {
      params: { path: { periodId } },
    });
  },

  // Transactions
  async listTransactions(periodId?: string) {
    return api.GET('/transactions', {
      params: periodId ? { query: { periodId } } : undefined,
    });
  },

  async createTransaction(transaction: {
    envelopeId: string;
    amount: number;
    description?: string;
    category?: string;
  }) {
    return api.POST('/transactions', {
      body: transaction,
    });
  },

  async getTransaction(transactionId: string) {
    return api.GET('/transactions/{transactionId}', {
      params: { path: { transactionId } },
    });
  },
};