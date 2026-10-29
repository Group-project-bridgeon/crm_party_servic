package dto

type CreateNoteRequest struct {
	NoteType  *string `json:"note_type"`
	Content   string  `json:"content"`
	CreatedBy *string `json:"created_by"`
}

type NoteResponse struct {
	ID        string  `json:"id"`
	PartyID   string  `json:"party_id"`
	NoteType  *string `json:"note_type,omitempty"`
	Content   string  `json:"content"`
	CreatedBy *string `json:"created_by,omitempty"`
	CreatedAt string  `json:"created_at"`
}
