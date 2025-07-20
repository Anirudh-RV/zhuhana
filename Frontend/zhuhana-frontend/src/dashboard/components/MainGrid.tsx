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

const data: StatCardProps[] = [];

interface Algorithm {
  id: string;
  scriptName: string;
  scriptUrl: string;
  order_domain: string;
  created_at: string;
  updated_at: string;
}

export default function MainGrid() {
  const { user, accessToken } = useAuth();
  const [rows, setRows] = useState<Algorithm[]>([]);
  const navigate = useNavigate();

  const handleRowDoubleClick = (params: any) => {
    const algorithmId = params.row.id;
    if (algorithmId) {
      navigate(`/code?algorithm_id=${algorithmId}`);
    }
  };

  useEffect(() => {
    const fetchAlgorithms = async () => {
      if (!user || !accessToken) {
        console.error("User not authenticated");
        return;
      }

      try {
        const response = await fetch(
          "http://localhost:8008/v1/user/algorithm/",
          {
            method: "GET",
            headers: {
              USER_TOKEN: accessToken,
            },
          }
        );

        if (!response.ok) throw new Error("Failed to fetch algorithms");

        const result = await response.json();
        console.log("✅ Fetch Success:", result);

        if (result?.status === 1) {
          setRows(result.user_algorithms);
        } else {
          console.warn("⚠️ Unexpected API response:", result.statusDescription);
        }
      } catch (err) {
        console.error("❌ Fetch error:", err);
      }
    };

    fetchAlgorithms();
  }, [user, accessToken]);

  const columns = [
    { field: "scriptName", headerName: "Script Name", flex: 1 },

    { field: "order_domain", headerName: "Domain", flex: 1 },
    { field: "created_at", headerName: "Created At", flex: 1 },
    { field: "updated_at", headerName: "Updated At", flex: 1 },
  ];

  return (
    <Box sx={{ width: "100%", maxWidth: { sm: "100%", md: "1700px" } }}>
      {/* cards */}
      <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
        Overview
      </Typography>
      <Grid
        container
        spacing={2}
        columns={12}
        sx={{ mb: (theme) => theme.spacing(2) }}
      >
        <Grid size={{ xs: 12, sm: 6, lg: 3 }}>
          <HighlightedCard />
        </Grid>
        {data.map((card, index) => (
          <Grid key={index} size={{ xs: 12, sm: 6, lg: 3 }}>
            <StatCard {...card} />
          </Grid>
        ))}
      </Grid>
      <Typography component="h2" variant="h6" sx={{ mb: 2 }}>
        Details
      </Typography>
      <Grid container spacing={2} columns={12}>
        <Grid size={{ xs: 12 }}>
          <CustomizedDataGrid
            rows={rows}
            columns={columns}
            onRowDoubleClick={handleRowDoubleClick}
          />
        </Grid>
        <Grid size={{ xs: 12, lg: 3 }}>
          <Stack
            gap={2}
            direction={{ xs: "column", sm: "row", lg: "column" }}
          ></Stack>
        </Grid>
      </Grid>
    </Box>
  );
}
