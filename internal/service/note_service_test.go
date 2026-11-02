package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
)

type mockNoteRepository struct {
	createPartyNoteCalled       bool
	createPartyNoteParams       db.CreatePartyNoteParams
	getPartyNoteByIDFunc        func(ctx context.Context, id pgtype.UUID) (db.PartyNote, error)
	listPartyNotesByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error)
	deletePartyNoteFunc         func(ctx context.Context, id pgtype.UUID) error
}

func (m *mockNoteRepository) CreatePartyNote(ctx context.Context, params db.CreatePartyNoteParams) (db.PartyNote, error) {
	m.createPartyNoteCalled = true
	m.createPartyNoteParams = params
	return db.PartyNote{
		ID:       pgtype.UUID{Valid: true},
		PartyID:  params.PartyID,
		NoteType: params.NoteType,
		Content:  params.Content,
	}, nil
}

func (m *mockNoteRepository) GetPartyNoteByID(ctx context.Context, id pgtype.UUID) (db.PartyNote, error) {
	if m.getPartyNoteByIDFunc != nil {
		return m.getPartyNoteByIDFunc(ctx, id)
	}
	return db.PartyNote{}, pgx.ErrNoRows
}

func (m *mockNoteRepository) ListPartyNotesByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error) {
	if m.listPartyNotesByPartyIDFunc != nil {
		return m.listPartyNotesByPartyIDFunc(ctx, partyID)
	}
	return []db.PartyNote{}, nil
}

func (m *mockNoteRepository) DeletePartyNote(ctx context.Context, id pgtype.UUID) error {
	if m.deletePartyNoteFunc != nil {
		return m.deletePartyNoteFunc(ctx, id)
	}
	return nil
}

func TestNote_Flow(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	noteRepo := &mockNoteRepository{}
	svc := NewNoteService(partyRepo, noteRepo)

	// Create
	note, err := svc.CreateNote(context.Background(), domain.CreateNoteInput{
		PartyID: partyUUID,
		Content: "Followed up with client regarding annual contract",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "Followed up with client regarding annual contract" {
		t.Fatalf("expected content match, got %s", note.Content)
	}

	// Empty content validation
	_, err = svc.CreateNote(context.Background(), domain.CreateNoteInput{
		PartyID: partyUUID,
		Content: "   ",
	})
	if !errors.Is(err, ErrEmptyNoteContent) {
		t.Fatalf("expected ErrEmptyNoteContent, got %v", err)
	}
}
