import apiClient from './api';
import type {Relationship} from "@/types/relationship.ts";

// 创建关系记录
export const createRelationship = async (data: Partial<Relationship>): Promise<Relationship> => {
  const response = await apiClient.post<Relationship>('/relationship', data);
  return response.data;
};

// 根据ID获取关系记录
export const getRelationshipById = async (id: number): Promise<Relationship> => {
  const response = await apiClient.get<Relationship>(`/relationship/${id}`);
  return response.data;
};

// 根据账号ID获取关系记录列表
export const getRelationshipsByAccountId = async (accountId: number): Promise<Relationship[]> => {
  const response = await apiClient.get<Relationship[]>(`/relationship/account/${accountId}`);
  return response.data;
};

// 获取所有关系记录
export const getAllRelationships = async (): Promise<Relationship[]> => {
  const response = await apiClient.get<Relationship[]>('/relationship');
  return response.data;
};

// 更新关系记录
export const updateRelationship = async (id: number, data: Partial<Relationship>): Promise<Relationship> => {
  const response = await apiClient.put<Relationship>(`/relationship/${id}`, data);
  return response.data;
};

// 删除关系记录
export const deleteRelationship = async (id: number): Promise<void> => {
  await apiClient.delete(`/relationship/${id}`);
};

// 分析关系记录
// export const analyzeRelationship = async (id: number): Promise<string> => {
//
//   const response = await apiClient.post<{ analyze: string }>(`/relationship/${id}/analyze`);
//   return response.data.analyze;
// };
export const analyzeRelationship = (id: number, onMessage: (msg: string) => void, onDone?: () => void) => {
  const url = `/api/relationship/${id}/analyze`;

  const eventSource = new EventSource(url);

  eventSource.onmessage = (event) => {
    onMessage(event.data);
  };

  eventSource.onerror = (err) => {
    console.error("SSE 错误:", err);
    eventSource.close();
    if (onDone) onDone();
  };

  // 如果后端发送了 event: done，可在此关闭连接
  eventSource.addEventListener("done", () => {
    eventSource.close();
    if (onDone) onDone();
  });

  // 返回 eventSource 供上层控制
  return eventSource;
};

// 根据选中的表创建关系记录
export const createRelationshipFromSelectedTables = async (
  accountId: number,
  selectedTables: any[],
  title: string
): Promise<Relationship> => {
  // 将选中的表信息转换为需要的格式
  const tablesData = JSON.stringify(selectedTables);
  
  const data = {
    account_id: accountId,
    title: title || `关系分析 - ${new Date().toLocaleString()}`,
    tables: tablesData,
    ddls: '', // 可以根据需要添加表结构信息
    result: '' // 初始结果为空
  };

  const response = await apiClient.post<Relationship>('/relationship', data);
  return response.data;
};