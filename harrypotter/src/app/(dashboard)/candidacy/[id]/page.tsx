"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Button,
  TextField,
  Paper,
  Divider,
  Link,
  CircularProgress,
  IconButton,
  Avatar,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Stack,
  FormControlLabel,
  Switch,
  Collapse,
} from "@mui/material";
import {
  OpenInNew as OpenInNewIcon,
  ExpandMore as ExpandMoreIcon,
} from "@mui/icons-material";
import { useTranslation } from "@/hooks/useTranslation";
import {
  GetCandidacyInfoRequest,
  GetCandidacyCommentsRequest,
  Candidacy,
  CandidacyComment,
  CandidacyState,
  HubInterview as Interview,
  InterviewType,
  InterviewState,
  GetHubInterviewsByCandidacyRequest as GetInterviewsByCandidacyRequest,
  AddInterviewRequest,
  TimeZone,
  validTimezones,
} from "@psankar/vetchi-typespec";
import { AddEmployerCandidacyCommentRequest } from "@psankar/vetchi-typespec/employer/candidacy";
import { config } from "@/config";
import Cookies from "js-cookie";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";

function CandidacyStateLabel({
  state,
  t,
}: {
  state: CandidacyState;
  t: (key: string) => string;
}) {
  let color:
    | "primary"
    | "secondary"
    | "error"
    | "info"
    | "success"
    | "warning" = "info";
  switch (state) {
    case "INTERVIEWING":
      color = "info";
      break;
    case "OFFERED":
      color = "warning";
      break;
    case "OFFER_ACCEPTED":
      color = "success";
      break;
    case "OFFER_DECLINED":
    case "CANDIDATE_UNSUITABLE":
    case "CANDIDATE_NOT_RESPONDING":
    case "CANDIDATE_WITHDREW":
    case "EMPLOYER_DEFUNCT":
      color = "error";
      break;
  }
  return (
    <Chip label={t(`candidacies.states.${state}`)} color={color} size="small" />
  );
}

function InterviewStateLabel({
  state,
  t,
}: {
  state: InterviewState;
  t: (key: string) => string;
}) {
  let color:
    | "primary"
    | "secondary"
    | "error"
    | "info"
    | "success"
    | "warning" = "info";

  switch (state) {
    case "SCHEDULED_INTERVIEW":
      color = "info";
      break;
    case "COMPLETED_INTERVIEW":
      color = "success";
      break;
    case "CANCELLED_INTERVIEW":
      color = "error";
      break;
  }
  return (
    <Chip label={t(`interviews.states.${state}`)} color={color} size="small" />
  );
}

