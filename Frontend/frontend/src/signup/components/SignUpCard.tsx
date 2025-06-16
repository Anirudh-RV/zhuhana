import * as React from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import MuiCard from "@mui/material/Card";
import Checkbox from "@mui/material/Checkbox";
import FormLabel from "@mui/material/FormLabel";
import FormControl from "@mui/material/FormControl";
import FormControlLabel from "@mui/material/FormControlLabel";
import MuiLink from "@mui/material/Link";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { styled } from "@mui/material/styles";
import { SitemarkIcon } from "./CustomIcons";
import { Link } from "react-router-dom";
import { MuiOtpInput } from "mui-one-time-password-input";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../AuthContext";

import {
  SIGN_UP_V1_INIT_ENDPOINT,
  SIGN_UP_V1_VERIFY_OTP_ENDPOINT,
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
  const [nameError, setNameError] = React.useState(false);
  const [nameErrorMessage, setNameErrorMessage] = React.useState("");
  const [otpSent, setOtpSent] = React.useState(false);
  const [otp, setOtp] = React.useState("");
  const [timer, setTimer] = React.useState(300); // optional: 5-min timer
  const [emailForOtp, setEmailForOtp] = React.useState("");
  const navigate = useNavigate();
  const { setAuth } = useAuth();

  const validateInputs = () => {
    const email = document.getElementById("email") as HTMLInputElement;
    const password = document.getElementById("password") as HTMLInputElement;
    const name = document.getElementById("name") as HTMLInputElement;

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

    if (!name.value || name.value.length < 1) {
      setNameError(true);
      setNameErrorMessage("Name is required.");
      isValid = false;
    } else {
      setNameError(false);
      setNameErrorMessage("");
    }

    return isValid;
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!validateInputs()) return;

    const data = new FormData(event.currentTarget);

    const payload = {
      emailId: data.get("email"),
      firstName: data.get("name"),
      lastName: "", // fill if you have it
      middleName: "",
      password: data.get("password"),
    };

    setEmailForOtp(payload.emailId as string);

    try {
      const res = await fetch(SIGN_UP_V1_INIT_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        setEmailForOtp(payload.emailId as string); // ✅ store email
        setOtpSent(true); // switch view to OTP
      } else {
        const error = await res.json();
        console.error("Signup error", error);
      }
    } catch (err) {
      console.error("Network error", err);
    }
  };

  return (
    <Card variant="outlined">
      <Box sx={{ display: { xs: "flex", md: "none" } }}>
        <SitemarkIcon />
      </Box>

      <Typography
        component="h1"
        variant="h4"
        sx={{ width: "100%", fontSize: "clamp(2rem, 10vw, 2.15rem)" }}
      >
        {otpSent ? "Verify OTP" : "Create your account"}
      </Typography>

      {!otpSent ? (
        <Box
          component="form"
          onSubmit={handleSubmit}
          noValidate
          sx={{
            display: "flex",
            flexDirection: "column",
            width: "100%",
            gap: 2,
          }}
        >
          <Typography
            component="h1"
            variant="h4"
            sx={{ width: "100%", fontSize: "clamp(2rem, 10vw, 2.15rem)" }}
          >
            Create your account
          </Typography>
          <FormControl>
            <FormLabel htmlFor="name">Full name</FormLabel>
            <TextField
              autoComplete="name"
              name="name"
              required
              fullWidth
              id="name"
              placeholder="Jon Snow"
              error={nameError}
              helperText={nameErrorMessage}
              color={nameError ? "error" : "primary"}
            />
          </FormControl>
          <FormControl>
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
              color={emailError ? "error" : "primary"}
            />
          </FormControl>
          <FormControl>
            <Box sx={{ display: "flex", justifyContent: "space-between" }}>
              <FormLabel htmlFor="password">Password</FormLabel>
            </Box>
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
              color={passwordError ? "error" : "primary"}
            />
          </FormControl>
          <FormControlLabel
            control={
              <Checkbox
                value="agreeToConditions"
                color="primary"
                defaultChecked
              />
            }
            label={
              <Box>
                <Typography
                  variant="caption"
                  color="text.secondary"
                  sx={{ textAlign: "center" }}
                >
                  Agree to our&nbsp;
                  <Link to="/sign-up" style={{ alignSelf: "center" }}>
                    Terms of Service
                  </Link>
                  &nbsp;and&nbsp;
                  <Link to="/sign-up" style={{ alignSelf: "center" }}>
                    Privacy Policy
                  </Link>
                  .
                </Typography>
              </Box>
            }
          />
          <Button type="submit" fullWidth variant="contained">
            Sign up
          </Button>
          <Typography sx={{ textAlign: "center" }}>
            Already have an account?{" "}
            <span>
              <Link to="/login" style={{ alignSelf: "center" }}>
                Log in
              </Link>
            </span>
          </Typography>
        </Box>
      ) : (
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
          <MuiOtpInput
            value={otp}
            onChange={setOtp}
            length={6}
            autoFocus
            TextFieldsProps={{ size: "small", sx: { width: "3rem", mx: 0.5 } }}
          />
          <Button
            variant="contained"
            onClick={() => {
              fetch(SIGN_UP_V1_VERIFY_OTP_ENDPOINT, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ emailId: emailForOtp, otp }),
              })
                .then(async (res) => {
                  if (res.ok) {
                    const data = await res.json();

                    // Store in localStorage
                    setAuth(data.user, data.accessToken);
                    navigate("/dashboard");

                    navigate("/dashboard");
                  } else {
                    const errData = await res.json();
                    alert(
                      errData.statusDescription || "OTP verification failed"
                    );
                  }
                })
                .catch((err) => {
                  console.error("Verification error:", err);
                  alert("An error occurred while verifying OTP");
                });
            }}
          >
            Verify OTP
          </Button>
          <Typography variant="caption" color="text.secondary">
            Didn’t receive OTP?{" "}
            <MuiLink
              onClick={() => {
                // Call resend OTP endpoint
              }}
              sx={{ cursor: "pointer" }}
            >
              Resend
            </MuiLink>
          </Typography>
        </Box>
      )}
    </Card>
  );
}
