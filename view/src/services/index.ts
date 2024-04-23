import { toSnakeCase } from "@/utils/converter";

export class CommonService {
  createPathWithQuery(path: string, query: RequestQuery): string {
    const queryParams: string[] = [];
    Object.entries(query).forEach(([key, value]) => {
      const snakeKey = toSnakeCase(key);
      if (Array.isArray(value) && value.length > 0) {
        queryParams.push(`${snakeKey}=${value.join(",")}`);
      } else {
        queryParams.push(`${snakeKey}=${value}`);
      }
    });
    const queryString = queryParams.join("&");
    return `${path}?${queryString}`;
  }

  createConfig(method: RequestMethod, body?: RequestBody): RequestInit {
    return {
      method,
      headers: {
        "Content-Type": "application/json",
      },
      cache: "no-store",
      body: body ? JSON.stringify(body) : undefined,
    };
  }

  returnUnknownError(): ErrResult {
    return {
      errorCode: "unknown",
      errorMessage: "an unknown error occured",
      ok: false,
    };
  }
}
