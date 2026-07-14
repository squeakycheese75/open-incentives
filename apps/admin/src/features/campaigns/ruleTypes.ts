export interface Condition {
  fact?: string;
  operator?: string;
  value?: unknown;
  all?: Condition[];
  any?: Condition[];
}

export interface Action {
  type: string;
  params?: Record<string, unknown>;
}

export interface Rule {
  id?: string;
  name?: string;
  conditions?: Condition;
  actions?: Action[];
}

// A campaign's `rules` field is a JSON array of Rule objects (see
// engine/model.go: `Rules []Rule`), not a single rule.
export type RuleList = Rule[];

export function isRuleList(value: unknown): value is RuleList {
  return Array.isArray(value);
}
