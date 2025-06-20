import { useState } from "react";
import { Box, Typography, TextField, IconButton, Paper } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";

type LLMPanelProps = {
  output: string;
  onSend: (input: string, onChunk: (token: string) => void) => void;
};

export default function LLMPanel({ output, onSend }: LLMPanelProps) {
  const [input, setInput] = useState("");

  const handleSend = () => {
    if (!input.trim()) return;
    onSend(input, () => {});
    setInput("");
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Typography variant="h6" gutterBottom>
        LLM Output
      </Typography>
      <Paper
        elevation={1}
        sx={{
          borderRadius: 2,
          p: 2,
          flexGrow: 1,
          overflowY: "auto",
          mb: 2,
          whiteSpace: "pre-wrap",
        }}
      >
        {output}
      </Paper>

      <Box
        component="form"
        onSubmit={(e) => {
          e.preventDefault();
          handleSend();
        }}
        sx={{
          display: "flex",
          gap: 1,
        }}
      >
        <TextField
          fullWidth
          placeholder="Ask the LLM..."
          size="small"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <IconButton type="submit" color="primary">
          <SendIcon />
        </IconButton>
      </Box>
    </Box>
  );
}
