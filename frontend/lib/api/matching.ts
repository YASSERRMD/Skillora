import api from "../api";

export interface SkillMatch {
  skill_id: string;
  skill_name: string;
  category_name: string;
  owner_id: string;
  proficiency_level: number;
  credit_value: number;
  similarity_score: number;
}

export interface MatchResponse {
  query: string;
  results: SkillMatch[];
}

export const searchSkills = async (query: string, limit = 10): Promise<MatchResponse> => {
  const { data } = await api.get<MatchResponse>("/api/v1/match", {
    params: { q: query, limit },
  });
  return data;
};
