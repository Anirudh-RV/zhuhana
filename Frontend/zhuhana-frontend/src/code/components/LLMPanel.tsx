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
import { useAuth } from "../../AuthContext";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import {
  CREATE_CHAT_SESSION_V1_ENDPOINT,
  GET_MESSAGES_V1_ENDPOINT,
} from "../../constants";

hljs.registerLanguage("python", python);

type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

export interface LLMPanelHandle {
  reset: () => void;
}

type LLMPanelProps = {
  onSend: (
    messages: Message[],
    onChunk: (token: string) => void,
    signal: AbortSignal
  ) => Promise<void>;
  messages: Message[];
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>;
  algorithmId: string | null;
  sessionId: string | null;
  setSessionId: React.Dispatch<React.SetStateAction<string | null>>;
  isNewSession: boolean;
  setIsNewSession: React.Dispatch<React.SetStateAction<boolean>>;
  inputBoxPlaceHolder: string;
};

export default function LLMPanel({
  onSend,
  messages,
  setMessages,
  algorithmId,
  sessionId,
  setSessionId,
  isNewSession,
  setIsNewSession,
  inputBoxPlaceHolder,
}: LLMPanelProps) {
  const [input, setInput] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const [showTypingIndicator, setShowTypingIndicator] = useState(false);
  const [copied, setCopied] = useState(false);
  const controllerRef = useRef<AbortController | null>(null);
  const bottomRef = useRef<HTMLDivElement | null>(null);
  const { user, accessToken } = useAuth();

  const [loadingMessages, setLoadingMessages] = useState(false);
  const [hasStartedChat, setHasStartedChat] = useState(false);

  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;
  const panelBgColor =
    resolvedMode === "dark" ? "background.paper" : "background.default";

  const [chatSessions, setChatSessions] = useState<any[]>([]);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleOpenChatHistory = async (
    event: React.MouseEvent<HTMLElement>
  ) => {
    setAnchorEl(event.currentTarget);

    try {
      const response = await fetch(
        `${CREATE_CHAT_SESSION_V1_ENDPOINT}?algorithm_id=${algorithmId}`,
        {
          headers: {
            "Content-Type": "application/json",
            ...(accessToken ? { USER_TOKEN: accessToken } : {}),
          },
        }
      );

      const data = await response.json();
      const sessions = data?.Result;
      if (Array.isArray(sessions)) {
        setChatSessions(sessions);
      } else {
        console.error("Invalid chat session response:", data);
      }
    } catch (err) {
      console.error("Failed to fetch chat sessions", err);
    }
  };

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
    const fetchSessionMessages = async () => {
      console.log("isNewSession: ", isNewSession);
      if (!sessionId || isNewSession) return;

      setLoadingMessages(true);
      setMessages([]);

      try {
        const response = await fetch(
          `${GET_MESSAGES_V1_ENDPOINT}?session_id=${sessionId}`,
          {
            headers: {
              "Content-Type": "application/json",
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
          }
        );

        const data = await response.json();
        const messageEntries = data?.Result;

        if (Array.isArray(messageEntries)) {
          const loadedMessages: Message[] = messageEntries.flatMap(
            (entry: any) => [
              {
                role: "user" as const,
                content: entry.user_message.replace(/^User:\s*/, ""),
              },
              {
                role: "assistant" as const,
                content: entry.system_message,
              },
            ]
          );
          setMessages(loadedMessages);
        } else {
          console.error("Invalid message response:", data);
        }
      } catch (error) {
        console.error("Failed to load session messages", error);
      } finally {
        setLoadingMessages(false);
      }
    };

    fetchSessionMessages();
  }, [sessionId, accessToken, setMessages]);

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
          justifyContent: "space-between",
          alignItems: "center",
          px: 2,
          py: 1,
          borderBottom: "1px solid",
          borderColor: "divider",
          backgroundColor: "background.paper",
        }}
      >
        <Typography variant="h6">Zhuhana AI</Typography>

        <Box sx={{ display: "flex", alignItems: "center" }}>
          <Tooltip title="New Chat">
            <IconButton
              aria-label="new chat"
              onClick={() => {
                setMessages([]);
                setSessionId(null);

                // Optional: also update the URL to remove ?session_id=
                const params = new URLSearchParams(window.location.search);
                params.delete("session_id");

                const newUrl = `${window.location.pathname}${
                  params.toString() ? `?${params.toString()}` : ""
                }${window.location.hash}`;

                window.history.replaceState({}, "", newUrl);
              }}
              disableRipple
              sx={{
                backgroundColor: "background.default",
                p: 1,
                mr: 0.5,
                "&:hover": {
                  backgroundColor: "action.hover",
                },
              }}
            >
              <ChatBubbleOutlineIcon />
            </IconButton>
          </Tooltip>

          <Tooltip title="Chat History">
            <IconButton
              aria-label="chat history"
              onClick={handleOpenChatHistory}
              disableRipple
              sx={{
                backgroundColor: "background.default",
                p: 1,
                "&:hover": {
                  backgroundColor: "action.hover",
                },
              }}
            >
              <HistoryIcon />
            </IconButton>
          </Tooltip>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={() => setAnchorEl(null)}
          >
            <MenuItem disabled>
              <Typography variant="subtitle2" fontWeight="bold">
                Chats
              </Typography>
            </MenuItem>
            {chatSessions.length > 0 ? (
              chatSessions.map((session) => (
                <MenuItem
                  key={session.id}
                  onClick={() => {
                    setAnchorEl(null);
                    setSessionId(session.id);
                    setIsNewSession(false);

                    const params = new URLSearchParams(window.location.search);
                    params.set("session_id", session.id);

                    const newUrl = `${
                      window.location.pathname
                    }?${params.toString()}${window.location.hash}`;
                    window.history.pushState({}, "", newUrl);
                  }}
                >
                  {session.title ??
                    new Date(session.created_at).toLocaleString()}
                </MenuItem>
              ))
            ) : (
              <MenuItem disabled>
                <Typography variant="body2" color="text.secondary">
                  No sessions found
                </Typography>
              </MenuItem>
            )}
          </Menu>
        </Box>
      </Box>

      {/* Chat Messages */}
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
        {messages.length === 0 ? (
          <Box
            sx={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              height: "100%",
            }}
          >
            <Typography
              variant="h4"
              sx={{ color: "text.secondary", fontWeight: 500 }}
            >
              Hi {user?.FirstName}, Let's get started!
            </Typography>
          </Box>
        ) : (
          messages
            .filter((msg) => msg.content.trim() !== "")
            .map((msg, idx) => (
              <Box
                key={idx}
                sx={{
                  display: "flex",
                  justifyContent:
                    msg.role === "user" ? "flex-end" : "flex-start",
                }}
              >
                <Box
                  sx={{
                    width: msg.role === "user" ? "75%" : "100%",
                    backgroundColor:
                      msg.role === "user" ? "action.selected" : "transparent",
                    px: 2,
                    py: 1,
                    borderRadius: 2,
                    fontSize: "0.95rem",
                    whiteSpace: "pre-wrap",
                    overflowWrap: "anywhere",
                    textAlign: "left",
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
            ))
        )}

        {isStreaming && showTypingIndicator && (
          <Box sx={{ display: "flex", gap: 1, px: 2, mt: 1 }}>
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

      {/* Input */}
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
          placeholder={inputBoxPlaceHolder}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              handleSend();
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

      <Snackbar
        open={copied}
        message="Copied"
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
        autoHideDuration={1000}
      />
    </Box>
  );
}
