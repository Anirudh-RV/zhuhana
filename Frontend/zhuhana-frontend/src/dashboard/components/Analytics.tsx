import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import CustomizedDataGrid from "./CustomizedDataGrid";
import HighlightedCard from "./HighlightedCard";
import StatCard from "./StatCard";
import type { StatCardProps } from "./StatCard";
import { useAuth } from "../../AuthContext";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { USER_PYTHON_ALGORITHMS_INFORMATION_V1_ENDPOINT } from "../../constants";

const data: StatCardProps[] = [];

export default function Analytics() {
  const { user, accessToken } = useAuth();
  const navigate = useNavigate();

  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      {/* cards */}
      <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
        Analytics
      </Typography>
      <Grid
        container
        spacing={2}
        columns={12}
        sx={{ mb: (theme) => theme.spacing(2) }}
      >
        <Grid size={{ xs: 12, sm: 6, lg: 3 }}></Grid>
        {data.map((card, index) => (
          <Grid key={index} size={{ xs: 12, sm: 6, lg: 3 }}>
            <StatCard {...card} />
          </Grid>
        ))}
      </Grid>
    </Box>
  );
}
