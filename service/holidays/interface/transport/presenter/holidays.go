package presenter

import (
	"context"

	"github.com/usk81/ashihara/service/holidays/core/domain/usecase"
)

type (
	HolidayBloc interface {
		Find(ctx context.Context, input usecase.HolidayFinderInput) (output *usecase.HolidayFinderOutput, err error)
		Import(ctx context.Context) (err error)
		Search(ctx context.Context, input usecase.HolidaySearcherInput) (output *usecase.HolidaySearcherOutput, err error)
	}

	holidayBlocImpl struct {
		findUsecase   usecase.HolidayFinder
		importUsecase usecase.HolidayImporter
		searchUsecase usecase.HolidaySearcher
	}
)

func NewHolidayBloc(
	findUsecase usecase.HolidayFinder,
	importUsecase usecase.HolidayImporter,
	searchUsecase usecase.HolidaySearcher,
) HolidayBloc {
	return &holidayBlocImpl{
		findUsecase:   findUsecase,
		importUsecase: importUsecase,
		searchUsecase: searchUsecase,
	}
}

func (b *holidayBlocImpl) Find(ctx context.Context, input usecase.HolidayFinderInput) (output *usecase.HolidayFinderOutput, err error) {
	return b.findUsecase.Execute(ctx, input)
}

func (b *holidayBlocImpl) Import(ctx context.Context) (err error) {
	return b.importUsecase.Execute(ctx)
}

func (b *holidayBlocImpl) Search(ctx context.Context, input usecase.HolidaySearcherInput) (output *usecase.HolidaySearcherOutput, err error) {
	return b.searchUsecase.Execute(ctx, input)
}
