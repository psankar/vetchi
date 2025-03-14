"use client";

import { config } from "@/config";
import { FeatureFlags } from "@/config/features";
import { useTranslation } from "@/hooks/useTranslation";
import Alert from "@mui/material/Alert";
import Autocomplete from "@mui/material/Autocomplete";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Chip from "@mui/material/Chip";
import Container from "@mui/material/Container";
import FormControl from "@mui/material/FormControl";
import FormControlLabel from "@mui/material/FormControlLabel";
import InputLabel from "@mui/material/InputLabel";
import MenuItem from "@mui/material/MenuItem";
import Paper from "@mui/material/Paper";
import Select from "@mui/material/Select";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  CreateOpeningRequest,
  EducationLevel,
  EducationLevels,
  GlobalCountryCode,
  Location,
  LocationStates,
  OpeningTag,
  OpeningTagID,
  OpeningType,
  OpeningTypes,
  OrgUserShort,
  validTimezones,
} from "@psankar/vetchi-typespec";
import countries from "@psankar/vetchi-typespec/common/countries.json";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import CircularProgress from "@mui/material/CircularProgress";

export default function CreateOpeningPage() {
  const { t } = useTranslation();
  const router = useRouter();

  // Form state
  const [title, setTitle] = useState("");
  const [positions, setPositions] = useState(1);
  const [jobDescription, setJobDescription] = useState("");
  const [employerNotes, setEmployerNotes] = useState("");
  const [selectedType, setSelectedType] = useState<OpeningType>(
    OpeningTypes.FULL_TIME
  );
  const [yoeMin, setYoeMin] = useState(0);
  const [yoeMax, setYoeMax] = useState(1);
  const [selectedEducation, setSelectedEducation] = useState<EducationLevel>(
    EducationLevels.UNSPECIFIED
  );
  const [selectedRecruiter, setSelectedRecruiter] =
    useState<OrgUserShort | null>(null);
  const [selectedHiringManager, setSelectedHiringManager] =
    useState<OrgUserShort | null>(null);
  const [selectedCostCenter, setSelectedCostCenter] = useState("");

  // Location state
  const [selectedLocations, setSelectedLocations] = useState<string[]>([]);
  const [selectedTimezones, setSelectedTimezones] = useState<string[]>([]);
  const [selectedCountries, setSelectedCountries] = useState<string[]>([]);
  const [isGloballyRemote, setIsGloballyRemote] = useState(false);

  // Loading and error state
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Data state
  const [orgUsers, setOrgUsers] = useState<OrgUserShort[]>([]);
  const [costCenters, setCostCenters] = useState<string[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [isLoadingTags, setIsLoadingTags] = useState(false);
  const [selectedTags, setSelectedTags] = useState<OpeningTag[]>([]);
  const [availableTags, setAvailableTags] = useState<OpeningTag[]>([]);
  const [isInitialDataLoaded, setIsInitialDataLoaded] = useState(false);

  // Initialize from sessionStorage on mount
  useEffect(() => {
    const saved = sessionStorage.getItem("isGloballyRemote");
    if (saved !== null) {
      setIsGloballyRemote(JSON.parse(saved));
    }
  }, []);

  // Persist isGloballyRemote state to sessionStorage
  useEffect(() => {
    sessionStorage.setItem(
      "isGloballyRemote",
      JSON.stringify(isGloballyRemote)
    );
  }, [isGloballyRemote]);

  const fetchOrgUsers = useCallback(
    async (searchPrefix?: string) => {
      try {
        const token = Cookies.get("session_token");
        if (!token) {
          router.push("/signin");
          return;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/employer/filter-org-users`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              prefix: searchPrefix,
              limit: 40,
            }),
          }
        );

        if (response.status === 401) {
          Cookies.remove("session_token");
          router.push("/signin");
          return;
        }

        if (!response.ok) {
          throw new Error(t("openings.fetchError"));
        }

        const data = await response.json();
        setOrgUsers(data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : t("openings.fetchError"));
      }
    },
    [router, t]
  );

  const fetchCostCenters = useCallback(async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-cost-centers`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({}),
        }
      );

      if (!response.ok) {
        throw new Error(t("openings.fetchCostCentersError"));
      }

      const data = await response.json();
      setCostCenters(data.map((cc: { name: string }) => cc.name) || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("openings.fetchCostCentersError")
      );
    }
  }, [router, t]);

  const fetchLocations = useCallback(async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-locations`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({}),
        }
      );

      if (!response.ok) {
        throw new Error(t("openings.fetchLocationsError"));
      }

      const data = await response.json();
      setLocations(data || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("openings.fetchLocationsError")
      );
    }
  }, [router, t]);

  const fetchTags = useCallback(
    async (searchPrefix?: string) => {
      try {
        setIsLoadingTags(true);
        const token = Cookies.get("session_token");
        if (!token) {
          router.push("/signin");
          return;
        }

        const response = await fetch(
          `${config.API_SERVER_PREFIX}/employer/filter-opening-tags`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              prefix: searchPrefix,
            }),
          }
        );

        if (!response.ok) {
          throw new Error(t("openings.fetchTagsError"));
        }

        const data = await response.json();
        setAvailableTags(data || []);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : t("openings.fetchTagsError")
        );
      } finally {
        setIsLoadingTags(false);
      }
    },
    [router, t]
  );

  useEffect(() => {
    if (!isInitialDataLoaded) {
      const loadInitialData = async () => {
        await Promise.all([
          fetchCostCenters(),
          fetchLocations(),
          fetchOrgUsers(),
        ]);
        setIsInitialDataLoaded(true);
      };
      loadInitialData();
    }
  }, [fetchCostCenters, fetchLocations, fetchOrgUsers, isInitialDataLoaded]);

  const handleTagChange = (
    event: React.SyntheticEvent,
    newValue: (OpeningTag | string)[]
  ) => {
    const processedTags = newValue.map((tag): OpeningTag => {
      if (typeof tag === "string") {
        // This is a new tag
        return {
          id: "" as OpeningTagID,
          name: tag.trim(),
        };
      }
      return tag;
    });

    if (processedTags.length > 3) {
      setError(t("openings.maxTagsError"));
      return;
    }

    setSelectedTags(processedTags);
  };

  const handleSave = async () => {
    try {
      setIsLoading(true);
      setError(null);

      // Validate required fields
      if (!title || title.length < 3 || title.length > 32) {
        setError(t("validation.title.lengthError"));
        return;
      }

      if (positions < 1 || positions > 20) {
        setError(t("validation.positions.range.1.20"));
        return;
      }

      if (
        !jobDescription ||
        jobDescription.length < 10 ||
        jobDescription.length > 1024
      ) {
        setError(t("validation.jobDescription.lengthError"));
        return;
      }

      if (employerNotes && employerNotes.length > 1024) {
        setError(t("validation.employerNotes.maxLength.1024"));
        return;
      }

      if (!selectedRecruiter || !selectedHiringManager) {
        setError(t("validation.roles.required"));
        return;
      }

      // Location validation
      if (
        !isGloballyRemote &&
        selectedLocations.length === 0 &&
        selectedTimezones.length === 0 &&
        selectedCountries.length === 0
      ) {
        setError(t("openings.locationRequiredError"));
        return;
      }

      if (selectedTags.length === 0) {
        setError(t("openings.tagsRequired"));
        return;
      }

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      // Separate existing tags and new tags
      const existingTags = selectedTags.filter((tag) => tag.id !== "");
      const newTags = selectedTags
        .filter((tag) => tag.id === "")
        .map((tag) => tag.name);

      const request: CreateOpeningRequest = {
        title,
        positions,
        jd: jobDescription,
        employer_notes: employerNotes || undefined,
        opening_type: selectedType,
        yoe_min: yoeMin,
        yoe_max: yoeMax,
        min_education_level: selectedEducation,
        recruiter: selectedRecruiter?.email || "",
        hiring_manager: selectedHiringManager?.email || "",
        cost_center_name: selectedCostCenter,
        location_titles: selectedLocations,
        remote_timezones:
          selectedTimezones.length > 0 ? selectedTimezones : undefined,
        remote_country_codes: isGloballyRemote
          ? [GlobalCountryCode]
          : selectedCountries.length > 0
          ? selectedCountries
          : [],
        tags: existingTags.map((tag) => tag.id),
        new_tags: newTags.length > 0 ? newTags : undefined,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/create-opening`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || t("openings.createError"));
      }

      router.push("/openings");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("openings.createError"));
    } finally {
      setIsLoading(false);
    }
  };

  // Add handlers for numeric inputs
  const handlePositionsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value, 10);
    setPositions(isNaN(value) ? 1 : Math.max(1, Math.min(20, value)));
  };

  const handleYoeMinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value, 10);
    setYoeMin(isNaN(value) ? 0 : Math.max(0, Math.min(100, value)));
  };

  const handleYoeMaxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value, 10);
    setYoeMax(isNaN(value) ? 1 : Math.max(1, Math.min(100, value)));
  };

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            {t("openings.createTitle")}
          </Typography>
          {error && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {error}
            </Alert>
          )}
        </Box>

        <Box component="form" noValidate sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            label={t("openings.openingTitle")}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            inputProps={{ minLength: 3, maxLength: 32 }}
            error={title.length > 0 && (title.length < 3 || title.length > 32)}
            helperText={
              title.length > 0 && (title.length < 3 || title.length > 32)
                ? t("validation.title.lengthError")
                : ""
            }
          />

          <TextField
            margin="normal"
            required
            fullWidth
            type="number"
            label={t("openings.positions")}
            value={positions}
            onChange={handlePositionsChange}
            inputProps={{ min: 1, max: 20 }}
            error={positions < 1 || positions > 20}
            helperText={
              positions < 1 || positions > 20
                ? t("validation.positions.range.1.20")
                : ""
            }
          />

          <TextField
            margin="normal"
            required
            fullWidth
            multiline
            rows={4}
            label={t("openings.jobDescription")}
            value={jobDescription}
            onChange={(e) => setJobDescription(e.target.value)}
            inputProps={{ minLength: 10, maxLength: 1024 }}
            error={
              jobDescription.length > 0 &&
              (jobDescription.length < 10 || jobDescription.length > 1024)
            }
            helperText={
              jobDescription.length > 0 &&
              (jobDescription.length < 10 || jobDescription.length > 1024)
                ? t("validation.jobDescription.lengthError")
                : ""
            }
          />

          <Autocomplete
            options={orgUsers}
            getOptionLabel={(option) => `${option.name} (${option.email})`}
            value={selectedRecruiter}
            onChange={(_, newValue) => setSelectedRecruiter(newValue)}
            onInputChange={(_, value) => fetchOrgUsers(value)}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.recruiter")}
              />
            )}
          />

          <Autocomplete
            options={orgUsers}
            getOptionLabel={(option) => `${option.name} (${option.email})`}
            value={selectedHiringManager}
            onChange={(_, newValue) => setSelectedHiringManager(newValue)}
            onInputChange={(_, value) => fetchOrgUsers(value)}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.hiringManager")}
              />
            )}
          />

          <Autocomplete
            options={costCenters}
            value={selectedCostCenter}
            onChange={(_, newValue) => setSelectedCostCenter(newValue || "")}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.costCenter")}
              />
            )}
          />

          <FormControl fullWidth margin="normal" required>
            <InputLabel>{t("openings.type")}</InputLabel>
            <Select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value as OpeningType)}
              label={t("openings.type")}
            >
              {Object.values(OpeningTypes).map((type) => (
                <MenuItem key={type} value={type}>
                  {t(`openings.types.${type}`)}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              type="number"
              label={t("openings.minYoe")}
              value={yoeMin}
              onChange={handleYoeMinChange}
              inputProps={{ min: 0, max: 100 }}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              type="number"
              label={t("openings.maxYoe")}
              value={yoeMax}
              onChange={handleYoeMaxChange}
              inputProps={{ min: 1, max: 100 }}
            />
          </Box>

          <FormControl fullWidth margin="normal">
            <InputLabel>{t("openings.minEducation")}</InputLabel>
            <Select
              value={selectedEducation}
              onChange={(e) =>
                setSelectedEducation(e.target.value as EducationLevel)
              }
              label={t("openings.minEducation")}
            >
              {Object.values(EducationLevels).map((level) => (
                <MenuItem key={level} value={level}>
                  {t(`openings.education.${level}`)}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControlLabel
            control={
              <Switch
                checked={isGloballyRemote}
                onChange={(e) => {
                  setIsGloballyRemote(e.target.checked);
                  if (e.target.checked) {
                    setSelectedCountries([]);
                    setSelectedTimezones([]);
                  }
                }}
                color="primary"
              />
            }
            label={t("openings.globallyRemote")}
            sx={{ mt: 2, mb: 1, display: "block" }}
          />

          <Autocomplete
            multiple
            options={countries}
            getOptionLabel={(option) => option.en}
            value={countries.filter((c) =>
              selectedCountries.includes(c.country_code)
            )}
            onChange={(_, newValue) =>
              setSelectedCountries(newValue.map((c) => c.country_code))
            }
            disabled={isGloballyRemote}
            renderInput={(params) => (
              <TextField
                {...params}
                margin="normal"
                label={t("openings.remoteCountries")}
                helperText={t("openings.remoteCountriesHelp")}
              />
            )}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip
                  label={option.en}
                  {...getTagProps({ index })}
                  key={option.country_code}
                />
              ))
            }
          />
          {FeatureFlags.TimezonesForCreateOpening && (
            <Autocomplete
              multiple
              options={Array.from(validTimezones)}
              value={selectedTimezones}
              onChange={(_, newValue) => setSelectedTimezones(newValue)}
              disabled={isGloballyRemote}
              renderInput={(params) => (
                <TextField
                  {...params}
                  margin="normal"
                  label={t("openings.remoteTimezones")}
                  helperText={t("openings.remoteTimezonesHelp")}
                />
              )}
              renderTags={(value, getTagProps) =>
                value.map((option, index) => (
                  <Chip
                    label={option}
                    {...getTagProps({ index })}
                    key={option}
                  />
                ))
              }
            />
          )}
          <Autocomplete
            multiple
            options={locations
              .filter((loc) => loc.state === LocationStates.ACTIVE)
              .map((loc) => loc.title)}
            value={selectedLocations}
            onChange={(_, newValue) => setSelectedLocations(newValue)}
            disabled={isGloballyRemote}
            renderInput={(params) => (
              <TextField
                {...params}
                margin="normal"
                label={t("openings.officeLocations")}
                helperText={t("openings.officeLocationsHelp")}
              />
            )}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip label={option} {...getTagProps({ index })} key={option} />
              ))
            }
          />

          <TextField
            margin="normal"
            fullWidth
            multiline
            rows={4}
            label={t("openings.employerNotes")}
            value={employerNotes}
            onChange={(e) => setEmployerNotes(e.target.value)}
            inputProps={{ maxLength: 1024 }}
            error={employerNotes.length > 1024}
            helperText={
              employerNotes.length > 1024
                ? t("validation.employerNotes.maxLength.1024")
                : ""
            }
          />

          <Autocomplete
            multiple
            options={availableTags}
            getOptionLabel={(option: OpeningTag | string) => {
              if (typeof option === "string") {
                return option;
              }
              return option.name;
            }}
            value={selectedTags}
            onChange={handleTagChange}
            onInputChange={(_, value) => {
              fetchTags(value);
            }}
            loading={isLoadingTags}
            freeSolo
            filterSelectedOptions
            renderOption={(props, option) => {
              // Extract key from props
              const { key, ...otherProps } = props;
              if (typeof option === "string") {
                return (
                  <li key={key} {...otherProps}>
                    {t("openings.createNewTag")}: "{option}"
                  </li>
                );
              }
              return (
                <li key={key} {...otherProps}>
                  {option.name}
                </li>
              );
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                required
                margin="normal"
                label={t("openings.tags")}
                placeholder={
                  selectedTags.length < 3
                    ? t("openings.selectOrCreateTag")
                    : t("openings.maxTagsReached")
                }
                helperText={t("openings.tagsHelp")}
                InputProps={{
                  ...params.InputProps,
                  endAdornment: (
                    <>
                      {isLoadingTags ? (
                        <CircularProgress color="inherit" size={20} />
                      ) : null}
                      {params.InputProps.endAdornment}
                    </>
                  ),
                }}
              />
            )}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip
                  label={typeof option === "string" ? option : option.name}
                  {...getTagProps({ index })}
                  key={
                    typeof option === "string"
                      ? option
                      : option.id || option.name
                  }
                />
              ))
            }
            filterOptions={(options, params) => {
              const filtered = options.filter((option) =>
                typeof option === "string"
                  ? option
                      .toLowerCase()
                      .includes(params.inputValue.toLowerCase())
                  : option.name
                      .toLowerCase()
                      .includes(params.inputValue.toLowerCase())
              );

              // Add option to create a new tag if input value exists and doesn't match exactly
              const inputValue = params.inputValue.trim();
              if (
                inputValue !== "" &&
                !options.some((option) =>
                  typeof option === "string"
                    ? option.toLowerCase() === inputValue.toLowerCase()
                    : option.name.toLowerCase() === inputValue.toLowerCase()
                )
              ) {
                filtered.push(inputValue);
              }

              return filtered;
            }}
            noOptionsText={
              selectedTags.length >= 3
                ? t("openings.maxTagsReached")
                : t("openings.createNewTag")
            }
          />

          <Box sx={{ mt: 4, display: "flex", gap: 2 }}>
            <Button variant="outlined" onClick={() => router.push("/openings")}>
              {t("common.cancel")}
            </Button>
            <Button
              variant="contained"
              onClick={handleSave}
              disabled={
                isLoading ||
                !title ||
                title.length < 3 ||
                title.length > 32 ||
                !positions ||
                positions < 1 ||
                positions > 20 ||
                !jobDescription ||
                jobDescription.length < 10 ||
                jobDescription.length > 1024 ||
                !selectedRecruiter ||
                !selectedHiringManager ||
                !selectedCostCenter ||
                !selectedType ||
                yoeMin < 0 ||
                yoeMax <= yoeMin ||
                yoeMin > 100 ||
                yoeMax > 100 ||
                (!isGloballyRemote &&
                  selectedLocations.length === 0 &&
                  selectedTimezones.length === 0 &&
                  selectedCountries.length === 0) ||
                selectedTags.length === 0
              }
            >
              {isLoading ? t("common.loading") : t("common.save")}
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}
