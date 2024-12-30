"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import Alert from "@mui/material/Alert";
import Paper from "@mui/material/Paper";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useTranslation } from "@/hooks/useTranslation";

export default function LoginPage() {
  const router = useRouter();
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      const response = await fetch(`${config.API_SERVER_PREFIX}/hub/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        switch (response.status) {
          case 401:
            throw new Error(t("auth.errors.invalidCredentials"));
          case 422:
            throw new Error(t("auth.errors.accountDisabled"));
          case 500:
          case 501:
          case 502:
          case 503:
          case 504:
            throw new Error(t("auth.errors.serverError"));
          default:
            throw new Error(t("auth.loginFailed"));
        }
      }

      const data = await response.json();
      // Store the token in a cookie and redirect to TFA page
      Cookies.set("tfa_token", data.token, { path: "/" });
      router.push("/tfa");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("auth.loginFailed"));
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Paper
          elevation={3}
          sx={{
            p: 4,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            width: "100%",
          }}
        >
          {/* Logo placeholder */}
          <Box
            sx={{
              width: 64,
              height: 64,
              bgcolor: "primary.main",
              borderRadius: 2,
              mb: 2,
            }}
          />
          <Typography component="h1" variant="h5">
            {t("common.login")}
          </Typography>
          <Box
            component="form"
            onSubmit={handleSubmit}
            noValidate
            sx={{ mt: 1 }}
          >
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label={t("common.email")}
              name="email"
              autoComplete="email"
              autoFocus
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label={t("common.password")}
              type="password"
              id="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              {t("common.login")}
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
}