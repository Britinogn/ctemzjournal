/** Setups — reference CRUD. Names are unique per user. */

export interface Setup {
  ID: string;
  UserID: string;
  Name: string;
  Rules: string | null;
  Invalidation: string | null;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface SetupCreate {
  name: string;
  rules?: string;
  invalidation?: string;
}

export interface SetupUpdate {
  name?: string;
  rules?: string;
  invalidation?: string;
}
