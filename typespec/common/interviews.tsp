union InterviewState {
    Scheduled: "SCHEDULED",
    Completed: "COMPLETED",
    CandidateWithdrew: "CANDIDATE_WITHDREW",
    EmployerWithdrew: "EMPLOYER_WITHDREW",
}

union CandidacyState {
    @doc("The candidate is being interviewed")
    Interviewing: "INTERVIEWING",

    @doc("The candidate has been offered a job")
    Offered: "OFFERED",

    @doc("The candidate has declined the offer")
    OfferDeclined: "OFFER_DECLINED",

    @doc("The candidate has accepted the offer")
    OfferAccepted: "OFFER_ACCEPTED",

    @doc("The Employer has deemed the candidate unsuitable for the job")
    CandidateUnsuitable: "CANDIDATE_UNSUITABLE",

    @doc("The candidate is not responding")
    CandidateNotResponding: "CANDIDATE_NOT_RESPONDING",

    @doc("The Employer is defunct")
    EmployerDefunct: "EMPLOYER_DEFUNCT",
}

union InterviewersDecision {
    StrongYes: "STRONG_YES",
    Yes: "YES",
    No: "NO",
    StrongNo: "STRONG_NO",
}
