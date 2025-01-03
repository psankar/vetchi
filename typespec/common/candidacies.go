package common

import "time"

type GetCandidacyCommentsRequest struct {
	CandidacyID string `json:"candidacy_id"`
}

type CandidacyComment struct {
	CommentID     string    `json:"comment_id"`
	CommenterName string    `json:"commenter_name"`
	CommenterType string    `json:"commenter_type"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}
