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

// Fetch all parent categories
export const getCategories = async (): Promise<Category[]> => {
  const { data } = await api.get("/api/v1/categories");
  return data;
};

// Fetch skills mapped to a specific category
export const getCategorySkills = async (categoryId: string): Promise<Skill[]> => {
  const { data } = await api.get(`/api/v1/categories/${categoryId}/skills`);
  return data;
};
