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

type NoteService struct {
	partyRepository repository.PartyRepository
	noteRepository  repository.NoteRepository
}

func NewNoteService(
	partyRepository repository.PartyRepository,
	noteRepository repository.NoteRepository,
) *NoteService {
	return &NoteService{
		partyRepository: partyRepository,
		noteRepository:  noteRepository,
	}
}

func (s *NoteService) CreateNote(
	ctx context.Context,
	input domain.CreateNoteInput,
) (db.PartyNote, error) {
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.PartyNote{}, ErrInvalidPartyID
	}

	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyNote{}, ErrPartyNotFound
		}
		return db.PartyNote{}, err
	}

	content := strings.TrimSpace(input.Content)
	if content == "" {
		return db.PartyNote{}, ErrEmptyNoteContent
	}

	var createdBy pgtype.UUID
	if input.CreatedBy != nil && *input.CreatedBy != "" {
		if err := createdBy.Scan(*input.CreatedBy); err != nil {
			return db.PartyNote{}, err
		}
	}

	params := db.CreatePartyNoteParams{
		PartyID:   partyID,
		NoteType:  input.NoteType,
		Content:   content,
		CreatedBy: createdBy,
	}

	return s.noteRepository.CreatePartyNote(ctx, params)
}

func (s *NoteService) GetNoteByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyNote, error) {
	note, err := s.noteRepository.GetPartyNoteByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.PartyNote{}, ErrNoteNotFound
		}
		return db.PartyNote{}, err
	}
	return note, nil
}

func (s *NoteService) ListNotesByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) ([]db.PartyNote, error) {
	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPartyNotFound
		}
		return nil, err
	}

	return s.noteRepository.ListPartyNotesByPartyID(ctx, partyID)
}

func (s *NoteService) DeleteNote(
	ctx context.Context,
	id pgtype.UUID,
) error {
	_, err := s.noteRepository.GetPartyNoteByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoteNotFound
		}
		return err
	}
	return s.noteRepository.DeletePartyNote(ctx, id)
}
