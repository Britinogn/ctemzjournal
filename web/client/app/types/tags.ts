/** Tags — reference CRUD, including mistake tags. Names unique per user. */

export type TagKind = 'mistake' | 'general';

export interface Tag {
  ID: string;
  UserID: string;
  Name: string;
  Kind: TagKind;
  CreatedAt: string;
}

export interface TagCreate {
  name: string;
  kind?: TagKind;
}

export interface TagUpdate {
  name?: string;
  kind?: TagKind;
}
