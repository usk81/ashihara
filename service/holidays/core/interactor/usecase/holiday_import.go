package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/usk81/ashihara/service/holidays/core/domain/entity"
	"github.com/usk81/ashihara/service/holidays/core/domain/errors"
	"github.com/usk81/ashihara/service/holidays/core/domain/repository"
	"github.com/usk81/ashihara/service/holidays/core/domain/usecase"
)

type (
	importHolidayImpl struct {
		crawler  repository.Crawler
		importer repository.Importer
		logger   *slog.Logger
	}
)

func ImportHoliday(
	crawler repository.Crawler,
	importer repository.Importer,
	logger *slog.Logger,
) usecase.HolidayImporter {
	return &importHolidayImpl{
		crawler:  crawler,
		importer: importer,
		logger:   logger,
	}
}

func (u *importHolidayImpl) Execute(ctx context.Context) (err error) {
	rs, err := u.importer.FindAllDefinitions(ctx)
	if err != nil {
		return
	}
	ds := map[string]int{}
	for _, r := range rs {
		ds[r.Name] = r.DifinitionID
	}

	entities, err := u.crawler.Crawl(ctx)
	if err != nil {
		u.logger.ErrorContext(ctx,
			"ImportHoliday.Execute",
			slog.String("action", "crawler.Crawl"),
			slog.Any("error", err),
		)
		return errors.NewCause(err, errors.CaseBackendError)
	}

	hs := make([]*entity.Holiday, 0, len(entities))
	for _, et := range entities {
		id, ok := ds[et.Name]
		if !ok {
			err = fmt.Errorf("%s does not found", et.Name)
			u.logger.ErrorContext(ctx,
				"ImportHoliday.Execute",
				slog.String("action", "crawler.Crawl"),
				slog.Any("error", err),
			)
			return errors.NewCause(err, errors.CaseBackendError)
		}
		hs = append(hs, &entity.Holiday{
			Date:         et.Date.Format("2006-01-02"),
			DifinitionID: id,
		})
	}
	for _, h := range hs {
		if err = u.importer.Create(ctx, h); err != nil {
			u.logger.ErrorContext(ctx,
				"ImportHoliday.Execute",
				slog.String("action", "importer.Create"),
				slog.Any("error", err),
			)
			return errors.NewCause(err, errors.CaseBackendError)
		}
	}
	return
}
