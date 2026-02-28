import * as React from "react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { apiClient } from "@/api/client"
import type { Transaction } from "@/types/api"

interface EditTransactionModalProps {
  transaction: Transaction | null
  envelopes: Array<{ id: string; name: string }>
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

export function EditTransactionModal({
  transaction,
  envelopes,
  open,
  onOpenChange,
  onSuccess,
}: EditTransactionModalProps) {
  const [isLoading, setIsLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const [envelopeId, setEnvelopeId] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [description, setDescription] = React.useState("")
  const [category, setCategory] = React.useState("")
  const [date, setDate] = React.useState("")

  React.useEffect(() => {
    if (open && transaction) {
      setEnvelopeId(transaction.envelopeId)
      // Convert cents to display amount, handle negative for expenses
      const displayAmount = (Math.abs(transaction.amount) / 100).toFixed(2)
      setAmount(displayAmount)
      setDescription(transaction.description || "")
      setCategory(transaction.category || "")
      // Extract YYYY-MM-DD from the ISO date string
      const localDate = transaction.date ? transaction.date.split("T")[0] : new Date().toISOString().split("T")[0]
      setDate(localDate)
      setError(null)
    }
  }, [open, transaction])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!transaction) return

    if (!envelopeId) {
      setError("Please select an envelope")
      return
    }

    const parsedAmount = parseFloat(amount)
    if (isNaN(parsedAmount) || parsedAmount <= 0) {
      setError("Please enter a valid positive amount")
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      // Assuming it's an expense if it was an expense before. If amount was > 0, we should keep it positive?
      // Actually, all envelope transactions created from CreateSpendingModal are expenses (negative).
      // If we are editing, we should probably keep the sign of the original transaction, or just assume it's an expense 
      // since the prompt implies modifying a spending transaction. Let's keep original sign:
      const sign = transaction.amount < 0 ? -1 : 1
      const amountInCents = Math.round(parsedAmount * 100) * sign
      
      const dateObj = new Date(date)
      dateObj.setHours(12, 0, 0, 0)
      const isoDate = dateObj.toISOString()

      const { error: apiError } = await apiClient.updateTransaction(transaction.id, {
        envelopeId,
        amount: amountInCents,
        description: description.trim() || undefined,
        date: isoDate,
        category: category.trim() || undefined,
      })

      if (apiError) {
        setError(apiError.message || "Failed to update transaction")
        return
      }

      onSuccess()
      onOpenChange(false)
    } catch (err) {
      setError("An unexpected error occurred")
      console.error(err)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Edit Transaction</DialogTitle>
            <DialogDescription>
              Modify the details of this transaction.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="edit-envelope">Envelope</Label>
              <select
                id="edit-envelope"
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                value={envelopeId}
                onChange={(e) => setEnvelopeId(e.target.value)}
                disabled={isLoading}
              >
                <option value="" disabled>Select an envelope</option>
                {envelopes.map((env) => (
                  <option key={env.id} value={env.id}>
                    {env.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-amount">Amount</Label>
              <Input
                id="edit-amount"
                type="number"
                step="0.01"
                min="0.01"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="0.00"
                disabled={isLoading}
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-description">Description</Label>
              <Input
                id="edit-description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="e.g. Grocery shopping"
                disabled={isLoading}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-category">Category (Optional)</Label>
              <Input
                id="edit-category"
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                placeholder="e.g. food, transport"
                disabled={isLoading}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-date">Date</Label>
              <Input
                id="edit-date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                disabled={isLoading}
                required
              />
            </div>
            {error && <p className="text-sm font-medium text-destructive">{error}</p>}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !amount || !envelopeId}>
              {isLoading ? "Saving..." : "Save Changes"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
