package repository

import (
	"context"
	"errors"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/location"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/google/uuid"
)

func (r *repository) ListLocations(filters map[string]interface{}) ([]*ent.Location, error) {
	ctx := context.Background()
	query := r.db.Location.Query()

	if city, ok := filters["city"].(string); ok && city != "" {
		query = query.Where(location.CityEQ(city))
	}

	locations, err := query.All(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to list locations")
		return nil, err
	}
	return locations, nil
}

func (r *repository) GetLocationByID(id string) (*ent.Location, error) {
	location, err := r.db.Location.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("location not found")
		}
		logger.Get().Error().Err(err).Msg("Failed to get location")
		return nil, err
	}
	return location, nil
}

func (r *repository) AddLocation(localityName, city, state, phoneNumber, country, pincode string) (*ent.Location, error) {
	location, err := r.db.Location.Create().
		SetID(uuid.New().String()).
		SetLocalityName(localityName).
		SetCity(city).
		SetState(state).
		SetPhoneNumber(phoneNumber).
		SetCountry(country).
		SetPincode(pincode).
		SetIsActive(true).
		Save(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create location")
		return nil, err
	}
	return location, nil
}

func (r *repository) SoftDeleteLocation(id string) error {
	_, err := r.db.Location.UpdateOneID(id).
		SetIsActive(false).
		Save(context.Background())
	return err
}

func (r *repository) GetAllUniqueCities() ([]string, error) {
	cities, err := r.db.Location.Query().Select(location.FieldCity).Strings(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get all unique cities")
		return nil, err
	}
	return cities, nil
}

func (r *repository) GetAllUniqueLocations() ([]string, error) {
	locations, err := r.db.Location.Query().Select(location.FieldLocalityName).Strings(context.Background())
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get all unique locations")
		return nil, err
	}
	return locations, nil
}
