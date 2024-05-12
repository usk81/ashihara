package repository

import (
	"context"
	"time"

	"github.com/usk81/ashihara/service/holidays/core/domain/entity"
)

type (
	HolidayEntity struct {
		Date time.Time
		Name string
	}

	DateRange struct {
		Gte *string
		Lte *string
	}

	SearchOption struct {
		Range  *DateRange
		Limit  int
		Offset int
	}

	DifinitionFinder interface {
		FindDefinition(ctx context.Context, id int) (output *entity.Holiday, err error)
	}

	DifinitionAllFinder interface {
		FindAllDefinitions(ctx context.Context) (output []*entity.Holiday, err error)
	}

	Searcher interface {
		Search(ctx context.Context, options SearchOption) (output []*entity.Holiday, err error)
	}

	Creater interface {
		Create(ctx context.Context, input *entity.Holiday) (err error)
	}

	Importer interface {
		DifinitionAllFinder
		Creater
	}

	Holidays interface {
		DifinitionFinder
		DifinitionAllFinder
		Searcher
		Creater
	}

	Crawler interface {
		Crawl(ctx context.Context) (output []*HolidayEntity, err error)
	}
)
