# Create Allocation Modal Proposal (Amended)

## Objective
Enable users to create Allocations (positive transactions) for a specific envelope via a modal window.

## Background
Currently, users can create envelopes but cannot easily add funds (allocations) to them from the UI. An allocation is essentially a transaction with a positive amount assigned to an envelope.

## Proposed Changes

### Frontend

1.  **API Client Refactor (`apps/frontend/src/api/client.ts`)**
    *   Refactor `apiClient.createTransaction` to use the auto-generated `components["schemas"]["CreateTransaction"]` type.
    *   Ensure the `date` field is properly handled and sent as an ISO string.

2.  **New Component: `CreateAllocationModal`**
    *   Location: `apps/frontend/src/components/dashboard/CreateAllocationModal.tsx`
    *   Props:
        *   `envelopes`: List of available envelopes (`id`, `name`).
        *   `defaultEnvelopeId`: (Optional) ID to pre-select.
        *   `onAllocationCreated`: Callback function to refresh data after success.
    *   Implementation:
        *   Use `Dialog` and other `shadcn/ui` components.
        *   **State:** `envelopeId`, `amount` (string for input), `description`, `date` (default to today).
        *   **Validation:** 
            *   Amount must be a valid positive number.
            *   Envelope must be selected.
        *   **Currency Conversion:** Use a robust conversion (e.g., `Math.round(parseFloat(amount) * 100)`) to cents.
        *   **Date Conversion:** Convert the input date to ISO string via `new Date(date).toISOString()`.

3.  **Dashboard Integration**
    *   Update `Sidebar.tsx` to include the `CreateAllocationModal` trigger.
    *   In `Dashboard.tsx`, implement a `refreshData` function that re-fetches the current period summary.
    *   Pass `refreshData` to `onAllocationCreated`.
    *   Pass the current `selectedCategory` (if it's an envelope) as `defaultEnvelopeId` to the modal.

### Backend
No changes required. The `POST /transactions` endpoint already supports creating positive transactions (allocations).

## Detailed Design

### API Call Details
- **Method:** `apiClient.createTransaction`
- **Payload:**
    - `envelopeId`: Selected ID.
    - `amount`: Calculated cents (positive for allocation).
    - `description`: User input.
    - `date`: ISO 8601 string.
    - `category`: Default to "Funding".

### State Management Strategy
To ensure data integrity, the frontend will **not** attempt to manually calculate the new period balances. Instead, it will re-fetch the `PeriodSummary` from the backend after a successful transaction. This ensures that `totalRemaining`, `spent`, and `projectedBalance` are always accurate according to the server's logic.

## Risks & Considerations
- **Floating Point Math:** Using `Math.round` handles standard decimal inputs correctly, but we should be mindful of precision.
- **Data Freshness:** Re-fetching adds a small network overhead but significantly improves reliability compared to manual state synchronization.