export default function CandidacyDetailPage() {
  const params = useParams();
  const candidacyId = params.id as string;
  const { t } = useTranslation();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [candidacy, setCandidacy] = useState<Candidacy | null>(null);
  const [comments, setComments] = useState<CandidacyComment[]>([]);
  const [newComment, setNewComment] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [interviews, setInterviews] = useState<Interview[]>([]);
  const [openAddInterview, setOpenAddInterview] = useState(false);
  const [showInterviews, setShowInterviews] = useState(true);
  const [expandedInterviews, setExpandedInterviews] = useState<
    Record<string, boolean>
  >({});

  // Initialize expanded state for new interviews
  useEffect(() => {
    const newExpandedState: Record<string, boolean> = {};
    interviews.forEach((interview) => {
      if (!(interview.interview_id in expandedInterviews)) {
        newExpandedState[interview.interview_id] = false;
      }
    });
    if (Object.keys(newExpandedState).length > 0) {
      setExpandedInterviews((prev) => ({ ...prev, ...newExpandedState }));
    }
  }, [interviews]);

  // Get user's timezone and find closest matching TimeZone enum value
  const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const userOffset = new Date().getTimezoneOffset();
  const offsetHours = Math.floor(Math.abs(userOffset) / 60);
  const offsetMinutes = Math.abs(userOffset) % 60;
  const offsetStr = `${userOffset <= 0 ? "+" : "-"}${offsetHours
    .toString()
    .padStart(2, "0")}${offsetMinutes.toString().padStart(2, "0")}`;

  // Find the closest matching timezone from validTimezones
  const defaultTimezone =
    Array.from(validTimezones).find((tz) => tz.includes(`GMT${offsetStr}`)) ||
    "UTC Coordinated Universal Time GMT+0000";

  const [newInterview, setNewInterview] = useState<{
    startTime: string;
    endTime: string;
    type: InterviewType;
    description: string;
    timezone: TimeZone;
  }>({
    startTime: "",
    endTime: "",
    type: "VIDEO_CALL",
    description: "",
    timezone: defaultTimezone,
  });

  const [allowPastDates, setAllowPastDates] = useState(false);
  const [use24HourFormat, setUse24HourFormat] = useState(true);

  // Handle localStorage in useEffect to avoid SSR issues
  useEffect(() => {
    const saved = localStorage.getItem("create_interview_24hour_format");
    if (saved !== null) {
      setUse24HourFormat(saved === "true");
    }
  }, []);

  // Update localStorage when time format preference changes
  useEffect(() => {
    localStorage.setItem(
      "create_interview_24hour_format",
      use24HourFormat.toString()
    );
  }, [use24HourFormat]);

  // Reset interview form with default timezone
  const resetInterviewForm = () => {
    setNewInterview({
      startTime: "",
      endTime: "",
      type: "VIDEO_CALL",
      description: "",
      timezone: defaultTimezone,
    });
    setAllowPastDates(false);
    // Don't reset the time format preference
  };

  // Fetch candidacy info, comments, and interviews
  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      // Fetch candidacy info
      const infoResponse = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-candidacy-info`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
          } as GetCandidacyInfoRequest),
        }
      );

      if (infoResponse.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!infoResponse.ok) throw new Error(t("candidacies.fetchError"));
      const candidacyData = await infoResponse.json();
      if (!candidacyData) {
        throw new Error(t("candidacies.fetchError"));
      }
      setCandidacy(candidacyData);

      // Fetch interviews
      const interviewsResponse = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-interviews-by-candidacy`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
          } as GetInterviewsByCandidacyRequest),
        }
      );

      if (!interviewsResponse.ok) throw new Error(t("interviews.fetchError"));
      const interviewsData = await interviewsResponse.json();
      setInterviews(interviewsData || []);

      // Fetch comments
      const commentsResponse = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-candidacy-comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
          } as GetCandidacyCommentsRequest),
        }
      );

      if (!commentsResponse.ok) throw new Error(t("candidacies.fetchError"));
      const commentsData = await commentsResponse.json();
      setComments(commentsData || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.serverError"));
    } finally {
      setLoading(false);
    }
  };

  // Add new comment
  const handleAddComment = async () => {
    if (!newComment.trim()) return;
    setSubmitting(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/add-candidacy-comment`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
            comment: newComment.trim(),
          } as AddEmployerCandidacyCommentRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) throw new Error(t("common.serverError"));
      setNewComment("");
      // Refresh comments
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.serverError"));
    } finally {
      setSubmitting(false);
    }
  };

  const handleAddInterview = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      // Convert local times to UTC for API
      const startDate = new Date(newInterview.startTime);
      const endDate = new Date(newInterview.endTime);

      // Get timezone offset from selected timezone
      const tzMatch = newInterview.timezone.match(/GMT([+-]\d{4})/);
      const tzOffset = tzMatch ? tzMatch[1] : "+0000";
      const tzHours = parseInt(tzOffset.slice(1, 3));
      const tzMinutes = parseInt(tzOffset.slice(3));
      const offsetMillis =
        (tzHours * 60 + tzMinutes) * 60 * 1000 * (tzOffset[0] === "+" ? -1 : 1);

      // Adjust dates to UTC
      const utcStartDate = new Date(startDate.getTime() + offsetMillis);
      const utcEndDate = new Date(endDate.getTime() + offsetMillis);

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/add-interview`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
            start_time: utcStartDate,
            end_time: utcEndDate,
            interview_type: newInterview.type,
            description: newInterview.description,
          } satisfies AddInterviewRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) throw new Error(t("interviews.addError"));

      setOpenAddInterview(false);
      resetInterviewForm();
      await fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.serverError"));
    }
  };

  // Fetch data on mount
  useEffect(() => {
    fetchData();
  }, []); // Empty dependency array means this runs once on mount

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">{error}</Typography>
        <Button variant="contained" onClick={fetchData} sx={{ mt: 2 }}>
          {t("common.retry")}
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("candidacies.viewCandidacy")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      {candidacy && (
        <Paper sx={{ p: 3, mb: 3 }}>
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              mb: 1,
            }}
          >
            <Typography variant="h6">{candidacy.applicant_name}</Typography>
            <CandidacyStateLabel state={candidacy.candidacy_state} t={t} />
          </Box>
          <Typography
            variant="subtitle1"
            gutterBottom
            sx={{ color: "text.secondary" }}
          >
            @{candidacy.applicant_handle}
          </Typography>
          <Divider sx={{ my: 2 }} />
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <Link
              href={`/openings/${candidacy.opening_id}`}
              target="_blank"
              rel="noopener noreferrer"
              sx={{
                color: "primary.main",
                textDecoration: "none",
                display: "flex",
                alignItems: "center",
                gap: 0.5,
              }}
            >
              <Typography variant="subtitle1">
                {candidacy.opening_title}
              </Typography>
              <OpenInNewIcon fontSize="small" />
            </Link>
          </Box>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            {candidacy.opening_description}
          </Typography>
        </Paper>
      )}

      {/* Interviews Section */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 1, mb: 2 }}>
          <Typography variant="h6">{t("interviews.title")}</Typography>
          <IconButton
            onClick={() => setShowInterviews(!showInterviews)}
            sx={{
              transform: showInterviews ? "rotate(0deg)" : "rotate(180deg)",
              transition: "transform 0.2s",
            }}
            size="small"
          >
            <ExpandMoreIcon />
          </IconButton>
          <Box sx={{ flex: 1 }} />
          <Button
            variant="contained"
            onClick={() => setOpenAddInterview(true)}
            size="small"
          >
            {t("interviews.addNew")}
          </Button>
        </Box>

        <Collapse in={showInterviews}>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            {interviews.length > 0 ? (
              interviews.map((interview) => (
                <Paper
                  key={interview.interview_id}
                  elevation={1}
                  sx={{
                    p: expandedInterviews[interview.interview_id] ? 3 : 2,
                    transition: "padding 0.2s",
                  }}
                >
                  {/* Header with start time and collapse control */}
                  <Box
                    sx={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center",
                      mb: expandedInterviews[interview.interview_id] ? 2 : 0,
                    }}
                  >
                    <Box
                      sx={{ display: "flex", alignItems: "baseline", gap: 1 }}
                    >
                      <Typography variant="body1" color="text.secondary">
                        {new Date(interview.start_time).toLocaleString(
                          "default",
                          {
                            weekday: "short",
                            year: "numeric",
                            month: "short",
                            day: "numeric",
                          }
                        )}
                      </Typography>
                      <Typography
                        variant="subtitle1"
                        sx={{ fontWeight: 500, color: "primary.main" }}
                      >
                        {new Date(interview.start_time).toLocaleTimeString(
                          "default",
                          {
                            hour: "2-digit",
                            minute: "2-digit",
                            hour12: undefined,
                          }
                        )}
                      </Typography>
                    </Box>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <InterviewStateLabel
                        state={interview.interview_state}
                        t={t}
                      />
                      <IconButton
                        onClick={() =>
                          setExpandedInterviews((prev) => ({
                            ...prev,
                            [interview.interview_id]:
                              !prev[interview.interview_id],
                          }))
                        }
                        sx={{
                          transform: expandedInterviews[interview.interview_id]
                            ? "rotate(180deg)"
                            : "rotate(0deg)",
                          transition: "transform 0.2s",
                        }}
                        size="small"
                      >
                        <ExpandMoreIcon />
                      </IconButton>
                    </Box>
                  </Box>

                  <Collapse in={expandedInterviews[interview.interview_id]}>
                    {/* Interview Type */}
                    <Typography
                      variant="subtitle1"
                      sx={{ mb: 2, fontWeight: 500 }}
                    >
                      {t(`interviews.types.${interview.interview_type}`)}
                    </Typography>

                    {/* Time section */}
                    <Box sx={{ mb: 3 }}>
                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 1 }}
                      >
                        <Typography variant="caption" color="text.secondary">
                          {t("interviews.endTime")}
                        </Typography>
                        <Typography>
                          {new Date(interview.end_time).toLocaleString(
                            "default",
                            {
                              weekday: "long",
                              year: "numeric",
                              month: "long",
                              day: "numeric",
                              hour: "2-digit",
                              minute: "2-digit",
                              hour12: undefined,
                            }
                          )}
                        </Typography>
                      </Box>
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ mt: 0.5, display: "block" }}
                      >
                        {Intl.DateTimeFormat().resolvedOptions().timeZone}
                      </Typography>
                    </Box>

                    {/* Interviewers section */}
                    <Box sx={{ mb: 2 }}>
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ mb: 1, display: "block" }}
                      >
                        {t("interviews.interviewers")}
                      </Typography>
                      <Table>
                        <TableBody>
                          <TableRow>
                            <TableCell>
                              <Box
                                sx={{
                                  display: "flex",
                                  flexDirection: "column",
                                }}
                              >
                                {interview.interviewers?.map(
                                  (interviewer, idx) => (
                                    <Box
                                      key={idx}
                                      sx={{
                                        display: "flex",
                                        alignItems: "center",
                                        gap: 1,
                                        mb: 0.5,
                                      }}
                                    >
                                      <Avatar
                                        sx={{ width: 24, height: 24 }}
                                        alt={interviewer}
                                      >
                                        {interviewer.charAt(0).toUpperCase()}
                                      </Avatar>
                                      <Typography variant="body2">
                                        {interviewer}
                                      </Typography>
                                    </Box>
                                  )
                                )}
                              </Box>
                            </TableCell>
                          </TableRow>
                        </TableBody>
                      </Table>
                    </Box>

                    {/* Description section */}
                    <Box sx={{ mb: 2 }}>
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ mb: 1, display: "block" }}
                      >
                        {t("interviews.description")}
                      </Typography>
                      <Typography sx={{ whiteSpace: "pre-wrap" }}>
                        {interview.description}
                      </Typography>
                    </Box>

                    {/* Actions */}
                    <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
                      <Button
                        size="small"
                        variant="outlined"
                        onClick={() =>
                          router.push(`/interviews/${interview.interview_id}`)
                        }
                      >
                        {t("interviews.manage")}
                      </Button>
                    </Box>
                  </Collapse>
                </Paper>
              ))
            ) : (
              <Paper sx={{ p: 3, textAlign: "center" }}>
                <Typography color="text.secondary">
                  {t("interviews.noInterviews")}
                </Typography>
              </Paper>
            )}
          </Box>
        </Collapse>
      </Paper>

      {/* Add Interview Dialog */}
      <Dialog
        open={openAddInterview}
        onClose={() => setOpenAddInterview(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>{t("interviews.addNew")}</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 2 }}>
            <FormControl fullWidth>
              <InputLabel>{t("interviews.type")}</InputLabel>
              <Select
                value={newInterview.type}
                label={t("interviews.type")}
                onChange={(e) =>
                  setNewInterview({
                    ...newInterview,
                    type: e.target.value as InterviewType,
                  })
                }
              >
                <MenuItem value="VIDEO_CALL">
                  {t("interviews.types.VIDEO_CALL")}
                </MenuItem>
                <MenuItem value="IN_PERSON">
                  {t("interviews.types.IN_PERSON")}
                </MenuItem>
                <MenuItem value="TAKE_HOME">
                  {t("interviews.types.TAKE_HOME")}
                </MenuItem>
              </Select>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel>{t("interviews.timezone")}</InputLabel>
              <Select
                value={newInterview.timezone}
                label={t("interviews.timezone")}
                onChange={(e) =>
                  setNewInterview({
                    ...newInterview,
                    timezone: e.target.value as TimeZone,
                  })
                }
              >
                {Array.from(validTimezones).map((tz) => (
                  <MenuItem key={tz} value={tz}>
                    {tz}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <FormControlLabel
              control={
                <Switch
                  checked={allowPastDates}
                  onChange={(e) => setAllowPastDates(e.target.checked)}
                />
              }
              label={t("interviews.allowPastDates")}
            />

            <FormControlLabel
              control={
                <Switch
                  checked={use24HourFormat}
                  onChange={(e) => setUse24HourFormat(e.target.checked)}
                />
              }
              label={t("interviews.use24HourFormat")}
            />

            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DateTimePicker
                label={t("interviews.startTime")}
                value={
                  newInterview.startTime
                    ? new Date(newInterview.startTime)
                    : null
                }
                onChange={(newValue: Date | null) => {
                  if (newValue) {
                    const startDate = newValue;
                    // Set end time to 1 hour after start time
                    const endDate = new Date(startDate);
                    endDate.setHours(startDate.getHours() + 1);

                    setNewInterview({
                      ...newInterview,
                      startTime: startDate.toISOString(),
                      endTime: endDate.toISOString(),
                    });
                  }
                }}
                views={["year", "month", "day", "hours", "minutes"]}
                ampm={!use24HourFormat}
                format={
                  use24HourFormat
                    ? "MMMM dd, yyyy HH:mm"
                    : "MMMM dd, yyyy hh:mm a"
                }
                minDateTime={allowPastDates ? undefined : new Date()}
                slotProps={{
                  textField: {
                    fullWidth: true,
                  },
                }}
              />

              <DateTimePicker
                label={t("interviews.endTime")}
                value={
                  newInterview.endTime ? new Date(newInterview.endTime) : null
                }
                onChange={(newValue: Date | null) => {
                  if (newValue) {
                    // Only update if end time is after start time
                    if (
                      newInterview.startTime &&
                      new Date(newValue) <= new Date(newInterview.startTime)
                    ) {
                      setError(t("interviews.endTimeBeforeStart"));
                      return;
                    }
                    setNewInterview({
                      ...newInterview,
                      endTime: newValue.toISOString(),
                    });
                  }
                }}
                views={["year", "month", "day", "hours", "minutes"]}
                ampm={!use24HourFormat}
                format={
                  use24HourFormat
                    ? "MMMM dd, yyyy HH:mm"
                    : "MMMM dd, yyyy hh:mm a"
                }
                minDateTime={
                  newInterview.startTime
                    ? new Date(newInterview.startTime)
                    : undefined
                }
                slotProps={{
                  textField: {
                    fullWidth: true,
                  },
                }}
              />
            </LocalizationProvider>

            <TextField
              label={t("interviews.description")}
              multiline
              rows={4}
              value={newInterview.description}
              onChange={(e) =>
                setNewInterview({
                  ...newInterview,
                  description: e.target.value,
                })
              }
              fullWidth
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAddInterview(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={handleAddInterview}
            variant="contained"
            disabled={
              !newInterview.startTime ||
              !newInterview.endTime ||
              !newInterview.description
            }
          >
            {t("common.add")}
          </Button>
        </DialogActions>
      </Dialog>

      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          {t("comments.title")}
        </Typography>

        {comments.length > 0 ? (
          <Box sx={{ mt: 3 }}>
            {[...comments]
              .sort(
                (a, b) =>
                  new Date(a.created_at).getTime() -
                  new Date(b.created_at).getTime()
              )
              .map((comment) => (
                <Box
                  key={comment.comment_id}
                  sx={{
                    display: "flex",
                    gap: 2,
                    mb: 3,
                    flexDirection:
                      comment.commenter_type === "ORG_USER"
                        ? "row"
                        : "row-reverse",
                  }}
                >
                  <Avatar
                    sx={{
                      width: 40,
                      height: 40,
                      bgcolor: (theme) =>
                        comment.commenter_type === "ORG_USER"
                          ? theme.palette.primary.main
                          : theme.palette.grey[400],
                    }}
                  >
                    {comment.commenter_name.charAt(0).toUpperCase()}
                  </Avatar>
                  <Box sx={{ flexGrow: 1 }}>
                    <Paper
                      sx={{
                        p: 2,
                        borderRadius: 2,
                        border: "1px solid",
                        borderColor: "divider",
                        position: "relative",
                        "&::before": {
                          content: '""',
                          position: "absolute",
                          ...(comment.commenter_type === "ORG_USER"
                            ? {
                                left: -8,
                                borderRight: (theme) =>
                                  `8px solid ${theme.palette.divider}`,
                              }
                            : {
                                right: -8,
                                borderLeft: (theme) =>
                                  `8px solid ${theme.palette.divider}`,
                              }),
                          top: 16,
                          width: 0,
                          height: 0,
                          borderTop: "8px solid transparent",
                          borderBottom: "8px solid transparent",
                        },
                      }}
                    >
                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                          mb: 1,
                        }}
                      >
                        <Typography
                          variant="subtitle2"
                          sx={{
                            fontWeight: "bold",
                            color: (theme) =>
                              comment.commenter_type === "ORG_USER"
                                ? theme.palette.primary.main
                                : theme.palette.text.primary,
                          }}
                        >
                          {comment.commenter_name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          {new Date(comment.created_at).toLocaleDateString(
                            undefined,
                            {
                              year: "numeric",
                              month: "short",
                              day: "2-digit",
                            }
                          )}{" "}
                          {new Date(comment.created_at).toLocaleTimeString(
                            undefined,
                            {
                              hour: "2-digit",
                              minute: "2-digit",
                            }
                          )}
                        </Typography>
                      </Box>
                      <Typography
                        sx={{
                          whiteSpace: "pre-wrap",
                          wordBreak: "break-word",
                        }}
                      >
                        {comment.content}
                      </Typography>
                    </Paper>
                  </Box>
                </Box>
              ))}
          </Box>
        ) : (
          <Typography color="text.secondary" sx={{ my: 2 }}>
            {t("comments.noComments")}
          </Typography>
        )}

        <Divider sx={{ my: 3 }} />

        <Box>
          <TextField
            fullWidth
            multiline
            rows={4}
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder={t("comments.addPlaceholder")}
            disabled={submitting}
          />
          <Box sx={{ display: "flex", justifyContent: "flex-end", mt: 2 }}>
            <Button
              variant="contained"
              onClick={handleAddComment}
              disabled={!newComment.trim() || submitting}
            >
              {submitting ? t("common.loading") : t("comments.add")}
            </Button>
          </Box>
        </Box>
      </Paper>
    </Box>
  );
}
