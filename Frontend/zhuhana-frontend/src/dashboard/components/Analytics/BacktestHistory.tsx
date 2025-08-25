import {useColorScheme, useTheme} from "@mui/material/styles";
import { Box, Typography } from "@mui/material";
import { DataGrid, GridColDef } from "@mui/x-data-grid";
import { LineChart, Line, ResponsiveContainer } from "recharts";
import {BacktestRunsList} from "./BacktestRunList";


export function BacktestHistorySection() {
    const theme = useTheme();

    const bestStrategy = {
        name: "Momentum",
        cagr: 15.3,
        sharpe: 1.60,
        mdd: -20.5,
        sortino: 2.00,
        equity: [100, 105, 115, 120, 130],
    };

    const equityData = bestStrategy.equity.map((val, idx) => ({ idx, val }));
    const { mode, systemMode } = useColorScheme();
    const resolvedMode = mode === "system" ? systemMode : mode;
    const isDark = resolvedMode === "dark";

    return (
        <Box>
        <Box
            sx={{
                p: 3,
                mb: 3,
                borderRadius: 3,
                bgcolor: isDark ? "grey.900" : "grey.100",
                border: "1px solid",
                borderColor: isDark ? "success.light" : "success.dark",
            }}
        >
            <Typography
                variant="h6"
                sx={{ fontWeight: 700, color: isDark ? "success.light" : "success.dark" }}
            >
                🏆 Highlighted Strategy: {bestStrategy.name}
            </Typography>
            <Typography
                variant="body2"
                sx={{ mb: 1, color: isDark ? "grey.300" : "grey.700" }}
            >
                CAGR: {bestStrategy.cagr}% | Sharpe: {bestStrategy.sharpe} | MDD:{" "}
                {bestStrategy.mdd}% | Sortino: {bestStrategy.sortino}
            </Typography>
            <Box sx={{ height: 80 }}>
                <ResponsiveContainer>
                    <LineChart data={equityData}>
                        <Line
                            type="monotone"
                            dataKey="val"
                            stroke={isDark ? "#66ff99" : "#2e7d32"} // 荧光绿 / 深绿
                            strokeWidth={2}
                            dot={false}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </Box>
        </Box>
        <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>
            Backtest Runs
        </Typography>
        <BacktestRunsList></BacktestRunsList>
        </Box>
    );
}
