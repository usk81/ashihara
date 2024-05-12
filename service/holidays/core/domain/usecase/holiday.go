package usecase

import (
	"context"

	"github.com/usk81/ashihara/service/holidays/core/domain/entity"
	"github.com/usk81/ashihara/shared/domain"
)

type (
	DateRange struct {
		From string
		To   string
	}

	HolidayFinderInput struct {
		ID int
	}

	HolidayFinderOutput struct {
		Holiday *entity.Holiday
	}

	HolidaySearcherInput struct {
		Fields    []string
		DateRange *DateRange
		Limit     int
		Offset    int
	}

	HolidaySearcherOutput struct {
		Holidays []*entity.Holiday
	}

	HolidayFinder = domain.Usecase[HolidayFinderInput, HolidayFinderOutput]

	HolidaySearcher = domain.Usecase[HolidaySearcherInput, HolidaySearcherOutput]

	HolidayImporter interface {
		Execute(ctx context.Context) (err error)
	}
)
