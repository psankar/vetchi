export const en = {
  common: {
    home: "Home",
    openings: "Openings",
    logout: "Logout",
    loading: "Loading...",
    error: "Error",
    retry: "Retry",
    costCenters: "Cost Centers",
    locations: "Locations",
    actions: "Actions",
    add: "Add",
    save: "Save",
    cancel: "Cancel",
    loadMore: "Load More",
    serverError: "Please try again after some time.",
    none: "None",
    back: "Back",
  },
  auth: {
    signin: "Sign In",
    domain: "Domain",
    email: "Email",
    password: "Password",
    tfa: "Two Factor Authentication",
    tfaCode: "Enter TFA Code",
    verify: "Verify",
    submit: "Submit",
    rememberMe: "Remember me",
    domainNotVerified:
      "Please add CNAME records to verify your domain with Vetchi.",
    domainVerifyPending:
      "Please ask your domain admin to check their email and complete the onboarding process.",
    accountDisabled: "Your account has been disabled.",
    serverError: "Please try again after some time.",
    invalidCredentials: "Invalid credentials.",
    unauthorized: "Unauthorized access",
  },
  dashboard: {
    welcome: "Welcome to your dashboard",
  },
  openings: {
    title: "Job Openings",
    noOpenings: "No openings found",
    create: "Create Opening",
    createTitle: "Create New Opening",
    openingTitle: "Title",
    positions: "Number of Positions",
    jobDescription: "Job Description",
    recruiter: "Recruiter",
    hiringManager: "Hiring Manager",
    costCenter: "Cost Center",
    type: "Opening Type",
    stateLabel: "State",
    minYoe: "Minimum Years of Experience",
    maxYoe: "Maximum Years of Experience",
    minEducation: "Minimum Education Level",
    employerNotes: "Employer Notes",
    remoteTimezones: "Remote Timezones",
    remoteTimezonesHelp:
      "Select the timezones where remote work is allowed. Leave empty if remote work is not allowed.",
    officeLocations: "Office Locations",
    officeLocationsHelp:
      "Select the office locations where this position is available. Leave empty if the position is fully remote.",
    showClosed: "Show Closed Openings",
    types: {
      FULL_TIME_OPENING: "Full Time",
      PART_TIME_OPENING: "Part Time",
      CONTRACT_OPENING: "Contract",
      INTERNSHIP_OPENING: "Internship",
      UNSPECIFIED_OPENING: "Unspecified",
    },
    education: {
      BACHELOR_EDUCATION: "Bachelor's Degree",
      MASTER_EDUCATION: "Master's Degree",
      DOCTORATE_EDUCATION: "Doctorate",
      NOT_MATTERS_EDUCATION: "Not Required",
      UNSPECIFIED_EDUCATION: "Unspecified",
    },
    state: {
      DRAFT_OPENING_STATE: "Draft",
      ACTIVE_OPENING_STATE: "Active",
      SUSPENDED_OPENING_STATE: "Suspended",
      CLOSED_OPENING_STATE: "Closed",
    },
    details: "Opening Details",
    id: "Opening ID",
    filledPositions: "Filled Positions",
    description: "Description",
    contacts: "Contact Information",
    actions: "Actions",
    publish: "Publish Opening",
    suspend: "Suspend Opening",
    reactivate: "Reactivate Opening",
    viewCandidacies: "View Candidacies",
    viewInterviews: "View Interviews",
    invalidStateTransition: "Invalid state transition",
    notFound: "Opening not found",
    fetchError: "Failed to fetch openings",
    createError: "Failed to create opening",
    fetchCostCentersError: "Failed to fetch cost centers",
    fetchLocationsError: "Failed to fetch locations",
    missingUserError: "Please select both recruiter and hiring manager",
    close: "Close Opening",
    closeConfirmTitle: "Close Opening",
    closeConfirmMessage:
      "Are you sure you want to close this opening? This action cannot be undone.",
    confirmClose: "Yes, Close Opening",
    stateChangeSuccess: "Opening state updated successfully",
  },
  costCenters: {
    title: "Cost Centers",
    addTitle: "Add Cost Center",
    editTitle: "Edit Cost Center",
    name: "Name",
    notes: "Notes",
    state: "State",
    add: "Add Cost Center",
    active: "Active",
    defunct: "Defunct",
    includeDefunct: "Show Defunct Cost Centers",
    noCostCenters:
      "No cost centers found. Click 'Add Cost Center' to create one.",
    fetchError: "Failed to fetch cost centers",
    addError: "Failed to add cost center",
    updateError: "Failed to update cost center",
    defunctError: "Failed to defunct cost center",
  },
  locations: {
    title: "Locations",
    addTitle: "Add Location",
    editTitle: "Edit Location",
    locationTitle: "Title",
    countryCode: "Country Code",
    countryCodeHelp: "Enter 3-letter ISO country code (e.g., USA, IND, GBR)",
    postalAddress: "Postal Address",
    postalCode: "Postal Code",
    mapUrl: "OpenStreetMap URL",
    cityAka: "Alternative City Names",
    cityAkaPlaceholder: "Enter alternative city name",
    state: "State",
    active: "Active",
    defunct: "Defunct",
    add: "Add Location",
    fetchError: "Failed to fetch locations",
    addError: "Failed to add location",
    updateError: "Failed to update location",
    defunctError: "Failed to defunct location",
    noLocations: "No locations found. Click 'Add Location' to create one.",
    includeDefunct: "Include Defunct Locations",
    viewMap: "View on Map",
  },
  "validation.title.length.3.32": "Title must be between 3 and 32 characters",
  "validation.positions.range.1.20":
    "Number of positions must be between 1 and 20",
  "validation.jobDescription.length.10.1024":
    "Job description must be between 10 and 1024 characters",
  "validation.employerNotes.maxLength.1024":
    "Employer notes must not exceed 1024 characters",
};
