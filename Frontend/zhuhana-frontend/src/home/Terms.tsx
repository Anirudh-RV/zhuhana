import { useEffect } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import AppTheme from "../shared-ui-theme/AppTheme";
import AppAppBar from "./components/AppAppBar";

export default function TermsPage(props: { disableCustomTheme?: boolean }) {
  useEffect(() => {
    document.title = "Zhuhana | Terms and Conditions";
  }, []);
  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />

      <AppAppBar />
    </AppTheme>
  );
}
