package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Hub Login", Ordered, func() {
	var db *pgxpool.Pool

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0007-hub-login-up.pgsql")
	})

	AfterAll(func() {
		seedDatabase(db, "0007-hub-login-down.pgsql")
		db.Close()
	})

	Describe("Hub Login Flow", func() {
		type loginTestCase struct {
			description   string
			request       vetchi.LoginRequest
			wantStatus    int
			wantErrFields []string
		}

		It("should handle login requests correctly", func() {
			testCases := []loginTestCase{
				{
					description: "valid credentials for active user",
					request: vetchi.LoginRequest{
						Email:    "active@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "invalid password for active user",
					request: vetchi.LoginRequest{
						Email:    "active@hub.example",
						Password: "WrongPassword123$",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "disabled user",
					request: vetchi.LoginRequest{
						Email:    "disabled@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "deleted user",
					request: vetchi.LoginRequest{
						Email:    "deleted@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "non-existent user",
					request: vetchi.LoginRequest{
						Email:    "nonexistent@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusInternalServerError,
				},
				{
					description: "invalid email format",
					request: vetchi.LoginRequest{
						Email:    "invalid-email",
						Password: "NewPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"email"},
				},
				{
					description: "empty password",
					request: vetchi.LoginRequest{
						Email:    "active@hub.example",
						Password: "",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"password"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				loginReqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				resp, err := http.Post(
					serverURL+"/hub/login",
					"application/json",
					bytes.NewBuffer(loginReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if len(tc.wantErrFields) > 0 {
					var validationErrors vetchi.ValidationErrors
					err = json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
					continue
				}

				if tc.wantStatus == http.StatusOK {
					var loginResp vetchi.LoginResponse
					err = json.NewDecoder(resp.Body).Decode(&loginResp)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(loginResp.Token).ShouldNot(BeEmpty())
				}
			}
		})

		type tfaTestCase struct {
			description   string
			request       vetchi.HubTFARequest
			wantStatus    int
			wantErrFields []string
		}

		It("should handle TFA flow correctly", func() {
			// First get a valid TFA token through login
			loginReqBody, err := json.Marshal(vetchi.LoginRequest{
				Email:    "active@hub.example",
				Password: "NewPassword123$",
			})
			Expect(err).ShouldNot(HaveOccurred())

			loginResp, err := http.Post(
				serverURL+"/hub/login",
				"application/json",
				bytes.NewBuffer(loginReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(loginResp.StatusCode).Should(Equal(http.StatusOK))

			var loginRespObj vetchi.LoginResponse
			err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			tfaToken := loginRespObj.Token

			// Get the TFA code from mailpit
			var messageID string
			for i := 0; i < 3; i++ {
				<-time.After(10 * time.Second)

				mailPitResp, err := http.Get(
					fmt.Sprintf(
						"%s/api/v1/search?query=to:active@hub.example subject:Vetchi Two Factor Authentication",
						mailPitURL,
					),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(mailPitResp.StatusCode).Should(Equal(http.StatusOK))

				body, err := io.ReadAll(mailPitResp.Body)
				Expect(err).ShouldNot(HaveOccurred())

				var mailPitRespObj MailPitResponse
				err = json.Unmarshal(body, &mailPitRespObj)
				Expect(err).ShouldNot(HaveOccurred())

				if len(mailPitRespObj.Messages) > 0 {
					messageID = mailPitRespObj.Messages[0].ID
					break
				}
			}
			Expect(messageID).ShouldNot(BeEmpty())

			// Get the email content
			mailResp, err := http.Get(
				mailPitURL + "/api/v1/message/" + messageID,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailResp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(mailResp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			re := regexp.MustCompile(`Code:\s*([0-9]+)`)
			matches := re.FindStringSubmatch(string(body))
			Expect(len(matches)).Should(BeNumerically(">=", 2))
			tfaCode := matches[1]

			testCases := []tfaTestCase{
				{
					description: "valid TFA token and code",
					request: vetchi.HubTFARequest{
						TFAToken:   tfaToken,
						TFACode:    tfaCode,
						RememberMe: false,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "invalid TFA token",
					request: vetchi.HubTFARequest{
						TFAToken:   "invalid-token",
						TFACode:    tfaCode,
						RememberMe: false,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid TFA code",
					request: vetchi.HubTFARequest{
						TFAToken:   tfaToken,
						TFACode:    "000000",
						RememberMe: false,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "empty TFA token",
					request: vetchi.HubTFARequest{
						TFAToken:   "",
						TFACode:    tfaCode,
						RememberMe: false,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"tfa_token"},
				},
				{
					description: "empty TFA code",
					request: vetchi.HubTFARequest{
						TFAToken:   tfaToken,
						TFACode:    "",
						RememberMe: false,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"tfa_code"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				tfaReqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				resp, err := http.Post(
					serverURL+"/hub/tfa",
					"application/json",
					bytes.NewBuffer(tfaReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if len(tc.wantErrFields) > 0 {
					var validationErrors vetchi.ValidationErrors
					err = json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
					continue
				}

				if tc.wantStatus == http.StatusOK {
					var tfaResp vetchi.HubTFAResponse
					err = json.NewDecoder(resp.Body).Decode(&tfaResp)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(tfaResp.SessionToken).ShouldNot(BeEmpty())
				}
			}

			// Clean up the email
			deleteReqBody, err := json.Marshal(MailPitDeleteRequest{
				IDs: []string{messageID},
			})
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"DELETE",
				mailPitURL+"/api/v1/messages",
				bytes.NewBuffer(deleteReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")

			deleteResp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(deleteResp.StatusCode).Should(Equal(http.StatusOK))
		})

		It("should handle remember me flag correctly", func() {
			// First get a valid TFA token through login
			loginReqBody, err := json.Marshal(vetchi.LoginRequest{
				Email:    "active@hub.example",
				Password: "NewPassword123$",
			})
			Expect(err).ShouldNot(HaveOccurred())

			loginResp, err := http.Post(
				serverURL+"/hub/login",
				"application/json",
				bytes.NewBuffer(loginReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(loginResp.StatusCode).Should(Equal(http.StatusOK))

			var loginRespObj vetchi.LoginResponse
			err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			tfaToken := loginRespObj.Token

			// Get the TFA code from mailpit
			var messageID string
			for i := 0; i < 3; i++ {
				<-time.After(10 * time.Second)

				mailPitResp, err := http.Get(
					fmt.Sprintf(
						"%s/api/v1/search?query=to:active@hub.example subject:Vetchi Two Factor Authentication",
						mailPitURL,
					),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(mailPitResp.StatusCode).Should(Equal(http.StatusOK))

				body, err := io.ReadAll(mailPitResp.Body)
				Expect(err).ShouldNot(HaveOccurred())

				var mailPitRespObj MailPitResponse
				err = json.Unmarshal(body, &mailPitRespObj)
				Expect(err).ShouldNot(HaveOccurred())

				if len(mailPitRespObj.Messages) > 0 {
					messageID = mailPitRespObj.Messages[0].ID
					break
				}
			}
			Expect(messageID).ShouldNot(BeEmpty())

			// Get the email content
			mailResp, err := http.Get(
				mailPitURL + "/api/v1/message/" + messageID,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailResp.StatusCode).Should(Equal(http.StatusOK))

			body, err := io.ReadAll(mailResp.Body)
			Expect(err).ShouldNot(HaveOccurred())

			re := regexp.MustCompile(`Code:\s*([0-9]+)`)
			matches := re.FindStringSubmatch(string(body))
			Expect(len(matches)).Should(BeNumerically(">=", 2))
			tfaCode := matches[1]

			// Test with remember_me flag
			tfaReqBody, err := json.Marshal(vetchi.HubTFARequest{
				TFAToken:   tfaToken,
				TFACode:    tfaCode,
				RememberMe: true,
			})
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.Post(
				serverURL+"/hub/tfa",
				"application/json",
				bytes.NewBuffer(tfaReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var tfaResp vetchi.HubTFAResponse
			err = json.NewDecoder(resp.Body).Decode(&tfaResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tfaResp.SessionToken).ShouldNot(BeEmpty())

			// Clean up the email
			deleteReqBody, err := json.Marshal(MailPitDeleteRequest{
				IDs: []string{messageID},
			})
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"DELETE",
				mailPitURL+"/api/v1/messages",
				bytes.NewBuffer(deleteReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")

			deleteResp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(deleteResp.StatusCode).Should(Equal(http.StatusOK))
		})
	})
})
