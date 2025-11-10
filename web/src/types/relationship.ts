export  interface Relationship {
  id: number;
  account_id: number;
  title: string;
  tables: string; // JSON格式存储所有表信息
  ddls: string;   // JSON格式存储所有表结构
  result: string; // 结果存储
  created_time: string;
  updated_time: string;
}

export interface TableInfo {
  name: string;
  comment: string;
  selected?: boolean;
}