package repository

import (
	"context"
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
	"github.com/jeffvo/event-sourcing/internal/models"
)

type ProductRepositoryInterface interface {
	CreateNewStream(context.Context, uuid.UUID, esdb.EventData) error
	AddEventToStream(context.Context, uuid.UUID, esdb.EventData) error
	GetLatestProductVersion(context.Context, uuid.UUID) (models.ProductStock, error)
}

type ProductRepository struct {
	conn *esdb.Client
}

func CreateProductRepository(db *esdb.Client) ProductRepositoryInterface {
	return &ProductRepository{
		conn: db,
	}
}

// instead of returning the list of events I have decided to create the latest product version in the repository
// this is because this could would be used in two places,
// but also we have to loop through all the events from the stream so we can close it in this scope
func (repo *ProductRepository) GetLatestProductVersion(ctx context.Context, id uuid.UUID) (models.ProductStock, error) {

	var db = repo.conn

	var stream, err = db.ReadStream(ctx, id.String(), esdb.ReadStreamOptions{
		Direction: esdb.Forwards,    // So we make sure we move forwards through the events
		From:      esdb.Revision(0), // Start with the first event
	}, 100) // Limit to 100 events

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	// The product model we want to return
	var product = models.ProductStock{}

	for {
		resolvedEvent, err := stream.Recv()

		// If we don't receive anything from the stream we will retrieve an error meaning we have no other events
		if err != nil {
			break
		}

		// This is its own variable for readability
		event := resolvedEvent.Event

		var productStock models.ProductStock
		err = json.Unmarshal(event.Data, &productStock)

		if err != nil {
			panic(err)
		}

		if event.EventType == "add" {
			product.Amount += productStock.Amount
		}

		if event.EventType == "remove" {
			product.Amount -= productStock.Amount
		}
	}

	return product, nil
}

func (repo *ProductRepository) CreateNewStream(ctx context.Context, id uuid.UUID, event esdb.EventData) error {
	var db = repo.conn
	var _, err = db.AppendToStream(ctx, id.String(), esdb.AppendToStreamOptions{}, event)

	if err != nil {
		return err
	}

	return err
}

func (repo *ProductRepository) AddEventToStream(ctx context.Context, id uuid.UUID, event esdb.EventData) error {
	var db = repo.conn

	readOperations := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}

	//We first read the last event so we now what the latest version is
	stream, err := db.ReadStream(context.Background(), id.String(), readOperations, 1)
	if err != nil {

		return err
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		return err
	}

	// Based on the last revision we will append the event at the end of the list
	var expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	_, err = db.AppendToStream(
		ctx,
		id.String(),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		event,
	)

	if err != nil {

		return err
	}

	return nil
}
