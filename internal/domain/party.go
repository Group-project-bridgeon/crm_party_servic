package domain

type PartyType string

const (
	PartyTypeOrganization PartyType = "ORGANIZATION"
	PartyTypeIndividual   PartyType = "INDIVIDUAL"
)

func (p PartyType) IsValid() bool {
	return p == PartyTypeOrganization ||
		p == PartyTypeIndividual
}