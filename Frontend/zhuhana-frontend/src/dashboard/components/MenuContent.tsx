import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import HomeRoundedIcon from "@mui/icons-material/HomeRounded";
import AnalyticsRoundedIcon from "@mui/icons-material/AnalyticsRounded";
import LockIcon from "@mui/icons-material/Lock";
import HelpRoundedIcon from "@mui/icons-material/HelpRounded";

export default function MenuContent({
  selectedPage,
  onSelectPage,
}: {
  selectedPage: "home" | "analytics" | "vault";
  onSelectPage: (page: "home" | "analytics" | "vault") => void;
}) {
  const mainListItems = [
    { text: "Home", icon: <HomeRoundedIcon />, key: "home" },
    { text: "Analytics", icon: <AnalyticsRoundedIcon />, key: "analytics" },
    { text: "Vault", icon: <LockIcon />, key: "vault" },
  ];

  return (
    <Stack sx={{ flexGrow: 1, p: 1, justifyContent: "space-between" }}>
      <List dense>
        {mainListItems.map((item) => (
          <ListItem key={item.key} disablePadding sx={{ display: "block" }}>
            <ListItemButton
              onClick={() => onSelectPage(item.key as any)}
              selected={selectedPage === item.key}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>

      <List dense>
        <ListItem disablePadding sx={{ display: "block" }}>
          <ListItemButton
            onClick={() => {
              window.open("mailto:support@zhuhana.com", "_blank");
            }}
          >
            <ListItemIcon>
              <HelpRoundedIcon />
            </ListItemIcon>
            <ListItemText primary="Feedback" />
          </ListItemButton>
        </ListItem>
      </List>
    </Stack>
  );
}
