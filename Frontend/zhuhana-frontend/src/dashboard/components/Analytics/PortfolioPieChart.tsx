import {PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend} from "recharts";
import { Box, Typography } from "@mui/material";

const data = [
    { name: "AAPL", value: 35 },
    { name: "MSFT", value: 40 },
    { name: "TSLA", value: 25 },
];

const COLORS = [
    "rgba(38,198,218,0.9)", // teal
    "rgba(156,39,176,0.9)", // purple
    "rgba(239,83,80,0.9)",  // red
];

export function PortfolioPieChart() {
    return (
        <Box sx={{ width: "100%", height: 300 }}>
            <ResponsiveContainer>
                <PieChart>
                    <Pie
                        data={data}
                        cx="50%"
                        cy="50%"
                        innerRadius={60}
                        outerRadius={100}
                        paddingAngle={5}
                        fill="url(#pieGradient)"
                        filter="url(#glow)"
                        dataKey="value"
                        label={({ name, percent }) => (
                                `${name} ${(percent * 100).toFixed(0)}%`
                        )}
                    >
                        {data.map((entry, index) => (
                            <Cell
                                key={`cell-${index}`}
                                fill={COLORS[index % COLORS.length]}
                                stroke="rgba(255,255,255,0.2)"
                                strokeWidth={2}
                            />
                        ))}
                    </Pie>
                    <Tooltip
                        contentStyle={{
                            // backgroundColor: "rgba(20,20,30,0.9)",
                            borderRadius: "8px",
                            border: "1px solid rgba(255,255,255,0.1)",
                            color: "#fff",
                        }}
                    />
                    <Legend
                        verticalAlign="top"
                        align="right"
                        layout="vertical"
                        wrapperStyle={{
                            fontSize: "0.8rem",
                            fontWeight: 500,
                        }}
                    />
                </PieChart>
            </ResponsiveContainer>
        </Box>
    );
}