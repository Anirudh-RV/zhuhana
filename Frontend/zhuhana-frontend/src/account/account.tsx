import { useEffect } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import AppTheme from "../shared-ui-theme/AppTheme";
import IconButton from "@mui/material/IconButton";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { useNavigate } from "react-router-dom";
import ColorModeIconDropdown from "../shared-ui-theme/ColorModeIconDropdown";

export default function Account(props: { disableCustomTheme?: boolean }) {
  useEffect(() => {
    document.title = "Zhuhana | Account";
    window.scrollTo({ top: 0, behavior: "smooth" });
  }, []);

  const navigate = useNavigate();

  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />
      <Box
        sx={(theme) => ({
          position: "relative",
          minHeight: "100vh",
          "&::before": {
            content: '""',
            position: "absolute",
            inset: 0,
            zIndex: -1,
            backgroundRepeat: "no-repeat",
            backgroundImage:
              "radial-gradient(ellipse 80% 25% at 50% 0%, hsl(210, 100%, 90%), transparent)",
            ...theme.applyStyles?.("dark", {
              backgroundImage:
                "radial-gradient(ellipse 80% 25% at 50% 0%, hsl(210, 100%, 16%), transparent)",
            }),
          },
        })}
      >
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            px: 4,
            pt: 3,
          }}
        >
          <IconButton
            size="small"
            onClick={() => {
              if (window.history.length > 1) {
                navigate(-1);
              } else {
                navigate("/dashboard");
              }
            }}
            sx={{
              "&:hover": {
                backgroundColor: "action.hover",
              },
            }}
          >
            <ArrowBackIcon fontSize="small" />
          </IconButton>
          {/* Replace this with your actual color picker component */}
          <ColorModeIconDropdown
            sx={{ position: "fixed", top: "1rem", right: "1rem" }}
          />
          ;
        </Box>
        <Container
          maxWidth="md"
          sx={{
            py: 12,
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            minHeight: "60vh",
            textAlign: "center",
          }}
        >
          <Box>
            <Typography variant="h3" gutterBottom>
              Account Page
            </Typography>
            <Typography variant="h5" color="text.secondary">
              Coming soon...
            </Typography>
          </Box>
        </Container>
      </Box>
    </AppTheme>
  );
}
