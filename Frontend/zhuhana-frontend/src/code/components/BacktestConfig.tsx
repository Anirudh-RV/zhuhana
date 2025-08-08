import {
  Typography,
  TextField,
  Stack,
  Box,
  Button,
  Autocomplete,
  InputAdornment,
  Tooltip,
} from "@mui/material";
import { LocalizationProvider, DatePicker } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import dayjs from "dayjs";
import { useColorScheme } from "@mui/material/styles";

export interface BacktestValues {
  instrument: string;
  timeDuration: string;
  frequencyType: string; // "1D" | "1W" | "1M" | "Custom"
  customFrequencyDays: number | "";
  startDate: dayjs.Dayjs | null;
  endDate: dayjs.Dayjs | null;
  portfolioSize: number;
}

export default function BacktestConfig({
  values,
  onChange,
}: {
  values: BacktestValues;
  onChange: (changes: Partial<BacktestValues>) => void;
}) {
  const timeDurations = [
    { value: "1D", label: "1 Day" },
    { value: "1M", label: "1 Month" },
    { value: "1Y", label: "1 Year" },
    { value: "Custom", label: "Custom Range" },
  ];

  const frequencies = [
    { value: "1D", label: "1 Day" },
    { value: "1W", label: "1 Week" },
    { value: "1M", label: "1 Month" },
    { value: "Custom", label: "Custom Frequency" },
  ];

  const instruments = ["SPY", "AAPL", "GOOGL", "MSFT", "AMZN"];

  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;
  const selectedColor = resolvedMode === "dark" ? "grey.700" : "grey.300";

  return (
    <Box sx={{ mt: 2 }}>
      {/* Instrument Selection */}
      <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
        Select Instrument
      </Typography>
      <Autocomplete
        fullWidth
        size="small"
        options={instruments}
        value={values.instrument}
        onChange={(_, newValue) => {
          if (newValue) onChange({ instrument: newValue });
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            size="small"
            sx={{
              "& .MuiInputBase-root": {
                height: 48,
                backgroundColor: "background.default",
              },
              "& .MuiInputBase-input": {
                padding: "12px 14px",
              },
              "& .MuiInputLabel-root": {
                top: -5,
              },
            }}
          />
        )}
        sx={{
          mb: 2,
          "& .MuiAutocomplete-endAdornment": {
            backgroundColor: "background.default",
            borderRadius: 0,
            height: "100%",
            display: "flex",
            alignItems: "center",
          },
        }}
      />

      {/* Time Duration */}
      <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
        Select Time Duration
      </Typography>
      <Box
        sx={{
          display: "flex",
          mb: 2,
          border: "1px solid",
          borderColor: "divider",
          borderRadius: "8px",
          overflow: "hidden",
          width: "100%",
        }}
      >
        {timeDurations.map(({ value, label }, index) => {
          const isSelected = values.timeDuration === value;
          return (
            <Tooltip key={value} title={label} arrow>
              <Button
                onClick={() => onChange({ timeDuration: value })}
                variant="text"
                sx={{
                  flex: 1,
                  textTransform: "none",
                  borderRadius: 0,
                  py: 3,
                  fontWeight: isSelected ? 600 : 400,
                  backgroundColor: isSelected
                    ? selectedColor
                    : "background.paper",
                  color: "text.primary",
                  borderRight:
                    index < timeDurations.length - 1 ? "1px solid" : "none",
                  borderColor: "divider",
                  "&:hover": {
                    backgroundColor: isSelected ? "grey.400" : "action.hover",
                  },
                }}
              >
                {value}
              </Button>
            </Tooltip>
          );
        })}
      </Box>

      {values.timeDuration === "Custom" && (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
          <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
            <DatePicker
              label="Start Date"
              value={values.startDate}
              onChange={(newValue) =>
                newValue && onChange({ startDate: newValue })
              }
              slotProps={{
                textField: {
                  fullWidth: true,
                  size: "small",
                },
              }}
            />
            <DatePicker
              label="End Date"
              value={values.endDate}
              onChange={(newValue) =>
                newValue && onChange({ endDate: newValue })
              }
              slotProps={{
                textField: {
                  fullWidth: true,
                  size: "small",
                },
              }}
            />
          </Stack>
        </LocalizationProvider>
      )}

      {/* Frequency */}
      {/* Frequency */}
      <Typography variant="h6" sx={{ mt: 3, mb: 1 }}>
        Select Frequency
      </Typography>
      <Box
        sx={{
          display: "flex",
          mb: 2,
          border: "1px solid",
          borderColor: "divider",
          borderRadius: "8px",
          overflow: "hidden",
          width: "100%",
        }}
      >
        {frequencies.map(({ value, label }, index) => {
          const isSelected = values.frequencyType === value;
          return (
            <Tooltip key={value} title={label} arrow>
              <Button
                onClick={() => onChange({ frequencyType: value })}
                variant="text"
                sx={{
                  flex: 1,
                  textTransform: "none",
                  borderRadius: 0,
                  py: 3,
                  fontWeight: isSelected ? 600 : 400,
                  backgroundColor: isSelected
                    ? selectedColor
                    : "background.paper",
                  color: "text.primary",
                  borderRight:
                    index < frequencies.length - 1 ? "1px solid" : "none",
                  borderColor: "divider",
                  "&:hover": {
                    backgroundColor: isSelected ? "grey.400" : "action.hover",
                  },
                }}
              >
                {value}
              </Button>
            </Tooltip>
          );
        })}
      </Box>

      {values.frequencyType === "Custom" && (
        <TextField
          fullWidth
          label="Frequency (Days)"
          type="number"
          size="small"
          value={values.customFrequencyDays}
          onChange={(e) =>
            onChange({ customFrequencyDays: Number(e.target.value) })
          }
          sx={{
            mb: 2,
            "& .MuiInputBase-root": {
              height: 48,
            },
            "& .MuiInputLabel-root": {
              color: "text.secondary",
            },
            "& .MuiInputLabel-root.Mui-focused": {
              color: "text.secondary",
            },
          }}
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
          value={values.portfolioSize}
          onChange={(e) => onChange({ portfolioSize: Number(e.target.value) })}
          InputProps={{
            startAdornment: <InputAdornment position="start">$</InputAdornment>,
            inputProps: {
              min: 0,
              step: 100,
            },
          }}
          sx={{
            "& .MuiInputBase-root": {
              height: 48,
            },
            "& .MuiInputLabel-root": {
              color: "text.secondary",
            },
            "& .MuiInputLabel-root.Mui-focused": {
              color: "text.secondary",
            },
          }}
        />
      </Stack>
    </Box>
  );
}
