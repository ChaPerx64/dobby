import createClient from 'openapi-fetch';
import type { paths } from './types';

// Create base client
export const api = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
});

let authToken: string | null = null;

// Middleware for JWT authentication
api.use({
  async onRequest({ request }) {
    if (authToken) {
      request.headers.set('Authorization', `Bearer ${authToken}`);
    }
    return request;
  },
});

export function setAuthToken(token: string) {
  authToken = token;
}

export function clearAuthToken() {
  authToken = null;
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