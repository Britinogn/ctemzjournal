/** Accounts — reference CRUD. Names are unique per user. */

export type AccountType = 'demo' | 'live';

/** V1 allows naira and USD. The API stores any currency string. */
export type AccountCurrency = 'USD' | 'NGN';

export interface Account {
  ID: string;
  UserID: string;
  Name: string;
  Type: AccountType;
  Currency: string;
  StartingBalance: number;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface AccountCreate {
  name: string;
  type: AccountType;
  currency?: AccountCurrency;
  starting_balance?: number;
}

export interface AccountUpdate {
  name?: string;
  type?: AccountType;
  currency?: AccountCurrency;
  starting_balance?: number;
}
