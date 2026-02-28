import * as React from "react"
import { Receipt } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { apiClient } from "@/api/client"

interface CreateSpendingModalProps {
  envelopes: Array<{ id: string; name: string }>
  defaultEnvelopeId?: string
  onSpendingCreated: () => void
}

export function CreateSpendingModal({
  envelopes,
  defaultEnvelopeId,
  onSpendingCreated,
}: CreateSpendingModalProps) {
  const [open, setOpen] = React.useState(false)
  const [isLoading, setIsLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const [envelopeId, setEnvelopeId] = React.useState(defaultEnvelopeId || "")
  const [amount, setAmount] = React.useState("")
  const [description, setDescription] = React.useState("")
  const [category, setCategory] = React.useState("")
  const [date, setDate] = React.useState(new Date().toISOString().split("T")[0])

  // Reset form when defaultEnvelopeId changes or modal opens
  React.useEffect(() => {
    if (open) {
      setEnvelopeId(defaultEnvelopeId || (envelopes.length > 0 ? envelopes[0].id : ""))
      setAmount("")
      setDescription("")
      setCategory("")
      setDate(new Date().toISOString().split("T")[0])
      setError(null)
    }
  }, [open, defaultEnvelopeId, envelopes])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

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
      // Currency Conversion: Convert to negative cents with rounding to avoid floating point issues
      const amountInCents = Math.round(parsedAmount * 100) * -1
      
      // Date Conversion: Convert to ISO string (keeping the time as noon to avoid timezone shifts)
      const dateObj = new Date(date)
      dateObj.setHours(12, 0, 0, 0)
      const isoDate = dateObj.toISOString()

      const { error: apiError } = await apiClient.createTransaction({
        envelopeId,
        amount: amountInCents,
        description: description.trim() || "Spending",
        date: isoDate,
        category: category.trim() || undefined,
      })

      if (apiError) {
        setError(apiError.message || "Failed to record spending")
        return
      }

      onSpendingCreated()
      setOpen(false)
    } catch (err) {
      setError("An unexpected error occurred")
      console.error(err)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full justify-start gap-2 px-3 py-2 h-auto font-normal text-muted-foreground hover:text-foreground"
        >
          <Receipt size={16} />
          Add Spending
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Add Spending</DialogTitle>
            <DialogDescription>
              Record an expense from an envelope. This will decrease its balance.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="envelope">Envelope</Label>
              <select
                id="envelope"
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
              <Label htmlFor="amount">Amount</Label>
              <Input
                id="amount"
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
              <Label htmlFor="description">Description (Optional)</Label>
              <Input
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="e.g. Grocery shopping"
                disabled={isLoading}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="category">Category (Optional)</Label>
              <Input
                id="category"
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                placeholder="e.g. food, transport"
                disabled={isLoading}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="date">Date</Label>
              <Input
                id="date"
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
            <Button type="button" variant="outline" onClick={() => setOpen(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !amount || !envelopeId}>
              {isLoading ? "Recording..." : "Record Spending"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
