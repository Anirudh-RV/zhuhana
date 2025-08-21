import { Card, CardContent, Typography, Box, Avatar } from "@mui/material";
import { red, green } from "@mui/material/colors";

interface WatchlistCardProps {
    symbol: string;
    name: string;
    price: number;
    change: number; // percentage
    iconUrl?: string;
}

export function WatchlistCard({ symbol, name, price, change, iconUrl }: WatchlistCardProps) {
    return (
        <Card
            sx={(theme) => ({
                borderRadius: 3,
                boxShadow: theme.palette.mode === "dark"
                    ? "0 0 15px rgba(0, 150, 255, 0.15)"
                    : "0 4px 12px rgba(0,0,0,0.1)",
                background: theme.palette.background.paper,
                transition: "transform 0.2s ease, box-shadow 0.2s ease",
                "&:hover": {
                    transform: "scale(1.03)",
                    boxShadow: theme.palette.mode === "dark"
                        ? "0 0 25px rgba(0, 150, 255, 0.3)"
                        : "0 6px 20px rgba(0,0,0,0.15)",
                },
                cursor: "pointer",
                minWidth: 200,
            })}
        >
            <CardContent>
                <Box display="flex" alignItems="center" gap={1}>
                    <Avatar
                        src={iconUrl}
                        sx={{
                            bgcolor: "#f7931a",
                            width: 32,
                            height: 32,
                            fontSize: "0.9rem",
                            fontWeight: 600,
                        }}
                    >
                        {symbol[0]}
                    </Avatar>
                    <Box>
                        <Typography variant="subtitle1" fontWeight={600}>
                            {symbol}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            {name}
                        </Typography>
                    </Box>
                </Box>

                <Typography variant="h6" mt={2}>
                    {price.toLocaleString()} <Typography component="span" variant="caption">USD</Typography>
                </Typography>

                <Typography
                    variant="body2"
                    mt={0.5}
                    sx={{ color: change >= 0 ? green[500] : red[400], fontWeight: 500 }}
                >
                    {change >= 0 ? "+" : ""}
                    {change.toFixed(2)}%
                </Typography>
            </CardContent>
        </Card>
    );
}