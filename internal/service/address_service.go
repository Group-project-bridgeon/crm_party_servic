package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/repository"
)

type AddressService struct {
	partyRepository   repository.PartyRepository
	addressRepository repository.AddressRepository
}

func NewAddressService(
	partyRepository repository.PartyRepository,
	addressRepository repository.AddressRepository,
) *AddressService {
	return &AddressService{
		partyRepository:   partyRepository,
		addressRepository: addressRepository,
	}
}

func (s *AddressService) CreateAddress(
	ctx context.Context,
	input domain.CreateAddressInput,
) (db.PartyAddress, error) {
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.PartyAddress{}, ErrInvalidPartyID
	}

	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyAddress{}, ErrPartyNotFound
		}
		return db.PartyAddress{}, err
	}

	addressType := strings.ToUpper(strings.TrimSpace(input.AddressType))
	if !domain.AddressType(addressType).IsValid() {
		return db.PartyAddress{}, ErrInvalidAddressType
	}

	addressLine1 := strings.TrimSpace(input.AddressLine1)
	if addressLine1 == "" {
		return db.PartyAddress{}, ErrEmptyAddressLine1
	}

	city := strings.TrimSpace(input.City)
	if city == "" {
		return db.PartyAddress{}, ErrEmptyCity
	}

	countryCode := strings.ToUpper(strings.TrimSpace(input.CountryCode))
	if len(countryCode) != 2 {
		return db.PartyAddress{}, ErrInvalidCountryCode
	}

	var latitude pgtype.Numeric
	if input.Latitude != nil {
		if *input.Latitude < -90 || *input.Latitude > 90 {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
		if err := latitude.Scan(fmt.Sprintf("%.7f", *input.Latitude)); err != nil {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
	}

	var longitude pgtype.Numeric
	if input.Longitude != nil {
		if *input.Longitude < -180 || *input.Longitude > 180 {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
		if err := longitude.Scan(fmt.Sprintf("%.7f", *input.Longitude)); err != nil {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
	}

	if input.IsPrimary {
		_ = s.addressRepository.ResetPrimaryPartyAddress(ctx, db.ResetPrimaryPartyAddressParams{
			PartyID:     partyID,
			AddressType: addressType,
		})
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	params := db.CreatePartyAddressParams{
		PartyID:      partyID,
		AddressType:  addressType,
		AddressName:  input.AddressName,
		AddressLine1: addressLine1,
		AddressLine2: input.AddressLine2,
		Landmark:     input.Landmark,
		City:         city,
		State:        input.State,
		PostalCode:   input.PostalCode,
		CountryCode:  countryCode,
		Latitude:     latitude,
		Longitude:    longitude,
		IsPrimary:    input.IsPrimary,
		IsActive:     isActive,
	}

	return s.addressRepository.CreatePartyAddress(ctx, params)
}

func (s *AddressService) GetAddressByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyAddress, error) {
	address, err := s.addressRepository.GetPartyAddressByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyAddress{}, ErrAddressNotFound
		}
		return db.PartyAddress{}, err
	}
	return address, nil
}

func (s *AddressService) ListAddressesByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
	addressType *string,
) ([]db.PartyAddress, error) {
	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPartyNotFound
		}
		return nil, err
	}

	var at *string
	if addressType != nil && strings.TrimSpace(*addressType) != "" {
		cleaned := strings.ToUpper(strings.TrimSpace(*addressType))
		at = &cleaned
	}

	params := db.ListPartyAddressesByPartyIDParams{
		PartyID:     partyID,
		AddressType: at,
	}

	return s.addressRepository.ListPartyAddressesByPartyID(ctx, params)
}

func (s *AddressService) UpdateAddress(
	ctx context.Context,
	input domain.UpdateAddressInput,
) (db.PartyAddress, error) {
	var addressID pgtype.UUID
	if err := addressID.Scan(input.ID); err != nil {
		return db.PartyAddress{}, ErrInvalidAddressID
	}

	existing, err := s.addressRepository.GetPartyAddressByID(ctx, addressID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyAddress{}, ErrAddressNotFound
		}
		return db.PartyAddress{}, err
	}

	addressType := strings.ToUpper(strings.TrimSpace(input.AddressType))
	if !domain.AddressType(addressType).IsValid() {
		return db.PartyAddress{}, ErrInvalidAddressType
	}

	addressLine1 := strings.TrimSpace(input.AddressLine1)
	if addressLine1 == "" {
		return db.PartyAddress{}, ErrEmptyAddressLine1
	}

	city := strings.TrimSpace(input.City)
	if city == "" {
		return db.PartyAddress{}, ErrEmptyCity
	}

	countryCode := strings.ToUpper(strings.TrimSpace(input.CountryCode))
	if len(countryCode) != 2 {
		return db.PartyAddress{}, ErrInvalidCountryCode
	}

	var latitude pgtype.Numeric
	if input.Latitude != nil {
		if *input.Latitude < -90 || *input.Latitude > 90 {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
		if err := latitude.Scan(fmt.Sprintf("%.7f", *input.Latitude)); err != nil {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
	}

	var longitude pgtype.Numeric
	if input.Longitude != nil {
		if *input.Longitude < -180 || *input.Longitude > 180 {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
		if err := longitude.Scan(fmt.Sprintf("%.7f", *input.Longitude)); err != nil {
			return db.PartyAddress{}, ErrInvalidCoordinates
		}
	}

	if input.IsPrimary {
		_ = s.addressRepository.ResetPrimaryPartyAddress(ctx, db.ResetPrimaryPartyAddressParams{
			PartyID:     existing.PartyID,
			AddressType: addressType,
		})
	}

	params := db.UpdatePartyAddressParams{
		ID:           addressID,
		AddressType:  addressType,
		AddressName:  input.AddressName,
		AddressLine1: addressLine1,
		AddressLine2: input.AddressLine2,
		Landmark:     input.Landmark,
		City:         city,
		State:        input.State,
		PostalCode:   input.PostalCode,
		CountryCode:  countryCode,
		Latitude:     latitude,
		Longitude:    longitude,
		IsPrimary:    input.IsPrimary,
		IsActive:     input.IsActive,
	}

	return s.addressRepository.UpdatePartyAddress(ctx, params)
}

func (s *AddressService) DeleteAddress(
	ctx context.Context,
	id pgtype.UUID,
) error {
	_, err := s.addressRepository.GetPartyAddressByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAddressNotFound
		}
		return err
	}
	return s.addressRepository.DeletePartyAddress(ctx, id)
}
