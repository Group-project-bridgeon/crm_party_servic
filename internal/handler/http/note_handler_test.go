package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

type mockNoteService struct {
	createNoteFunc         func(ctx context.Context, input domain.CreateNoteInput) (db.PartyNote, error)
	getNoteByIDFunc        func(ctx context.Context, id pgtype.UUID) (db.PartyNote, error)
	listNotesByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error)
	deleteNoteFunc         func(ctx context.Context, id pgtype.UUID) error
}

func (m *mockNoteService) CreateNote(ctx context.Context, input domain.CreateNoteInput) (db.PartyNote, error) {
	if m.createNoteFunc != nil {
		return m.createNoteFunc(ctx, input)
	}
	return db.PartyNote{}, nil
}

func (m *mockNoteService) GetNoteByID(ctx context.Context, id pgtype.UUID) (db.PartyNote, error) {
	if m.getNoteByIDFunc != nil {
		return m.getNoteByIDFunc(ctx, id)
	}
	return db.PartyNote{}, nil
}

func (m *mockNoteService) ListNotesByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error) {
	if m.listNotesByPartyIDFunc != nil {
		return m.listNotesByPartyIDFunc(ctx, partyID)
	}
	return []db.PartyNote{}, nil
}

func (m *mockNoteService) DeleteNote(ctx context.Context, id pgtype.UUID) error {
	if m.deleteNoteFunc != nil {
		return m.deleteNoteFunc(ctx, id)
	}
	return nil
}

func TestNoteHandler_Flow(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockNoteService{
		createNoteFunc: func(ctx context.Context, input domain.CreateNoteInput) (db.PartyNote, error) {
			return db.PartyNote{
				ID:      partyID,
				PartyID: partyID,
				Content: input.Content,
			}, nil
		},
		listNotesByPartyIDFunc: func(ctx context.Context, pID pgtype.UUID) ([]db.PartyNote, error) {
			return []db.PartyNote{
				{ID: partyID, PartyID: pID, Content: "Note 1"},
			}, nil
		},
		deleteNoteFunc: func(ctx context.Context, id pgtype.UUID) error {
			return nil
		},
	}
	handler := NewNoteHandler(mockSvc)
	app.Post("/api/v1/parties/:party_id/notes", handler.CreateNote)
	app.Get("/api/v1/parties/:party_id/notes", handler.ListNotesByPartyID)
	app.Delete("/api/v1/notes/:note_id", handler.DeleteNote)

	// Create
	body, _ := json.Marshal(dto.CreateNoteRequest{
		Content: "Customer requested pricing catalog",
	})
	req := httptest.NewRequest("POST", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// List
	reqList := httptest.NewRequest("GET", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/notes", nil)
	respList, _ := app.Test(reqList)
	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respList.StatusCode)
	}

	// Delete
	reqDel := httptest.NewRequest("DELETE", "/api/v1/notes/550e8400-e29b-41d4-a716-446655440000", nil)
	respDel, _ := app.Test(reqDel)
	if respDel.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respDel.StatusCode)
	}
}
