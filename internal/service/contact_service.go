package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/repository"
)

type ContactService struct {
	partyRepository   repository.PartyRepository
	contactRepository repository.ContactRepository
}

func NewContactService(
	partyRepository repository.PartyRepository,
	contactRepository repository.ContactRepository,
) *ContactService {
	return &ContactService{
		partyRepository:   partyRepository,
		contactRepository: contactRepository,
	}
}

// ----------------- Direct Contacts -----------------

func (s *ContactService) CreateContact(
	ctx context.Context,
	input domain.CreateContactInput,
) (db.PartyContact, error) {
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.PartyContact{}, ErrInvalidPartyID
	}

	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyContact{}, ErrPartyNotFound
		}
		return db.PartyContact{}, err
	}

	contactType := strings.ToUpper(strings.TrimSpace(input.ContactType))
	if !domain.ContactType(contactType).IsValid() {
		return db.PartyContact{}, ErrInvalidContactType
	}

	contactValue := strings.TrimSpace(input.ContactValue)
	if contactValue == "" {
		return db.PartyContact{}, ErrEmptyContactValue
	}

	if input.IsPrimary {
		_ = s.contactRepository.ResetPrimaryPartyContact(ctx, db.ResetPrimaryPartyContactParams{
			PartyID:     partyID,
			ContactType: contactType,
		})
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	params := db.CreatePartyContactParams{
		PartyID:      partyID,
		ContactType:  contactType,
		ContactValue: contactValue,
		Label:        input.Label,
		IsPrimary:    input.IsPrimary,
		IsActive:     isActive,
	}

	return s.contactRepository.CreatePartyContact(ctx, params)
}

func (s *ContactService) GetContactByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyContact, error) {
	contact, err := s.contactRepository.GetPartyContactByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyContact{}, ErrContactNotFound
		}
		return db.PartyContact{}, err
	}
	return contact, nil
}

func (s *ContactService) ListContactsByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
	contactType *string,
) ([]db.PartyContact, error) {
	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPartyNotFound
		}
		return nil, err
	}

	var ct *string
	if contactType != nil && strings.TrimSpace(*contactType) != "" {
		cleaned := strings.ToUpper(strings.TrimSpace(*contactType))
		ct = &cleaned
	}

	params := db.ListPartyContactsByPartyIDParams{
		PartyID:     partyID,
		ContactType: ct,
	}

	return s.contactRepository.ListPartyContactsByPartyID(ctx, params)
}

func (s *ContactService) UpdateContact(
	ctx context.Context,
	input domain.UpdateContactInput,
) (db.PartyContact, error) {
	var contactID pgtype.UUID
	if err := contactID.Scan(input.ID); err != nil {
		return db.PartyContact{}, ErrInvalidContactID
	}

	existing, err := s.contactRepository.GetPartyContactByID(ctx, contactID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyContact{}, ErrContactNotFound
		}
		return db.PartyContact{}, err
	}

	contactType := strings.ToUpper(strings.TrimSpace(input.ContactType))
	if !domain.ContactType(contactType).IsValid() {
		return db.PartyContact{}, ErrInvalidContactType
	}

	contactValue := strings.TrimSpace(input.ContactValue)
	if contactValue == "" {
		return db.PartyContact{}, ErrEmptyContactValue
	}

	if input.IsPrimary {
		_ = s.contactRepository.ResetPrimaryPartyContact(ctx, db.ResetPrimaryPartyContactParams{
			PartyID:     existing.PartyID,
			ContactType: contactType,
		})
	}

	params := db.UpdatePartyContactParams{
		ID:           contactID,
		ContactType:  contactType,
		ContactValue: contactValue,
		Label:        input.Label,
		IsPrimary:    input.IsPrimary,
		IsActive:     input.IsActive,
	}

	return s.contactRepository.UpdatePartyContact(ctx, params)
}

func (s *ContactService) DeleteContact(
	ctx context.Context,
	id pgtype.UUID,
) error {
	_, err := s.contactRepository.GetPartyContactByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrContactNotFound
		}
		return err
	}
	return s.contactRepository.DeletePartyContact(ctx, id)
}

// ----------------- Contact Persons -----------------

func (s *ContactService) CreateContactPerson(
	ctx context.Context,
	input domain.CreateContactPersonInput,
) (db.PartyContactPerson, error) {
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.PartyContactPerson{}, ErrInvalidPartyID
	}

	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyContactPerson{}, ErrPartyNotFound
		}
		return db.PartyContactPerson{}, err
	}

	firstName := strings.TrimSpace(input.FirstName)
	if firstName == "" {
		return db.PartyContactPerson{}, ErrEmptyFirstName
	}

	if input.IsPrimary {
		_ = s.contactRepository.ResetPrimaryPartyContactPerson(ctx, partyID)
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	params := db.CreatePartyContactPersonParams{
		PartyID:     partyID,
		FirstName:   firstName,
		MiddleName:  input.MiddleName,
		LastName:    input.LastName,
		Designation: input.Designation,
		Department:  input.Department,
		Email:       input.Email,
		Phone:       input.Phone,
		Mobile:      input.Mobile,
		IsPrimary:   input.IsPrimary,
		IsActive:    isActive,
		Notes:       input.Notes,
	}

	return s.contactRepository.CreatePartyContactPerson(ctx, params)
}

func (s *ContactService) GetContactPersonByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyContactPerson, error) {
	person, err := s.contactRepository.GetPartyContactPersonByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyContactPerson{}, ErrContactPersonNotFound
		}
		return db.PartyContactPerson{}, err
	}
	return person, nil
}

func (s *ContactService) ListContactPersonsByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) ([]db.PartyContactPerson, error) {
	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPartyNotFound
		}
		return nil, err
	}

	return s.contactRepository.ListPartyContactPersonsByPartyID(ctx, partyID)
}

func (s *ContactService) UpdateContactPerson(
	ctx context.Context,
	input domain.UpdateContactPersonInput,
) (db.PartyContactPerson, error) {
	var personID pgtype.UUID
	if err := personID.Scan(input.ID); err != nil {
		return db.PartyContactPerson{}, ErrInvalidContactPersonID
	}

	existing, err := s.contactRepository.GetPartyContactPersonByID(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyContactPerson{}, ErrContactPersonNotFound
		}
		return db.PartyContactPerson{}, err
	}

	firstName := strings.TrimSpace(input.FirstName)
	if firstName == "" {
		return db.PartyContactPerson{}, ErrEmptyFirstName
	}

	if input.IsPrimary {
		_ = s.contactRepository.ResetPrimaryPartyContactPerson(ctx, existing.PartyID)
	}

	params := db.UpdatePartyContactPersonParams{
		ID:          personID,
		FirstName:   firstName,
		MiddleName:  input.MiddleName,
		LastName:    input.LastName,
		Designation: input.Designation,
		Department:  input.Department,
		Email:       input.Email,
		Phone:       input.Phone,
		Mobile:      input.Mobile,
		IsPrimary:   input.IsPrimary,
		IsActive:    input.IsActive,
		Notes:       input.Notes,
	}

	return s.contactRepository.UpdatePartyContactPerson(ctx, params)
}

func (s *ContactService) DeleteContactPerson(
	ctx context.Context,
	id pgtype.UUID,
) error {
	_, err := s.contactRepository.GetPartyContactPersonByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrContactPersonNotFound
		}
		return err
	}
	return s.contactRepository.DeletePartyContactPerson(ctx, id)
}
