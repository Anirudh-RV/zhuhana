import { useState } from "react";
import { Box, IconButton, Toolbar, Typography } from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import CheckIcon from "@mui/icons-material/Check";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import { useTheme, useColorScheme } from "@mui/material/styles";

export type TerminalLine = { text: string; type: "info" | "success" | "error" };

export default function TerminalPanel({
  terminalOutput,
  isLoadingPyodide,
  onRunCode,
  onAskAI,
}: {
  terminalOutput: TerminalLine[];
  isLoadingPyodide: boolean;
  onRunCode: () => void;
  onAskAI: (errorMessage: string) => void;
}) {
  const theme = useTheme();
  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;

  const [copied, setCopied] = useState(false);

  const lineCount = terminalOutput.length;
  const maxLines = 20; // cap height if too many lines

  const height = Math.min(lineCount * 22 + 40, 400);

  return (
    <Box
      sx={{
        flexGrow: 1,
        display: "flex",
        borderRadius: 1,
        flexDirection: "column",
        border: `1px solid ${theme.palette.divider}`,
        backgroundColor: "background.default",
        overflow: "hidden",
        minHeight: 0,
      }}
    >
      <Toolbar
        variant="dense"
        sx={{
          backgroundColor: resolvedMode === "dark" ? "#222" : "#e0e0e0",
          borderBottom: `1px solid ${theme.palette.divider}`,
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          px: 1,
        }}
      >
        <Typography
          variant="subtitle2"
          sx={{ color: resolvedMode === "dark" ? "#bbb" : "#333" }}
        >
          Terminal
        </Typography>

        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          {/* Copy Button with background */}
        </Box>
      </Toolbar>

      <Box
        sx={{
          flexGrow: 1,
          overflowY: "auto",
          overflowX: "auto",
          whiteSpace: "pre-wrap",
          wordBreak: "break-word",
          p: 1,
          fontFamily: "monospace",
        }}
      >
        {terminalOutput.map((line, index) => (
          <Box
            key={index}
            sx={{ display: "flex", alignItems: "center", gap: 1 }}
          >
            <Typography
              sx={{
                color:
                  line.type === "success"
                    ? "#0f0"
                    : line.type === "error"
                    ? "#f55"
                    : resolvedMode === "dark"
                    ? "#aaa"
                    : "#222",
                fontSize: "0.875rem",
                whiteSpace: "pre-wrap",
                flexGrow: 1,
              }}
            >
              {line.text}
            </Typography>

            {line.type === "error" && (
              <IconButton
                size="small"
                onClick={() => onAskAI(line.text)}
                title="Ask AI about this error"
              >
                <PlayArrowIcon fontSize="small" />
              </IconButton>
            )}
          </Box>
        ))}
      </Box>
    </Box>
  );
}
