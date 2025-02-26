export interface Client {
  ID: string;
  IP: string;
}

export interface Metrics {
  disk_total: number;
  disk_free: number;
  memory_total: number;
  memory_available: number;
  processor: string;
  os: string;
}
