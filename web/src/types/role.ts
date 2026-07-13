import type { Id } from "./common";

export interface RoleRecord {
  id: Id;
  code?: string;
  name: string;
  description?: string;
  sort?: number;
  enabled?: boolean;
  permissionIDs?: Id[];
}

export interface RolePayload extends Partial<RoleRecord> {
  id?: Id;
  permissionIDs?: Id[];
  [key: string]: unknown;
}
