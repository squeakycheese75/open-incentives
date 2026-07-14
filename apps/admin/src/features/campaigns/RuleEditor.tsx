import { useMemo } from "react";

import { RuleSummary } from "./RuleSummary";
import { isRuleList } from "./ruleTypes";

// MVP: a plain JSON textarea + a rendered read-only summary. Deliberately
// avoids a heavier code-editor dependency (Monaco/CodeMirror). This is
// structured as a self-contained "mode" so a future visual condition/action
// builder can replace just the editing surface below without touching the
// API layer or the parent form - it only ever needs to produce/consume the
// same raw JSON string via `value`/`onChange`.
interface RuleEditorProps {
  value: string;
  onChange: (next: string) => void;
  error?: string | null;
}

export function RuleEditor({ value, onChange, error }: RuleEditorProps) {
  const parsed = useMemo(() => {
    try {
      const json = JSON.parse(value);
      return isRuleList(json) ? json : null;
    } catch {
      return null;
    }
  }, [value]);

  return (
    <div className="space-y-2">
      <label htmlFor="rules-json" className="block text-sm font-medium text-gray-700">
        Rules (JSON)
      </label>
      <textarea
        id="rules-json"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        rows={14}
        spellCheck={false}
        className="w-full rounded border border-gray-300 bg-gray-900 p-3 font-mono text-xs text-gray-100 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
      />

      {error && <p className="text-sm text-red-600">{error}</p>}

      {parsed && <RuleSummary rules={parsed} />}
    </div>
  );
}
