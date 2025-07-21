import { useEffect } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import AppTheme from "../shared-ui-theme/AppTheme";
import AppAppBar from "./components/AppAppBar";

export default function ContactPage(props: { disableCustomTheme?: boolean }) {
  useEffect(() => {
    document.title = "Contact Zhuhana";
  }, []);
  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />

      <AppAppBar />
    </AppTheme>
  );
}
