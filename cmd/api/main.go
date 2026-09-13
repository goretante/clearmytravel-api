package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/goretante/clearmytravel-api/internal/compliance"
	"github.com/goretante/clearmytravel-api/internal/config"
	"github.com/goretante/clearmytravel-api/internal/db"
	apphttp "github.com/goretante/clearmytravel-api/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// Compliance dependencies
	visaRepository := compliance.NewDBVisaRuleRepository(queries)

	visaEvaluator := compliance.NewRepositoryVisaEvaluator(
		visaRepository,
	)

	etaRepository := compliance.NewDBETARuleRepository(queries)
	etaEvaluator := compliance.NewRepositoryETAEvaluator(etaRepository)

	stayRepository := compliance.NewDBStayRuleRepository(queries)
	stayEvaluator := compliance.NewRepositoryStayDurationEvaluator(stayRepository)

	complianceService := compliance.NewService(
		compliance.PassportValidityEvaluator{
			BufferDays: 6,
		},
		visaEvaluator,
		etaEvaluator,
		stayEvaluator,
	)

	router := apphttp.NewRouter(queries, complianceService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("server listening on :8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
