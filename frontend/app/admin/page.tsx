"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { getLLMProviders } from "@/lib/api/admin";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  LayoutDashboard,
  Settings,
  Users,
  Cpu,
  Database,
  ChevronRight,
  Activity,
  Shield
} from "lucide-react";
import Link from "next/link";

export default function AdminDashboardPage() {
  const { data: providers, isLoading } = useQuery({
    queryKey: ["admin-llm-providers"],
    queryFn: getLLMProviders,
  });

  const stats = [
    {
      title: "LLM Providers",
      value: providers?.length || 0,
      icon: Cpu,
      color: "text-indigo-600",
      bgColor: "bg-indigo-50",
      href: "/admin/llm",
    },
    {
      title: "Active Models",
      value: providers?.filter(p => p.is_active).length || 0,
      icon: Activity,
      color: "text-emerald-600",
      bgColor: "bg-emerald-50",
      href: "/admin/llm",
    },
    {
      title: "System Status",
      value: "Operational",
      icon: Shield,
      color: "text-blue-600",
      bgColor: "bg-blue-50",
      href: "#",
    },
  ];

  return (
    <div className="container mx-auto p-6 max-w-7xl font-sans min-h-screen bg-zinc-50 dark:bg-black">
      <div className="mb-10">
        <h1 className="text-4xl font-extrabold tracking-tight text-zinc-900 dark:text-white flex items-center gap-3">
          <LayoutDashboard className="h-10 w-10 text-indigo-600" />
          Admin Dashboard
        </h1>
        <p className="text-zinc-500 dark:text-zinc-400 mt-2 text-lg">
          Manage your Skillora platform configuration and settings.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
        {stats.map((stat) => (
          <Link key={stat.title} href={stat.href}>
            <Card className="border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 hover:shadow-lg transition-all cursor-pointer">
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-zinc-600 dark:text-zinc-400">
                      {stat.title}
                    </p>
                    <p className="text-3xl font-bold text-zinc-900 dark:text-white mt-2">
                      {stat.value}
                    </p>
                  </div>
                  <div className={`p-3 rounded-xl ${stat.bgColor}`}>
                    <stat.icon className={`h-6 w-6 ${stat.color}`} />
                  </div>
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-10">
        <Card className="border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50">
          <CardHeader>
            <CardTitle className="text-xl font-bold flex items-center gap-2">
              <Settings className="h-5 w-5 text-indigo-500" />
              Configuration
            </CardTitle>
            <CardDescription>Manage platform settings and integrations</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <Link href="/admin/llm" className="block">
              <Button variant="outline" className="w-full justify-between h-12 rounded-xl border-zinc-200 hover:border-indigo-300 hover:bg-indigo-50/50 transition-all">
                <span className="flex items-center gap-2">
                  <Cpu className="h-4 w-4" />
                  LLM Providers
                </span>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </Link>
            <Button variant="outline" className="w-full justify-between h-12 rounded-xl border-zinc-200 opacity-50" disabled>
              <span className="flex items-center gap-2">
                <Database className="h-4 w-4" />
                Database Migrations
              </span>
              <ChevronRight className="h-4 w-4" />
            </Button>
          </CardContent>
        </Card>

        <Card className="border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50">
          <CardHeader>
            <CardTitle className="text-xl font-bold flex items-center gap-2">
              <Users className="h-5 w-5 text-indigo-500" />
              User Management
            </CardTitle>
            <CardDescription>Manage users and permissions</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <Button variant="outline" className="w-full justify-between h-12 rounded-xl border-zinc-200 opacity-50" disabled>
              <span className="flex items-center gap-2">
                <Shield className="h-4 w-4" />
                Admin Privileges
              </span>
              <ChevronRight className="h-4 w-4" />
            </Button>
            <Button variant="outline" className="w-full justify-between h-12 rounded-xl border-zinc-200 opacity-50" disabled>
              <span className="flex items-center gap-2">
                <Users className="h-4 w-4" />
                All Users
              </span>
              <ChevronRight className="h-4 w-4" />
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity */}
      <Card className="border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50">
        <CardHeader>
          <CardTitle className="text-xl font-bold flex items-center gap-2">
            <Activity className="h-5 w-5 text-indigo-500" />
            System Overview
          </CardTitle>
          <CardDescription>Current platform configuration status</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl">
              <div>
                <p className="font-semibold text-zinc-900 dark:text-white">Authentication</p>
                <p className="text-sm text-zinc-500">Google OAuth2 with JWT tokens</p>
              </div>
              <Badge className="bg-emerald-500/10 text-emerald-500 border-none">Active</Badge>
            </div>
            <div className="flex items-center justify-between p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl">
              <div>
                <p className="font-semibold text-zinc-900 dark:text-white">LLM Router</p>
                <p className="text-sm text-zinc-500">
                  {providers?.length || 0} providers configured
                </p>
              </div>
              <Badge className={providers && providers.length > 0 ? "bg-emerald-500/10 text-emerald-500 border-none" : "bg-amber-500/10 text-amber-500 border-none"}>
                {providers && providers.length > 0 ? "Operational" : "Setup Required"}
              </Badge>
            </div>
            <div className="flex items-center justify-between p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-xl">
              <div>
                <p className="font-semibold text-zinc-900 dark:text-white">Admin Access</p>
                <p className="text-sm text-zinc-500">Role-based access control enabled</p>
              </div>
              <Badge className="bg-emerald-500/10 text-emerald-500 border-none">Active</Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
