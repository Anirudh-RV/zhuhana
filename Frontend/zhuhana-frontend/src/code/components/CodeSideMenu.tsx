import * as React from "react";
import { Box, Button, IconButton, Typography, Tooltip } from "@mui/material";
import BacktestConfig from "./BacktestConfig";
import PaperTradeConfig from "./PaperTradingConfig";
import LiveTradeConfig from "./LiveTradingConfig";
import MenuIcon from "@mui/icons-material/Menu";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";
import { useColorScheme } from "@mui/material/styles";

const executionModes = ["Backtest", "Paper Trade", "Live Trade"];

export default function CodeSideMenu({ onClose }: { onClose?: () => void }) {
  const [mode, setMode] = React.useState("Backtest");
  const { mode: themeMode, systemMode: themeSystemMode } = useColorScheme();
  const resolvedThemeMode =
    themeMode === "system" ? themeSystemMode : themeMode;
  const selectedColor = resolvedThemeMode === "dark" ? "grey.700" : "grey.300";

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
      <Box sx={{ overflowY: "auto", maxHeight: "100%", mt: 1 }}>
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
