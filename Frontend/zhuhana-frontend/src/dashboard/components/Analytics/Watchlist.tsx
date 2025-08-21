import { Box, Typography } from "@mui/material";
import { WatchlistCard } from "./WatchlistCard";

const mockWatchlist = [
    {
        symbol: "BTCUSD",
        name: "Bitcoin",
        price: 113724.1,
        change: -0.47,
        iconUrl: "https://cryptologos.cc/logos/bitcoin-btc-logo.png",
    },
    {
        symbol: "ETHUSD",
        name: "Ethereum",
        price: 3900.55,
        change: +2.15,
        iconUrl: "https://cryptologos.cc/logos/ethereum-eth-logo.png",
    },
    {
        symbol: "AAPL",
        name: "Apple",
        price: 172.33,
        change: +1.25,
        iconUrl: "https://logo.clearbit.com/apple.com",
    },
    {
        symbol: "MSFT",
        name: "Microsoft",
        price: 310.12,
        change: -0.85,
        iconUrl: "https://logo.clearbit.com/microsoft.com",
    },
];

export function WatchlistSection() {
    return (
        <Box>
            <Box
                sx={{
                    display: "flex",
                    gap: 2,
                    overflowX: "auto",
                    pb: 1,
                    "&::-webkit-scrollbar": {
                        height: 6,
                    },
                    "&::-webkit-scrollbar-thumb": {
                        backgroundColor: "rgba(120,120,120,0.3)",
                        borderRadius: 4,
                    },
                }}
            >
                {mockWatchlist.map((item) => (
                    <WatchlistCard key={item.symbol} {...item} />
                ))}
            </Box>
        </Box>
    );
}