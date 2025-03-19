package granger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (g *Granger) scoreApplications(quit <-chan struct{}) {
	defer g.wg.Done()

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			g.log.Dbg("Resume scoring job received quit signal")
			return
		case <-ticker.C:
			g.log.Dbg("Starting resume scoring job")
			if err := g.processApplicationsForScoring(context.Background()); err != nil {
				g.log.Err(
					"Failed to process applications for scoring",
					"error",
					err,
				)
			}
		}
	}
}

func (g *Granger) processApplicationsForScoring(ctx context.Context) error {
	// Get all active scoring models
	models, err := g.db.GetActiveApplicationScoringModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active scoring models: %w", err)
	}
	g.log.Dbg("Got active scoring models", "count", len(models))

	// Get openings with unscored applications in APPLIED state
	openings, err := g.db.GetOpeningsWithUnscoredApplications(ctx)
	if err != nil {
		return err
	}
	g.log.Dbg("Got openings with unscored applications", "count", len(openings))

	// Process each opening
	for _, opening := range openings {
		// Get job description for this opening
		jd, err := g.db.GetOpeningJD(ctx, opening.EmployerID, opening.ID)
		if err != nil {
			g.log.Dbg("failed to get job description", "err", err)
			continue
		}
		g.log.Dbg("Got job description", "jd_length", len(jd))

		// Get unscore applications for this opening (max 10 at a time)
		applications, err := g.db.GetUnscoredApplicationsForOpening(
			ctx,
			opening.EmployerID,
			opening.ID,
			vetchi.MaxApplicationsToScorePerBatch,
		)
		if err != nil {
			g.log.Dbg("failed to get unscored applications", "err", err)
			continue
		}
		g.log.Dbg("Got unscored applications", "count", len(applications))

		if len(applications) == 0 {
			continue
		}

		// Process this batch of applications
		err = g.scoreApplicationBatch(ctx, applications, jd, models)
		if err != nil {
			g.log.Dbg("failed to score application batch", "err", err)
			continue
		}
	}

	return nil
}

func (g *Granger) scoreApplicationBatch(
	ctx context.Context,
	applications []db.ApplicationForScoring,
	jd string,
	models []db.ApplicationScoringModel,
) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Get S3 bucket name from environment variable
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		g.log.Err("S3_BUCKET environment variable not set")
		return fmt.Errorf("S3_BUCKET environment variable not set")
	}

	// Collect all scores for all applications to save in a single transaction
	var allScores []db.ApplicationScore

	for _, app := range applications {
		// Format fileurl as expected by sortinghat: s3://bucket/key
		fileurl := fmt.Sprintf("s3://%s/%s", bucket, app.ResumeSHA)
		g.log.Dbg("Scoring resume", "application", app.ID, "fileurl", fileurl)

		// Build request URL with query parameters
		apiURL := "http://sortinghat:8080/score-resumes-jd"
		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			g.log.Err("failed to create request", "err", err)
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Add query parameters
		q := req.URL.Query()
		q.Add("fileurl", fileurl)
		q.Add("job_description", jd)
		req.URL.RawQuery = q.Encode()

		// Execute request
		resp, err := client.Do(req)
		if err != nil {
			g.log.Err("failed to call sortinghat API", "err", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			g.log.Err(
				"sortinghat API",
				"application",
				app.ID,
				"status",
				resp.Status,
			)
			continue
		}

		// Parse response
		var scoreResponse []resumeScoreResponse
		if err := json.NewDecoder(resp.Body).Decode(&scoreResponse); err != nil {
			g.log.Err("failed to decode sortinghat response", "err", err)
			continue
		}

		if len(scoreResponse) == 0 {
			g.log.Err("empty response from sortinghat", "application", app.ID)
			continue
		}

		scores := scoreResponse[0].CompatibilityScores
		g.log.Dbg("Got scores", "application", app.ID, "scores", scores)

		// Process scores for this application
		for _, model := range models {
			// Map sortinghat model names to our model names
			// We have to do this because the model_names in the response don't match our database exactly
			var scoreKey string
			if model.ModelName == "sentence-transformers-all-MiniLM-L6-v2" {
				scoreKey = "sentence-transformers"
			} else if model.ModelName == "tfidf-1.0" {
				scoreKey = "tfidf"
			} else {
				g.log.Err("unknown model name", "model_name", model.ModelName)
				continue
			}

			score, ok := scores[scoreKey]
			if !ok {
				g.log.Dbg("no score for model", "model", model.ModelName)
				continue
			}

			// Convert to integer score
			intScore := int(score)

			allScores = append(allScores, db.ApplicationScore{
				ApplicationID: app.ID,
				ModelName:     model.ModelName,
				Score:         intScore,
			})
		}
	}

	// Save all scores for all applications in a single transaction
	if len(allScores) > 0 {
		err := g.db.SaveApplicationScores(ctx, allScores)
		if err != nil {
			g.log.Dbg("failed to save application scores", "err", err)
			return fmt.Errorf("failed to save application scores: %w", err)
		}
		g.log.Dbg("Saved all application scores", "count", len(allScores))
	} else {
		g.log.Dbg("No scores to save")
	}

	return nil
}
