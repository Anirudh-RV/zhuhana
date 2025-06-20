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
        width: "20%",
        height: "100vh",
        p: 2,
        backgroundColor: "background.paper",
        display: "flex",
        flexDirection: "column",
        borderRight: "1px solid",
        borderColor: "divider",
      }}
    >
      {/* Header */}
      <Stack direction="row" alignItems="center" spacing={1} mb={1}>
        <IconButton>
          <ArrowBackIcon />
        </IconButton>
        <Typography variant="subtitle1" fontWeight="bold">
          Strategy - NewAlgorithm
        </Typography>
      </Stack>

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
