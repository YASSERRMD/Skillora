import api from "../api";

export interface Category {
  id: string;
  name: string;
  slug: string;
}

export interface Skill {
  id: string;
  category_id: string;
  name: string;
  description: string;
}
export interface UserSkillDetail {
  user_id: string;
  skill_id: string;
  skill_name: string;
  category_name: string;
  proficiency_level: number;
  credit_value: number;
  is_verified: boolean;
}


export const getCategories = async (): Promise<Category[]> => {
  const { data } = await api.get("/api/v1/categories");
  return data;
};

// Fetch skills mapped to a specific category
export const getCategorySkills = async (categoryId: string): Promise<Skill[]> => {
  const { data } = await api.get(`/api/v1/categories/${categoryId}/skills`);
  return data;
};
