import api from "../api";

export interface LLMProvider {
  id: string;
  provider_name: string;
  model_name: string;
  use_case: string;
  priority: number;
  is_active: boolean;
  created_at: string;
}

export interface AddProviderPayload {
  provider_name: string;
  model_name: string;
  api_key: string;
  use_case: string;
  priority: number;
}

export const getLLMProviders = async (): Promise<LLMProvider[]> => {
  const { data } = await api.get<LLMProvider[]>("/api/v1/admin/llm-providers");
  return data;
};

export const addLLMProvider = async (payload: AddProviderPayload): Promise<LLMProvider> => {
  const { data } = await api.post<LLMProvider>("/api/v1/admin/llm-providers", payload);
  return data;
};
