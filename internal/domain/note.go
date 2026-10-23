package domain

type CreateNoteInput struct {
	PartyID   string
	NoteType  *string
	Content   string
	CreatedBy *string
}
