import { useEffect } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import AppTheme from "../shared-ui-theme/AppTheme";
import AppAppBar from "./components/AppAppBar";

export default function PrivacyPage(props: { disableCustomTheme?: boolean }) {
  useEffect(() => {
    document.title = "Zhuhana | Privacy";
  }, []);
  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />

      <AppAppBar />
    </AppTheme>
  );
}
