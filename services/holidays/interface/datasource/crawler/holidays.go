package crawler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/usk81/ashihara/services/holidays/core/domain/repository"
	"github.com/usk81/ashihara/shared/utils/jst"
)

type (
	holidayImpl struct {
		client *http.Client
		logger *slog.Logger
	}

	holiday struct {
		Date string `csv:"国民の祝日・休日月日"`
		Name string `csv:"国民の祝日・休日名称"`
	}
)

const (
	holidayCsvURL = "https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv"
)

func Holidays(c *http.Client, logger *slog.Logger) repository.Crawler {
	return &holidayImpl{
		client: c,
		logger: logger,
	}
}

func (d *holidayImpl) Crawl(ctx context.Context) (output []*repository.HolidayEntity, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, holidayCsvURL, http.NoBody)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"holidays.Crawl",
			slog.String("action", "http.NewRequestWithContext"),
			slog.Any("error", err),
		)
		return
	}

	res, err := d.client.Do(req)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"holidays.Crawl",
			slog.String("action", "client.Do"),
			slog.Any("error", err),
		)
		return
	}

	rb := res.Body
	defer func() {
		_ = rb.Close()
	}()

	r := transform.NewReader(rb, japanese.ShiftJIS.NewDecoder())
	hs := []*holiday{}
	if err = gocsv.Unmarshal(r, &hs); err != nil {
		d.logger.ErrorContext(
			ctx,
			"holidays.Crawl",
			slog.String("action", "gocsv.Unmarshal"),
			slog.Any("error", err),
		)
		return
	}

	output = make([]*repository.HolidayEntity, 0, len(hs))
	for _, h := range hs {
		if h.Name == "休日（祝日扱い）" {
			h.Name = "休日"
		} else if h.Name == "体育の日（スポーツの日）" {
			h.Name = "体育の日"
		}

		dt, err := jst.Parse("2006/1/2", h.Date)
		if err != nil {
			d.logger.ErrorContext(
				ctx,
				"holidays.Crawl",
				slog.String("action", "jst.Parse"),
				slog.String("date", h.Date),
				slog.Any("error", err),
			)
			return nil, err
		}

		output = append(output, &repository.HolidayEntity{
			Date: dt,
			Name: h.Name,
		})
	}
	return
}
