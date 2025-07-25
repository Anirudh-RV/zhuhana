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
import { Link, useNavigate } from "react-router-dom";
import { MuiOtpInput } from "mui-one-time-password-input";
import { useAuth } from "../../AuthContext";
import {
  SIGN_UP_V1_INIT_ENDPOINT,
  SIGN_UP_V1_VERIFY_OTP_ENDPOINT,
} from "../../constants";
import CircularProgress from "@mui/material/CircularProgress";

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

export default function SignUpCard() {
  const [emailError, setEmailError] = React.useState(false);
  const [emailErrorMessage, setEmailErrorMessage] = React.useState("");
  const [passwordError, setPasswordError] = React.useState(false);
  const [passwordErrorMessage, setPasswordErrorMessage] = React.useState("");
  const [nameError, setNameError] = React.useState(false);
  const [nameErrorMessage, setNameErrorMessage] = React.useState("");
  const [otpSent, setOtpSent] = React.useState(false);
  const [otp, setOtp] = React.useState("");
  const [emailForOtp, setEmailForOtp] = React.useState("");
  const [passwordForOtp, setPasswordForOtp] = React.useState("");
  const [nameForOtp, setNameForOtp] = React.useState("");
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [isSubmittingOtp, setIsSubmittingOtp] = React.useState(false);

  const [otpError, setOtpError] = React.useState(false);
  const [otpErrorMessage, setOtpErrorMessage] = React.useState("");

  const [resendCount, setResendCount] = React.useState(0);
  const [resendCooldown, setResendCooldown] = React.useState(0);

  React.useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(
        () => setResendCooldown(resendCooldown - 1),
        1000
      );
      return () => clearTimeout(timer);
    }
  }, [resendCooldown]);

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

    setIsSubmitting(true);

    const data = new FormData(event.currentTarget);

    const fullName = data.get("name")?.toString().trim() || "";
    const nameParts = fullName.split(/\s+/);

    let firstName = "",
      middleName = "",
      lastName = "";

    if (nameParts.length === 1) {
      firstName = nameParts[0];
    } else if (nameParts.length === 2) {
      [firstName, lastName] = nameParts;
    } else if (nameParts.length >= 3) {
      firstName = nameParts[0];
      lastName = nameParts[nameParts.length - 1];
      middleName = nameParts.slice(1, -1).join(" ");
    }

    const emailId = data.get("email") as string;
    const password = data.get("password") as string;

    const payload = {
      emailId,
      firstName,
      middleName,
      lastName,
      password,
    };

    setEmailForOtp(emailId);
    setPasswordForOtp(password);
    setNameForOtp(fullName);

    try {
      const res = await fetch(SIGN_UP_V1_INIT_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        setOtpSent(true);
        setResendCount(0);
        setResendCooldown(60);
      } else {
        const error = await res.json();
        console.error("Signup error", error);
      }
    } catch (err) {
      console.error("Network error", err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleResendOtp = async () => {
    if (resendCount >= 2 || resendCooldown > 0) return;

    try {
      const fullName = nameForOtp.trim();
      const nameParts = fullName.split(/\s+/);
      let firstName = "",
        middleName = "",
        lastName = "";

      if (nameParts.length === 1) {
        firstName = nameParts[0];
      } else if (nameParts.length === 2) {
        [firstName, lastName] = nameParts;
      } else {
        firstName = nameParts[0];
        lastName = nameParts[nameParts.length - 1];
        middleName = nameParts.slice(1, -1).join(" ");
      }

      const payload = {
        emailId: emailForOtp,
        firstName,
        middleName,
        lastName,
        password: passwordForOtp,
      };

      const res = await fetch(SIGN_UP_V1_INIT_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        setResendCount((prev) => prev + 1);
        setResendCooldown(60);
        setOtpError(false);
        setOtpErrorMessage("");
      } else {
        setOtpError(true);
        setOtpErrorMessage("Failed to resend OTP. Please try again.");
      }
    } catch (err) {
      setOtpError(true);
      setOtpErrorMessage("Network error while resending OTP.");
    }
  };

  const handleOtpVerification = async () => {
    setIsSubmittingOtp(true);
    setOtpError(false);
    setOtpErrorMessage("");

    try {
      const res = await fetch(SIGN_UP_V1_VERIFY_OTP_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ emailId: emailForOtp, otp }),
      });

      if (res.ok) {
        const data = await res.json();
        setAuth(data.user, data.accessToken);
        navigate("/dashboard");
      } else {
        const errData = await res.json();
        setOtpError(true);
        setOtpErrorMessage(
          errData.statusDescription || "Invalid OTP, please try again."
        );
      }
    } catch (err) {
      setOtpError(true);
      setOtpErrorMessage("Network error. Please try again.");
    } finally {
      setIsSubmittingOtp(false);
    }
  };

  return (
    <Card variant="outlined">
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
          sx={{ display: "flex", flexDirection: "column", gap: 2 }}
        >
          <FormControl>
            <FormLabel htmlFor="name">Full name</FormLabel>
            <TextField
              id="name"
              name="name"
              required
              fullWidth
              placeholder="Jon Snow"
              error={nameError}
              helperText={nameErrorMessage}
            />
          </FormControl>

          <FormControl>
            <FormLabel htmlFor="email">Email</FormLabel>
            <TextField
              id="email"
              name="email"
              type="email"
              required
              fullWidth
              placeholder="your@email.com"
              error={emailError}
              helperText={emailErrorMessage}
            />
          </FormControl>

          <FormControl>
            <FormLabel htmlFor="password">Password</FormLabel>
            <TextField
              id="password"
              name="password"
              type="password"
              required
              fullWidth
              placeholder="••••••"
              error={passwordError}
              helperText={passwordErrorMessage}
            />
          </FormControl>

          <FormControlLabel
            control={<Checkbox defaultChecked />}
            label={
              <Typography variant="caption" color="text.secondary">
                Agree to our <Link to="/terms">Terms of Service</Link> and{" "}
                <Link to="/privacy">Privacy Policy</Link>.
              </Typography>
            }
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={isSubmitting}
            sx={{
              color: isSubmitting ? "common.white" : undefined,
            }}
          >
            {isSubmitting ? (
              <CircularProgress size={24} color="primary" />
            ) : (
              "Sign up"
            )}
          </Button>

          <Typography sx={{ textAlign: "center" }}>
            Already have an account?{" "}
            <Link to="/login" style={{ alignSelf: "center" }}>
              Log in
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
              error: otpError,
              sx: {
                mx: 0.5,
                "& .MuiInputBase-root": {
                  width: "3rem",
                  height: "4rem",
                  borderRadius: "0.5rem",
                },
                "& input": {
                  width: "100%",
                  height: "100%",
                  textAlign: "center",
                  fontSize: "2rem",
                  fontWeight: 700,
                  lineHeight: 1.2,
                  padding: 0,
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

          <Button
            variant="contained"
            onClick={handleOtpVerification}
            disabled={isSubmittingOtp}
          >
            {isSubmittingOtp ? (
              <CircularProgress size={24} color="primary" />
            ) : (
              "Verify OTP"
            )}
          </Button>

          <Typography
            variant="caption"
            color="text.secondary"
            sx={{ textAlign: "center" }}
          >
            Didn't receive OTP?{" "}
            <MuiLink
              onClick={handleResendOtp}
              sx={{
                cursor:
                  resendCooldown > 0 || resendCount >= 2
                    ? "not-allowed"
                    : "pointer",
                pointerEvents:
                  resendCooldown > 0 || resendCount >= 2 ? "none" : "auto",
                color:
                  resendCooldown > 0 || resendCount >= 2
                    ? "text.disabled"
                    : "primary.main",
              }}
            >
              {resendCooldown > 0
                ? `Resend in ${resendCooldown}s`
                : resendCount >= 2
                ? "Resend limit reached"
                : "Resend"}
            </MuiLink>
          </Typography>
        </Box>
      )}
    </Card>
  );
}
