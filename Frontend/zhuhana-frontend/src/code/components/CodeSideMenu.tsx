import * as React from "react";
import { Box, Button, IconButton, Typography, Tooltip } from "@mui/material";
import BacktestConfig from "./BacktestConfig";
import PaperTradeConfig from "./PaperTradingConfig";
import LiveTradeConfig from "./LiveTradingConfig";
import MenuIcon from "@mui/icons-material/Menu";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import { useColorScheme } from "@mui/material/styles";
import TerminalPanel, { TerminalLine } from "./TerminalPanel";
import { BacktestValues } from "./BacktestConfig";
import dayjs from "dayjs";
import { useAuth } from "../../AuthContext";
import { START_USER_PYTHON_ALGORITHM_V1_ENDPOINT } from "../../constants";

const executionModes = ["Backtest", "Paper Trade", "Live Trade"];

interface CodeSideMenuProps {
  onClose?: () => void;
  terminalOutput: TerminalLine[];
  isLoadingPyodide: boolean;
  onRunCode: () => void;
}

export default function CodeSideMenu({
  onClose,
  terminalOutput,
  isLoadingPyodide,
  onRunCode,
}: CodeSideMenuProps) {
  const [mode, setMode] = React.useState("Backtest");
  const { mode: themeMode, systemMode: themeSystemMode } = useColorScheme();
  const resolvedThemeMode =
    themeMode === "system" ? themeSystemMode : themeMode;
  const selectedColor = resolvedThemeMode === "dark" ? "grey.700" : "grey.300";
  const { user, accessToken } = useAuth();

  const [backtestValues, setBacktestValues] = React.useState<BacktestValues>({
    instrument: "SPY",
    timeDuration: "1Y",
    frequencyType: "1D",
    customFrequencyDays: "",
    startDate: null, // means not chosen
    endDate: null, // means not chosen
    portfolioSize: 10000,
  });

  const handleRunBacktest = async () => {
    const params = new URLSearchParams(window.location.search);
    const algorithmID = params.get("algorithm_id"); // adjust name if different
    if (!algorithmID) {
      console.error("No algorithmID in URL");
      return;
    }

    // Determine dates
    let startDate = backtestValues.startDate;
    let endDate = backtestValues.endDate;

    if (!startDate || !endDate) {
      endDate = dayjs();
      const unit = backtestValues.timeDuration.slice(-1); // Y, M, W, D
      const amount = parseInt(backtestValues.timeDuration.slice(0, -1));
      startDate = endDate.subtract(
        amount,
        unit === "Y"
          ? "year"
          : unit === "M"
          ? "month"
          : unit === "W"
          ? "week"
          : "day"
      );
    }

    let frequencySeconds: number;

    if (backtestValues.frequencyType === "1D") {
      frequencySeconds = 86400;
    } else if (backtestValues.frequencyType === "1W") {
      frequencySeconds = 604800;
    } else if (backtestValues.frequencyType === "1M") {
      frequencySeconds = 2592000;
    } else if (backtestValues.frequencyType === "Custom") {
      // Custom: days → seconds
      frequencySeconds = Number(backtestValues.customFrequencyDays) * 86400;
    } else {
      // Fallback to daily
      frequencySeconds = 86400;
    }

    const body = {
      algorithmID,
      market: "NYSEARCA",
      symbol: backtestValues.instrument,
      start_time: startDate.toISOString(),
      end_time: endDate.toISOString(),
      frequency: frequencySeconds,
      portfolio_size: backtestValues.portfolioSize,
    };

    try {
      const res = await fetch(START_USER_PYTHON_ALGORITHM_V1_ENDPOINT, {
        method: "POST",
        headers: {
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
    } catch (err) {
      console.error("Failed to start backtest", err);
    }
  };

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
        minHeight: 0,
      }}
    >
      <Box
        sx={{
          overflowY: "auto",
          flexGrow: 1,
          minHeight: 0,
          pr: 1,
          mt: 1,
        }}
      >
        {/* Header */}
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            mb: 2,
          }}
        >
          <Typography variant="h5">Configuration Panel</Typography>
          {onClose && (
            <IconButton size="small" onClick={onClose}>
              <MenuIcon fontSize="small" />
            </IconButton>
          )}
        </Box>

        {/* Mode Selector Title */}
        <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
          <Typography variant="h6" sx={{ mr: 0.5 }}>
            Select Execution Mode
          </Typography>
          <Tooltip
            title="⚠️ Paper Trading and Live Trading coming soon!"
            placement="bottom"
            slotProps={{
              tooltip: {
                sx: {
                  backgroundColor:
                    resolvedThemeMode === "dark" ? "#333" : "#eee",
                  color: resolvedThemeMode === "dark" ? "#fff" : "#000",
                  boxShadow: 3,
                  opacity: 1,
                },
              },
            }}
          >
            <InfoOutlinedIcon
              sx={{ fontSize: 18, color: "text.secondary", cursor: "default" }}
            />
          </Tooltip>
        </Box>

        {/* Custom Toggle Button Replacement */}
        <Box
          sx={{
            display: "flex",
            mb: 2,
            border: "1px solid",
            borderColor: "divider",
            borderRadius: "8px", // No rounding
            overflow: "hidden", // Remove border overlaps
            width: "100%",
          }}
        >
          {executionModes.map((m, index) => {
            const isSelected = mode === m;
            const isDisabled = m === "Live Trade" || m === "Paper Trade";

            return (
              <Button
                key={m}
                onClick={() => setMode(m)}
                disabled={isDisabled}
                variant="text"
                sx={{
                  flex: 1,
                  textTransform: "none",
                  borderRadius: 0, // Boxy edges
                  py: 4,
                  fontWeight: isSelected ? 600 : 400,
                  backgroundColor: isSelected
                    ? selectedColor
                    : "background.paper",
                  color: isSelected ? "text.primary" : "text.primary",
                  borderRight:
                    index < executionModes.length - 1 ? "1px solid" : "none",
                  borderColor: "divider",
                  "&:hover": {
                    backgroundColor: isSelected
                      ? "primary.dark"
                      : "action.hover",
                  },
                }}
              >
                {m}
              </Button>
            );
          })}
        </Box>

        {/* Conditional Configs */}
        {mode === "Backtest" && (
          <BacktestConfig
            values={backtestValues}
            onChange={(changes) =>
              setBacktestValues((prev) => ({ ...prev, ...changes }))
            }
          />
        )}

        {mode === "Paper Trade" && <PaperTradeConfig />}
        {mode === "Live Trade" && <LiveTradeConfig />}

        <Button
          variant="contained"
          fullWidth
          sx={{ mt: 2 }}
          onClick={() => {
            if (mode === "Backtest") handleRunBacktest();
          }}
        >
          {mode}
        </Button>
      </Box>
    </Box>
  );
}
