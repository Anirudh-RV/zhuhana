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
import BacktestConfig from "./BacktestConfig";
import PaperTradeConfig from "./PaperTradingConfig";
import LiveTradeConfig from "./LiveTradingConfig";
import MenuIcon from "@mui/icons-material/Menu";
import Tooltip from "@mui/material/Tooltip";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import { useColorScheme } from "@mui/material/styles";
import ArrowBackIosIcon from "@mui/icons-material/ArrowBackIos";

const executionModes = ["Backtest", "Paper Trade", "Live Trade"];

export default function CodeSideMenu({ onClose }: { onClose?: () => void }) {
  const [mode, setMode] = React.useState("Backtest");
  const { mode: themeMode, systemMode: themeSystemMode } = useColorScheme();
  const resolvedThemeMode =
    themeMode === "system" ? themeSystemMode : themeMode;

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
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            mb: 2,
          }}
        >
          <Typography variant="h5">Configuration</Typography>
          {onClose && (
            <IconButton size="small" onClick={onClose}>
              <MenuIcon fontSize="small" />
            </IconButton>
          )}
        </Box>

        <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
          <Typography variant="h6" sx={{ mr: 0.5 }}>
            Select Execution Mode
          </Typography>
          <Tooltip
            title="⚠️ Live Trading and Paper Trading coming soon!"
            placement="bottom"
            slotProps={{
              tooltip: {
                sx: {
                  backgroundColor:
                    resolvedThemeMode === "dark" ? "#333" : "#eee",
                  color: resolvedThemeMode === "dark" ? "#fff" : "#000",
                  boxShadow: 3,
                  opacity: 1, // Ensures it's fully opaque
                },
              },
            }}
          >
            <InfoOutlinedIcon
              sx={{ fontSize: 18, color: "text.secondary", cursor: "default" }}
            />
          </Tooltip>
        </Box>

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
