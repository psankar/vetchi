"use client";

import { useEffect, useState, use } from "react";
import {
  Box,
  Paper,
  Typography,
  Alert,
  CircularProgress,
  Button,
  Grid,
  Divider,
  Card,
  CardContent,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
} from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import {
  OpeningState,
  OpeningStates,
} from "@psankar/vetchi-typespec/common/openings";

interface Opening {
  id: string;
  title: string;
  positions: number;
  filled_positions: number;
  jd: string;
  recruiter: {
    name: string;
    email: string;
  };
  hiring_manager: {
    name: string;
    email: string;
  };
  cost_center_name: string;
  opening_type: string;
  state: OpeningState;
  created_at: string;
  last_updated_at: string;
}

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

export default function OpeningDetail({ params }: PageProps) {
  const { id } = use(params);
  const [opening, setOpening] = useState<Opening | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [showCloseConfirm, setShowCloseConfirm] = useState(false);
  const { t } = useTranslation();
  const router = useRouter();

  useEffect(() => {
    let isMounted = true;

    const fetchOpening = async () => {
      try {
        const sessionToken = Cookies.get("session_token");
        if (!sessionToken) {
          if (isMounted) {
            setError(t("auth.unauthorized"));
            setIsLoading(false);
          }
          return;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/employer/get-opening`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${sessionToken}`,
            },
            body: JSON.stringify({ id }),
          }
        );

        if (!isMounted) return;

        if (response.status === 200) {
          const data = await response.json();
          setOpening(data);
        } else if (response.status === 401) {
          setError(t("auth.unauthorized"));
        } else {
          setError(t("common.error"));
        }
      } catch (err) {
        if (isMounted) {
          setError(t("common.error"));
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    fetchOpening();

    return () => {
      isMounted = false;
    };
  }, [id]);

  const handleStateChange = async (toState: OpeningState) => {
    if (!opening) return;

    try {
      const sessionToken = Cookies.get("session_token");
      if (!sessionToken) {
        setError(t("auth.unauthorized"));
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/change-opening-state`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${sessionToken}`,
          },
          body: JSON.stringify({
            opening_id: opening.id,
            from_state: opening.state,
            to_state: toState,
          }),
        }
      );

      if (response.status === 200) {
        // Refresh the opening data
        setOpening((prev) => (prev ? { ...prev, state: toState } : null));
      } else if (response.status === 401) {
        setError(t("auth.unauthorized"));
      } else if (response.status === 409) {
        setError(t("openings.invalidStateTransition"));
      } else {
        setError(t("common.error"));
      }
    } catch (err) {
      setError(t("common.error"));
    }
  };

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", my: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  if (!opening) {
    return <Alert severity="info">{t("openings.notFound")}</Alert>;
  }

  return (
    <Box sx={{ width: "100%", p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Box>
          <Typography variant="h4">{opening.title}</Typography>
          <Typography variant="subtitle1" color="textSecondary" sx={{ mt: 1 }}>
            {t(`openings.state.${opening.state.toLowerCase()}`)}
          </Typography>
        </Box>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                {t("openings.details")}
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="subtitle2">
                    {t("openings.id")}
                  </Typography>
                  <Typography>{opening.id}</Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2">
                    {t("openings.positions")}
                  </Typography>
                  <Typography>{opening.positions}</Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2">
                    {t("openings.filledPositions")}
                  </Typography>
                  <Typography>{opening.filled_positions}</Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2">
                    {t("openings.costCenter")}
                  </Typography>
                  <Typography>{opening.cost_center_name}</Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2">
                    {t("openings.type")}
                  </Typography>
                  <Typography>{opening.opening_type}</Typography>
                </Grid>
                <Grid item xs={12}>
                  <Typography variant="subtitle2" sx={{ mt: 2 }}>
                    {t("openings.description")}
                  </Typography>
                  <Typography sx={{ whiteSpace: "pre-wrap" }}>
                    {opening.jd}
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                {t("openings.contacts")}
              </Typography>
              <Typography variant="subtitle2">
                {t("openings.recruiter")}
              </Typography>
              <Typography>{opening.recruiter.name}</Typography>
              <Typography color="textSecondary">
                {opening.recruiter.email}
              </Typography>

              <Box sx={{ my: 2 }}>
                <Divider />
              </Box>

              <Typography variant="subtitle2">
                {t("openings.hiringManager")}
              </Typography>
              <Typography>{opening.hiring_manager.name}</Typography>
              <Typography color="textSecondary">
                {opening.hiring_manager.email}
              </Typography>
            </CardContent>
          </Card>

          <Card sx={{ mt: 2 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                {t("openings.actions")}
              </Typography>
              <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                {opening.state === OpeningStates.DRAFT && (
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => handleStateChange(OpeningStates.ACTIVE)}
                  >
                    {t("openings.publish")}
                  </Button>
                )}
                {opening.state === OpeningStates.ACTIVE && (
                  <Button
                    variant="contained"
                    color="warning"
                    onClick={() => handleStateChange(OpeningStates.SUSPENDED)}
                  >
                    {t("openings.suspend")}
                  </Button>
                )}
                {opening.state === OpeningStates.SUSPENDED && (
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => handleStateChange(OpeningStates.ACTIVE)}
                  >
                    {t("openings.reactivate")}
                  </Button>
                )}
                {opening.state !== OpeningStates.CLOSED && (
                  <Button
                    variant="contained"
                    color="error"
                    onClick={() => setShowCloseConfirm(true)}
                  >
                    {t("openings.close")}
                  </Button>
                )}
                <Button
                  variant="outlined"
                  onClick={() =>
                    router.push(`/openings/${opening.id}/candidacies`)
                  }
                >
                  {t("openings.viewCandidacies")}
                </Button>
                <Button
                  variant="outlined"
                  onClick={() =>
                    router.push(`/openings/${opening.id}/interviews`)
                  }
                >
                  {t("openings.viewInterviews")}
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Dialog
        open={showCloseConfirm}
        onClose={() => setShowCloseConfirm(false)}
      >
        <DialogTitle>{t("openings.closeConfirmTitle")}</DialogTitle>
        <DialogContent>
          <DialogContentText>
            {t("openings.closeConfirmMessage")}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowCloseConfirm(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={() => {
              handleStateChange(OpeningStates.CLOSED);
              setShowCloseConfirm(false);
            }}
            color="error"
            variant="contained"
          >
            {t("openings.confirmClose")}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
