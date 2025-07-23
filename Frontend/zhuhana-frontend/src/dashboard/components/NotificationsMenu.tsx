// NotificationMenu.tsx (updated)
import React, { useEffect, useState } from "react";
import {
  Menu,
  MenuItem,
  IconButton,
  Badge,
  ListItemText,
  Typography,
  Divider,
} from "@mui/material";
import NotificationsRoundedIcon from "@mui/icons-material/NotificationsRounded";
import { useAuth } from "../../AuthContext";
import {
  GET_NOTIFICATIONS_V1_INIT_ENDPOINT,
  READ_NOTIFICATIONS_V1_INIT_ENDPOINT,
} from "../../constants";

interface Notification {
  ID: string;
  Title: string;
  Message: string;
  Link?: string;
  Read: boolean;
  CreatedAt: string;
}

export default function NotificationMenu() {
  const { accessToken } = useAuth(); // Ensure this is available in your context
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const unreadCount = notifications.filter((n) => !n.Read).length;

  const fetchNotifications = async () => {
    try {
      const res = await fetch(GET_NOTIFICATIONS_V1_INIT_ENDPOINT, {
        headers: {
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
        },
      });

      const data = await res.json();
      if (data.status === 1) {
        setNotifications(data.notifications);
      } else {
        console.error("Failed to fetch notifications:", data.statusDescription);
      }
    } catch (error) {
      console.error("Error fetching notifications:", error);
    }
  };

  const markNotificationsAsRead = async (ids: string[]) => {
    if (ids.length === 0) return;

    try {
      await fetch(READ_NOTIFICATIONS_V1_INIT_ENDPOINT, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
        },
        body: JSON.stringify({ ids }),
      });

      // Optimistically update state
      setNotifications((prev) =>
        prev.map((n) => (ids.includes(n.ID) ? { ...n, Read: true } : n))
      );
    } catch (error) {
      console.error("Error marking notifications as read:", error);
    }
  };

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);

    // On open, mark unread notifications as read
    const unreadIDs = notifications.filter((n) => !n.Read).map((n) => n.ID);
    markNotificationsAsRead(unreadIDs);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  useEffect(() => {
    fetchNotifications();
  }, []);

  return (
    <>
      <IconButton onClick={handleClick} color="inherit">
        <Badge badgeContent={unreadCount} color="error">
          <NotificationsRoundedIcon />
        </Badge>
      </IconButton>
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleClose}
        PaperProps={{
          sx: { width: 320 },
        }}
      >
        <Typography variant="h6" sx={{ px: 2, py: 1 }}>
          Notifications
        </Typography>
        <Divider />
        {notifications.length === 0 && (
          <MenuItem disabled>
            <ListItemText primary="No notifications" />
          </MenuItem>
        )}
        {notifications.map((notification) => (
          <MenuItem
            key={notification.ID}
            onClick={handleClose}
            sx={{ whiteSpace: "normal", alignItems: "flex-start" }}
          >
            <ListItemText
              primary={
                <Typography fontWeight={notification.Read ? "normal" : "bold"}>
                  {notification.Title}
                </Typography>
              }
              secondary={
                <Typography variant="body2" color="text.secondary">
                  {notification.Message}
                </Typography>
              }
            />
          </MenuItem>
        ))}
      </Menu>
    </>
  );
}
