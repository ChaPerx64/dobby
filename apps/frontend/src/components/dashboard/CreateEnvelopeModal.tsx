import * as React from "react"
import { Plus } from "lucide-react"
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
import type { components } from "@/api/types"

type Envelope = components["schemas"]["Envelope"]

interface CreateEnvelopeModalProps {
  onEnvelopeCreated: (envelope: Envelope) => void
}

export function CreateEnvelopeModal({ onEnvelopeCreated }: CreateEnvelopeModalProps) {
  const [open, setOpen] = React.useState(false)
  const [name, setName] = React.useState("")
  const [isLoading, setIsLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) return

    setIsLoading(true)
    setError(null)

    try {
      const { data, error: apiError, response } = await apiClient.createEnvelope({
        name: name.trim(),
      })

      if (apiError) {
        if (response.status === 409) {
          setError("An envelope with this name already exists")
        } else {
          setError(apiError.message || "Failed to create envelope")
        }
        return
      }

      if (data) {
        onEnvelopeCreated(data)
        setOpen(false)
        setName("")
      }
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
        <Button variant="ghost" size="sm" className="w-full justify-start gap-2 px-3 py-2 h-auto font-normal text-muted-foreground hover:text-foreground">
          <Plus size={16} />
          Add Envelope
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Envelope</DialogTitle>
            <DialogDescription>
              Add a new budget envelope to track your spending.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Groceries, Rent, Utilities"
                autoFocus
                disabled={isLoading}
              />
              {error && <p className="text-sm font-medium text-destructive">{error}</p>}
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setOpen(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !name.trim()}>
              {isLoading ? "Creating..." : "Create Envelope"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
