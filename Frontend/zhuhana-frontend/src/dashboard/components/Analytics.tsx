import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import StatCard from "./StatCard";
import PortfolioSection from "./Analytics/Portfolio";
import type { StatCardProps } from "./StatCard";
import { useAuth } from "../../AuthContext";
import { useNavigate, useSearchParams } from "react-router-dom";
import { WatchlistSection } from "./Analytics/Watchlist";

const overviewData: StatCardProps[] = [
    {
        id: "total-return",
        title: "Total Return",
        value: "23.5%",
        interval: "Since inception",
        trend: "up",
        data: [10, 20, 25, 30, 28, 35, 40],
    },
    {
        id: "annualized-return",
        title: "Annualized Return",
        value: "12.1%",
        interval: "CAGR",
        trend: "up",
        data: [5, 10, 12, 13, 14, 16],
    },
];

export default function Analytics() {
    const { user, accessToken } = useAuth();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const algorithmRunId = searchParams.get("algorithm_run_id");

    return (
        <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" }, p: 2 }}>
            {/* ---------- Overview ---------- */}
            <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
                Analytics Overview
            </Typography>
            <Grid container spacing={2} columns={12} sx={{ mb: (theme) => theme.spacing(2) }} >
                {/*<Grid size={{ xs: 12, sm: 6, lg: 3 }}></Grid>*/}
                {overviewData.map((card, index) => (
                    <Grid key={index} size={{ xs: 12, sm: 6, lg: 6 }}> <StatCard {...card} /> </Grid>
                ))}
            </Grid>

            {/* ---------- Section 1: Watchlist ---------- */}
            <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
                Watchlist
            </Typography>
            <Box sx={{ mb: 4, p: 2, border: "1px solid", borderColor: "divider", borderRadius: 2 }}>
                {/* TODO: Replace with DataGrid or Table */}
                <WatchlistSection></WatchlistSection>
            </Box>

            {/* ---------- Section 2: Portfolio ---------- */}
            <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
                Portfolio
            </Typography>
            <Box sx={{ mb: 4, p: 2, border: "1px solid", borderColor: "divider", borderRadius: 2 }}>
                <PortfolioSection></PortfolioSection>
            </Box>

            {/* ---------- Section 3: Backtest History ---------- */}
            <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
                Backtest History
            </Typography>
            <Box sx={{ mb: 4, p: 2, border: "1px solid", borderColor: "divider", borderRadius: 2 }}>
                {/* TODO: Replace with equity curve + drawdown + trade history */}
                <Typography variant="body2" color="text.secondary">
                    Equity curve + drawdown chart + trade history placeholder
                </Typography>
            </Box>
        </Box>
    );
}