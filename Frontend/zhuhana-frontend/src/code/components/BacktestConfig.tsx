import { useState } from "react";
import {
  Typography,
  TextField,
  Stack,
  ToggleButton,
  ToggleButtonGroup,
  MenuItem,
  Box,
} from "@mui/material";

export default function BacktestConfig() {
  const [timeDuration, setTimeDuration] = useState("Custom");
  const [frequency, setFrequency] = useState("15");

  return (
    <Box sx={{ mt: 2 }}>
      <Typography variant="h5" gutterBottom>
        Backtest Configuration
      </Typography>

      {/* Instrument Selection */}
      <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
        Select Instrument
      </Typography>
      <TextField
        fullWidth
        select
        defaultValue="AAPL"
        size="small"
        sx={{ mb: 1 }}
      >
        <MenuItem value="AAPL">AAPL</MenuItem>
        <MenuItem value="GOOGL">GOOGL</MenuItem>
        <MenuItem value="TSLA">TSLA</MenuItem>
        <MenuItem value="MSFT">MSFT</MenuItem>
        <MenuItem value="AMZN">AMZN</MenuItem>
      </TextField>

      {/* Time Duration */}
      <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
        Time Duration
      </Typography>
      <ToggleButtonGroup
        value={timeDuration}
        exclusive
        onChange={(e, val) => val && setTimeDuration(val)}
        fullWidth
        sx={{ mb: 2 }}
      >
        <ToggleButton value="1D">1D</ToggleButton>
        <ToggleButton value="1W">1W</ToggleButton>
        <ToggleButton value="1Y">1Y</ToggleButton>
        <ToggleButton value="Custom">Custom</ToggleButton>
      </ToggleButtonGroup>

      {timeDuration === "Custom" && (
        <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
          <TextField
            fullWidth
            label="Start Date"
            type="date"
            size="small"
            slotProps={{
              inputLabel: {
                shrink: true,
              },
            }}
          />
          <TextField
            fullWidth
            label="End Date"
            type="date"
            size="small"
            slotProps={{
              inputLabel: {
                shrink: true,
              },
            }}
          />
        </Stack>
      )}

      {/* Frequency */}
      <Typography variant="h6" sx={{ mt: 3, mb: 1 }}>
        Select Frequency (Minute)
      </Typography>
      <ToggleButtonGroup
        value={frequency}
        exclusive
        onChange={(e, val) => val && setFrequency(val)}
        sx={{ mb: 2 }}
      >
        <ToggleButton value="15">15</ToggleButton>
        <ToggleButton value="30">30</ToggleButton>
        <ToggleButton value="60">60</ToggleButton>
        <ToggleButton value="Custom">Custom</ToggleButton>
      </ToggleButtonGroup>

      {frequency === "Custom" && (
        <TextField
          fullWidth
          label="Frequency (seconds)"
          type="number"
          size="small"
          defaultValue="60"
          sx={{ mb: 1 }}
        />
      )}

      {/* Portfolio Details */}
      <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
        Portfolio Details
      </Typography>
      <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
        <TextField
          fullWidth
          label="Portfolio Size"
          type="number"
          size="small"
          defaultValue="100000"
        />
        <TextField
          fullWidth
          label="Risk Appetite"
          type="number"
          size="small"
          defaultValue="1000"
        />
      </Stack>
    </Box>
  );
}
