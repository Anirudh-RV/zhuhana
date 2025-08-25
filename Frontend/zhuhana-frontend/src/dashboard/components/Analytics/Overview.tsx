import Grid from "@mui/material/Grid";
import StatCard, {type StatCardProps} from "../StatCard";
import TrendingUpIcon from "@mui/icons-material/TrendingUp";
import TrendingDownIcon from "@mui/icons-material/TrendingDown";
import BarChartIcon from "@mui/icons-material/BarChart";
import PercentIcon from "@mui/icons-material/Percent";
import AttachMoneyIcon from "@mui/icons-material/AttachMoney";
import ShieldIcon from "@mui/icons-material/Shield";
const overviewData: StatCardProps[] = [
    {
        id: "total-return",
        title: "Total Return",
        value: "23.5%",
        interval: "Since inception",
        trend: "up",
        icon: <TrendingUpIcon color="success" />,
        data: [10, 20, 25, 30, 28, 35, 40],
    },
    {
        id: "annualized-return",
        title: "Annualized Return",
        value: "12.1%",
        interval: "CAGR",
        trend: "up",
        icon: <TrendingUpIcon color="success" />,
        data: [5, 10, 12, 13, 14, 16],
    },
    {
        id: "max-drawdown",
        title: "Max Drawdown",
        value: "-15.2%",
        interval: "Worst peak-to-trough",
        trend: "down",
        icon: <TrendingDownIcon color="error" />,
        data: [0, -5, -8, -12, -15, -10, -12],
    },
    {
        id: "volatility",
        title: "Volatility",
        value: "12.5%",
        interval: "Std dev of returns",
        trend: "neutral",
        icon: <BarChartIcon color="primary" />,
        data: [10, 12, 14, 11, 13, 12, 12.5],
    },
    {
        id: "win-rate",
        title: "Win Rate",
        value: "58.0%",
        interval: "Profitable trades",
        trend: "up",
        icon: <PercentIcon color="success" />,
        data: [50, 52, 55, 57, 59, 58],
    },
    {
        id: "avg-pnl",
        title: "Avg P&L / Trade",
        value: "$120.35",
        interval: "Per executed trade",
        trend: "up",
        icon: <AttachMoneyIcon color="success" />,
        data: [50, 60, 80, 100, 110, 120],
    },
    {
        id: "sharpe",
        title: "Sharpe Ratio",
        value: "1.45",
        interval: "Return / Volatility",
        trend: "up",
        icon: <TrendingUpIcon color="primary" />,
        data: [0.8, 1.0, 1.2, 1.3, 1.4, 1.45],
    },
    {
        id: "sortino",
        title: "Sortino Ratio",
        value: "1.88",
        interval: "Return / Downside risk",
        trend: "up",
        icon: <ShieldIcon color="success" />,
        data: [1.0, 1.2, 1.4, 1.6, 1.7, 1.88],
    },
];

export function OverviewSection() {
    return (
        //
        <Grid container spacing={2} columns={12} sx={{ mb: (theme) => theme.spacing(2) }} >
            {/*<Grid size={{ xs: 12, sm: 6, lg: 3 }}></Grid>*/}
            {overviewData.slice(0, 1).map((card, index) => (
                <Grid key={index} size={{ xs: 12, sm: 12, lg: 12 }}> <StatCard {...card} /> </Grid>
            ))}

            {overviewData.slice(1, 3).map((card, index) => (
                <Grid key={index + 1} size={{ xs: 12, sm: 6, lg: 6 }}>
                    <StatCard {...card} />
                </Grid>
            ))}

            {overviewData.slice(3, 6).map((card, index) => (
                <Grid key={index + 3} size={{ xs: 12, sm: 4, lg: 4 }}>
                    <StatCard {...card} />
                </Grid>
            ))}

            {overviewData.slice(6,8).map((card, index) => (
            <Grid key={index + 6} size={{ xs: 12, sm: 6, lg: 6 }}>
                <StatCard {...card} />
            </Grid>
            ))}
        </Grid>
    )
}