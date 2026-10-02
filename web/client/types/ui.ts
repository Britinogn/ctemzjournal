/**
 * Shared UI primitives. Styling comes only from design tokens —
 * components never contain hex values.
 */
import type { Component } from 'vue';

export type Theme = 'light' | 'dark' | 'system';

export interface NavItem {
  label: string;
  to: string;
  /** Hugeicons component (from @hugeicons/vue). */
  icon: Component;
  active?: boolean;
}

export interface StatCard {
  label: string;
  /** Formatted display value, e.g. "+$7,928" or "60%". */
  value: string;
  /** Semantic tone — never color alone (pair with sign/arrow). */
  tone?: 'neutral' | 'profit' | 'loss' | 'warning';
  sub?: string;
  icon?: Component;
}

export interface ConfirmOptions {
  title: string;
  message: string;
  confirmLabel?: string;
  danger?: boolean;
}

export interface TableColumn<T> {
  key: string;
  label: string;
  /** mono applies JetBrains Mono with tabular figures (prices, R). */
  mono?: boolean;
  render?: (row: T) => string;
}
