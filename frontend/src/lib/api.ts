export type User = { id: string; role: "admin" | "ta" | "student" };
export type Person = User & { email: string };
export type Classroom = { id: string; name: string; members?: Person[] };
export type Lab = {
  id: string;
  title: string;
  statement: string;
  classId?: string;
  requirements?: string[];
};
export type Submission = {
  id: string;
  labId: string;
  title?: string;
  createdAt?: string;
  status: string;
  verdict: string;
  detail?: string;
  source?: string;
};

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
  }
}

// Preserve useful errors even when a reverse proxy sends an HTML response.
export async function api<T>(
  path: string,
  body?: object,
  key?: string,
): Promise<T> {
  let response: Response;

  try {
    response = await fetch(`/api${path}`, {
      method: body ? "POST" : "GET",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json",
        ...(key ? { "Idempotency-Key": key } : {}),
      },
      body: body ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError(
      "Cannot connect to the server. Check your connection and retry.",
      0,
    );
  }

  if (!response.ok) {
    const problem = response.headers.get("content-type")?.includes("json")
      ? await response.json().catch(() => null)
      : null;

    throw new ApiError(
      problem?.detail ||
        `The server could not complete this request (${response.status}).`,
      response.status,
    );
  }

  return response.status === 204 ? (undefined as T) : response.json();
}
