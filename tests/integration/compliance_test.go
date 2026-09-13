package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/goretante/clearmytravel-api/internal/compliance"
	"github.com/goretante/clearmytravel-api/internal/db"
	apphttp "github.com/goretante/clearmytravel-api/internal/http"
	dbpkg "github.com/jackc/pgx/v5/pgxpool"
)

type complianceResponse struct {
	Status       string `json:"status"`
	Requirements []struct {
		Type    string `json:"type"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"requirements"`
}

func setupIntegrationTest(t *testing.T) (*dbpkg.Pool, http.Handler) {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := dbpkg.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create db pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping database: %v", err)
	}

	queries := db.New(pool)

	// Visa
	visaRepository := compliance.NewDBVisaRuleRepository(queries)
	visaEvaluator := compliance.NewRepositoryVisaEvaluator(
		visaRepository,
	)

	// ETA
	etaRepository := compliance.NewDBETARuleRepository(queries)
	etaEvaluator := compliance.NewRepositoryETAEvaluator(
		etaRepository,
	)

	// Stay duration
	stayRepository := compliance.NewDBStayRuleRepository(queries)
	stayEvaluator := compliance.NewRepositoryStayDurationEvaluator(
		stayRepository,
	)

	// Compliance service
	complianceService := compliance.NewService(
		compliance.PassportValidityEvaluator{
			BufferDays: 6,
		},
		visaEvaluator,
		etaEvaluator,
		stayEvaluator,
	)

	// HTTP router
	router := apphttp.NewRouter(
		queries,
		complianceService,
	)

	return pool, router
}

func checkCompliance(
	t *testing.T,
	router http.Handler,
	body string,
) complianceResponse {
	t.Helper()

	req := httptest.NewRequest(
		http.MethodPost,
		"/compliance/check",
		bytes.NewBufferString(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	var response complianceResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf(
			"decode response: %v",
			err,
		)
	}

	return response
}

func assertOverallStatus(
	t *testing.T,
	response complianceResponse,
	expected string,
) {
	t.Helper()

	if response.Status != expected {
		t.Fatalf(
			"expected overall status %s, got %s",
			expected,
			response.Status,
		)
	}
}

func assertRequirementStatus(
	t *testing.T,
	response complianceResponse,
	requirementType string,
	expectedStatus string,
) {
	t.Helper()

	for _, requirement := range response.Requirements {
		if requirement.Type != requirementType {
			continue
		}

		if requirement.Status != expectedStatus {
			t.Fatalf(
				"requirement %s: expected %s, got %s",
				requirementType,
				expectedStatus,
				requirement.Status,
			)
		}

		return
	}

	t.Fatalf(
		"requirement %s not found",
		requirementType,
	)
}

func TestComplianceCheckIntegration_USA_AllOK(t *testing.T) {
	pool, router := setupIntegrationTest(t)
	defer pool.Close()

	body := `{
		"nationality_code": "HRV",
		"destination_code": "USA",
		"departure_date": "2026-10-01",
		"return_date": "2026-10-20",
		"passport_expiry": "2027-06-01",
		"visa_status": "not_required",
		"eta_status": "obtained"
	}`

	response := checkCompliance(
		t,
		router,
		body,
	)

	assertOverallStatus(
		t,
		response,
		"cleared",
	)

	assertRequirementStatus(
		t,
		response,
		"passport_validity",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"visa",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"eta",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"stay_duration",
		"ok",
	)
}

func TestComplianceCheckIntegration_USA_ETAMissing(t *testing.T) {
	pool, router := setupIntegrationTest(t)
	defer pool.Close()

	body := `{
		"nationality_code": "HRV",
		"destination_code": "USA",
		"departure_date": "2026-10-01",
		"return_date": "2026-10-20",
		"passport_expiry": "2027-06-01",
		"visa_status": "not_required",
		"eta_status": "missing"
	}`

	response := checkCompliance(
		t,
		router,
		body,
	)

	assertOverallStatus(
		t,
		response,
		"action_required",
	)

	assertRequirementStatus(
		t,
		response,
		"passport_validity",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"visa",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"eta",
		"problem",
	)

	assertRequirementStatus(
		t,
		response,
		"stay_duration",
		"ok",
	)
}

func TestComplianceCheckIntegration_CHN_2026_VisaNotRequired(
	t *testing.T,
) {
	pool, router := setupIntegrationTest(t)
	defer pool.Close()

	body := `{
		"nationality_code": "HRV",
		"destination_code": "CHN",
		"departure_date": "2026-10-01",
		"return_date": "2026-10-20",
		"passport_expiry": "2027-06-01",
		"visa_status": "not_required",
		"eta_status": "not_required"
	}`

	response := checkCompliance(
		t,
		router,
		body,
	)

	assertOverallStatus(
		t,
		response,
		"cleared",
	)

	assertRequirementStatus(
		t,
		response,
		"passport_validity",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"visa",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"eta",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"stay_duration",
		"ok",
	)
}

func TestComplianceCheckIntegration_CHN_2027_VisaRequired(
	t *testing.T,
) {
	pool, router := setupIntegrationTest(t)
	defer pool.Close()

	body := `{
		"nationality_code": "HRV",
		"destination_code": "CHN",
		"departure_date": "2027-04-01",
		"return_date": "2027-04-20",
		"passport_expiry": "2028-01-01",
		"visa_status": "missing",
		"eta_status": "not_required"
	}`

	response := checkCompliance(
		t,
		router,
		body,
	)

	assertOverallStatus(
		t,
		response,
		"action_required",
	)

	assertRequirementStatus(
		t,
		response,
		"passport_validity",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"visa",
		"problem",
	)

	assertRequirementStatus(
		t,
		response,
		"eta",
		"ok",
	)

	assertRequirementStatus(
		t,
		response,
		"stay_duration",
		"ok",
	)
}
