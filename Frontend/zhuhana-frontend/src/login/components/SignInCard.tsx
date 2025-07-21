import * as React from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import MuiCard from "@mui/material/Card";
import FormLabel from "@mui/material/FormLabel";
import FormControl from "@mui/material/FormControl";
import MuiLink from "@mui/material/Link";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { styled } from "@mui/material/styles";
import { useNavigate, Link } from "react-router-dom";
import { MuiOtpInput } from "mui-one-time-password-input";
import { useAuth } from "../../AuthContext";
import ForgotPassword from "./ForgotPassword";

import {
  LOGIN_V1_VERIFY_PASSWORD_ENDPOINT,
  LOGIN_V1_VERIFY_OTP_ENDPOINT,
} from "../../constants";

const Card = styled(MuiCard)(({ theme }) => ({
  display: "flex",
  flexDirection: "column",
  alignSelf: "center",
  width: "100%",
  padding: theme.spacing(4),
  gap: theme.spacing(2),
  boxShadow:
    "hsla(220, 30%, 5%, 0.05) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.05) 0px 15px 35px -5px",
  [theme.breakpoints.up("sm")]: {
    width: "450px",
  },
  ...theme.applyStyles("dark", {
    boxShadow:
      "hsla(220, 30%, 5%, 0.5) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.08) 0px 15px 35px -5px",
  }),
}));

export default function SignInCard() {
  const [emailError, setEmailError] = React.useState(false);
  const [emailErrorMessage, setEmailErrorMessage] = React.useState("");
  const [passwordError, setPasswordError] = React.useState(false);
  const [passwordErrorMessage, setPasswordErrorMessage] = React.useState("");
  const [otpSent, setOtpSent] = React.useState(false);
  const [otp, setOtp] = React.useState("");
  const [emailForOtp, setEmailForOtp] = React.useState("");
  const [forgotPasswordOpen, setForgotPasswordOpen] = React.useState(false);

  const [otpError, setOtpError] = React.useState(false);
  const [otpErrorMessage, setOtpErrorMessage] = React.useState("");

  const openForgotPasswordDialog = () => setForgotPasswordOpen(true);
  const closeForgotPasswordDialog = () => setForgotPasswordOpen(false);

  const navigate = useNavigate();
  const { setAuth } = useAuth();

  const validateInputs = () => {
    const email = document.getElementById("email") as HTMLInputElement;
    const password = document.getElementById("password") as HTMLInputElement;

    let isValid = true;

    if (!email.value || !/\S+@\S+\.\S+/.test(email.value)) {
      setEmailError(true);
      setEmailErrorMessage("Please enter a valid email address.");
      isValid = false;
    } else {
      setEmailError(false);
      setEmailErrorMessage("");
    }

    if (!password.value || password.value.length < 6) {
      setPasswordError(true);
      setPasswordErrorMessage("Password must be at least 6 characters long.");
      isValid = false;
    } else {
      setPasswordError(false);
      setPasswordErrorMessage("");
    }

    return isValid;
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!validateInputs()) return;

    const data = new FormData(event.currentTarget);
    const payload = {
      emailId: data.get("email"),
      password: data.get("password"),
    };

    setEmailForOtp(payload.emailId as string);

    try {
      const res = await fetch(LOGIN_V1_VERIFY_PASSWORD_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        setOtpSent(true);
        setPasswordError(false);
        setPasswordErrorMessage("");
      } else {
        const error = await res.json();

        setPasswordError(true);
        setPasswordErrorMessage("Login Error, incorrect email or password");
      }
    } catch (err) {
      alert("Network error");
    }
  };

  const handleOtpVerification = async () => {
    setOtpError(false);
    setOtpErrorMessage("");

    try {
      const res = await fetch(LOGIN_V1_VERIFY_OTP_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ emailId: emailForOtp, otp }),
      });

      if (res.ok) {
        const data = await res.json();
        setAuth(data.user, data.accessToken);
        navigate("/dashboard");
      } else {
        const error = await res.json();
        setOtpError(true);
        setOtpErrorMessage("Invalid OTP, resend and try again.");
      }
    } catch (err) {
      setOtpError(true);
      setOtpErrorMessage("Network error. Please try again.");
    }
  };

  return (
    <Card variant="outlined">
      <Box sx={{ display: { xs: "flex", md: "none" } }}></Box>
      <Typography
        component="h1"
        variant="h4"
        sx={{ width: "100%", fontSize: "clamp(2rem, 10vw, 2.15rem)" }}
      >
        {otpSent ? "Verify OTP" : "Log in"}
      </Typography>

      {!otpSent ? (
        <Box
          component="form"
          onSubmit={handleSubmit}
          noValidate
          sx={{ display: "flex", flexDirection: "column", gap: 2 }}
        >
          <FormLabel htmlFor="email">Email</FormLabel>
          <TextField
            error={emailError}
            helperText={emailErrorMessage}
            id="email"
            type="email"
            name="email"
            placeholder="your@email.com"
            autoComplete="email"
            required
            fullWidth
            variant="outlined"
          />

          <FormLabel htmlFor="password">Password</FormLabel>
          <TextField
            error={passwordError}
            helperText={passwordErrorMessage}
            name="password"
            placeholder="••••••"
            type="password"
            id="password"
            autoComplete="current-password"
            required
            fullWidth
            variant="outlined"
          />

          <Typography
            variant="body2"
            sx={{
              textAlign: "right",
              cursor: "pointer",
              color: "primary.main",
            }}
            onClick={openForgotPasswordDialog}
          >
            Forgot password?
          </Typography>

          <Button type="submit" fullWidth variant="contained">
            Log in
          </Button>

          <Typography sx={{ textAlign: "center" }}>
            Don't have an account?{" "}
            <Link to="/signup" style={{ alignSelf: "center" }}>
              Sign up
            </Link>
          </Typography>
        </Box>
      ) : (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
          <MuiOtpInput
            value={otp}
            onChange={(value) => {
              setOtp(value);
              setOtpError(false);
              setOtpErrorMessage("");
            }}
            length={6}
            autoFocus
            TextFieldsProps={{
              size: "small",
              error: otpError,
              sx: {
                width: "3rem",
                mx: 0.5,
                input: {
                  color: "text.primary",
                  backgroundColor: "background.default",
                },
              },
            }}
          />
          {otpError && (
            <Typography
              variant="caption"
              color="error"
              sx={{ textAlign: "center", mt: -1 }}
            >
              {otpErrorMessage}
            </Typography>
          )}
          <Button variant="contained" onClick={handleOtpVerification}>
            Verify OTP
          </Button>
          <Typography variant="caption" color="text.secondary">
            Didn't receive OTP?{" "}
            <MuiLink
              sx={{ cursor: "pointer" }}
              onClick={() => alert("Resend OTP endpoint not implemented yet")}
            >
              Resend
            </MuiLink>
          </Typography>
        </Box>
      )}
      <ForgotPassword
        open={forgotPasswordOpen}
        handleClose={closeForgotPasswordDialog}
      />
    </Card>
  );
}
