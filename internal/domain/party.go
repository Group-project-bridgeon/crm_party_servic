package domain

type PartyType string

const (
	PartyTypeOrganization PartyType = "ORGANIZATION"
	PartyTypeIndividual   PartyType = "INDIVIDUAL"
)

func (p PartyType) IsValid() bool {
	switch p {
	case PartyTypeOrganization, PartyTypeIndividual:
		return true
	default:
		return false
	}
}
