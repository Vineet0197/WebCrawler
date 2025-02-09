package api

import (
	"net/http"

	"github.com/Vineet0197/shoppin-webcrawler/internal/services"
	"github.com/Vineet0197/shoppin-webcrawler/pkg/utils"
	"github.com/labstack/echo"
)

func AcceptDomainAndURL(ctx echo.Context, domains []string, ks *services.KafkaStorage) error {
	// Accept Domains and produce/post to Kafka
	// Validate and Produce to Kafka
	for _, domain := range domains {
		if !utils.IsValidDomain(domain) {
			ctx.Logger().Errorf("Invalid domain: %s", domain)
			continue
		}

		err := ks.Produce(domain)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "Failed to produce message to Kafka")
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "success"})
}
