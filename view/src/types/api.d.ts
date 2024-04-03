type RequestMethod = "GET" | "POST" | "PATCH" | "DELETE";

type RequestBody = Record<string, string | number>;

type RequestQuery = Record<string, string | number | string[] | boolean | null>;

type ErrRespBody = {
  errorCode: string;
  errorMessage: string;
};

type ErrResult = {
  errorCode: string;
  errorMessage: string;
  ok: false;
};
