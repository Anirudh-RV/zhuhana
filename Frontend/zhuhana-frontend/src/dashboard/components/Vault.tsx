import {
  Box,
  Typography,
  Button,
  CircularProgress,
  Stack,
  TextField,
  IconButton,
} from "@mui/material";
import { useEffect, useState } from "react";
import { useAuth } from "../../AuthContext";
import { GET_SECRET_KEYS_V1_ENDPOINT } from "../../constants";
import AddIcon from "@mui/icons-material/Add";
import EditableSecretsList from "./EditableSecretsList";

export default function Vault() {
  const { accessToken } = useAuth();
  const [keys, setKeys] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchKeys = async () => {
    try {
      setLoading(true);
      const res = await fetch(GET_SECRET_KEYS_V1_ENDPOINT, {
        headers: { ...(accessToken ? { USER_TOKEN: accessToken } : {}) },
      });
      const data = await res.json();
      if (data.status === 1) setKeys(data.keys);
    } catch (err) {
      console.error("Failed to fetch secret keys:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchKeys();
  }, [accessToken]);

  return (
    <Box sx={{ width: "100%", maxWidth: "75%", mx: "auto", p: 2 }}>
      <Stack
        direction="row"
        justifyContent="space-between"
        alignItems="center"
        mb={2}
      >
        <Typography variant="h5">Secrets Vault</Typography>
      </Stack>

      {loading ? (
        <CircularProgress />
      ) : (
        <EditableSecretsList
          keys={keys}
          accessToken={accessToken}
          onUpdated={fetchKeys}
        />
      )}
    </Box>
  );
}
