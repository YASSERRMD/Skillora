"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { searchSkills, SkillMatch } from "@/lib/api/matching";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Search, Sparkles, Users, Star, Coins } from "lucide-react";
import Link from "next/link";

export default function MarketplacePage() {
  const [query, setQuery] = React.useState("");
  const [activeQuery, setActiveQuery] = React.useState("");

  const { data, isFetching, isError } = useQuery({
    queryKey: ["matches", activeQuery],
    queryFn: () => searchSkills(activeQuery),
    enabled: !!activeQuery,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setActiveQuery(query.trim());
  };

  const proficiencyLabel = (level: number) => {
    const labels = ["", "Beginner", "Intermediate", "Advanced", "Expert", "Master"];
    return labels[level] || "Unknown";
  };

  const proficiencyColor = (level: number) => {
    const colors = ["", "bg-slate-500", "bg-blue-500", "bg-indigo-500", "bg-purple-500", "bg-amber-500"];
    return colors[level] || "bg-gray-400";
  };

  return (
    <div className="container mx-auto p-6 max-w-5xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight mb-2">Skill Marketplace</h1>
        <p className="text-muted-foreground">
          Search the knowledge economy. Our AI oracle finds the best matches for what you want to learn.
        </p>
      </div>

      <form onSubmit={handleSearch} className="flex gap-3 mb-8">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            id="marketplace-search"
            type="text"
            placeholder="e.g. machine learning, UI design, guitar..."
            className="pl-9"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <Button type="submit" disabled={!query.trim() || isFetching}>
          <Sparkles className="h-4 w-4 mr-2" />
          {isFetching ? "Searching..." : "Find Matches"}
        </Button>
      </form>

      {isError && (
        <div className="text-center py-16 text-destructive">
          The matching engine encountered an error. Please try again.
        </div>
      )}

      {data && data.results.length === 0 && (
        <div className="text-center py-16">
          <Sparkles className="h-10 w-10 text-muted-foreground/40 mx-auto mb-4" />
          <h3 className="text-lg font-medium mb-2">No matches found</h3>
          <p className="text-sm text-muted-foreground">
            Try a different query. Our AI will find semantic matches even if the exact words differ.
          </p>
        </div>
      )}

      {!activeQuery && !isFetching && (
        <div className="text-center py-16">
          <div className="w-16 h-16 rounded-full bg-indigo-100 dark:bg-indigo-950/30 flex items-center justify-center mx-auto mb-4">
            <Search className="h-7 w-7 text-indigo-500" />
          </div>
          <h3 className="text-lg font-medium mb-2">Discover What You Want to Learn</h3>
          <p className="text-sm text-muted-foreground max-w-sm mx-auto">
            Enter any topic above and our semantic AI engine will surface the most relevant skill holders in the community.
          </p>
        </div>
      )}

      {data && data.results.length > 0 && (
        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Showing <span className="font-medium text-foreground">{data.results.length}</span> top matches for "{data.query}"
          </p>
          {data.results.map((match: SkillMatch) => (
            <Card key={`${match.skill_id}-${match.owner_id}`} className="border-muted/60 shadow-sm hover:shadow-md transition-shadow duration-200">
              <CardContent className="p-5">
                <div className="flex flex-col md:flex-row justify-between gap-4">
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 flex-wrap">
                      <h3 className="font-semibold text-lg leading-none">{match.skill_name}</h3>
                      <Badge variant="secondary" className="text-xs">{match.category_name}</Badge>
                      <Badge className={`text-white text-xs ${proficiencyColor(match.proficiency_level)}`}>
                        {proficiencyLabel(match.proficiency_level)}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Users className="h-3.5 w-3.5" /> {match.owner_id.slice(0, 8)}...
                      </span>
                      <span className="flex items-center gap-1">
                        <Star className="h-3.5 w-3.5 text-amber-500" />
                        {(match.similarity_score * 100).toFixed(1)}% match
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center justify-between md:justify-end gap-4">
                    <div className="flex items-center gap-1.5">
                      <Coins className="h-4 w-4 text-indigo-500" />
                      <span className="font-bold text-indigo-600 dark:text-indigo-400">{match.credit_value}</span>
                      <span className="text-sm text-muted-foreground">credits</span>
                    </div>
                    <Link href={`/dashboard`}>
                      <Button size="sm" variant="outline">
                        Propose Barter
                      </Button>
                    </Link>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
