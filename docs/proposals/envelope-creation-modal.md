# Proposal: Envelope Creation Modal

## Overview
Add a button to the frontend dashboard that opens a modal window for creating new budget envelopes. The goal is to provide a seamless user experience with immediate feedback and robust error handling.

## Proposed Changes

### 1. Update API Client
*   Add the `createEnvelope` method to `apps/frontend/src/api/client.ts` to support the `POST /envelopes` endpoint defined in the OpenAPI spec.
*   Ensure the client method correctly propagates `409 Conflict` errors (for duplicate names) so the UI can handle them specifically.

### 2. Add UI Primitives
*   Install `@radix-ui/react-dialog` for accessible modal functionality.
*   Install `@radix-ui/react-label` for accessible form labeling.
*   Create `apps/frontend/src/components/ui/button.tsx` using `class-variance-authority` (cva) to support variants (default, outline, ghost, destructive) and sizes, following existing Tailwind patterns.
*   Create `apps/frontend/src/components/ui/dialog.tsx` as a wrapper around Radix UI Dialog.
*   Create `apps/frontend/src/components/ui/input.tsx` for the envelope name field, ensuring it supports `forwardRef` for form library compatibility.
*   Create `apps/frontend/src/components/ui/label.tsx` using Radix Label primitive.

### 3. Create Envelope Modal Component
*   Implement `CreateEnvelopeModal.tsx` in `apps/frontend/src/components/dashboard/`.
*   **State Management:**
    *   Manage local form state for the envelope name.
    *   Manage `isLoading` state during API submission.
    *   Manage `error` state for displaying validation or API errors.
*   **Behavior:**
    *   Auto-focus the input field when the modal opens.
    *   Disable the "Create" button and show a loading indicator (e.g., "Creating...") while the request is in-flight.
    *   Close the modal automatically upon successful creation.
    *   Display specific error messages for `409 Conflict` (e.g., "An envelope with this name already exists") and generic messages for other failures.

### 4. Integrate into Sidebar
*   Update `apps/frontend/src/components/dashboard/Sidebar.tsx` to include a "Add Envelope" button (using the `Plus` icon from `lucide-react`).
*   **Placement:** Add the button at the top of the category list or in a distinct section for better visibility.
*   **Data Flow:**
    *   Pass a callback `onEnvelopeCreated(newEnvelope: Envelope)` from the parent `Dashboard` component.
    *   Trigger the modal when the "Add Envelope" button is clicked.
    *   On success, call `onEnvelopeCreated` with the new envelope object.
    *   The parent component should optimistically update the `categories` list with the new envelope (initialized with 0 allocation/spent) to avoid a full page reload or flicker.

## Verification Plan
*   **Manual Test**: Verify the "Add Envelope" button appears in the sidebar and is accessible via keyboard navigation.
*   **Manual Test**: Confirm the modal opens with focus on the input field.
*   **Manual Test (Happy Path)**: Successfully create a new envelope and verify it immediately appears in the sidebar list without a page reload.
*   **Manual Test (Error Path)**: Attempt to create an envelope with an existing name and verify that a specific error message is displayed without closing the modal.
*   **Manual Test (Loading)**: Verify that the submit button is disabled during the API request to prevent double submissions.
