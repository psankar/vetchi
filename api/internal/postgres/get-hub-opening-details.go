package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetHubOpeningDetails(
	ctx context.Context,
	req hub.GetHubOpeningDetailsRequest,
) (hub.HubOpeningDetails, error) {
	query := `
	SELECT
		o.id as opening_id_within_company,
		d.domain_name as company_domain,
		e.company_name as company_name,
		o.title as job_title,
		o.jd as jd,
		o.opening_type,
		o.yoe_min,
		o.yoe_max,
		o.min_education_level,
		o.salary_min,
		o.salary_max,
		o.salary_currency,
		o.created_at,
		o.pagination_key,
		o.state,
		hm.name as hiring_manager_name,
		hu_hm.handle as hiring_manager_vetchi_handle,
		r.name as recruiter_name
	FROM openings o
		JOIN employers e ON o.employer_id = e.id
		JOIN domains d ON d.employer_id = e.id
		LEFT JOIN org_users hm ON o.hiring_manager = hm.id
		LEFT JOIN hub_users_official_emails hue_hm ON hm.email = hue_hm.official_email
		LEFT JOIN hub_users hu_hm ON hue_hm.hub_user_id = hu_hm.id
		LEFT JOIN org_users r ON o.recruiter = r.id
	WHERE o.id = $1
		AND d.domain_name = $2
`

	var details hub.HubOpeningDetails
	var salaryMin, salaryMax *float64
	var salaryCurrency *string
	var hiringManagerHandle *string

	err := p.pool.QueryRow(
		ctx,
		query,
		req.OpeningIDWithinCompany,
		req.CompanyDomain,
	).Scan(
		&details.OpeningIDWithinCompany,
		&details.CompanyDomain,
		&details.CompanyName,
		&details.JobTitle,
		&details.JD,
		&details.OpeningType,
		&details.YoeMin,
		&details.YoeMax,
		&details.EducationLevel,
		&salaryMin,
		&salaryMax,
		&salaryCurrency,
		&details.CreatedAt,
		&details.PaginationKey,
		&details.State,
		&details.HiringManagerName,
		&hiringManagerHandle,
		&details.RecruiterName,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return hub.HubOpeningDetails{}, db.ErrNoOpening
		}
		p.log.Err("failed to get opening details", "error", err)
		return hub.HubOpeningDetails{}, db.ErrInternal
	}

	// Build salary if all components are present
	if salaryMin != nil && salaryMax != nil && salaryCurrency != nil {
		// TODO: We should raise an error if one or two conditions only fail
		details.Salary = &common.Salary{
			MinAmount: *salaryMin,
			MaxAmount: *salaryMax,
			Currency:  common.Currency(*salaryCurrency),
		}
	}

	// Set hiring manager handle if present
	if hiringManagerHandle != nil {
		details.HiringManagerVetchiHandle = hiringManagerHandle
	}

	return details, nil
}
