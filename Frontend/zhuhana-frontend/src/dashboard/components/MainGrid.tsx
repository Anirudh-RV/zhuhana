import Grid from "@mui/material/Grid";
import Stack from "@mui/material/Stack";
import {
  Box,
  Collapse,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Paper,
} from "@mui/material";

import HighlightedCard from "./HighlightedCard";
import StatCard from "./StatCard";
import type { StatCardProps } from "./StatCard";
import { useAuth } from "../../AuthContext";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  USER_PYTHON_ALGORITHMS_INFORMATION_V1_ENDPOINT,
  USER_PYTHON_ALGORITHM_RUNS_V1_ENDPOINT,
} from "../../constants";
import { KeyboardArrowDown, KeyboardArrowUp } from "@mui/icons-material";
import { useTheme, useColorScheme } from "@mui/material/styles";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import Tooltip from "@mui/material/Tooltip";

const data: StatCardProps[] = [];

interface Algorithm {
  id: string;
  scriptName: string;
  scriptUrl: string;
  order_domain: string;
  created_at: string;
  updated_at: string;
}

interface Run {
  ID: string;
  OrderDomain: string;
  Status: string;
  Market: string;
  Symbol: string;
  StartTime: string;
  EndTime: string;
  PortfolioSize: number;
}

export default function MainGrid() {
  const { user, accessToken } = useAuth();
  const [rows, setRows] = useState<Algorithm[]>([]);
  const navigate = useNavigate();

  const [algorithms, setAlgorithms] = useState<Algorithm[]>([]);
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [runsMap, setRunsMap] = useState<Record<string, Run[]>>({});

  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;

  useEffect(() => {
    const fetchAlgorithms = async () => {
      if (!user || !accessToken) return;

      const res = await fetch(USER_PYTHON_ALGORITHMS_INFORMATION_V1_ENDPOINT, {
        method: "GET",
        headers: { USER_TOKEN: accessToken },
      });
      const json = await res.json();
      if (json?.status === 1) {
        setAlgorithms(json.user_algorithms);
      }
    };
    fetchAlgorithms();
  }, [user, accessToken]);

  useEffect(() => {
    const fetchAlgorithms = async () => {
      if (!user || !accessToken) {
        console.error("User not authenticated");
        return;
      }

      try {
        const response = await fetch(
          USER_PYTHON_ALGORITHMS_INFORMATION_V1_ENDPOINT,
          {
            method: "GET",
            headers: {
              USER_TOKEN: accessToken,
            },
          }
        );

        if (!response.ok) throw new Error("Failed to fetch algorithms");

        const result = await response.json();

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
    { field: "scriptName", headerName: "Algorithm Name", flex: 1 },

    { field: "order_domain", headerName: "Domain", flex: 1 },
    { field: "created_at", headerName: "Created At", flex: 1 },
    { field: "updated_at", headerName: "Updated At", flex: 1 },
  ];

  const handleExpand = async (algoId: string) => {
    const newExpandedIds = new Set(expandedIds);

    if (expandedIds.has(algoId)) {
      newExpandedIds.delete(algoId);
    } else {
      newExpandedIds.add(algoId);
      if (!runsMap[algoId]) {
        try {
          const res = await fetch(
            `${USER_PYTHON_ALGORITHM_RUNS_V1_ENDPOINT}?algorithm_id=${algoId}`,
            {
              headers: {
                ...(accessToken ? { USER_TOKEN: accessToken } : {}),
              },
            }
          );
          const json = await res.json();
          if (json?.status === 1) {
            setRunsMap((prev) => ({
              ...prev,
              [algoId]: json.user_algorithm_runs,
            }));
          }
        } catch (e) {
          console.error("❌ Failed to fetch runs:", e);
        }
      }
    }

    setExpandedIds(newExpandedIds);
  };

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
          <Box sx={{ width: "100%", maxWidth: "100%" }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Algorithm Details
            </Typography>

            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Sl. No</TableCell>
                    <TableCell>Algorithm Name</TableCell>
                    <TableCell>Domain</TableCell>
                    <TableCell>Created At</TableCell>
                    <TableCell>Updated At</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {algorithms.map((algo, index) => (
                    <>
                      <TableRow
                        key={algo.id}
                        sx={{
                          "&:hover": {
                            cursor: "pointer",
                          },
                        }}
                        onClick={() => handleExpand(algo.id)}
                      >
                        <TableCell>{index + 1}</TableCell>
                        <TableCell>{algo.scriptName}</TableCell>
                        <TableCell>{algo.order_domain}</TableCell>
                        <TableCell>
                          {new Date(algo.created_at).toLocaleString()}
                        </TableCell>
                        <TableCell>
                          {new Date(algo.updated_at).toLocaleString()}
                        </TableCell>
                        <TableCell align="right">
                          <Tooltip title="Open Algorithm">
                            <IconButton
                              onClick={(e) => {
                                e.stopPropagation();
                                navigate(`/code?algorithm_id=${algo.id}`);
                              }}
                              sx={{ border: 0 }}
                            >
                              <OpenInNewIcon />
                            </IconButton>
                          </Tooltip>
                        </TableCell>
                      </TableRow>

                      <TableRow>
                        <TableCell
                          colSpan={6}
                          sx={{
                            p: 0,
                            backgroundColor:
                              resolvedMode === "dark"
                                ? "#000000"
                                : "background.default",
                          }}
                        >
                          <Collapse
                            in={expandedIds.has(algo.id)}
                            timeout="auto"
                            unmountOnExit
                          >
                            <Box sx={{ margin: 2 }}>
                              {runsMap[algo.id]?.length ? (
                                <Table size="small">
                                  <TableHead>
                                    <TableRow>
                                      <TableCell>Order Domain</TableCell>
                                      <TableCell>Status</TableCell>
                                      <TableCell>Market</TableCell>
                                      <TableCell>Symbol</TableCell>
                                      <TableCell>Start</TableCell>
                                      <TableCell>End</TableCell>
                                      <TableCell>Portfolio Size</TableCell>
                                    </TableRow>
                                  </TableHead>
                                  <TableBody>
                                    {runsMap[algo.id].map((run) => (
                                      <TableRow key={run.ID}>
                                        <TableCell>{run.OrderDomain}</TableCell>
                                        <TableCell>{run.Status}</TableCell>
                                        <TableCell>{run.Market}</TableCell>
                                        <TableCell>{run.Symbol}</TableCell>
                                        <TableCell>
                                          {new Date(
                                            run.StartTime
                                          ).toLocaleDateString()}
                                        </TableCell>
                                        <TableCell>
                                          {new Date(
                                            run.EndTime
                                          ).toLocaleDateString()}
                                        </TableCell>
                                        <TableCell>
                                          {run.PortfolioSize}
                                        </TableCell>
                                        <TableCell align="right">
                                          <Tooltip title="Open Run Analytics">
                                            <IconButton
                                              onClick={() => {
                                                const url = new URL(
                                                  window.location.href
                                                );
                                                url.hash = "analytics";
                                                url.searchParams.set(
                                                  "algorithm_run_id",
                                                  run.ID
                                                );
                                                window.location.href =
                                                  url.toString();
                                              }}
                                              sx={{
                                                border: 0,
                                              }}
                                            >
                                              <OpenInNewIcon />
                                            </IconButton>
                                          </Tooltip>
                                        </TableCell>
                                      </TableRow>
                                    ))}
                                  </TableBody>
                                </Table>
                              ) : (
                                <Typography
                                  variant="body2"
                                  sx={{ color: "gray" }}
                                >
                                  No runs found.
                                </Typography>
                              )}
                            </Box>
                          </Collapse>
                        </TableCell>
                      </TableRow>
                    </>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
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
