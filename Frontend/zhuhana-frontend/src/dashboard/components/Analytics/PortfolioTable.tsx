import {Paper, Table, TableBody, TableCell, TableHead, TableRow} from "@mui/material";
import {useColorScheme, useTheme} from "@mui/material/styles";

const portfolioData = [
    { symbol: "AAPL", position: 100, avgCost: 150, currentPrice: 170, pnl: 2000, weight: 0.35 },
    { symbol: "MSFT", position: 50, avgCost: 250, currentPrice: 310, pnl: 3000, weight: 0.40 },
    { symbol: "TSLA", position: 30, avgCost: 700, currentPrice: 680, pnl: -600, weight: 0.25 },
];


export function PortfolioTable() {

    const theme = useTheme();
    const { mode, systemMode } = useColorScheme();
    const resolvedMode = mode === "system" ? systemMode : mode;

    return (
    <Paper sx={{ mb: 3 }} >
        <Table
            size="medium"
            sx={{
                "& th": {
                    backgroundColor: "rgba(255,255,255,0.05)",
                    // color: "#9ca3af",
                    fontWeight: 600,
                    textTransform: "uppercase",
                    fontSize: "0.75rem",
                    letterSpacing: "0.05em",
                },
                "& td": {
                    borderBottom: "1px solid rgba(255,255,255,0.05)",
                    fontFamily: "monospace",
                    fontSize: "0.9rem",
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
    )
}
