export type OperationType =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage';

export interface CalculateRequest {
  operation: OperationType | string;
  a: number;
  b?: number;
}

export interface CalculateResponse {
  operation: string;
  a: number;
  b?: number;
  result: number;
  formatted: string;
}

export interface ApiErrorResponse {
  success: boolean;
  error: string;
  code?: string;
}

export interface HistoryItem {
  id: string;
  expression: string;
  result: number;
  timestamp: Date;
}
