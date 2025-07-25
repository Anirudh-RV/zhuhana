import {
  Box,
  Stack,
  TextField,
  Typography,
  Button,
  IconButton,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import { useState } from "react";
import { useEffect, useRef } from "react";
import DeleteForeverIcon from "@mui/icons-material/DeleteForever";

interface Props {
  keys: string[];
  accessToken: string | null;
  onUpdated: () => void;
}

interface SecretRow {
  id?: string;
  key: string;
  value: string;
  fetched: boolean;
  editing: boolean;
  isNew?: boolean;
  hovering?: boolean;
  hideGetValueAt?: number;
}

export default function EditableSecretsList({
  keys,
  accessToken,
  onUpdated,
}: Props) {
  const [rows, setRows] = useState<SecretRow[]>(
    (keys ?? []).map((key) => ({
      key,
      value: "********",
      fetched: false,
      editing: false,
    }))
  );

  const keyInputRef = useRef<HTMLInputElement | null>(null);

  const handleAddNewRow = () => {
    setRows((prev) => [
      ...prev,
      {
        key: "",
        value: "",
        fetched: false,
        editing: true, // ✅ immediately open value input
        isNew: true,
      },
    ]);

    setTimeout(() => {
      keyInputRef.current?.focus();
    }, 0);
  };

  const handleDelete = async (row: SecretRow) => {
    if (!row.id) {
      // try fetching the ID first
      const res = await fetch(
        `http://localhost:8004/v1/user/secret/?key=${row.key}`,
        {
          headers: { ...(accessToken ? { USER_TOKEN: accessToken } : {}) },
        }
      );
      const data = await res.json();
      if (data.status !== 1) return;

      row.id = data.userSecret.ID;
    }

    await fetch("http://localhost:8004/v1/user/secret/", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        ...(accessToken ? { USER_TOKEN: accessToken } : {}),
      },
      body: JSON.stringify({ id: row.id }),
    });

    setRows((prev) => prev.filter((r) => r.key !== row.key));
  };

  const handleGetValue = async (key: string) => {
    const res = await fetch(
      `http://localhost:8004/v1/user/secret/?key=${key}`,
      {
        headers: { ...(accessToken ? { USER_TOKEN: accessToken } : {}) },
      }
    );
    const data = await res.json();
    if (data.status === 1) {
      const expiryTime = Date.now() + 30_000;
      const { Value, ID } = data.userSecret;

      setRows((prev) =>
        prev.map((r) =>
          r.key === key
            ? {
                ...r,
                value: Value || "",
                id: ID, // ← Save ID
                fetched: true,
                hideGetValueAt: expiryTime,
              }
            : r
        )
      );

      setTimeout(() => {
        setRows((prev) =>
          prev.map((r) =>
            r.key === key
              ? {
                  ...r,
                  hideGetValueAt: undefined,
                  fetched: false,
                  value: "********",
                }
              : r
          )
        );
      }, 30_000); // ← fix timeout to 30s
    }
  };

  const handleSave = async (row: SecretRow) => {
    const res = await fetch("http://localhost:8004/v1/user/secret/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(accessToken ? { USER_TOKEN: accessToken } : {}),
      },
      body: JSON.stringify({
        key: row.key,
        value: row.value,
      }),
    });
    const data = await res.json();
    if (data.status === 1) {
      onUpdated();
    }
  };

  const updateRow = (index: number, changes: Partial<SecretRow>) => {
    setRows((prev) =>
      prev.map((r, i) => (i === index ? { ...r, ...changes } : r))
    );
  };

  return (
    <Stack spacing={2} sx={{ width: "100%" }}>
      {rows.map((row, index) => (
        <Stack
          key={index}
          direction="row"
          spacing={2}
          alignItems="center"
          sx={{
            borderBottom: "1px solid #ccc",
            pb: 1,
            width: "100%",
            flexWrap: "nowrap",
          }}
        >
          {row.isNew ? (
            <TextField
              label="Key"
              value={row.key}
              inputRef={row.isNew ? keyInputRef : undefined}
              onChange={(e) => updateRow(index, { key: e.target.value })}
              size="small"
              sx={{ width: 200 }}
            />
          ) : (
            <Typography variant="body1" sx={{ width: 200 }}>
              {row.key}
            </Typography>
          )}

          {/* Editable Value Field */}
          {row.editing ? (
            <TextField
              size="small"
              value={row.value}
              autoFocus={!row.isNew}
              onChange={(e) => updateRow(index, { value: e.target.value })}
              onBlur={() => {
                if (!row.isNew) {
                  updateRow(index, { editing: false });
                }
              }}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === "Escape") {
                  updateRow(index, { editing: false });
                }
              }}
              sx={{
                flexGrow: 1,
                "& input": {
                  fontSize: "0.95rem",
                },
              }}
            />
          ) : (
            <Box
              onMouseEnter={() => updateRow(index, { hovering: true })}
              onMouseLeave={() => updateRow(index, { hovering: false })}
              onClick={() =>
                updateRow(index, {
                  editing: true,
                  value:
                    !row.fetched && row.value === "********" ? "" : row.value,
                })
              }
              sx={{
                flexGrow: 1,
                cursor: "text",
                borderRadius: "4px",
                padding: "4px 8px",
                border: "1px solid transparent",
                "&:hover": {
                  border: "1px solid #1976d2",
                },
                transition: "border 0.2s",
              }}
            >
              <Typography
                variant="body2"
                sx={{
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                }}
              >
                {row.fetched
                  ? row.value
                  : row.value && row.value !== ""
                  ? row.value
                  : row.isNew
                  ? ""
                  : "********"}
              </Typography>
            </Box>
          )}

          {/* GET VALUE Button */}
          <Box sx={{ width: 100, display: "flex", justifyContent: "center" }}>
            {!row.fetched && !row.isNew && !row.hideGetValueAt ? (
              <Button size="small" onClick={() => handleGetValue(row.key)}>
                Get Value
              </Button>
            ) : (
              <Box sx={{ width: 0, height: 0 }} />
            )}
          </Box>

          {/* SAVE Button */}
          <Box sx={{ width: 80, display: "flex", justifyContent: "center" }}>
            <Button
              size="small"
              onClick={() => {
                handleSave(row);
                updateRow(index, {
                  editing: false,
                  isNew: false,
                  fetched: true,
                });
              }}
              disabled={!row.key || !row.value}
              variant="contained"
              sx={{
                "&.Mui-disabled": {
                  backgroundColor: "#ccc",
                  color: "#666",
                },
              }}
            >
              Save
            </Button>
          </Box>
          <Box sx={{ width: 60, display: "flex", justifyContent: "center" }}>
            {!row.isNew && (
              <IconButton
                disableRipple
                onClick={() => handleDelete(row)}
                sx={{
                  color: "error.main",
                  backgroundColor: "transparent !important",
                  border: 0,
                }}
              >
                <DeleteForeverIcon />
              </IconButton>
            )}
          </Box>
        </Stack>
      ))}

      {/* Add new key-value button */}
      <Box>
        <IconButton onClick={handleAddNewRow}>
          <AddIcon />
        </IconButton>
      </Box>
    </Stack>
  );
}
