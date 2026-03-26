import api from "../api";

export interface UserProfile {
  id: string;
  full_name: string;
  email: string;
  avatar_url: string;
  bio: string;
}

export interface UpdateProfilePayload {
  full_name?: string;
  bio?: string;
}

export const getMyProfile = async (): Promise<UserProfile> => {
  const { data } = await api.get<UserProfile>("/api/v1/users/me");
  return data;
};

export const updateMyProfile = async (payload: UpdateProfilePayload): Promise<UserProfile> => {
  const { data } = await api.put<UserProfile>("/api/v1/users/me", payload);
  return data;
};
