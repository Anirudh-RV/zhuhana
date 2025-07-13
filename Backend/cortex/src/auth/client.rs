use bytes::Bytes;
use futures_util::{Stream, TryStreamExt};
use reqwest::Client;
use std::io;

/// System prompt that tells the model to format code in Markdown with Python fences
const SYSTEM_PROMPT: &str = r#"You are a helpful assistant. Always respond in Markdown.
When showing code, use triple backticks and specify the language, for example, for python, you would write:

```python
def greet(name):
    print(f"Hello, {name}!")
```

"#;

pub async fn query_ollama_stream(
    user_prompt: String,
) -> Result<impl Stream<Item = Result<Bytes, io::Error>>, reqwest::Error> {
    let full_prompt = format!(
    "{system}\n\nUser: {user}",
    system = SYSTEM_PROMPT,
    user = user_prompt
    );

    let client = Client::new();
    let payload = serde_json::json!({
        "model": "llama3:8b-instruct-q4_0",
        "prompt": full_prompt,
        "stream": true
    });

    let res = client
        .post("http://ollama:11434/api/generate")
        .json(&payload)
        .send()
        .await?;

    let byte_stream = res
        .bytes_stream()
        .map_ok(Bytes::from)
        .map_err(|e| io::Error::new(io::ErrorKind::Other, e));

    Ok(byte_stream)
}
