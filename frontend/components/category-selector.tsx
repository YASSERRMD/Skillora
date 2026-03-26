"use client";

import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { getCategories, getCategorySkills, Category, Skill } from "@/lib/api/skills";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Props {
  onCategoryChange: (category: Category | null) => void;
  onSkillChange: (skill: Skill | null) => void;
}

export function CategorySelector({ onCategoryChange, onSkillChange }: Props) {
  const [selectedCatId, setSelectedCatId] = React.useState<string>("");

  // Fetch parent categories
  const { data: categories, isLoading: isLoadingCats } = useQuery({
    queryKey: ["categories"],
    queryFn: getCategories,
  });

  // Fetch children skills when subset category is populated
  const { data: skills, isLoading: isLoadingSkills } = useQuery({
    queryKey: ["skills", selectedCatId],
    queryFn: () => getCategorySkills(selectedCatId),
    enabled: !!selectedCatId,
  });

  const handleCategoryChange = (id: string | null) => {
    setSelectedCatId(id || "");
    onSkillChange(null); // Reset child form control
    if (categories && id) {
      const cat = categories.find((c) => c.id === id);
      onCategoryChange(cat || null);
    } else {
      onCategoryChange(null);
    }
  };

  const handleSkillChange = (id: string | null) => {
    if (skills && id) {
      const skill = skills.find((s) => s.id === id);
      onSkillChange(skill || null);
    } else {
      onSkillChange(null);
    }
  };

  return (
    <div className="flex flex-col space-y-4 md:flex-row md:space-y-0 md:space-x-4 w-full">
      <div className="flex-1 flex flex-col space-y-2">
        <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
          Domain Category
        </label>
        <Select onValueChange={handleCategoryChange} value={selectedCatId}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder={isLoadingCats ? "Loading..." : "Select a category"} />
          </SelectTrigger>
          <SelectContent>
            {categories?.map((c) => (
              <SelectItem key={c.id} value={c.id}>
                {c.name}
              </SelectItem>
            ))}
            {(!categories || categories.length === 0) && !isLoadingCats && (
              <SelectItem value="none" disabled>
                No categories found
              </SelectItem>
            )}
          </SelectContent>
        </Select>
      </div>

      <div className="flex-1 flex flex-col space-y-2">
        <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
          Specific Skill
        </label>
        <Select
          onValueChange={handleSkillChange}
          disabled={!selectedCatId || isLoadingSkills}
        >
          <SelectTrigger className="w-full">
            <SelectValue
              placeholder={
                !selectedCatId
                  ? "Select a category first"
                  : isLoadingSkills
                  ? "Loading skills..."
                  : "Select a skill"
              }
            />
          </SelectTrigger>
          <SelectContent>
            {skills?.map((s) => (
              <SelectItem key={s.id} value={s.id}>
                {s.name}
              </SelectItem>
            ))}
            {skills && skills.length === 0 && !isLoadingSkills && (
              <SelectItem value="none" disabled>
                No skills found in this category
              </SelectItem>
            )}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
