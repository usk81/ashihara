package mysql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"github.com/usk81/ashihara/service/holidays/core/domain/entity"
	"github.com/usk81/ashihara/service/holidays/core/domain/repository"
	"github.com/usk81/ashihara/shared/interface/datasource/mysql"
	"github.com/usk81/ashihara/shared/utils/jst"
)

type (
	holidaysImpl struct {
		dbType DBType
		db     mysql.DB
		logger *slog.Logger
	}

	holiday struct {
		DifinitionID int        `db:"difinition_id"`
		Date         string     `db:"date"`
		CreatedAt    *time.Time `db:"created_at"`
		UpdatedAt    *time.Time `db:"updated_at"`
		DeletedAt    *time.Time `db:"deleted_at"`
	}

	definition struct {
		ID          int        `db:"id"`
		Name        string     `db:"name_ja"`
		Summary     string     `db:"summary_ja"`
		Description string     `db:"description_ja"`
		CreatedAt   *time.Time `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
	}

	holidayWithDefinition struct {
		// holiday
		DifinitionID int        `db:"difinition_id"`
		Date         string     `db:"date"`
		CreatedAt    *time.Time `db:"created_at"`
		UpdatedAt    *time.Time `db:"updated_at"`
		DeletedAt    *time.Time `db:"deleted_at"`

		// definition
		Name        string `db:"name_ja"`
		Summary     string `db:"summary_ja"`
		Description string `db:"description_ja"`
	}
)

func HolidaysReader(db *sql.DB, logger *slog.Logger) repository.Holidays {
	return holidays(db, Reader, logger)
}

func HolidaysWriter(db *sql.DB, logger *slog.Logger) repository.Holidays {
	return holidays(db, Writer, logger)
}

func holidays(
	db *sql.DB,
	dbType DBType,
	logger *slog.Logger,
) repository.Holidays {
	return &holidaysImpl{
		dbType: dbType,
		db:     sqlx.NewDb(db, mysql.DriverName),
		logger: logger,
	}
}

func fromDefinition(d *definition) entity.Holiday {
	return entity.Holiday{
		Name:         d.Name,
		DifinitionID: d.ID,
		Summary:      &d.Summary,
		Description:  &d.Description,
	}
}

func (d *holidaysImpl) FindDefinition(ctx context.Context, id int) (output *entity.Holiday, err error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("holiday_difinitions")
	sb.Where(sb.E("id", id))
	q, args := sb.Build()
	d.logger.DebugContext(
		ctx,
		"Holidays.FindDefinition",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
			slog.Any("args", args),
		),
	)
	var v definition
	if err = d.db.GetContext(ctx, &v, q, args...); err != nil {
		d.logger.ErrorContext(
			ctx,
			"Holidays.FindDefinition",
			slog.String("dbType", string(d.dbType)),
			slog.String("action", "db.GetContext"),
			slog.Group(
				"request",
				slog.String("query", q),
				slog.Any("args", args),
			),
			slog.Any("error", err),
		)
		return nil, err
	}

	entity := fromDefinition(&v)
	return &entity, nil
}

func (d *holidaysImpl) FindAllDefinitions(ctx context.Context) (output []*entity.Holiday, err error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("holiday_difinitions")
	q, _ := sb.Build()
	d.logger.DebugContext(
		ctx,
		"Holidays.FindAllDefinitions",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
		),
	)
	var ds []*definition
	if err = d.db.SelectContext(ctx, &ds, q); err != nil {
		d.logger.ErrorContext(
			ctx,
			"Holidays.FindAllDefinitions",
			slog.String("dbType", string(d.dbType)),
			slog.String("action", "db.SelectContext"),
			slog.Group(
				"request",
				slog.String("query", q),
			),
			slog.Any("error", err),
		)
		return
	}

	output = make([]*entity.Holiday, 0, len(ds))
	for _, v := range ds {
		r := fromDefinition(v)
		output = append(output, &r)
	}
	return
}

func (d *holidaysImpl) Search(ctx context.Context, options repository.SearchOption) (output []*entity.Holiday, err error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"holidays.*",
		"holiday_difinitions.name_ja",
		"holiday_difinitions.summary_ja",
		"holiday_difinitions.description_ja",
	)
	sb.From("holidays")
	sb.Join("holiday_difinitions", "holidays.difinition_id = holiday_difinitions.id")
	if options.Limit > 0 {
		sb.Limit(options.Limit)
	}
	if options.Offset > 0 {
		sb.Offset(options.Offset)
	}
	if options.Range != nil {
		if options.Range.Gte != nil && options.Range.Lte != nil {
			sb.Where(sb.And(
				sb.GE("holidays.date", options.Range.Gte),
				sb.LE("holidays.date", options.Range.Lte),
			))
		} else if options.Range.Gte != nil {
			sb.Where(sb.GE("holidays.date", options.Range.Gte))
		} else if options.Range.Lte != nil {
			sb.Where(sb.LE("holidays.date", options.Range.Lte))
		}
	}
	sb.OrderBy("holidays.date DESC")
	q, args := sb.Build()

	d.logger.DebugContext(
		ctx,
		"Holidays.Search",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
			slog.Any("args", args),
		),
	)

	var hs []*holidayWithDefinition
	if err = d.db.SelectContext(ctx, &hs, q, args...); err != nil {
		d.logger.ErrorContext(
			ctx,
			"Holidays.Search",
			slog.String("dbType", string(d.dbType)),
			slog.String("action", "db.SelectContext"),
			slog.Group(
				"request",
				slog.String("query", q),
				slog.Any("args", args),
			),
			slog.Any("error", err),
		)
		return
	}

	output = make([]*entity.Holiday, len(hs))
	for i, h := range hs {
		r := entity.Holiday{
			Date:         h.Date,
			Name:         h.Name,
			DifinitionID: h.DifinitionID,
			Summary:      &h.Summary,
			Description:  &h.Description,
		}
		output[i] = &r
	}
	return
}

func (d *holidaysImpl) Create(ctx context.Context, input *entity.Holiday) (err error) {
	if d.dbType != Writer {
		err = repository.ErrNoPermission
		d.logger.ErrorContext(
			ctx,
			"Holidays.Create",
			slog.String("dbType", string(d.dbType)),
			slog.Any("error", err),
		)
		return
	}
	tt := jst.Now()
	ib := sqlbuilder.
		NewStruct(new(holiday)).
		InsertIgnoreInto("holidays", holiday{
			DifinitionID: input.DifinitionID,
			Date:         input.Date,
			CreatedAt:    &tt,
			UpdatedAt:    &tt,
			DeletedAt:    &tt,
		})
	q, args := ib.Build()
	d.logger.DebugContext(
		ctx,
		"Holidays.Create",
		slog.String("dbType", string(d.dbType)),
		slog.Group(
			"request",
			slog.String("query", q),
			slog.Any("args", args),
		),
	)
	if _, err = d.db.ExecContext(ctx, q, args...); err != nil {
		d.logger.ErrorContext(
			ctx,
			"Holidays.Create",
			slog.String("dbType", string(d.dbType)),
			slog.String("action", "db.ExecContext"),
			slog.Group(
				"request",
				slog.String("query", q),
				slog.Any("args", args),
			),
			slog.Any("error", err),
		)
	}
	return
}
