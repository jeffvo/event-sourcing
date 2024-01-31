package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/jeffvo/event-sourcing/internal/dto"
	"github.com/jeffvo/event-sourcing/internal/models"
	repository "github.com/jeffvo/event-sourcing/internal/repositories"
)

type ProductUsecaseInterface interface {
	NewProductInStock(ctx context.Context, product *dto.NewProductStockDTO) (*models.ProductStockId, error)
	AdjustStock(ctx context.Context, product *dto.UpdateProductStockDto) error
	GetProductStockById(ctx context.Context, id uuid.UUID) (*models.ProductStock, error)
}

type ProductUsecase struct {
	repo repository.ProductRepositoryInterface
}

func CreateProductUsecase(repo repository.ProductRepositoryInterface) ProductUsecaseInterface {
	return &ProductUsecase{
		repo: repo,
	}
}

func (usecase *ProductUsecase) NewProductInStock(ctx context.Context, product *dto.NewProductStockDTO) (*models.ProductStockId, error) {
	product.Action = models.Action(models.Add)
	var event = dto.ToEvent(product)
	var err = usecase.repo.CreateNewStream(ctx, product.ID, *event)
	if err != nil {
		return nil, err
	}

	return &models.ProductStockId{ID: product.ID}, nil
}

// Future improvement would be to make sure if we two request happen at the same time they both don't remove the stock
// Race condition could potentially happen here
func (usecase *ProductUsecase) AdjustStock(ctx context.Context, product *dto.UpdateProductStockDto) error {
	var event = dto.ToPutEvent(product)

	// In order to do a remove function we first have to make sure we have enough stock to even remove it
	if product.Action == models.Remove {
		var productStock, err = usecase.repo.GetLatestProductVersion(ctx, product.ID)
		if err != nil {
			panic(err)
		}

		// Check the max and return error if the stock isn't enough
		if productStock.Amount-product.Amount < 0 {
			var max = strconv.FormatInt(productStock.Amount, 10)
			return errors.New("We don't have enough stock the maximum amount is " + max)
		}
	}

	var err = usecase.repo.AddEventToStream(ctx, product.ID, *event)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *ProductUsecase) GetProductStockById(ctx context.Context, id uuid.UUID) (*models.ProductStock, error) {
	var product, err = usecase.repo.GetLatestProductVersion(ctx, id)

	if err != nil {
		return nil, err
	}

	return &product, err
}
