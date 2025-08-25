import { useEffect, useRef } from "react";
import {
    Box,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    Paper,
    Grid,
} from "@mui/material";
import { PieChart, Pie, Cell, Tooltip, Legend } from "recharts";
import {BaselineSeries, createChart, ISeriesApi} from "lightweight-charts";
import {useTheme,useColorScheme} from "@mui/material/styles";
import {PortfolioPieChart} from "./PortfolioPieChart";
import {PortfolioTable} from "./PortfolioTable";

export default function PortfolioSection() {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const theme = useTheme();
    const { mode, systemMode } = useColorScheme();
    const resolvedMode = mode === "system" ? systemMode : mode;

    useEffect(() => {
        if (!chartContainerRef.current) return;

        const chart = createChart(chartContainerRef.current, {
            width: chartContainerRef.current.clientWidth,
            height: 250,
            layout: {
                textColor: "#ccc",
                background: { color: "transparent" },
            },
            grid: {
                vertLines: { color: "#444" },
                horzLines: { color: "#444" },
            },
        });

        const baselineSeries: ISeriesApi<"Baseline"> = chart.addSeries(BaselineSeries, {
            baseValue: { type: "price", price: 10000 },
            topFillColor1: "rgba(38,198,218,0.28)",
            topFillColor2: "rgba(38,198,218,0.05)",
            topLineColor: "rgba(38,198,218,1)",
            bottomFillColor1: "rgba(239,83,80,0.05)",
            bottomFillColor2: "rgba(239,83,80,0.28)",
            bottomLineColor: "rgba(239,83,80,1)",
        });

        const data = [
            { time: "2024-01-01", value: 10000 },
            { time: "2024-02-01", value: 10800 },
            { time: "2024-03-01", value: 11200 },
            { time: "2024-04-01", value: 10700 },
            { time: "2024-05-01", value: 12000 },
        ];
        baselineSeries.setData(data);
        chart.timeScale().fitContent();

        const resizeObserver = new ResizeObserver(() => {
            chart.applyOptions({ width: chartContainerRef.current!.clientWidth });
        });
        resizeObserver.observe(chartContainerRef.current);

        return () => {
            resizeObserver.disconnect();
            chart.remove();
        };
    }, []);

    return (
        <Box sx={{ mt: 0 }}>
            <Grid container spacing={2} alignItems="flex-start">
                <Grid size={{xs: 12, md: 5}}>
                    <Typography variant="h6" gutterBottom color="text.secondary" sx={{
                        fontSize: "0.9rem",
                        fontWeight: 600,
                        textTransform: "uppercase",
                        letterSpacing: "0.025em",
                        mb: 1,
                    }}>
                        Portfolio Allocation
                    </Typography>
                    <PortfolioPieChart />
                </Grid>

                <Grid size={{xs: 12, md: 7}}>
                    <Typography gutterBottom variant="subtitle1"
                                sx={(theme) => ({
                                    fontSize: "0.9rem",
                                    fontWeight: 600,
                                    textTransform: "uppercase",
                                    letterSpacing: "0.025em",
                                    mb: 1,
                                    color: resolvedMode === "dark"
                                        ? "rgba(255,255,255,0.7)"
                                        : "rgba(0,0,0,1)"
                                })}
                    >
                        Holdings
                    </Typography>
                    <PortfolioTable />
                </Grid>
            </Grid>
            <Grid size={{xs:12, md:6}}>
                <Typography variant="subtitle1" gutterBottom color="text.secondary"
                            sx={{
                                fontSize: "0.9rem",
                                fontWeight: 600,
                                textTransform: "uppercase",
                                letterSpacing: "0.025em",
                                mb: 1,
                    }}>
                    Portfolio Equity Curve
                </Typography>
                <Box ref={chartContainerRef} />
            </Grid>
        </Box>
    );
}