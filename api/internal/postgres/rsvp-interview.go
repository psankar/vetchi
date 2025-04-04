package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (p *PG) HubRSVPInterview(
	ctx context.Context,
	req hub.HubRSVPInterviewRequest,
) error {
	const (
		interviewNotFound     = "interview_not_found"
		interviewInvalidState = "interview_invalid_state"
		success               = "success"
	)

	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return db.ErrInternal
	}

	query := `
WITH interview_hub_user AS (
    SELECT i.id as interview_id, a.hub_user_id
    FROM interviews i
    JOIN candidacies c ON i.candidacy_id = c.id
    JOIN applications a ON c.application_id = a.id
),
updated_interview AS (
    UPDATE interviews i
    SET candidate_rsvp = $1
    FROM interview_hub_user ihu
    WHERE i.id = $2
      AND ihu.interview_id = i.id
      AND ihu.hub_user_id = $3
      AND i.interview_state = $4
    RETURNING i.*
)
SELECT
    CASE
        WHEN NOT EXISTS (
            SELECT 1 FROM interview_hub_user 
            WHERE interview_id = $2 AND hub_user_id = $3
        ) THEN $5
        WHEN NOT EXISTS (
            SELECT 1 FROM interviews i
            JOIN interview_hub_user ihu ON i.id = ihu.interview_id
            WHERE i.id = $2 
            AND ihu.hub_user_id = $3 
            AND i.interview_state = $4
        ) THEN $6
        WHEN EXISTS (SELECT 1 FROM updated_interview) THEN $7
    END AS result
`

	var result string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.RSVPStatus,
		req.InterviewID,
		hubUser.ID,
		common.ScheduledInterviewState,
		interviewNotFound,
		interviewInvalidState,
		success,
	).Scan(&result)
	if err != nil {
		p.log.Err("failed to update interview rsvp status", "error", err)
		return db.ErrInternal
	}

	switch result {
	case interviewNotFound:
		p.log.Dbg("interview not found", "interview_id", req.InterviewID)
		return db.ErrNoInterview
	case interviewInvalidState:
		p.log.Dbg("interview invalid state", "interview_id", req.InterviewID)
		return db.ErrInvalidInterviewState
	default:
		p.log.Dbg(
			"rsvp status updated",
			"interview_id",
			req.InterviewID,
			"hub_user_id",
			hubUser.ID,
			"rsvp",
			req.RSVPStatus,
		)
		return nil
	}
}

func (p *PG) EmployerRSVPInterview(
	ctx context.Context,
	req common.RSVPInterviewRequest,
) error {
	const (
		interviewNotFound     = "interview_not_found"
		interviewInvalidState = "interview_invalid_state"
		success               = "success"
	)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
WITH updated_interviewer AS (
    UPDATE interview_interviewers
    SET rsvp_status = $1
    WHERE interview_id = $2
      AND interviewer_id = $3
      AND EXISTS (
        SELECT 1 FROM interviews i
        WHERE i.id = interview_interviewers.interview_id
        AND i.interview_state = $4
      )
    RETURNING *
)
SELECT
    CASE
        WHEN NOT EXISTS (SELECT 1 FROM interviews WHERE id = $2) THEN $5
        WHEN NOT EXISTS (SELECT 1 FROM interviews WHERE id = $2 AND interview_state = $4) THEN $6
        WHEN EXISTS (SELECT 1 FROM updated_interviewer) THEN $7
        ELSE $5
    END AS result
`

	var result string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.RSVPStatus,
		req.InterviewID,
		orgUser.ID,
		common.ScheduledInterviewState,
		interviewNotFound,
		interviewInvalidState,
		success,
	).Scan(&result)
	if err != nil {
		p.log.Err("failed to update interview rsvp status", "error", err)
		return db.ErrInternal
	}

	switch result {
	case interviewNotFound:
		p.log.Dbg("interview not found", "interview_id", req.InterviewID)
		return db.ErrNoInterview
	case interviewInvalidState:
		p.log.Dbg("interview invalid state", "interview_id", req.InterviewID)
		return db.ErrInvalidInterviewState
	default:
		p.log.Dbg(
			"rsvp status updated",
			"interview_id",
			req.InterviewID,
			"org_user_id",
			orgUser.ID,
			"rsvp",
			req.RSVPStatus,
		)
		return nil
	}
}
