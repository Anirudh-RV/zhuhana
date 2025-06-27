import { useEffect, useRef, useState } from "react";
import {
  Box,
  Typography,
  TextField,
  IconButton,
  Paper,
  Tooltip,
} from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import StopIcon from "@mui/icons-material/Stop";
import ReactMarkdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";

type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

type LLMPanelProps = {
  onSend: (
    input: string,
    onChunk: (token: string) => void,
    signal: AbortSignal
  ) => Promise<void>;
};

export default function LLMPanel({ onSend }: LLMPanelProps) {
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const controllerRef = useRef<AbortController | null>(null);
  const bottomRef = useRef<HTMLDivElement | null>(null);

  const handleSend = () => {
    if (!input.trim() || isStreaming) return;

    const userMessage: Message = { role: "user", content: input };
    const botMessage: Message = { role: "assistant", content: "" };

    setMessages((prev) => [...prev, userMessage, botMessage]);

    const controller = new AbortController();
    controllerRef.current = controller;
    setIsStreaming(true);

    let currentBotResponse = "";

    onSend(
      input,
      (token: string) => {
        currentBotResponse += token;
        setMessages((prev) =>
          prev.map((msg, i) =>
            i === prev.length - 1
              ? { ...msg, content: currentBotResponse }
              : msg
          )
        );
      },
      controller.signal
    )
      .catch((err) => {
        if (err.name !== "AbortError") {
          console.error("Streaming error:", err);
          setMessages((prev) => [
            ...prev,
            {
              role: "system",
              content: `---\n❌ ${err.message}`,
            },
          ]);
        }
      })
      .finally(() => {
        setIsStreaming(false);
        controllerRef.current = null;
        // ✅ Always add a divider after assistant message
        setMessages((prev) => [...prev, { role: "system", content: "---" }]);
      });

    setInput("");
  };

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const CodeBlock = ({ inline, className, children }: any) => {
    return !inline ? (
      <Box
        component="pre"
        sx={{
          backgroundColor: "#1e1e1e",
          color: "#f8f8f2",
          borderRadius: 1,
          padding: 2,
          overflowX: "auto",
          fontSize: "0.875rem",
          fontFamily: "monospace",
        }}
      >
        <Box component="code" className={className}>
          {children}
        </Box>
      </Box>
    ) : (
      <code
        style={{
          backgroundColor: "#2e2e2e",
          padding: "2px 4px",
          borderRadius: "4px",
          fontFamily: "monospace",
        }}
      >
        {children}
      </code>
    );
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Typography variant="h6" gutterBottom>
        LLM Chat
      </Typography>

      <Paper
        elevation={1}
        sx={{
          borderRadius: 2,
          p: 2,
          flexGrow: 1,
          overflowY: "auto",
          mb: 2,
          display: "flex",
          flexDirection: "column",
          gap: 2,
          backgroundColor: "#0A0F1A",
        }}
      >
        {messages.map((msg, idx) => (
          <Box key={idx}>
            {msg.role === "system" && msg.content.trim() === "---" ? (
              <Box
                sx={{
                  width: "100%",
                  borderTop: "1px solid #333",
                  my: 1,
                }}
              />
            ) : msg.role === "system" ? (
              <Box
                sx={{
                  textAlign: "center",
                  color: "#888",
                  fontSize: "0.8rem",
                  py: 1,
                  whiteSpace: "pre-wrap",
                }}
              >
                <ReactMarkdown>{msg.content}</ReactMarkdown>
              </Box>
            ) : (
              <Box
                sx={{
                  alignSelf: msg.role === "user" ? "flex-end" : "flex-start",
                  maxWidth: "100%",
                  width: "100%",
                  display: "flex",
                  justifyContent:
                    msg.role === "user" ? "flex-end" : "flex-start",
                }}
              >
                {msg.role === "user" ? (
                  <Box
                    sx={{
                      backgroundColor: "#0C1018",
                      color: "#ffffff",
                      px: 2,
                      py: 1,
                      borderRadius: 2,
                      maxWidth: "75%",
                      boxShadow: "0 0 4px rgba(255, 255, 255, 0.05)",
                    }}
                  >
                    <Typography
                      variant="caption"
                      sx={{ fontWeight: "bold", color: "#888" }}
                    >
                      You
                    </Typography>
                    <Typography variant="body1">{msg.content}</Typography>
                  </Box>
                ) : (
                  <Box
                    sx={{
                      backgroundColor: "transparent",
                      color: "#ddd",
                      fontSize: "0.95rem",
                      maxWidth: "100%",
                      px: 1,
                    }}
                  >
                    <ReactMarkdown
                      rehypePlugins={[rehypeHighlight]}
                      components={{ code: CodeBlock }}
                    >
                      {msg.content}
                    </ReactMarkdown>
                  </Box>
                )}
              </Box>
            )}
          </Box>
        ))}
        <div ref={bottomRef} />
      </Paper>

      <Box
        component="form"
        onSubmit={(e) => {
          e.preventDefault();
          handleSend();
        }}
        sx={{
          display: "flex",
          alignItems: "flex-end",
          gap: 1,
        }}
      >
        <TextField
          fullWidth
          multiline
          minRows={1}
          maxRows={6}
          placeholder="Ask the LLM..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          disabled={isStreaming}
        />
        <Tooltip title={isStreaming ? "Stop" : "Send"}>
          <IconButton
            color={isStreaming ? "secondary" : "primary"}
            onClick={() => {
              if (isStreaming) {
                controllerRef.current?.abort();
                setIsStreaming(false);
              } else {
                handleSend();
              }
            }}
          >
            {isStreaming ? <StopIcon /> : <SendIcon />}
          </IconButton>
        </Tooltip>
      </Box>
    </Box>
  );
}
