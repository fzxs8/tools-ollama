export interface ModelParams {
  temperature: number;
  topP: number;
  context: number;
  numPredict: number;
  topK: number;
  repeatPenalty: number;
  outputMode: 'stream' | 'blocking';
}

export interface OnlineModel {
  model_identifier: string;
  namespace: string | null;
  model_name: string;
  model_type: string;
  description: string;
  capability: string | null;
  labels: string[];
  pulls: number;
  tags: number;
  last_updated: string;
  last_updated_str: string;
  url: string;
}

export interface DownloadProgress {
  model: string;
  status: string;
  percentage: number;
  notification?: any;
}
