export type InterviewState = 'SCHEDULED' | 'COMPLETED' | 'CANCELLED';

export const InterviewStates = {
    SCHEDULED: 'SCHEDULED' as InterviewState,
    COMPLETED: 'COMPLETED' as InterviewState,
    CANCELLED: 'CANCELLED' as InterviewState,
} as const;

export type CandidacyState = 'INTERVIEWING' | 'OFFERED' | 'OFFER_DECLINED' | 'OFFER_ACCEPTED' | 'CANDIDATE_UNSUITABLE' | 'CANDIDATE_NOT_RESPONDING' | 'CANDIDATE_WITHDREW' | 'EMPLOYER_DEFUNCT';

export const CandidacyStates = {
    INTERVIEWING: 'INTERVIEWING' as CandidacyState,
    OFFERED: 'OFFERED' as CandidacyState,
    OFFER_DECLINED: 'OFFER_DECLINED' as CandidacyState,
    OFFER_ACCEPTED: 'OFFER_ACCEPTED' as CandidacyState,
    CANDIDATE_UNSUITABLE: 'CANDIDATE_UNSUITABLE' as CandidacyState,
    CANDIDATE_NOT_RESPONDING: 'CANDIDATE_NOT_RESPONDING' as CandidacyState,
    CANDIDATE_WITHDREW: 'CANDIDATE_WITHDREW' as CandidacyState,
    EMPLOYER_DEFUNCT: 'EMPLOYER_DEFUNCT' as CandidacyState,
} as const;

export type InterviewersDecision = 'STRONG_YES' | 'YES' | 'NEUTRAL' | 'NO' | 'STRONG_NO';

export const InterviewersDecisions = {
    STRONG_YES: 'STRONG_YES' as InterviewersDecision,
    YES: 'YES' as InterviewersDecision,
    NEUTRAL: 'NEUTRAL' as InterviewersDecision,
    NO: 'NO' as InterviewersDecision,
    STRONG_NO: 'STRONG_NO' as InterviewersDecision,
} as const;

export type RSVPStatus = 'YES' | 'NO' | 'NOT_SET';

export const RSVPStatuses = {
    YES: 'YES' as RSVPStatus,
    NO: 'NO' as RSVPStatus,
    NOT_SET: 'NOT_SET' as RSVPStatus,
} as const;

export interface RSVPInterviewRequest {
    interview_id: string;
    rsvp_status: RSVPStatus;
}

export function isValidInterviewState(state: string): state is InterviewState {
    return Object.values(InterviewStates).includes(state as InterviewState);
}

export function isValidCandidacyState(state: string): state is CandidacyState {
    return Object.values(CandidacyStates).includes(state as CandidacyState);
}

export function isValidInterviewersDecision(decision: string): decision is InterviewersDecision {
    return Object.values(InterviewersDecisions).includes(decision as InterviewersDecision);
}

export function isValidRSVPStatus(status: string): status is RSVPStatus {
    return Object.values(RSVPStatuses).includes(status as RSVPStatus);
} 