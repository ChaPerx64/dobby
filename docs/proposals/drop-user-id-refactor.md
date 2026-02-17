# Proposal: Drop `user_id` from Financial Entities

## Objective
Decouple financial data (envelopes, transactions, and periods) from specific user identities. This simplifies the system for household-wide management where financial records are shared regardless of which user created them. The `users` table and user management endpoints will be preserved for authentication and profile purposes.

## Proposed Changes

### 1. Database Schema
Create a new migration to:
- Drop foreign key constraints: `fk_transactions_user` and `fk_envelopes_user`.
- Drop indices: `idx_transactions_user_id` and `idx_envelopes_user_id`.
- Remove the `user_id` column from the `transactions` and `envelopes` tables.
- **Note**: The `users` table itself will **not** be dropped.

### 2. Backend Model Refactor (`apps/backend/internal/service/models.go`)
- Remove the `UserID` field from the `Envelope` struct.
- Remove the `UserID` field from the `Transaction` struct.
- Keep the `User` struct as is.

### 3. Service Interfaces and Filters (`apps/backend/internal/service/interfaces.go`)
- **FinanceService Interface**:
    - Update `RecordTransaction` signature: remove `userID uuid.UUID`.
    - Update `CreateEnvelope` signature: remove `userID uuid.UUID`.
- **TransactionFilter Struct**:
    - Remove `UserID` field (if present) to ensure filtering logic is aligned with the new schema.

### 4. Persistence Layer (`apps/backend/internal/adapters/persistence/postgres.go`)
- Update `SaveEnvelope` and `ListEnvelopes` SQL queries and scanning logic to exclude `user_id`.
- Update `SaveTransaction` and `ListTransactions` SQL queries and scanning logic to exclude `user_id`.
- Update `TransactionFilter` handling in `ListTransactions` to remove any `user_id` filtering logic.
- Keep `SaveUser` and `GetUser` implementations for identity purposes.

### 5. Service Implementation (`apps/backend/internal/service/dobbyFinancier.go`)
- Update `RecordTransaction` and `CreateEnvelope` to match the new interface signatures.
- Remove logic that assigns `UserID` to transactions or envelopes.
- Ensure `GetPeriodSummary` and other aggregation logic no longer rely on `UserID` for grouping or filtering.

### 6. API Handlers and Security (`apps/backend/internal/adapters/api/handlers.go`)
- Update `CreateEnvelope` and `CreateTransaction` handlers to stop extracting `userID` from context.
- **Legacy Removal**: Delete the `getUserID` helper method and the hardcoded fallback UUID (`000...001`). This ensures the system explicitly moves away from "guest" mode for financial records.
- **Context Handling**: Verify that `GetUserID(ctx)` is only used for authentication/identity endpoints (like `/me`) and not for authorization of financial data access.

### 7. Frontend API and Types (`apps/frontend/src/api/types.ts`)
- Re-generate frontend types using `openapi-typescript` after the OpenAPI schema is updated.
- Verify that `Transaction` and `Envelope` interfaces in TypeScript no longer contain `userId`.
- Audit frontend components (e.g., `SpendingChart.tsx`) to ensure they don't attempt to access or display `userId`.

### 8. Verification Plan
1. **Database**: Run migrations and verify `envelopes` and `transactions` tables no longer have `user_id`.
2. **Backend Build**: Run `go generate ./...` and ensure the backend compiles without `UserID` references.
3. **Frontend Build**: Run `npm run build` (or equivalent) in `apps/frontend` to catch any broken type references.
4. **Tests**: Update service and persistence tests to remove `userID` expectations.
5. **Integration**: Verify that the `/me` endpoint still works and that transactions can be created without a user association.
