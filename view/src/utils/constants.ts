export const AUTH_ENDPOINTS = {
  login: "/auth-api/login",
} as const;

export const CLIENT_ENDPOINTS = {
  admin: "/client-api/admin",
  images: "/client-api/images",
} as const;

export const PAGE_ENDPOINTS = {
  home: "/",
  createImage: "/create-image",
  privacyPolicy: "/privacy-policy",
  termsOfService: "/terms-of-service",
} as const;

export const MAX_IMAGES_FETCH_COUNT = 9;
export const IMAGE_SIZE = 300;
export const MAX_KEYWORD_LENGTH = 50;

export const ACTIVE_TAB_ID = {
  latest: "latest",
  popular: "popular",
  favorite: "favorite",
} as const;

export const PATCH_IMAGE_REQUEST_TYPE = {
  used: "used",
  report: "report",
  confirm: "confirm",
} as const;
