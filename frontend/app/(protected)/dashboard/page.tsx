"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { useQuery } from "@tanstack/react-query";
import { getMySkills } from "@/lib/api/user-skills";
import { getMyBarters, getCreditBalance } from "@/lib/api/barter";
import { Button } from "@/components/ui/button";
import { Plus, Sparkles, Coins, Handshake, Star, ShoppingBag } from "lucide-react";
import Link from "next/link";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

const proficiencyLabel = (level: number) => {
  return ["", "Beginner", "Intermediate", "Advanced", "Expert", "Master"][level] || "Unknown";
};

const proficiencyColors: Record<number, string> = {
  1: "bg-slate-500",
  2: "bg-blue-500",
  3: "bg-indigo-500",
  4: "bg-purple-500",
  5: "bg-amber-500",
};

export default function DashboardPage() {
  const { user } = useAuthStore();

  const { data: skills, isLoading: skillsLoading } = useQuery({
    queryKey: ["my-skills"],
    queryFn: getMySkills,
  });

  const { data: barters, isLoading: bartersLoading } = useQuery({
    queryKey: ["barters"],
    queryFn: getMyBarters,
  });

  const { data: balance } = useQuery({
    queryKey: ["credit-balance"],
    queryFn: getCreditBalance,
  });

  const activeBarters = barters?.filter((b) => b.status === "pending" || b.status === "accepted") ?? [];

  return (
    <div className="container mx-auto p-6 max-w-5xl">
      {/* Welcome Banner */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-indigo-600 to-purple-700 text-white p-8 mb-8">
        <div className="relative z-10">
          <p className="text-indigo-200 text-sm font-medium mb-1">Welcome back</p>
          <h1 className="text-3xl font-bold mb-4">
            {user?.full_name?.split(" ")[0] || "Trader"} 👋
          </h1>
          <div className="flex flex-wrap gap-4">
            <div className="flex items-center gap-2 bg-white/10 rounded-lg px-4 py-2">
              <Coins className="h-4 w-4 text-amber-300" />
              <span className="font-semibold">{balance ?? "—"} credits</span>
            </div>
            <div className="flex items-center gap-2 bg-white/10 rounded-lg px-4 py-2">
              <Star className="h-4 w-4 text-indigo-200" />
              <span className="font-semibold">{skills?.length ?? 0} verified skills</span>
            </div>
            <div className="flex items-center gap-2 bg-white/10 rounded-lg px-4 py-2">
              <Handshake className="h-4 w-4 text-indigo-200" />
              <span className="font-semibold">{activeBarters.length} active barters</span>
            </div>
          </div>
        </div>
        <div className="absolute -right-12 -bottom-12 w-48 h-48 rounded-full bg-white/5" />
        <div className="absolute -right-4 -top-12 w-32 h-32 rounded-full bg-white/5" />
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-8">
        <Link href="/dashboard/add-skill" className="group">
          <div className="border border-muted/60 rounded-xl p-4 hover:border-indigo-300 hover:bg-indigo-50 dark:hover:bg-indigo-950/20 transition-all duration-200 cursor-pointer h-full flex flex-col items-center text-center gap-2">
            <div className="w-10 h-10 rounded-full bg-indigo-100 dark:bg-indigo-900/40 flex items-center justify-center group-hover:bg-indigo-200 dark:group-hover:bg-indigo-800/40 transition-colors">
              <Plus className="h-5 w-5 text-indigo-600 dark:text-indigo-400" />
            </div>
            <p className="text-sm font-medium">Offer a Skill</p>
          </div>
        </Link>
        <Link href="/marketplace" className="group">
          <div className="border border-muted/60 rounded-xl p-4 hover:border-purple-300 hover:bg-purple-50 dark:hover:bg-purple-950/20 transition-all duration-200 cursor-pointer h-full flex flex-col items-center text-center gap-2">
            <div className="w-10 h-10 rounded-full bg-purple-100 dark:bg-purple-900/40 flex items-center justify-center group-hover:bg-purple-200 dark:group-hover:bg-purple-800/40 transition-colors">
              <ShoppingBag className="h-5 w-5 text-purple-600 dark:text-purple-400" />
            </div>
            <p className="text-sm font-medium">Find Skills</p>
          </div>
        </Link>
        <Link href="/barters" className="group">
          <div className="border border-muted/60 rounded-xl p-4 hover:border-green-300 hover:bg-green-50 dark:hover:bg-green-950/20 transition-all duration-200 cursor-pointer h-full flex flex-col items-center text-center gap-2">
            <div className="w-10 h-10 rounded-full bg-green-100 dark:bg-green-900/40 flex items-center justify-center group-hover:bg-green-200 dark:group-hover:bg-green-800/40 transition-colors">
              <Handshake className="h-5 w-5 text-green-600 dark:text-green-400" />
            </div>
            <p className="text-sm font-medium">My Barters</p>
          </div>
        </Link>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* My Verified Skills */}
        <Card className="col-span-1 shadow-sm border-muted/60">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>My Verified Skills</CardTitle>
              <CardDescription>AI-appraised competencies on your profile.</CardDescription>
            </div>
            <Link href="/dashboard/add-skill">
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <Plus className="h-4 w-4" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {skillsLoading && (
              <div className="space-y-3">
                {[1, 2].map((i) => <div key={i} className="h-12 rounded-lg bg-muted animate-pulse" />)}
              </div>
            )}
            {!skillsLoading && (!skills || skills.length === 0) && (
              <div className="flex flex-col items-center justify-center p-8 text-center bg-muted/30 rounded-lg border border-dashed border-muted">
                <Sparkles className="h-6 w-6 text-muted-foreground mb-3" />
                <p className="text-sm text-muted-foreground mb-3">No skills verified yet.</p>
                <Link href="/dashboard/add-skill">
                  <Button variant="outline" size="sm">Add Your First Skill</Button>
                </Link>
              </div>
            )}
            <div className="space-y-3">
              {skills?.map((s) => (
                <div key={s.skill_id} className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-muted/50">
                  <div>
                    <p className="text-sm font-medium">{s.skill_name}</p>
                    <p className="text-xs text-muted-foreground">{s.category_name}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge className={`${proficiencyColors[s.proficiency_level] || "bg-gray-500"} text-white text-xs`}>
                      {proficiencyLabel(s.proficiency_level)}
                    </Badge>
                    <span className="text-xs font-medium text-indigo-600 dark:text-indigo-400">{s.credit_value}c</span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Barters */}
        <Card className="col-span-1 shadow-sm border-muted/60">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Recent Barters</CardTitle>
              <CardDescription>Your active knowledge exchanges.</CardDescription>
            </div>
            <Link href="/barters">
              <Button variant="ghost" size="sm" className="text-xs h-8">View All</Button>
            </Link>
          </CardHeader>
          <CardContent>
            {bartersLoading && (
              <div className="space-y-3">
                {[1, 2].map((i) => <div key={i} className="h-12 rounded-lg bg-muted animate-pulse" />)}
              </div>
            )}
            {!bartersLoading && (!barters || barters.length === 0) && (
              <div className="flex flex-col items-center justify-center p-8 text-center bg-muted/30 rounded-lg border border-dashed border-muted">
                <p className="text-sm text-muted-foreground">No barters yet.</p>
                <Link href="/marketplace" className="mt-3">
                  <Button variant="outline" size="sm">Explore Marketplace</Button>
                </Link>
              </div>
            )}
            <div className="space-y-3">
              {barters?.slice(0, 4).map((b) => (
                <div key={b.id} className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-muted/50">
                  <div>
                    <p className="text-xs font-mono text-muted-foreground">{b.id.slice(0, 8)}...</p>
                    <p className="text-xs text-muted-foreground mt-0.5">{b.credit_amount} credits</p>
                  </div>
                  <Badge className={
                    b.status === "pending" ? "bg-amber-500 text-white" :
                    b.status === "accepted" ? "bg-blue-500 text-white" :
                    b.status === "completed" ? "bg-green-500 text-white" :
                    "bg-red-500 text-white"
                  }>{b.status}</Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
