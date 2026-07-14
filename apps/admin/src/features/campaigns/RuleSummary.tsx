import type { Action, Condition, RuleList } from "./ruleTypes";

function describeCondition(condition: Condition | undefined): string {
  if (!condition) return "no conditions";

  if (condition.all) {
    return `ALL of: ${condition.all.map(describeCondition).join("; ")}`;
  }

  if (condition.any) {
    return `ANY of: ${condition.any.map(describeCondition).join("; ")}`;
  }

  if (condition.fact && condition.operator) {
    return `${condition.fact} ${condition.operator} ${JSON.stringify(condition.value)}`;
  }

  return "unknown condition";
}

function describeAction(action: Action): string {
  const params = action.params ? JSON.stringify(action.params) : "{}";
  return `${action.type}: ${params}`;
}

export function RuleSummary({ rules }: { rules: RuleList }) {
  if (rules.length === 0) {
    return <div className="rounded border border-gray-200 bg-gray-50 p-3 text-xs text-gray-500">No rules</div>;
  }

  return (
    <div className="space-y-2">
      {rules.map((rule, i) => (
        <div key={rule.id ?? i} className="space-y-1.5 rounded border border-gray-200 bg-gray-50 p-3 text-xs text-gray-600">
          {rule.name && <p className="font-medium text-gray-800">{rule.name}</p>}
          <p>
            <span className="font-medium text-gray-700">Conditions: </span>
            {describeCondition(rule.conditions)}
          </p>
          <p>
            <span className="font-medium text-gray-700">Actions: </span>
            {rule.actions && rule.actions.length > 0 ? rule.actions.map(describeAction).join(", ") : "none"}
          </p>
        </div>
      ))}
    </div>
  );
}
