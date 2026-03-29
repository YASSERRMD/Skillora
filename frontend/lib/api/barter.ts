import api from "../api";

export interface BarterTransaction {
  id: string;
  initiator_id: string;
  receiver_id: string;
  initiator_skill_id: string;
  receiver_skill_id: string;
  credit_amount: number;
  status: "pending" | "accepted" | "completed" | "cancelled";
  created_at: string;
  updated_at: string;
}

export interface Milestone {
  id: string;
  barter_id: string;
  title: string;
  description: string;
  credit_portion: number;
  status: "pending" | "completed" | "approved";
  created_at: string;
  updated_at: string;
}

export interface ProposeBarterPayload {
  receiver_id: string;
  initiator_skill_id: string;
  receiver_skill_id: string;
  credit_amount: number;
}

export const getMyBarters = async (): Promise<BarterTransaction[]> => {
  const { data } = await api.get<BarterTransaction[]>("/api/v1/barters");
  return data;
};

export const proposeBarter = async (payload: ProposeBarterPayload): Promise<BarterTransaction> => {
  const { data } = await api.post<BarterTransaction>("/api/v1/barters", payload);
  return data;
};

export const updateBarterStatus = async (barterID: string, status: "accepted" | "cancelled"): Promise<void> => {
  await api.patch(`/api/v1/barters/${barterID}/status`, { status });
};

export const getCreditBalance = async (): Promise<number> => {
  const { data } = await api.get<{ credit_balance: number }>("/api/v1/barters/balance");
  return data.credit_balance;
};

export const getBarterMilestones = async (barterID: string): Promise<Milestone[]> => {
  const { data } = await api.get<Milestone[]>(`/api/v1/barters/${barterID}/milestones`);
  return data;
};

export const completeMilestone = async (id: string): Promise<void> => {
  await api.post(`/api/v1/milestones/${id}/complete`);
};

export const approveMilestone = async (id: string): Promise<void> => {
  await api.post(`/api/v1/milestones/${id}/approve`);
};
