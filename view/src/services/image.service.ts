import { CommonService } from "@/services";
import { CLIENT_ENDPOINTS } from "@/utils/constants";

export class ImageService extends CommonService {
  async getImages(query: GetImagesQuery): Promise<GetImagesResult> {
    try {
      const path = this.createPathWithQuery(CLIENT_ENDPOINTS.images, query);
      const config = this.createConfig("GET");
      const res = await fetch(path, config);
      if (!res.ok) {
        const body: ErrRespBody = await res.json();
        return { ...body, ok: false };
      }
      const body: GetImagesRespBody = await res.json();
      return { ...body, ok: true };
    } catch {
      return this.returnUnknownError();
    }
  }

  async postImage(body: PostImageReqBody): Promise<PostImageResult> {
    try {
      const config = this.createConfig("POST", body);
      const res = await fetch(CLIENT_ENDPOINTS.images, config);
      if (!res.ok) {
        const body: ErrRespBody = await res.json();
        return { ...body, ok: false };
      }
      const resBody: PostImageRespBody = await res.json();
      return { ...resBody, ok: true };
    } catch {
      return this.returnUnknownError();
    }
  }

  async patchImage(
    id: string,
    body: PatchImageReqBody
  ): Promise<PatchImageResult> {
    try {
      const config = this.createConfig("PATCH", body);
      const res = await fetch(CLIENT_ENDPOINTS.images + "/" + id, config);
      if (!res.ok) {
        const body: ErrRespBody = await res.json();
        return { ...body, ok: false };
      }
      return { ok: true };
    } catch {
      return this.returnUnknownError();
    }
  }

  async deleteImage(id: string): Promise<DeleteImageResult> {
    try {
      const config = this.createConfig("DELETE");
      const res = await fetch(CLIENT_ENDPOINTS.images + "/" + id, config);
      if (!res.ok) {
        const body: ErrRespBody = await res.json();
        return { ...body, ok: false };
      }
      return { ok: true };
    } catch {
      return this.returnUnknownError();
    }
  }
}
