package hermione

import (
	"fmt"
	"net/http"

	app "github.com/psankar/vetchi/api/internal/hermione/applications"
	"github.com/psankar/vetchi/api/internal/hermione/candidacy"
	"github.com/psankar/vetchi/api/internal/hermione/colleagues"
	"github.com/psankar/vetchi/api/internal/hermione/costcenter"
	ea "github.com/psankar/vetchi/api/internal/hermione/employerauth"
	ha "github.com/psankar/vetchi/api/internal/hermione/hubauth"
	he "github.com/psankar/vetchi/api/internal/hermione/hubemp"
	ho "github.com/psankar/vetchi/api/internal/hermione/hubopenings"
	"github.com/psankar/vetchi/api/internal/hermione/interview"
	"github.com/psankar/vetchi/api/internal/hermione/locations"
	"github.com/psankar/vetchi/api/internal/hermione/openings"
	"github.com/psankar/vetchi/api/internal/hermione/orgusers"
	pp "github.com/psankar/vetchi/api/internal/hermione/profilepage"
	wh "github.com/psankar/vetchi/api/internal/hermione/workhistory"
	"github.com/psankar/vetchi/typespec/common"
)

func (h *Hermione) Run() error {
	// Authentication related endpoints
	http.HandleFunc("/employer/get-onboard-status", ea.GetOnboardStatus(h))
	http.HandleFunc("/employer/set-onboard-password", ea.SetOnboardPassword(h))
	http.HandleFunc("/employer/signin", ea.EmployerSignin(h))
	http.HandleFunc("/employer/tfa", ea.EmployerTFA(h))

	// CostCenter related endpoints
	h.mw.Protect(
		"/employer/add-cost-center",
		costcenter.AddCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-centers",
		costcenter.GetCostCenters(h),
		[]common.OrgUserRole{
			common.Admin,
			common.CostCentersCRUD,
			common.CostCentersViewer,
		},
	)
	h.mw.Protect(
		"/employer/defunct-cost-center",
		costcenter.DefunctCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/rename-cost-center",
		costcenter.RenameCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/update-cost-center",
		costcenter.UpdateCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-center",
		costcenter.GetCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersViewer},
	)

	// Location related endpoints
	h.mw.Protect(
		"/employer/add-location",
		locations.AddLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/defunct-location",
		locations.DefunctLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/get-locations",
		locations.GetLocations(h),
		[]common.OrgUserRole{
			common.Admin,
			common.LocationsCRUD,
			common.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/get-location",
		locations.GetLocation(h),
		[]common.OrgUserRole{
			common.Admin,
			common.LocationsCRUD,
			common.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/rename-location",
		locations.RenameLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/update-location",
		locations.UpdateLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)

	// OrgUser related endpoints
	h.mw.Protect(
		"/employer/add-org-user",
		orgusers.AddOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/update-org-user",
		orgusers.UpdateOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/disable-org-user",
		orgusers.DisableOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/enable-org-user",
		orgusers.EnableOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/filter-org-users",
		orgusers.FilterOrgUsers(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OrgUsersCRUD,
			common.OrgUsersViewer,
		},
	)
	http.HandleFunc("/employer/signup-org-user", orgusers.SignupOrgUser(h))

	// Opening related endpoints
	h.mw.Protect(
		"/employer/create-opening",
		openings.CreateOpening(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/get-opening",
		openings.GetOpening(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/filter-openings",
		openings.FilterOpenings(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/update-opening",
		openings.UpdateOpening(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/get-opening-watchers",
		openings.GetOpeningWatchers(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/add-opening-watchers",
		openings.AddOpeningWatchers(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/remove-opening-watcher",
		openings.RemoveOpeningWatcher(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/change-opening-state",
		openings.ChangeOpeningState(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)

	// Opening tags related endpoints
	h.mw.Protect(
		"/employer/filter-opening-tags",
		he.FilterOpeningTags(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)

	// Application related endpoints
	h.mw.Protect(
		"/employer/get-applications",
		app.GetApplications(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
		},
	)

	h.mw.Protect(
		"/employer/get-resume",
		app.GetResume(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
		},
	)

	h.mw.Protect(
		"/employer/set-application-color-tag",
		app.SetApplicationColorTag(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	h.mw.Protect(
		"/employer/remove-application-color-tag",
		app.RemoveApplicationColorTag(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	h.mw.Protect(
		"/employer/shortlist-application",
		app.ShortlistApplication(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	h.mw.Protect(
		"/employer/reject-application",
		app.RejectApplication(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	// Used by employer - Candidacies
	h.mw.Protect(
		"/employer/add-candidacy-comment",
		candidacy.EmployerAddComment(h),
		[]common.OrgUserRole{common.Any},
	)

	h.mw.Protect(
		"/employer/get-candidacy-comments",
		candidacy.EmployerGetComments(h),
		[]common.OrgUserRole{common.Any},
	)

	h.mw.Protect(
		"/employer/filter-candidacy-infos",
		candidacy.FilterCandidacyInfos(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.Any},
	)

	h.mw.Protect(
		"/employer/get-candidacy-info",
		candidacy.GetEmployerCandidacyInfo(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.Any},
	)
	h.mw.Protect(
		"/employer/offer-to-candidate",
		candidacy.OfferToCandidate(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	// Used by employer - Interviews
	h.mw.Protect(
		"/employer/add-interview",
		interview.AddInterview(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)
	h.mw.Protect(
		"/employer/add-interviewer",
		interview.AddInterviewer(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	h.mw.Protect(
		"/employer/remove-interviewer",
		interview.RemoveInterviewer(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	h.mw.Protect(
		"/employer/rsvp-interview",
		interview.EmployerRSVPInterview(h),
		[]common.OrgUserRole{common.Any},
	)
	h.mw.Protect(
		"/employer/get-interviews-by-opening",
		interview.GetInterviewsByOpening(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.Any},
	)
	h.mw.Protect(
		"/employer/get-interviews-by-candidacy",
		interview.GetEmployerInterviewsByCandidacy(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.Any},
	)
	h.mw.Protect(
		"/employer/get-assessment",
		interview.EmployerGetAssessment(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.Any},
	)
	h.mw.Protect(
		"/employer/get-interview-details",
		interview.GetInterviewDetails(h),
		[]common.OrgUserRole{common.Any},
	)
	h.mw.Protect(
		"/employer/put-assessment",
		interview.EmployerPutAssessment(h),
		[]common.OrgUserRole{common.Any},
	)

	wrap := func(fn http.Handler) http.Handler {
		return h.mw.HubWrap(fn)
	}
	// Hub related endpoints
	http.HandleFunc("/hub/login", ha.Login(h))
	http.HandleFunc("/hub/tfa", ha.HubTFA(h))
	http.Handle("/hub/get-my-handle", wrap(ha.GetMyHandle(h)))
	http.HandleFunc("/hub/logout", ha.Logout(h))

	http.HandleFunc("/hub/forgot-password", ha.ForgotPassword(h))
	http.HandleFunc("/hub/reset-password", ha.ResetPassword(h))
	http.Handle("/hub/change-password", wrap(ha.ChangePassword(h)))

	// Official Email related endpoints
	http.Handle("/hub/add-official-email", wrap(pp.AddOfficialEmail(h)))
	http.Handle(
		"/hub/verify-official-email",
		wrap(pp.VerifyOfficialEmail(h)),
	)
	http.Handle("/hub/trigger-verification", wrap(pp.TriggerVerification(h)))
	http.Handle(
		"/hub/delete-official-email",
		wrap(pp.DeleteOfficialEmail(h)),
	)
	http.Handle("/hub/my-official-emails", wrap(pp.MyOfficialEmails(h)))

	// ProfilePage related endpoints
	http.Handle("/hub/get-bio", wrap(pp.GetBio(h)))
	http.Handle("/hub/update-bio", wrap(pp.UpdateBio(h)))
	http.Handle("/hub/upload-profile-picture", wrap(pp.UploadProfilePicture(h)))
	http.Handle("/hub/profile-picture/", wrap(pp.GetProfilePicture(h)))
	http.Handle("/hub/remove-profile-picture", wrap(pp.RemoveProfilePicture(h)))

	http.Handle("/hub/find-openings", wrap(ho.FindHubOpenings(h)))
	http.Handle("/hub/filter-opening-tags", wrap(he.FilterOpeningTags(h)))
	http.Handle("/hub/get-opening-details", wrap(ho.GetOpeningDetails(h)))
	http.Handle("/hub/apply-for-opening", wrap(ho.ApplyForOpening(h)))
	http.Handle("/hub/my-applications", wrap(app.MyApplications(h)))
	http.Handle("/hub/withdraw-application", wrap(app.WithdrawApplication(h)))
	http.Handle("/hub/add-candidacy-comment", wrap(candidacy.HubAddComment(h)))
	http.Handle(
		"/hub/get-candidacy-comments",
		wrap(candidacy.HubGetComments(h)),
	)
	http.Handle("/hub/get-my-candidacies", wrap(candidacy.MyCandidacies(h)))
	http.Handle(
		"/hub/get-candidacy-info",
		wrap(candidacy.GetHubCandidacyInfo(h)),
	)
	http.Handle(
		"/hub/get-interviews-by-candidacy",
		wrap(interview.GetHubInterviewsByCandidacy(h)),
	)
	http.Handle("/hub/rsvp-interview", wrap(interview.HubRSVPInterview(h)))
	http.Handle("/hub/filter-employers", wrap(he.FilterEmployers(h)))

	// WorkHistory related endpoints
	http.Handle("/hub/add-work-history", wrap(wh.AddWorkHistory(h)))
	http.Handle("/hub/delete-work-history", wrap(wh.DeleteWorkHistory(h)))
	http.Handle("/hub/list-work-history", wrap(wh.ListWorkHistory(h)))
	http.Handle("/hub/update-work-history", wrap(wh.UpdateWorkHistory(h)))

	// Colleague related endpoints
	http.Handle("/hub/connect-colleague", wrap(colleagues.ConnectColleague(h)))
	http.Handle("/hub/unlink-colleague", wrap(colleagues.UnlinkColleague(h)))
	http.Handle(
		"/hub/my-colleague-approvals",
		wrap(colleagues.MyColleagueApprovals(h)),
	)
	http.Handle("/hub/my-colleague-seeks", wrap(colleagues.MyColleagueSeeks(h)))
	http.Handle("/hub/approve-colleague", wrap(colleagues.ApproveColleague(h)))
	http.Handle("/hub/reject-colleague", wrap(colleagues.RejectColleague(h)))

	port := fmt.Sprintf(":%d", h.Config().Port)
	return http.ListenAndServe(port, nil)
}
