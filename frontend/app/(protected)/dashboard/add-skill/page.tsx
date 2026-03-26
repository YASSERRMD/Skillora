"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import { CategorySelector } from "@/components/category-selector";
import { Category, Skill } from "@/lib/api/skills";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Sparkles, AlertCircle, CheckCircle2, ArrowLeft } from "lucide-react";
import api from "@/lib/api";
import Link from "next/link";

interface AppraisalResult {
  is_valid_skill: boolean;
  proficiency: number;
  credit_value: number;
  reasoning: string;
}

interface AppraisePayload {
  category_id: string;
  skill_id: string;
  description: string;
  category_name: string;
  skill_name: string;
}

export default function AddSkillPage() {
  const router = useRouter();
  const [category, setCategory] = React.useState<Category | null>(null);
  const [skill, setSkill] = React.useState<Skill | null>(null);
  const [description, setDescription] = React.useState("");

  const appraiseMutation = useMutation({
    mutationFn: async (payload: AppraisePayload) => {
      const { data } = await api.post<AppraisalResult>("/api/v1/skills/appraise", payload);
      return data;
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!category || !skill || description.length < 20) return;

    appraiseMutation.mutate({
      category_id: category.id,
      skill_id: skill.id,
      description: description,
      category_name: category.name,
      skill_name: skill.name,
    });
  };

  return (
    <div className="container mx-auto p-6 max-w-3xl">
      <div className="mb-6">
        <Link href="/dashboard" className="text-sm text-muted-foreground hover:text-primary flex items-center gap-1 w-fit">
          <ArrowLeft className="w-4 h-4" /> Back to Dashboard
        </Link>
      </div>

      <Card className="shadow-sm border-muted/60 relative overflow-hidden">
        {appraiseMutation.isPending && (
          <div className="absolute inset-0 bg-background/80 backdrop-blur-sm z-10 flex flex-col items-center justify-center">
            <Sparkles className="w-8 h-8 text-indigo-500 animate-pulse mb-4" />
            <p className="text-lg font-medium animate-pulse">AI Agent is appraising your skill...</p>
          </div>
        )}

        <CardHeader>
          <CardTitle>Offer a New Skill</CardTitle>
          <CardDescription>
            Select a category and describe your experience. Our AI Oracle will evaluate your proficiency and estimate its market value.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {!appraiseMutation.isSuccess && (
            <form onSubmit={handleSubmit} className="space-y-6">
              <CategorySelector onCategoryChange={setCategory} onSkillChange={setSkill} />

              <div className="space-y-2">
                <Label htmlFor="description">Experience & Portfolio</Label>
                <Textarea
                  id="description"
                  placeholder="Describe your background, years of experience, projects you've completed, and what you can teach others..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="min-h-[120px]"
                  required
                  minLength={20}
                />
                <p className="text-xs text-muted-foreground">
                  Minimum 20 characters. The more detailed your proof of work, the higher your proficiency rating.
                </p>
              </div>

              {appraiseMutation.isError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Appraisal Failed</AlertTitle>
                  <AlertDescription>
                    {appraiseMutation.error instanceof Error
                      ? appraiseMutation.error.message
                      : "The AI agent could not process your request at this time."}
                  </AlertDescription>
                </Alert>
              )}

              <Button
                type="submit"
                className="w-full bg-indigo-600 hover:bg-indigo-700 text-white"
                disabled={!category || !skill || description.length < 20 || appraiseMutation.isPending}
              >
                <Sparkles className="w-4 h-4 mr-2" />
                Submit for AI Appraisal
              </Button>
            </form>
          )}

          {appraiseMutation.isSuccess && appraiseMutation.data && (
            <div className="space-y-6 mt-2">
              {appraiseMutation.data.is_valid_skill ? (
                <Alert className="border-green-500/50 bg-green-500/10 text-green-700 dark:text-green-400">
                  <CheckCircle2 className="h-5 w-5 text-green-600 dark:text-green-400" />
                  <AlertTitle className="text-green-800 dark:text-green-300">Skill Verified!</AlertTitle>
                  <AlertDescription>
                    Your skill has been approved and added to your portfolio.
                  </AlertDescription>
                </Alert>
              ) : (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Skill Denied</AlertTitle>
                  <AlertDescription>
                    The AI oracle did not find sufficient evidence of teachable capability.
                  </AlertDescription>
                </Alert>
              )}

              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 rounded-lg bg-muted/50 border border-muted flex flex-col items-center justify-center text-center">
                  <p className="text-sm text-muted-foreground mb-1">Proficiency Level</p>
                  <p className="text-3xl font-bold">{appraiseMutation.data.proficiency} <span className="text-lg text-muted-foreground font-normal">/ 5</span></p>
                </div>
                <div className="p-4 rounded-lg bg-muted/50 border border-muted flex flex-col items-center justify-center text-center">
                  <p className="text-sm text-muted-foreground mb-1">Estimated Value</p>
                  <p className="text-3xl font-bold text-indigo-600 dark:text-indigo-400">{appraiseMutation.data.credit_value} <span className="text-lg text-muted-foreground font-normal">credits</span></p>
                </div>
              </div>

              <div className="p-4 rounded-lg bg-indigo-50 dark:bg-indigo-950/20 border border-indigo-100 dark:border-indigo-900/50">
                <p className="text-sm font-semibold text-indigo-900 dark:text-indigo-300 mb-2 flex items-center gap-2">
                  <Sparkles className="w-4 h-4" /> Oracle Reasoning
                </p>
                <p className="text-sm text-indigo-800 dark:text-indigo-400 italic">
                  "{appraiseMutation.data.reasoning}"
                </p>
              </div>

              <div className="flex gap-4 pt-4">
                <Button className="flex-1" variant="outline" onClick={() => appraiseMutation.reset()}>
                  Offer Another Skill
                </Button>
                <Button className="flex-1" onClick={() => router.push("/dashboard")}>
                  Return to Portfolio
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
