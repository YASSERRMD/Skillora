import api from "../api";
import { UserSkillDetail } from "./skills";

export const getMySkills = async (): Promise<UserSkillDetail[]> => {
  const { data } = await api.get<UserSkillDetail[]>("/api/v1/users/skills");
  return data;
};
