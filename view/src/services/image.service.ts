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
}
