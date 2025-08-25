import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import StatCard from "./StatCard";
import PortfolioSection from "./Analytics/Portfolio";
import type { StatCardProps } from "./StatCard";
import { useAuth } from "../../AuthContext";
import { useNavigate, useSearchParams } from "react-router-dom";
import { WatchlistSection } from "./Analytics/Watchlist";
import { OverviewSection } from "./Analytics/Overview";
import { BacktestHistorySection } from "./Analytics/BacktestHistory";


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
            <OverviewSection></OverviewSection>


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
                <BacktestHistorySection></BacktestHistorySection>
            </Box>
        </Box>
    );
}