import { CommonService } from "@/services";
import { AUTH_ENDPOINTS } from "@/utils/constants";

type RespBody = {
  redirectURL: string;
};

type ErrRespBody = {
  errorCode: string;
  errorMessage: string;
};

type PostLoginSuccessResult = RespBody & { ok: true };

type PostLoginErrorResult = ErrRespBody & { ok: false };

type PostLoginResult = PostLoginSuccessResult | PostLoginErrorResult;

export class LoginService extends CommonService {
  async postLogin(
    username: string,
    password: string
  ): Promise<PostLoginResult> {
    try {
      const formData = new URLSearchParams();
      formData.append("username", username);
      formData.append("password", password);
      const result = await fetch(AUTH_ENDPOINTS.login, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData.toString(),
      });
      if (!result.ok) {
        const body: ErrRespBody = await result.json();
        return { ...body, ok: false };
      }
      const body: RespBody = await result.json();
      return { ...body, ok: true };
    } catch {
      return this.returnUnknownError();
    }
  }
}
