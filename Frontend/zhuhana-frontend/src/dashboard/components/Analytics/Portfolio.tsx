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
import {useTheme} from "@mui/material/styles";
import {PortfolioPieChart} from "./PortfolioPieChart";

const portfolioData = [
    { symbol: "AAPL", position: 100, avgCost: 150, currentPrice: 170, pnl: 2000, weight: 0.35 },
    { symbol: "MSFT", position: 50, avgCost: 250, currentPrice: 310, pnl: 3000, weight: 0.40 },
    { symbol: "TSLA", position: 30, avgCost: 700, currentPrice: 680, pnl: -600, weight: 0.25 },
];


export default function PortfolioSection() {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const theme = useTheme();

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
                    <Typography gutterBottom variant="subtitle1" color="text.secondary"
                                sx={{
                                    fontSize: "0.9rem",
                                    fontWeight: 600,
                                    textTransform: "uppercase",
                                    letterSpacing: "0.05em",
                                }}
                    >
                        Holdings
                    </Typography>
                    <Paper sx={{ mb: 3 }} >
                        <Table
                            size="medium"
                            sx={{
                                "& th": {
                                    backgroundColor: "rgba(255,255,255,0.05)",
                                    color: "#9ca3af",
                                    fontWeight: 600,
                                    textTransform: "uppercase",
                                    fontSize: "0.75rem",
                                    letterSpacing: "0.05em",
                                },
                                "& td": {
                                    borderBottom: "1px solid rgba(255,255,255,0.05)",
                                    fontFamily: "monospace",
                                    fontSize: "0.9rem",
                                    color: "#e5e7eb",
                                },
                                "& tr:hover": {
                                    background: "rgba(59,130,246,0.08)",
                                    boxShadow: "0 0 12px rgba(59,130,246,0.4)",
                                    transition: "0.2s ease",
                                },
                            }}
                        >
                            <TableHead>
                                <TableRow>
                                    <TableCell>Symbol</TableCell>
                                    <TableCell>Position</TableCell>
                                    <TableCell>Avg Cost</TableCell>
                                    <TableCell>Current Price</TableCell>
                                    <TableCell>P&L</TableCell>
                                    <TableCell>Weight</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {portfolioData.map((row, idx) => (
                                    <TableRow key={idx}>
                                        <TableCell>{row.symbol}</TableCell>
                                        <TableCell>{row.position}</TableCell>
                                        <TableCell>${row.avgCost}</TableCell>
                                        <TableCell>${row.currentPrice}</TableCell>
                                        <TableCell
                                            sx={{ color: row.pnl >= 0 ? "success.main" : "error.main" }}
                                        >
                                            {row.pnl}
                                        </TableCell>
                                        <TableCell>{(row.weight * 100).toFixed(1)}%</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </Paper>
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