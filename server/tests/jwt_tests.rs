use hbb_common::tokio;
use hbbs::jwt;

#[test]
fn test_generate_token() {
    std::env::set_var("RUSTDESK_API_JWT_KEY", "testjwt");
    let token = jwt::generate_token(1, 3600).unwrap();
    assert!(!token.is_empty(), "Generated token should not be empty");
}

#[tokio::test]
async fn test_verify_token() {
    std::env::set_var("RUSTDESK_API_JWT_KEY", "testjwt");
    let token = jwt::generate_token(1, 2).unwrap();
    // let token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MzY4NzA4NjF9.u5pxmwNMrYUwtkspF1FuZj-R5ANAR9WT9_dMHuQhV0Y";

    // hbb_common::sleep(3f32).await;
    // println!("Token : {:?}, now: {:?}", token, chrono::Utc::now().timestamp());
    let result = jwt::verify_token(&token);
    assert!(result.is_ok(), "Token should be valid");
}

#[test]
fn test_reject_invalid_claims() {
    use jsonwebtoken::{encode, Algorithm, EncodingKey, Header};
    use serde_json::json;

    std::env::set_var("RUSTDESK_API_JWT_KEY", "testjwt");
    let future = chrono::Utc::now().timestamp() + 3600;
    let cases = [
        (Algorithm::HS256, json!({"user_id": 0, "exp": future})),
        (Algorithm::HS256, json!({"user_id": 1, "exp": 0})),
        (Algorithm::HS256, json!({"user_id": 1})),
        (Algorithm::HS384, json!({"user_id": 1, "exp": future})),
    ];
    for (algorithm, claims) in cases {
        let token = encode(
            &Header::new(algorithm),
            &claims,
            &EncodingKey::from_secret(b"testjwt"),
        )
        .unwrap();
        assert!(jwt::verify_token(&token).is_err());
    }
}
