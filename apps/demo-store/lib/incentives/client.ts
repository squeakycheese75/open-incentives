import "server-only";

import type { EvaluateApiResponse, EvaluateRequest } from "./types";

const EVALUATE_TIMEOUT_MS = 4000;

export class IncentivesApiError extends Error {
  constructor(
    message: string,
    public readonly status: "unauthorized" | "unavailable" | "invalid_request",
  ) {
    super(message);
    this.name = "IncentivesApiError";
  }
}

export interface IncentivesClient {
  evaluate(input: EvaluateRequest): Promise<EvaluateApiResponse>;
}

export class HttpIncentivesClient implements IncentivesClient {
  constructor(
    private readonly baseUrl: string,
    private readonly apiKey: string,
  ) {}

  async evaluate(input: EvaluateRequest): Promise<EvaluateApiResponse> {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), EVALUATE_TIMEOUT_MS);

    let response: Response;
    try {
      response = await fetch(`${this.baseUrl}/v1/evaluate`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `ApiKey ${this.apiKey}`,
        },
        body: JSON.stringify(input),
        signal: controller.signal,
      });
    } catch (err) {
      if (err instanceof Error && err.name === "AbortError") {
        throw new IncentivesApiError("evaluation request timed out", "unavailable");
      }
      throw new IncentivesApiError("incentives API unreachable", "unavailable");
    } finally {
      clearTimeout(timeout);
    }

    if (response.status === 401) {
      throw new IncentivesApiError("invalid API key", "unauthorized");
    }

    if (response.status === 400) {
      throw new IncentivesApiError("invalid evaluation request", "invalid_request");
    }

    if (!response.ok) {
      throw new IncentivesApiError(
        `incentives API returned status ${response.status}`,
        "unavailable",
      );
    }

    return (await response.json()) as EvaluateApiResponse;
  }
}

export function createIncentivesClient(): IncentivesClient {
  const baseUrl = process.env.INCENTIVES_API_URL;
  const apiKey = process.env.INCENTIVES_API_KEY;

  if (!baseUrl || !apiKey) {
    throw new IncentivesApiError(
      "INCENTIVES_API_URL and INCENTIVES_API_KEY must be set",
      "unavailable",
    );
  }

  return new HttpIncentivesClient(baseUrl, apiKey);
}
