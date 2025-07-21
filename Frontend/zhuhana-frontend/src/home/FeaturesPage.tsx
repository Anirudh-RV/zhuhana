import { useEffect } from "react";
import { useLocation } from "react-router-dom";
import CssBaseline from "@mui/material/CssBaseline";
import Divider from "@mui/material/Divider";
import AppTheme from "../shared-ui-theme/AppTheme";
import AppAppBar from "./components/AppAppBar";
import Hero from "./components/Hero";
import Highlights from "./components/Highlights";
import Pricing from "./components/Pricing";
import Features from "./components/Features";
import FAQ from "./components/FAQ";
import Footer from "./components/Footer";

export default function FeaturesPage(props: { disableCustomTheme?: boolean }) {
  const location = useLocation();

  useEffect(() => {
    document.title = "Zhuhana | Features";

    setTimeout(() => {
      if (location.hash) {
        const el = document.querySelector(location.hash);
        if (el) {
          el.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      } else {
        // Fallback scroll to top if no hash
        window.scrollTo({ top: 0, behavior: "smooth" });
      }
    }, 100);
  }, [location]);

  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />
      <AppAppBar />
      <div>
        <Features />
        <Divider />
        <div id="highlights">
          <Highlights />
        </div>
        <Divider />
        <Footer />
      </div>
    </AppTheme>
  );
}
