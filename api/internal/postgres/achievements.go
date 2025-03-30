package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
)

func (pg *PG) AddAchievement(
	ctx context.Context,
	addAchievementReq hub.AddAchievementRequest,
) (string, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return "", err
	}

	var id uuid.UUID
	err = pg.pool.QueryRow(ctx, `
INSERT INTO achievements (hub_user_id, title, description, url, at)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    id
`, hubUserID, addAchievementReq.Title, addAchievementReq.Description, addAchievementReq.URL, addAchievementReq.At).Scan(&id)
	if err != nil {
		pg.log.Err("failed to add achievement", "error", err)
		return "", err
	}

	pg.log.Dbg("achievement added", "id", id)
	return id.String(), nil
}

func (pg *PG) ListAchievements(
	ctx context.Context,
	listAchievementsReq hub.ListAchievementsRequest,
) ([]common.Achievement, error) {
	// If a handle is provided, check if it exists
	if listAchievementsReq.Handle != nil {
		var exists bool
		err := pg.pool.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM hub_users WHERE handle = $1)",
			string(*listAchievementsReq.Handle),
		).Scan(&exists)
		if err != nil {
			pg.log.Err("failed to check if handle exists", "error", err)
			return nil, db.ErrInternal
		}

		if !exists {
			pg.log.Dbg("notfound",
				"handle", string(*listAchievementsReq.Handle),
			)
			return nil, db.ErrNoHubUser
		}
	}

	var queryParams []interface{}
	var queryConditions []string
	var paramCounter int

	// Build the base query
	baseQuery := `
SELECT id, title, description, url, at, achievement_type as type
FROM achievements
WHERE `

	// Handle user selection: if handle is provided, use it; otherwise, use logged-in user
	if listAchievementsReq.Handle != nil {
		paramCounter++
		queryConditions = append(
			queryConditions,
			fmt.Sprintf(
				"hub_user_id = (SELECT id FROM hub_users WHERE handle = $%d)",
				paramCounter,
			),
		)
		queryParams = append(queryParams, string(*listAchievementsReq.Handle))
	} else {
		// Get logged-in user ID
		hubUserID, err := getHubUserID(ctx)
		if err != nil {
			pg.log.Err("failed to get hub user ID", "error", err)
			return nil, err
		}
		paramCounter++
		queryConditions = append(queryConditions, fmt.Sprintf("hub_user_id = $%d", paramCounter))
		queryParams = append(queryParams, hubUserID)
	}

	// Filter by type if provided
	if listAchievementsReq.Type != "" {
		paramCounter++
		queryConditions = append(
			queryConditions,
			fmt.Sprintf("achievement_type = $%d", paramCounter),
		)
		queryParams = append(queryParams, listAchievementsReq.Type)
	}

	// Combine the conditions
	query := baseQuery + strings.Join(queryConditions, " AND ")

	// Execute the query
	rows, err := pg.pool.Query(ctx, query, queryParams...)
	if err != nil {
		pg.log.Err("failed to list achievements", "error", err)
		return nil, err
	}
	defer rows.Close()

	achievements := []common.Achievement{}
	for rows.Next() {
		var achievement common.Achievement
		err = rows.Scan(
			&achievement.ID,
			&achievement.Title,
			&achievement.Description,
			&achievement.URL,
			&achievement.At,
			&achievement.Type,
		)
		if err != nil {
			pg.log.Err("failed to scan achievement", "error", err)
			return nil, db.ErrInternal
		}
		achievements = append(achievements, achievement)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error during rows iteration", "error", err)
		return nil, db.ErrInternal
	}

	pg.log.Dbg("achievements listed", "count", len(achievements))
	return achievements, nil
}

func (pg *PG) DeleteAchievement(
	ctx context.Context,
	deleteAchievementReq hub.DeleteAchievementRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	res, err := pg.pool.Exec(
		ctx,
		`DELETE FROM achievements WHERE id = $1::uuid AND hub_user_id = $2`,
		deleteAchievementReq.ID,
		hubUserID,
	)
	if err != nil {
		pg.log.Err("failed to delete achievement", "error", err)
		return db.ErrInternal
	}

	if res.RowsAffected() == 0 {
		pg.log.Dbg("achievement not found for deletion",
			"achievement_id", deleteAchievementReq.ID,
			"hub_user_id", hubUserID)
		return db.ErrNoAchievement
	}

	pg.log.Dbg("achievement deleted successfully",
		"achievement_id", deleteAchievementReq.ID,
		"hub_user_id", hubUserID)
	return nil
}

func (pg *PG) ListHubUserAchievements(
	ctx context.Context,
	listHubUserAchievementsReq employer.ListHubUserAchievementsRequest,
) ([]common.Achievement, error) {
	var exists bool
	err := pg.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM hub_users WHERE handle = $1)",
		listHubUserAchievementsReq.Handle,
	).Scan(&exists)
	if err != nil {
		pg.log.Err("failed to check if handle exists", "error", err)
		return nil, db.ErrInternal
	}

	if !exists {
		pg.log.Dbg("notfound",
			"handle", listHubUserAchievementsReq.Handle,
		)
		return nil, db.ErrNoHubUser
	}

	var queryParams []interface{}
	queryParams = append(queryParams, listHubUserAchievementsReq.Handle)

	achievementsQuery := `
SELECT id, title, description, url, at, achievement_type as type
FROM achievements
WHERE hub_user_id = (SELECT id FROM hub_users WHERE handle = $1)
`
	// If type is provided, add filter condition
	if listHubUserAchievementsReq.Type != "" {
		achievementsQuery += `AND achievement_type = $2`
		queryParams = append(queryParams, listHubUserAchievementsReq.Type)
	}

	rows, err := pg.pool.Query(ctx, achievementsQuery, queryParams...)
	if err != nil {
		pg.log.Err("failed to list hub user achievements", "error", err)
		return nil, err
	}
	defer rows.Close()

	achievements := []common.Achievement{}
	for rows.Next() {
		var achievement common.Achievement
		err = rows.Scan(
			&achievement.ID,
			&achievement.Title,
			&achievement.Description,
			&achievement.URL,
			&achievement.At,
			&achievement.Type,
		)
		if err != nil {
			pg.log.Err("failed to scan achievement", "error", err)
			return nil, db.ErrInternal
		}
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}
