package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
)

func (p *PG) GetEmployerCandidacyComments(
	ctx context.Context,
	empGetCommentsReq common.GetCandidacyCommentsRequest,
) ([]common.CandidacyComment, error) {
	candidacyComments := make([]common.CandidacyComment, 0)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return candidacyComments, db.ErrInternal
	}

	query := `
WITH access_check AS (
	SELECT EXISTS (
		SELECT 1 FROM candidacies 
		WHERE id = $1 AND employer_id = $2
	) as has_access
)
SELECT 
	cc.id,
	COALESCE(ou.name, hu.full_name) as commenter_name,
	cc.author_type,
	cc.comment_text as content,
	cc.created_at
FROM candidacy_comments cc
LEFT JOIN org_users ou ON cc.org_user_id = ou.id
LEFT JOIN hub_users hu ON cc.hub_user_id = hu.id,
access_check
WHERE cc.candidacy_id = $1
AND access_check.has_access
ORDER BY cc.created_at DESC
`

	rows, err := p.pool.Query(ctx, query,
		empGetCommentsReq.CandidacyID,
		orgUser.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to query candidacy comments", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var comment common.CandidacyComment
		err := rows.Scan(
			&comment.CommentID,
			&comment.CommenterName,
			&comment.CommenterType,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			p.log.Err("failed to scan candidacy comment", "error", err)
			return nil, db.ErrInternal
		}

		candidacyComments = append(candidacyComments, comment)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating over rows", "error", err)
		return nil, db.ErrInternal
	}

	return candidacyComments, nil
}

func (p *PG) GetHubCandidacyComments(
	ctx context.Context,
	hubGetCommentsReq common.GetCandidacyCommentsRequest,
) ([]common.CandidacyComment, error) {
	candidacyComments := make([]common.CandidacyComment, 0)

	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return candidacyComments, db.ErrInternal
	}

	query := `
WITH access_check AS (
	SELECT EXISTS (
		SELECT 1 FROM candidacies c
		JOIN applications a ON c.application_id = a.id
		WHERE c.id = $1 AND a.hub_user_id = $2
	) as has_access
)
SELECT 
	cc.id,
	COALESCE(ou.name, hu.full_name) as commenter_name,
	cc.author_type,
	cc.comment_text as content,
	cc.created_at
FROM candidacy_comments cc
LEFT JOIN org_users ou ON cc.org_user_id = ou.id
LEFT JOIN hub_users hu ON cc.hub_user_id = hu.id,
access_check
WHERE cc.candidacy_id = $1
AND access_check.has_access
ORDER BY cc.created_at DESC
`

	rows, err := p.pool.Query(ctx, query,
		hubGetCommentsReq.CandidacyID,
		hubUser.ID,
	)
	if err != nil {
		p.log.Err("failed to query candidacy comments", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var comment common.CandidacyComment
		err := rows.Scan(
			&comment.CommentID,
			&comment.CommenterName,
			&comment.CommenterType,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			p.log.Err("failed to scan candidacy comment", "error", err)
			return nil, db.ErrInternal
		}

		candidacyComments = append(candidacyComments, comment)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating over rows", "error", err)
		return nil, db.ErrInternal
	}

	p.log.Dbg("got candidacy comments", "comments", candidacyComments)
	return candidacyComments, nil
}
