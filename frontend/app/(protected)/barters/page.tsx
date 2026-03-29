"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { 
  getMyBarters, 
  getCreditBalance, 
  updateBarterStatus, 
  getBarterMilestones,
  completeMilestone,
  approveMilestone,
  BarterTransaction, 
  Milestone
} from "@/lib/api/barter";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { 
  Coins, 
  Clock, 
  CheckCircle2, 
  XCircle, 
  ArrowRight, 
  Handshake, 
  ChevronDown, 
  ChevronUp, 
  Play,
  Check
} from "lucide-react";
import { format } from "date-fns";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

const statusConfig: Record<string, { label: string; color: string; icon: React.ReactNode }> = {
  pending: { label: "Pending", color: "bg-amber-500", icon: <Clock className="h-3 w-3" /> },
  accepted: { label: "Accepted", color: "bg-blue-500", icon: <CheckCircle2 className="h-3 w-3" /> },
  completed: { label: "Completed", color: "bg-green-500", icon: <CheckCircle2 className="h-3 w-3" /> },
  cancelled: { label: "Cancelled", color: "bg-red-500", icon: <XCircle className="h-3 w-3" /> },
};

export default function BartersPage() {
  const [expandedId, setExpandedId] = React.useState<string | null>(null);

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
        <Card className="flex items-center gap-3 px-5 py-3 border-indigo-200 dark:border-indigo-800 bg-white dark:bg-zinc-950 shadow-sm">
          <Coins className="h-5 w-5 text-indigo-500" />
          <div>
            <p className="text-xs text-muted-foreground font-medium uppercase tracking-wider">Balance</p>
            <p className="text-xl font-bold text-indigo-600 dark:text-indigo-400 leading-tight">{balance ?? "0"} <span className="text-sm font-normal text-muted-foreground">credits</span></p>
          </div>
        </Card>
      </div>

      <div className="space-y-4">
        {barters?.map((b: BarterTransaction) => (
          <BarterItem 
            key={b.id} 
            barter={b} 
            isExpanded={expandedId === b.id} 
            onToggle={() => setExpandedId(expandedId === b.id ? null : b.id)}
            onStatusChange={handleStatusChange}
          />
        ))}
      </div>
    </div>
  );
}

function BarterItem({ barter: b, isExpanded, onToggle, onStatusChange }: any) {
  const cfg = statusConfig[b.status];
  
  const { data: milestones, refetch: refetchMilestones } = useQuery({
    queryKey: ["milestones", b.id],
    queryFn: () => getBarterMilestones(b.id),
    enabled: isExpanded,
  });

  const handleAction = async (id: string, action: "complete" | "approve") => {
    try {
      if (action === "complete") await completeMilestone(id);
      else await approveMilestone(id);
      toast.success(`Milestone updated.`);
      refetchMilestones();
    } catch {
      toast.error("Action failed.");
    }
  };

  return (
    <Card className={cn("overflow-hidden transition-all duration-200 border-muted/60", isExpanded && "ring-1 ring-indigo-500/20 shadow-md")}>
      <CardContent className="p-0">
        <div className="p-5 flex flex-col md:flex-row justify-between gap-4">
          <div className="flex-1 space-y-2">
            <div className="flex items-center gap-3 flex-wrap">
              <Badge className={cn(cfg.color, "text-white text-[10px] uppercase font-bold px-2")}>
                {cfg.label}
              </Badge>
              <span className="text-xs text-muted-foreground">
                Started {format(new Date(b.created_at), "MMM d, yyyy")}
              </span>
            </div>
            
            <h3 className="font-semibold text-lg flex items-center gap-2">
              Exchange: <span className="text-indigo-600 dark:text-indigo-400">{b.credit_amount} Credits</span>
            </h3>

            <div className="flex items-center gap-2 text-sm text-muted-foreground bg-muted/30 p-2 rounded-lg w-fit">
              <span className="font-mono text-[10px]">{b.initiator_id.slice(0, 8)}</span>
              <ArrowRight className="h-3 w-3" />
              <span className="font-mono text-[10px]">{b.receiver_id.slice(0, 8)}</span>
            </div>
          </div>

          <div className="flex items-center gap-2 self-end md:self-center">
            {b.status === "pending" && (
              <>
                <Button size="sm" variant="outline" className="text-green-600 h-9 px-4" onClick={() => onStatusChange(b.id, "accepted")}>
                  Accept
                </Button>
                <Button size="sm" variant="ghost" className="text-red-500 h-9" onClick={() => onStatusChange(b.id, "cancelled")}>
                   Decline
                </Button>
              </>
            )}
            <Button size="sm" variant="ghost" className="h-9 px-2" onClick={onToggle}>
              {isExpanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </Button>
          </div>
        </div>

        {isExpanded && (
          <div className="border-t border-muted bg-zinc-50/50 dark:bg-zinc-900/50 p-5 space-y-4">
            <div className="flex items-center justify-between">
               <h4 className="text-sm font-bold uppercase tracking-widest text-muted-foreground flex items-center gap-2">
                 <Handshake className="h-4 w-4" /> Study Track & Milestones
               </h4>
               <Badge variant="outline" className="bg-white dark:bg-black font-mono text-[10px]">AI-Generated Plan</Badge>
            </div>
            
            <div className="space-y-3">
              {milestones?.map((m: Milestone, idx: number) => (
                <div key={m.id} className="relative pl-6 py-1">
                  {/* Vertical connector line */}
                  {idx !== milestones.length - 1 && (
                    <div className="absolute left-[11px] top-6 bottom-[-12px] w-[1px] bg-muted-foreground/20" />
                  )}
                  
                  {/* Status node */}
                  <div className={cn(
                    "absolute left-0 top-1.5 w-[22px] h-[22px] rounded-full border-2 flex items-center justify-center bg-white dark:bg-black z-10",
                    m.status === "approved" ? "border-green-500 text-green-500" :
                    m.status === "completed" ? "border-amber-500 text-amber-500" :
                    "border-muted text-muted-foreground"
                  )}>
                    {m.status === "approved" ? <Check className="h-3 w-3" /> : <div className="w-1.5 h-1.5 rounded-full bg-current" />}
                  </div>

                  <div className="flex justify-between gap-4 group">
                    <div>
                      <h5 className="text-sm font-semibold">{m.title}</h5>
                      <p className="text-xs text-muted-foreground leading-relaxed mt-0.5">{m.description}</p>
                      <div className="flex items-center gap-2 mt-2">
                         <span className="text-[10px] font-bold text-indigo-500 bg-indigo-50 dark:bg-indigo-950/40 px-1.5 py-0.5 rounded">
                           {m.credit_portion} CR
                         </span>
                         <span className="text-[10px] capitalize text-muted-foreground">{m.status}</span>
                      </div>
                    </div>

                    <div className="shrink-0 transition-opacity flex items-center gap-1">
                      {m.status === "pending" && (
                        <Button size="xs" variant="outline" className="h-7 text-[10px] gap-1 px-2" onClick={() => handleAction(m.id, "complete")}>
                           <Play className="h-2.5 w-2.5" /> Submit Step
                        </Button>
                      )}
                      {m.status === "completed" && (
                        <Button size="xs" variant="outline" className="h-7 text-[10px] gap-1 px-2 border-green-500 text-green-600 hover:bg-green-50" onClick={() => handleAction(m.id, "approve")}>
                           <CheckCircle2 className="h-2.5 w-2.5" /> Approve & Release
                        </Button>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
