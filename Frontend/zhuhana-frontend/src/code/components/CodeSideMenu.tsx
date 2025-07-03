import * as React from "react";
import {
  Box,
  Button,
  IconButton,
  Stack,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from "@mui/material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import BacktestConfig from "./BacktestConfig";
import PaperTradeConfig from "./PaperTradingConfig";
import LiveTradeConfig from "./LiveTradingConfig";

const executionModes = ["Backtest", "Paper Trade", "Live Trade"];

export default function CodeSideMenu() {
  const [mode, setMode] = React.useState("Backtest");

  return (
    <Box
      sx={{
        flex: 1,
        height: "100%",
        p: 2,
        backgroundColor: "background.paper",
        display: "flex",
        flexDirection: "column",
        borderRight: "1px solid",
        borderColor: "divider",
      }}
    >
      {/* Bottom Content */}
      <Box sx={{ overflowY: "auto", maxHeight: "100%", mt: 1 }}>
        <Typography variant="h6" gutterBottom>
          Select Execution Mode
        </Typography>
        <ToggleButtonGroup
          value={mode}
          exclusive
          onChange={(_e, newMode) => newMode && setMode(newMode)}
          fullWidth
          sx={{ mb: 2 }}
        >
          {executionModes.map((m) => (
            <ToggleButton
              key={m}
              value={m}
              disabled={m === "Live Trade" || m === "Paper Trade"}
            >
              {m}
            </ToggleButton>
          ))}
        </ToggleButtonGroup>

        {/* Step 3: Conditional Components */}
        {mode === "Backtest" && <BacktestConfig />}
        {mode === "Paper Trade" && <PaperTradeConfig />}
        {mode === "Live Trade" && <LiveTradeConfig />}

        <Button variant="contained" fullWidth sx={{ mt: 2 }}>
          {mode}
        </Button>
      </Box>
    </Box>
  );
}
