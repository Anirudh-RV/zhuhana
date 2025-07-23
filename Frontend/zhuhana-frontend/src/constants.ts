export const ENV = "DEV";
export const UASAM_ENDPOINT = "http://localhost:8002";
export const GOVERNOR_ENDPOINT = "http://localhost:8008";
export const CORTEX_ENDPOINT = "http://localhost:3000";

export const SIGN_UP_V1_INIT_ENDPOINT =
  UASAM_ENDPOINT + "/v1/user/sign-up/init/";

export const GET_NOTIFICATIONS_V1_INIT_ENDPOINT =
  UASAM_ENDPOINT + "/v1/notification/list/";

export const READ_NOTIFICATIONS_V1_INIT_ENDPOINT =
  UASAM_ENDPOINT + "/v1/notification/read/";

export const PASSWORD_RESET_V1_INIT_ENDPOINT =
  UASAM_ENDPOINT + "/v1/user/reset-password/init/";

export const PASSWORD_RESET_V1_RESET_ENDPOINT =
  UASAM_ENDPOINT + "/v1/user/reset-password/reset/";

export const SIGN_UP_V1_VERIFY_OTP_ENDPOINT =
  UASAM_ENDPOINT + "/v1/user/sign-up/verify-otp/";

export const LOGIN_V1_VERIFY_PASSWORD_ENDPOINT =
  UASAM_ENDPOINT + "/v1/user/login/verify-password/";

export const LOGIN_V1_VERIFY_OTP_ENDPOINT =
  UASAM_ENDPOINT + "/v1/user/login/verify-otp/";

export const USER_PYTHON_ALGORITHM_UPLOAD_V1_ENDPOINT =
  GOVERNOR_ENDPOINT + "/v1/user/algorithm/python/upload/";

export const USER_PYTHON_ALGORITHMS_INFORMATION_V1_ENDPOINT =
  GOVERNOR_ENDPOINT + "/v1/user/algorithm/";

export const USER_PYTHON_ALGORITHM_INFORMATION_V1_ENDPOINT =
  GOVERNOR_ENDPOINT + "/v1/user/algorithm/info/";

export const CREATE_CHAT_SESSION_V1_ENDPOINT = CORTEX_ENDPOINT + "/v1/session/";
export const ASK_LLM_V1_ENDPOINT = CORTEX_ENDPOINT + "/v1/ask/";
export const GET_MESSAGES_V1_ENDPOINT = CORTEX_ENDPOINT + "/v1/messages/";
