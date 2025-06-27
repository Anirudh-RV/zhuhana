use bytes::Bytes;
use futures_util::{Stream, StreamExt, TryStreamExt};
use http_body_util::BodyExt;
use reqwest::Client;
use std::io;
use tokio_stream::wrappers::ReceiverStream;

pub async fn query_ollama_stream(
    prompt: String,
) -> Result<impl Stream<Item = Result<Bytes, io::Error>>, reqwest::Error> {
    let client = Client::new();
    let payload = serde_json::json!({
        "model": "llama3:8b-instruct-q4_0",
        "prompt": prompt,
        "stream": true
    });

    let res = client
        .post("http://ollama:11434/api/generate")
        .json(&payload)
        .send()
        .await?;

    let byte_stream = res
        .bytes_stream()
        .map_ok(|chunk| Bytes::from(chunk))
        .map_err(|e| io::Error::new(io::ErrorKind::Other, e));

    Ok(byte_stream)
}
