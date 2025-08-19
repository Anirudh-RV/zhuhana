import {
  Box,
  Typography,
  TextField,
  IconButton,
  Card,
  CardContent,
  CardActions,
  Button,
  CircularProgress,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import { useEffect, useState } from "react";
import { useAuth } from "../../AuthContext";
import LLMPanel from "../../code/components/LLMPanel";
import {
  GET_NEWS_ARTICLE_V1_ENDPOINT,
  ASK_LLM_V1_ENDPOINT,
  CREATE_CHAT_SESSION_V1_ENDPOINT,
} from "../../constants";
import { useSearchParams } from "react-router-dom";
import { useRef } from "react";
import ClearIcon from "@mui/icons-material/Clear";
import InputAdornment from "@mui/material/InputAdornment";

export type Article = {
  article_id: string;
  title: string;
  link: string;
  keywords?: string[];
  creator?: string[];
  description?: string;
  content?: string;
  pubDate: string;
  pubDateTZ: string;
  image_url?: string;
  video_url?: string | null;
  source_id: string;
  source_name: string;
  source_priority: number;
  source_url: string;
  source_icon: string;
  language: string;
  country?: string[];
  category?: string[];
  sentiment?: string;
  sentiment_stats?: string;
  ai_tag?: string;
  ai_region?: string;
  ai_org?: string;
  ai_summary?: string;
  ai_content?: string;
  duplicate: boolean;
};

type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

export default function Discovery() {
  const { user, accessToken } = useAuth();
  const [inputQuery, setInputQuery] = useState("");
  const [query, setQuery] = useState("");

  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(false);
  const [isAppending, setIsAppending] = useState(false);

  const userID: string = user?.ID ?? "id";

  const [sessionId, setSessionId] = useState<string | null>(null);
  const [isNewSession, setIsNewSession] = useState(false);

  const handleSearch = async ({
    append = false,
    customPage,
  }: { append?: boolean; customPage?: string | null } = {}) => {
    const queryToUse = query.trim() !== "" ? query : "Latest Financial News";

    if (!queryToUse.trim() && !customPage) return;

    if (append) setIsAppending(true);
    else setLoading(true);

    const pageParam = customPage ?? null;

    try {
      const endpoint = new URL(GET_NEWS_ARTICLE_V1_ENDPOINT);
      endpoint.searchParams.set("query", queryToUse);
      if (pageParam) {
        endpoint.searchParams.set("page", pageParam);
      }

      const response = await fetch(endpoint.toString(), {
        headers: {
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
        },
      });

      const json = await response.json();

      if (json.status === 1) {
        const newArticles = json.data.results;
        const nextPage = json.data.nextPage || null;

        setNextPageToken(nextPage);

        setArticles((prev: Article[]): Article[] => {
          if (append) {
            const existingIds = new Set(prev.map((a: Article) => a.article_id));
            const filteredNew = newArticles.filter(
              (a: Article) => !existingIds.has(a.article_id)
            );
            return [...prev, ...filteredNew];
          }
          return newArticles;
        });
      }
    } catch (error) {
      console.error("Error fetching articles:", error);
    } finally {
      if (append) setIsAppending(false);
      else setLoading(false);
    }
  };

  useEffect(() => {
    setQuery("Latest Financial News");
    setInputQuery("");
    handleSearch({ append: false });
  }, []);

  const [searchParams, setSearchParams] = useSearchParams();
  const [llmMessages, setLlmMessages] = useState<Message[]>([]);
  const newsFeedRef = useRef<HTMLDivElement>(null);
  const scrollTimeout = useRef<NodeJS.Timeout | null>(null);
  const lastFetchedPage = useRef<number>(1);
  const [nextPageToken, setNextPageToken] = useState<string | null>(null);

  useEffect(() => {
    const handleScroll = () => {
      if (scrollTimeout.current || loading || !nextPageToken) return;

      const div = newsFeedRef.current;
      if (!div) return;

      const nearBottom =
        div.scrollTop + div.clientHeight >= div.scrollHeight - 200;
      if (nearBottom) {
        handleSearch({ append: true, customPage: nextPageToken });
        scrollTimeout.current = setTimeout(() => {
          scrollTimeout.current = null;
        }, 500);
      }
    };

    const currentRef = newsFeedRef.current;
    if (currentRef) currentRef.addEventListener("scroll", handleScroll);

    return () => {
      if (currentRef) currentRef.removeEventListener("scroll", handleScroll);
      if (scrollTimeout.current) {
        clearTimeout(scrollTimeout.current);
        scrollTimeout.current = null;
      }
    };
  }, [loading, nextPageToken]);

  const handleSendToLLM = async (
    messages: Message[],
    onChunk: (token: string) => void,
    signal: AbortSignal
  ) => {
    try {
      const prompt = messages
        .map(
          (msg) =>
            `${msg.role === "user" ? "User" : "Assistant"}: ${msg.content}`
        )
        .join("\n");

      const latestUserMessage = [...messages]
        .reverse()
        .find((msg) => msg.role === "user");

      if (!latestUserMessage) {
        onChunk("[Error: No user message found]");
        return;
      }

      // Local variable to ensure correct sessionId usage
      let currentSessionId = sessionId;

      if (!currentSessionId) {
        const createSessionResponse = await fetch(
          CREATE_CHAT_SESSION_V1_ENDPOINT,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
            body: JSON.stringify({
              algorithm_id: user?.ID,
              title: latestUserMessage.content.slice(0, 60),
            }),
          }
        );

        if (!createSessionResponse.ok) {
          const errorText = await createSessionResponse.text();
          throw new Error(`Failed to create session: ${errorText}`);
        }

        const sessionData = await createSessionResponse.json();
        currentSessionId = sessionData.Result.id;

        setSessionId(currentSessionId);
        setIsNewSession(true);

        searchParams.set("session_id", currentSessionId!);
        setSearchParams(searchParams);
        window.location.hash = "discovery";
      }

      const queryParams = new URLSearchParams({
        q: prompt,
        current_user_q: latestUserMessage.content,
        ...(currentSessionId ? { session_id: currentSessionId } : {}),
      });

      const response = await fetch(`${ASK_LLM_V1_ENDPOINT}?${queryParams}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
        },
        signal,
      });

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();

      if (!reader) {
        onChunk("[Error: No stream]");
        return;
      }

      while (true) {
        if (signal.aborted) {
          try {
            await reader.cancel();
          } catch {}
          throw new DOMException("Aborted", "AbortError");
        }

        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        const lines = chunk
          .split("\n")
          .filter((line: string) => line.trim() !== "");

        for (const line of lines) {
          try {
            const data = JSON.parse(line);
            if (data.done) return;
            onChunk(data.response || "");
          } catch {
            onChunk(chunk);
          }
        }

        if (isNewSession) setIsNewSession(false);
      }
    } catch (err: any) {
      if (err.name !== "AbortError") {
        onChunk(`[LLM Error]: ${err.message}`);
      }
    }
  };

  useEffect(() => {
    // Run the search once on load with default topic
    const runInitialSearch = async () => {
      setLoading(true);
      try {
        const response = await fetch(
          `${GET_NEWS_ARTICLE_V1_ENDPOINT}?query=${encodeURIComponent(
            "Latest financial news"
          )}`,
          {
            headers: {
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
          }
        );
        const json = await response.json();
        if (json.status === 1) {
          setArticles(json.data.results);
          setNextPageToken(json.data.nextPage || null);
        }
      } catch (error) {
        console.error("Error fetching articles:", error);
      } finally {
        setLoading(false);
      }
    };

    runInitialSearch();
  }, []);

  const triggerSearch = () => {
    setArticles([]);
    setNextPageToken(null);
    setQuery(inputQuery); // this updates the real query used in search
    handleSearch({ append: false });
  };

  return (
    <Box
      sx={{
        px: { xs: 2, sm: 3, md: 4 }, // padding on smaller screens too
        pt: 2,
        maxWidth: { sm: "100%", md: "1700px" },
        mx: "auto",
        width: "100%",
      }}
    >
      {/* Search bar */}
      <Box sx={{ display: "flex", gap: 1, mb: 2 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Search for news, tickers or companies..."
          value={inputQuery}
          onChange={(e) => setInputQuery(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              triggerSearch();
            }
          }}
          slotProps={{
            input: {
              endAdornment: inputQuery ? (
                <InputAdornment position="end">
                  <IconButton
                    size="small"
                    onClick={() => setInputQuery("")} // <-- only clears input
                  >
                    <ClearIcon fontSize="small" />
                  </IconButton>
                </InputAdornment>
              ) : null,
            },
          }}
        />

        <IconButton
          onClick={() => {
            setArticles([]);
            setNextPageToken(null);
            handleSearch({ append: false });
          }}
        >
          <SearchIcon />
        </IconButton>
      </Box>

      {/* Grid layout */}
      <Box
        sx={{
          display: "flex",
          height: "calc(100vh - 130px)", // adjust if your header height is different
        }}
      >
        {/* News Feed (Left Panel) */}
        <Box
          ref={newsFeedRef}
          sx={{
            flex: 1,
            pr: 1,
            overflowY: "auto",
            display: "flex",
            flexDirection: "column",
          }}
        >
          {loading ? (
            <Box sx={{ width: "100vh", mx: "auto" }}>
              <CircularProgress />
            </Box>
          ) : articles.length === 0 ? (
            <Box sx={{ width: "100%", mx: "auto" }}>
              <Typography>No articles found.</Typography>
            </Box>
          ) : (
            <>
              <Box sx={{ maxWidth: "100%", mx: "auto", width: "100%" }}>
                <Box
                  sx={{
                    display: "grid",
                    gridTemplateColumns:
                      "repeat(auto-fill, minmax(300px, 1fr))",
                    gap: 2,
                  }}
                >
                  {articles.map((article) => (
                    <Card
                      key={article.article_id}
                      sx={{
                        display: "flex",
                        flexDirection: "column",
                        height: "100%",
                      }}
                    >
                      {/* Preview Image with fallback */}
                      <Box
                        sx={{
                          height: 180,
                          overflow: "hidden",
                          position: "relative",
                        }}
                      >
                        <img
                          src={article.image_url || "/images/news.png"}
                          onError={(e) => {
                            e.currentTarget.onerror = null; // prevent infinite loop
                            e.currentTarget.src = "/images/news.png"; // ✅ no `url(...)`
                          }}
                          alt="preview"
                          width="100%"
                          height="100%"
                          style={{
                            objectFit: "cover",
                            width: "100%",
                            height: "100%",
                            display: "block",
                          }}
                        />
                      </Box>

                      <CardContent sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" gutterBottom>
                          {article.title}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {article.description
                            ? article.description.length > 150
                              ? `${article.description.slice(0, 150)}...`
                              : article.description
                            : ""}
                        </Typography>

                        {article.source_name && (
                          <Typography variant="caption" color="text.secondary">
                            Source: {article.source_name}
                          </Typography>
                        )}
                        <br />
                        <Typography variant="caption" color="text.secondary">
                          Published at:{" "}
                          {new Date(article.pubDate).toLocaleString()}
                        </Typography>
                      </CardContent>
                      <CardActions>
                        <Button
                          size="small"
                          href={article.link}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          Read more
                        </Button>
                      </CardActions>
                    </Card>
                  ))}
                </Box>
              </Box>
              {isAppending && (
                <Box textAlign="center" my={2}>
                  <CircularProgress size={24} />
                </Box>
              )}
            </>
          )}
        </Box>

        {/* Drag handle & toggle button */}
        {/* LLM Panel (Right Panel) */}
        <Box
          sx={{
            width: "30%",
            height: "100%",
            overflowY: "auto",
          }}
        >
          <LLMPanel
            onSend={handleSendToLLM}
            messages={llmMessages}
            setMessages={setLlmMessages}
            algorithmId={userID}
            sessionId={sessionId}
            setSessionId={setSessionId}
            isNewSession={isNewSession}
            setIsNewSession={setIsNewSession}
            inputBoxPlaceHolder="Ask about the News..."
          />
        </Box>
      </Box>
    </Box>
  );
}
