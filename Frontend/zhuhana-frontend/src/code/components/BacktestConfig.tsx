import { useState } from "react";
import { Typography, TextField, Stack, Box, Button } from "@mui/material";
import { Autocomplete } from "@mui/material";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import dayjs from "dayjs";
import { DatePicker } from "@mui/x-date-pickers";
import InputAdornment from "@mui/material/InputAdornment";
import { useColorScheme } from "@mui/material/styles";
import { Tooltip } from "@mui/material";

export default function BacktestConfig() {
  const timeDurations = [
    { value: "1D", label: "1 Day" },
    { value: "1M", label: "1 Month" },
    { value: "1Y", label: "1 Year" },
    { value: "Custom", label: "Custom Range" },
  ];

  const [timeDuration, setTimeDuration] = useState("1Y");

  const frequencies = [
    { value: "1D", label: "1 Day" },
    { value: "1W", label: "1 Week" },
    { value: "1M", label: "1 Month" },
    { value: "Custom", label: "Custom Frequency" },
  ];
  const [frequency, setFrequency] = useState("1D");

  const instruments = ["SPY", "AAPL", "GOOGL", "MSFT", "AMZN"];
  const [instrument, setInstrument] = useState("SPY");
  const [startDate, setStartDate] = useState(dayjs().subtract(1, "year"));
  const [endDate, setEndDate] = useState(dayjs());

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
        value={instrument}
        onChange={(_, newValue) => {
          if (newValue) setInstrument(newValue);
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            size="small"
            sx={{
              // Match the container height
              "& .MuiInputBase-root": {
                height: 48,
                backgroundColor: "background.default", // matches theme background
              },
              // Match text input padding
              "& .MuiInputBase-input": {
                padding: "12px 14px",
              },
              // Optional: match label size
              "& .MuiInputLabel-root": {
                top: -5,
              },
            }}
          />
        )}
        sx={{
          mb: 2,
          "& .MuiAutocomplete-endAdornment": {
            backgroundColor: "background.default", // match input background
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
          const isSelected = timeDuration === value;
          return (
            <Tooltip key={value} title={label} arrow>
              <Button
                onClick={() => setTimeDuration(value)}
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

      {timeDuration === "Custom" && (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
          <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
            <DatePicker
              label="Start Date"
              value={startDate}
              onChange={(newValue) => newValue && setStartDate(newValue)}
              slotProps={{
                textField: {
                  fullWidth: true,
                  size: "small",
                },
              }}
            />
            <DatePicker
              label="End Date"
              value={endDate}
              onChange={(newValue) => newValue && setEndDate(newValue)}
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
          const isSelected = frequency === value;
          return (
            <Tooltip key={value} title={label} arrow>
              <Button
                onClick={() => setFrequency(value)}
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

      {frequency === "Custom" && (
        <TextField
          fullWidth
          label="Frequency (Days)"
          type="number"
          size="small"
          defaultValue={60}
          sx={{
            mb: 2,
            "& .MuiInputBase-root": {
              height: 48, // You can adjust this
            },
            "& .MuiInputLabel-root": {
              color: "text.secondary",
            },
            "& .MuiInputLabel-root.Mui-focused": {
              color: "text.secondary", // Prevents label from turning blue
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
          defaultValue={10000}
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">$</InputAdornment>
              ),
              inputProps: {
                min: 0,
                step: 100,
              },
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
