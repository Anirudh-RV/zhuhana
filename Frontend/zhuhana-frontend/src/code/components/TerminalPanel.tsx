// components/TerminalPanel.tsx
import { Box, IconButton, Toolbar, Typography } from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import { useTheme, useColorScheme } from "@mui/material/styles";

export type TerminalLine = { text: string; type: "info" | "success" | "error" };

export default function TerminalPanel({
  terminalOutput,
  isLoadingPyodide,
  onRunCode,
  onCopyTerminal,
}: {
  terminalOutput: TerminalLine[];
  isLoadingPyodide: boolean;
  onRunCode: () => void;
  onCopyTerminal: () => void;
}) {
  const theme = useTheme();
  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;
  console.log("resolvedMode: " + resolvedMode);

  return (
    <Box
      sx={{
        flexGrow: 1,
        display: "flex",
        flexDirection: "column",
        border: `1px solid ${theme.palette.divider}`,
        borderRadius: 1,
        backgroundColor: resolvedMode === "dark" ? "#000000" : "#f5f5f5",
        minHeight: "60px",
        overflow: "hidden",
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
          <IconButton
            onClick={onCopyTerminal}
            size="small"
            sx={{ color: resolvedMode === "dark" ? "#ccc" : "#333" }}
          >
            <ContentCopyIcon fontSize="small" />
          </IconButton>

          <IconButton
            onClick={onRunCode}
            size="small"
            disabled={isLoadingPyodide}
            sx={{
              color: isLoadingPyodide
                ? theme.palette.grey[500]
                : theme.palette.success.main,
              "&:hover": {
                backgroundColor: isLoadingPyodide
                  ? "transparent"
                  : theme.palette.success.dark,
              },
            }}
          >
            <PlayArrowIcon />
          </IconButton>
        </Box>
      </Toolbar>

      <Box
        sx={{
          flexGrow: 1,
          p: 1,
          overflowY: "auto",
          overflowX: "auto",
          whiteSpace: "pre-wrap",
          wordBreak: "break-word",
          maxWidth: "100%",
        }}
      >
        {terminalOutput.map((line, index) => (
          <Typography
            key={index}
            sx={{
              color:
                line.type === "success"
                  ? "#0f0"
                  : line.type === "error"
                  ? "#f55"
                  : resolvedMode === "dark"
                  ? "#aaa"
                  : "#222",
              fontFamily: "monospace",
              fontSize: "0.875rem",
            }}
          >
            {line.text}
          </Typography>
        ))}
      </Box>
    </Box>
  );
}
