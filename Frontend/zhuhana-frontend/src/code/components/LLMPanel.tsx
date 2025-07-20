import { useEffect, useRef, useState } from "react";
import {
  Box,
  Typography,
  IconButton,
  Tooltip,
  InputBase,
  Snackbar,
} from "@mui/material";

import SendIcon from "@mui/icons-material/Send";
import StopIcon from "@mui/icons-material/Stop";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import { useColorScheme } from "@mui/material/styles";
import CheckIcon from "@mui/icons-material/Check";
import ReactMarkdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import hljs from "highlight.js/lib/core";
import python from "highlight.js/lib/languages/python";
import ChatBubbleOutlineIcon from "@mui/icons-material/ChatBubbleOutline";
import HistoryIcon from "@mui/icons-material/History";

hljs.registerLanguage("python", python);

type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

type LLMPanelProps = {
  onSend: (
    messages: Message[],
    onChunk: (token: string) => void,
    signal: AbortSignal
  ) => Promise<void>;
  onClose?: () => void;
  messages: Message[];
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>;
};

export default function LLMPanel({
  onSend,
  onClose,
  messages,
  setMessages,
}: LLMPanelProps) {
  const [input, setInput] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const [showTypingIndicator, setShowTypingIndicator] = useState(false);
  const [copied, setCopied] = useState(false);
  const controllerRef = useRef<AbortController | null>(null);
  const bottomRef = useRef<HTMLDivElement | null>(null);

  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;
  const panelBgColor =
    resolvedMode === "dark" ? "#161B26" : "background.default";

  const handleSend = () => {
    if (!input.trim() || isStreaming) return;

    const userMessage: Message = { role: "user", content: input };
    const botMessage: Message = { role: "assistant", content: "" };
    const updatedMessages = [...messages, userMessage, botMessage];
    setMessages(updatedMessages);

    const controller = new AbortController();
    controllerRef.current = controller;
    setIsStreaming(true);
    setShowTypingIndicator(true);

    let currentBotResponse = "";

    onSend(
      updatedMessages.slice(0, -1),
      (token: string) => {
        setShowTypingIndicator(false);
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
        setMessages((prev) => [...prev, { role: "system", content: "---" }]);
      });

    setInput("");
  };

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, showTypingIndicator]);

  const CodeBlock = ({ inline, className, children }: any) => {
    const [copied, setCopied] = useState(false);
    const codeRef = useRef<HTMLElement>(null);

    const { mode, systemMode } = useColorScheme();
    const resolvedMode = mode === "system" ? systemMode : mode;

    const handleCopy = () => {
      const text = codeRef.current?.innerText || "";
      navigator.clipboard.writeText(text.trim()).then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 1500);
      });
    };

    if (inline) {
      return (
        <code
          style={{
            backgroundColor: "#e0e0e0",
            padding: "2px 4px",
            borderRadius: "4px",
            fontFamily: "monospace",
          }}
        >
          {children}
        </code>
      );
    }

    return (
      <Box
        sx={{
          position: "relative",
          borderRadius: 1,
          overflow: "hidden",
          border: "1px solid",
          borderColor: "divider",
          mb: 2,
        }}
      >
        <Box
          onClick={handleCopy}
          sx={{
            position: "absolute",
            top: 8,
            right: 8,
            display: "flex",
            alignItems: "center",
            gap: 0.5,
            px: 1,
            py: 0.25,
            borderRadius: 1,
            cursor: "pointer",
            backgroundColor:
              resolvedMode === "dark" ? "background.paper" : "#ddd",
            color: resolvedMode === "dark" ? "#fff" : "#000",
            transition: "background-color 0.2s",
            "&:hover": {
              backgroundColor: resolvedMode === "dark" ? "#444" : "#ccc",
            },
          }}
        >
          {copied ? (
            <>
              <CheckIcon fontSize="small" />
              <Typography variant="body2" fontSize="0.75rem">
                Copied
              </Typography>
            </>
          ) : (
            <>
              <ContentCopyIcon fontSize="small" />
              <Typography variant="body2" fontSize="0.75rem">
                Copy
              </Typography>
            </>
          )}
        </Box>

        <pre style={{ margin: 0, padding: "1rem" }}>
          <code className={className} ref={codeRef}>
            {children}
          </code>
        </pre>
      </Box>
    );
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        height: "100%",
        backgroundColor: "background.default",
        color: "text.primary",
      }}
    >
      {/* Header */}
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          px: 2,
          py: 1,
          borderBottom: "1px solid",
          borderColor: "divider",
          backgroundColor: "background.paper",
        }}
      >
        <Tooltip title="New Chat">
          <IconButton
            aria-label="new chat"
            disableRipple
            sx={{
              border: "none",
              backgroundColor: "transparent",
              p: 1, // padding to keep the click area large enough
              "&:hover": {
                backgroundColor: "action.hover", // optional: subtle hover effect
              },
            }}
          >
            <ChatBubbleOutlineIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Chat History">
          <IconButton
            aria-label="new chat"
            disableRipple
            sx={{
              border: "none",
              backgroundColor: "transparent",
              p: 1, // padding to keep the click area large enough
              "&:hover": {
                backgroundColor: "action.hover", // optional: subtle hover effect
              },
            }}
          >
            <HistoryIcon />
          </IconButton>
        </Tooltip>
        <Typography variant="h6" sx={{ ml: 1 }}>
          Zhuhana AI
        </Typography>
      </Box>

      {/* Chat Area */}
      <Box
        sx={{
          flexGrow: 1,
          overflowY: "auto",
          p: 2,
          display: "flex",
          flexDirection: "column",
          backgroundColor: panelBgColor,
          gap: 2,
        }}
      >
        {messages.map((msg, idx) => (
          <Box
            key={idx}
            sx={{
              display: "flex",
              justifyContent: msg.role === "user" ? "flex-end" : "flex-start",
              width: "100%",
            }}
          >
            <Box
              sx={{
                width: msg.role === "user" ? "75%" : "100%", // user gets bubble, assistant spans full
                backgroundColor:
                  msg.role === "user" ? "action.selected" : "transparent",
                px: 2,
                py: 1,
                borderRadius: 2,
                fontSize: "0.95rem",
                whiteSpace: "pre-wrap",
                overflowWrap: "anywhere",
                textAlign: "left", // ensures left justification
              }}
            >
              <ReactMarkdown
                rehypePlugins={[rehypeHighlight]}
                components={{ code: CodeBlock }}
              >
                {msg.content}
              </ReactMarkdown>
            </Box>
          </Box>
        ))}

        {isStreaming && showTypingIndicator && (
          <Box
            sx={{
              display: "flex",
              gap: 1,
              alignItems: "center",
              px: 2,
              mt: 1,
            }}
          >
            {[0, 1, 2].map((i) => (
              <Box
                key={i}
                sx={{
                  width: 8,
                  height: 8,
                  borderRadius: "50%",
                  backgroundColor: "text.secondary",
                  animation: "typing 1.2s infinite ease-in-out both",
                  animationDelay: `${i * 0.2}s`,
                  "@keyframes typing": {
                    "0%, 80%, 100%": { opacity: 0 },
                    "40%": { opacity: 1 },
                  },
                }}
              />
            ))}
          </Box>
        )}

        <div ref={bottomRef} />
      </Box>

      {/* Input Section */}
      <Box
        component="form"
        onSubmit={(e) => {
          e.preventDefault();
          handleSend();
        }}
        sx={{
          px: 1,
          py: 1,
          borderTop: "1px solid",
          borderColor: "divider",
          backgroundColor: "background.paper",
          display: "flex",
          alignItems: "flex-end",
          gap: 1,
        }}
      >
        <InputBase
          fullWidth
          multiline
          minRows={1}
          maxRows={6}
          placeholder="Ask the LLM..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault(); // prevent newline
              handleSend(); // trigger send
            }
          }}
          disabled={isStreaming}
          sx={{
            px: 2,
            py: 1.5,
            borderRadius: 2,
            backgroundColor: "background.default",
            color: "text.primary",
            fontSize: "0.95rem",
            fontFamily: "monospace",
            lineHeight: 1.6,
          }}
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

      {/* Copy Snackbar */}
      <Snackbar
        open={copied}
        message="Copied"
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
        autoHideDuration={1000}
      />
    </Box>
  );
}
