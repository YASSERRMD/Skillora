"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import Link from "next/link";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function DashboardPage() {
  const { user } = useAuthStore();

  return (
    <div className="container mx-auto p-6 max-w-5xl">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Welcome back, {user?.full_name?.split(" ")[0] || "Trader"}
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage your knowledge portfolio and discover new barters.
          </p>
        </div>
        <Link href="/dashboard/add-skill">
          <Button className="flex items-center gap-2">
            <Plus className="w-4 h-4" />
            Offer a Skill
          </Button>
        </Link>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card className="col-span-1 shadow-sm border-muted/60">
          <CardHeader>
            <CardTitle>My Verified Skills</CardTitle>
            <CardDescription>
              Your AI-appraised competencies available for barter.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {/* Will be populated in Phase 26 with real data from GetUserSkills */}
            <div className="flex flex-col items-center justify-center p-8 text-center bg-muted/30 rounded-lg border border-dashed border-muted">
              <p className="text-sm text-muted-foreground mb-4">
                You haven't added any skills to your portfolio yet.
              </p>
              <Link href="/dashboard/add-skill">
                <Button variant="outline" size="sm">
                  Add Your First Skill
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-1 shadow-sm border-muted/60">
          <CardHeader>
            <CardTitle>Recent Barters</CardTitle>
            <CardDescription>
              Active and past knowledge exchanges.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {/* Will be populated with barter transaction history */}
            <div className="flex flex-col items-center justify-center p-8 text-center bg-muted/30 rounded-lg border border-dashed border-muted">
              <p className="text-sm text-muted-foreground">
                No active barters. Explore the marketplace to find learning opportunities.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
