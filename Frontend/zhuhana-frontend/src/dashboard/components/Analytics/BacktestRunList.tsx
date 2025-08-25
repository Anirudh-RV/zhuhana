import { Box, Typography, Link } from "@mui/material";
import { DataGrid, GridColDef } from "@mui/x-data-grid";

const backtestRuns = [
    {
        id: "run_001",
        strategy: "SMA Crossover",
        start: "2024-01-01",
        end: "2024-06-30",
        capital: 100000,
        return: 23.5,
        cagr: 12.1,
        mdd: -15.2,
        sharpe: 1.45,
    },
    {
        id: "run_002",
        strategy: "Momentum",
        start: "2024-03-01",
        end: "2024-07-31",
        capital: 50000,
        return: 30.0,
        cagr: 15.3,
        mdd: -20.5,
        sharpe: 1.60,
    },
    {
        id: "run_003",
        strategy: "Mean Reversion",
        start: "2024-02-01",
        end: "2024-05-30",
        capital: 75000,
        return: 18.0,
        cagr: 10.2,
        mdd: -12.0,
        sharpe: 1.20,
    },
];

const runColumns: GridColDef[] = [
    {
        field: "id",
        headerName: "Run ID",
        flex: 1,
        renderCell: (params) => (
            <Link
                href={`/backtest/${params.value}`}
                underline="hover"
                sx={{
                    fontWeight: 600,
                    color: "primary.main",
                    "&:hover": {
                        color: "primary.dark",
                    },
                }}
            >
                {params.value}
            </Link>
        ),
    },
    { field: "strategy", headerName: "Strategy", flex: 1 },
    { field: "start", headerName: "Start Date", flex: 1 },
    { field: "end", headerName: "End Date", flex: 1 },
    {
        field: "capital",
        headerName: "Initial Capital",
        flex: 1,
        valueFormatter: (p) =>
            p.value != null
                ? new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "USD",
                    maximumFractionDigits: 0,
                }).format(p.value as number)
                : "-",
    },
    { field: "return", headerName: "Total Return %", flex: 1 },
    { field: "cagr", headerName: "CAGR %", flex: 1 },
    { field: "mdd", headerName: "Max DD %", flex: 1 },
    { field: "sharpe", headerName: "Sharpe", flex: 1 },
];

export function BacktestRunsList() {
    return (
        <Box sx={{ mt: 4 }}>
            <div style={{ height: 400, width: "100%" }}>
                <DataGrid
                    rows={backtestRuns}
                    columns={runColumns}
                    pageSize={5}
                    disableRowSelectionOnClick
                    sx={{
                        "& .MuiDataGrid-columnHeaders": {
                            fontWeight: "bold",
                        },
                    }}
                />
            </div>
        </Box>
    );
}
