package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
	"strings"

	"github.com/radon1/pg-stat-test-task/internal/repositories"
)

const filterSplitSeparator = ","

type container struct {
	logger       *zerolog.Logger
	repositories *repositories.Container
}

func Register(logger *zerolog.Logger, repositories *repositories.Container, fiberApp *fiber.App) {
	c := &container{
		logger:       logger,
		repositories: repositories,
	}

	fiberApp.Get("/queries-stat", c.findQueriesStat)
}

func (c *container) findQueriesStat(ctx *fiber.Ctx) error {
	var (
		limit       = 20
		offset      = 0
		filterParam = ctx.Query("filter")
	)

	queryLimitParam := ctx.Query("limit")
	if len(queryLimitParam) > 0 {
		queryLimit, err := strconv.Atoi(queryLimitParam)
		if err != nil {
			c.logger.Error().Err(err).Msg("invalid parse limit from the query")
			return ctx.Status(http.StatusBadRequest).JSON(map[string]string{
				"error": "invalid params",
			})
		}
		limit = queryLimit
	}

	queryOffsetParam := ctx.Query("offset")
	if len(queryOffsetParam) > 0 {
		queryOffset, err := strconv.Atoi(queryOffsetParam)
		if err != nil {
			c.logger.Error().Err(err).Msg("invalid parse offset from the query")
			return ctx.Status(http.StatusBadRequest).JSON(map[string]string{
				"error": "invalid params",
			})
		}
		offset = queryOffset
	}

	var queryFilters []string
	if len(filterParam) > 0 {
		queryFilters = strings.Split(filterParam, filterSplitSeparator)
	}
	stats, err := c.repositories.PGStat.FindQueriesStat(ctx.Context(), limit, offset, queryFilters)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to find queries stats")
		return ctx.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": "something is wrong",
		})
	}

	return ctx.Status(http.StatusOK).JSON(stats)
}
