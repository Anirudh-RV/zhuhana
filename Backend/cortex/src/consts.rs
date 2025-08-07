use std::env;

pub fn user_authentication_endpoint() -> String {
    let base = env::var("UASAM_URL").expect("UASAM_URL is not set");
    format!("{}/v1/user/authenticate/", base)
}
