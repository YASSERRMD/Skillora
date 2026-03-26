"use client";

import * as React from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { LogOut, User as UserIcon } from "lucide-react";
import api from "@/lib/api";
import { useAuthStore } from "@/store/useAuthStore";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function LayoutHeader() {
  const { user, setUser, clearUser } = useAuthStore();

  // Fetch the authenticated user's profile on mount to hydrate the Zustand store
  const { isLoading, isError } = useQuery({
    queryKey: ["currentUser"],
    queryFn: async () => {
      const res = await api.get("/api/v1/users/me");
      const userData = res.data;
      setUser(userData);
      return userData;
    },
    retry: false,
  });

  const handleLogout = () => {
    // In a real app, you would also call a DELETE /api/v1/auth/logout endpoint to clear the HttpOnly cookie.
    // Since Phase 8 only mints cookies but doesn't define a logout endpoint yet, we just clear cookies
    // by manually setting expiration or relying on the frontend state wipe for now.
    // A robust way to clear HttpOnly cookies is to let the backend do it. We will redirect to /login.
    document.cookie = "skillora_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    clearUser();
    window.location.href = "/login";
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center pl-4 pr-6">
        <Link href="/dashboard" className="flex items-center space-x-2 mr-6">
          <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center shadow">
            <span className="text-white font-bold text-sm">S</span>
          </div>
          <span className="hidden font-bold sm:inline-block">Skillora</span>
        </Link>
        
        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <nav className="flex items-center">
            {/* Desktop Navigation Links placeholder */}
          </nav>
          
          {isLoading ? (
            <div className="w-8 h-8 rounded-full bg-muted animate-pulse" />
          ) : user && !isError ? (
            <DropdownMenu>
              <DropdownMenuTrigger className="rounded-full outline-none focus-visible:ring-2 focus-visible:ring-ring">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={user.avatar_url || ""} alt={user.full_name} />
                    <AvatarFallback className="bg-primary/10 text-primary">
                      {user.full_name?.charAt(0).toUpperCase() || "U"}
                    </AvatarFallback>
                  </Avatar>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel className="font-normal">
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm font-medium leading-none">{user.full_name}</p>
                    <p className="text-xs leading-none text-muted-foreground">
                      {user.email}
                    </p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="cursor-pointer" onClick={() => window.location.href = "/dashboard"}>
                    <UserIcon className="mr-2 h-4 w-4" />
                    <span>Profile</span>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="text-destructive focus:text-destructive cursor-pointer" onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : null}
        </div>
      </div>
    </header>
  );
}
