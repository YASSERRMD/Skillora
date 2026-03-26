"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { getMyBarters, getCreditBalance, updateBarterStatus, BarterTransaction } from "@/lib/api/barter";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Coins, Clock, CheckCircle2, XCircle, ArrowRight, Handshake } from "lucide-react";
import { format } from "date-fns";
import { toast } from "sonner";

const statusConfig: Record<string, { label: string; color: string; icon: React.ReactNode }> = {
  pending: { label: "Pending", color: "bg-amber-500", icon: <Clock className="h-3 w-3" /> },
  accepted: { label: "Accepted", color: "bg-blue-500", icon: <CheckCircle2 className="h-3 w-3" /> },
  completed: { label: "Completed", color: "bg-green-500", icon: <CheckCircle2 className="h-3 w-3" /> },
  cancelled: { label: "Cancelled", color: "bg-red-500", icon: <XCircle className="h-3 w-3" /> },
};

export default function BartersPage() {
  const { data: barters, isLoading, refetch } = useQuery({
    queryKey: ["barters"],
    queryFn: getMyBarters,
  });

  const { data: balance } = useQuery({
    queryKey: ["credit-balance"],
    queryFn: getCreditBalance,
  });

  const handleStatusChange = async (id: string, status: "accepted" | "cancelled") => {
    try {
      await updateBarterStatus(id, status);
      toast.success(`Barter ${status} successfully.`);
      refetch();
    } catch {
      toast.error("Failed to update barter status.");
    }
  };

  return (
    <div className="container mx-auto p-6 max-w-5xl">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">My Barters</h1>
          <p className="text-muted-foreground mt-1">Manage your knowledge exchange agreements.</p>
        </div>
        <Card className="flex items-center gap-3 px-5 py-3 border-indigo-200 dark:border-indigo-800 bg-indigo-50 dark:bg-indigo-950/20">
          <Coins className="h-5 w-5 text-indigo-500" />
          <div>
            <p className="text-xs text-muted-foreground">Credit Balance</p>
            <p className="text-xl font-bold text-indigo-600 dark:text-indigo-400">{balance ?? "—"}</p>
          </div>
        </Card>
      </div>

      {isLoading && (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-28 rounded-xl bg-muted animate-pulse" />
          ))}
        </div>
      )}

      {!isLoading && barters && barters.length === 0 && (
        <div className="flex flex-col items-center justify-center py-20 text-center">
          <Handshake className="h-12 w-12 text-muted-foreground/40 mb-4" />
          <h3 className="text-lg font-medium mb-2">No Barters Yet</h3>
          <p className="text-sm text-muted-foreground">
            Head to the Marketplace to find skill holders and propose your first knowledge exchange.
          </p>
        </div>
      )}

      <div className="space-y-4">
        {barters?.map((b: BarterTransaction) => {
          const cfg = statusConfig[b.status];
          return (
            <Card key={b.id} className="border-muted/60 shadow-sm">
              <CardContent className="p-5">
                <div className="flex flex-col md:flex-row justify-between gap-4">
                  <div className="space-y-2">
                    <div className="flex items-center gap-3">
                      <Badge className={`${cfg.color} text-white text-xs flex items-center gap-1`}>
                        {cfg.icon} {cfg.label}
                      </Badge>
                      <span className="text-sm text-muted-foreground">
                        {format(new Date(b.created_at), "MMM d, yyyy")}
                      </span>
                    </div>
                    <div className="flex items-center gap-2 text-sm">
                      <span className="font-mono text-xs text-muted-foreground">{b.initiator_id.slice(0, 8)}...</span>
                      <ArrowRight className="h-3 w-3 text-muted-foreground" />
                      <span className="font-mono text-xs text-muted-foreground">{b.receiver_id.slice(0, 8)}...</span>
                    </div>
                    <div className="flex items-center gap-1.5 text-sm">
                      <Coins className="h-3.5 w-3.5 text-indigo-500" />
                      <span className="font-semibold text-indigo-600 dark:text-indigo-400">{b.credit_amount} credits</span>
                    </div>
                  </div>
                  {b.status === "pending" && (
                    <div className="flex items-center gap-2">
                      <Button size="sm" variant="outline" className="text-green-600 border-green-300 hover:bg-green-50 dark:hover:bg-green-950/20"
                        onClick={() => handleStatusChange(b.id, "accepted")}>
                        Accept
                      </Button>
                      <Button size="sm" variant="outline" className="text-red-600 border-red-300 hover:bg-red-50 dark:hover:bg-red-950/20"
                        onClick={() => handleStatusChange(b.id, "cancelled")}>
                        Decline
                      </Button>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
