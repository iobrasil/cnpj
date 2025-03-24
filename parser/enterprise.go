package parser

type EnterpriseData struct {
	persist func(map[uint64]any) error

	BasicCNPJ                string
	CorporateName            string
	LegalNature              string
	ResponsibleQualification string
	SocialCapital            string
	CompanySize              string
	FederativeEntity         string
}

var _ Parser = (*EnterpriseData)(nil)

func NewEnterprise(persist func(map[uint64]any) error) Parser {
	return &EnterpriseData{persist: persist}
}

func (d *EnterpriseData) ID(reading []byte, indexes []int) uint64 {
	return hash(reading[indexes[0]:indexes[1]])
}

func (d *EnterpriseData) Size() int {
	return 14
}

func (d *EnterpriseData) Data(reading []byte, indexes []int) Parser {
	return &EnterpriseData{
		BasicCNPJ:                string(reading[indexes[0]:indexes[1]]),
		CorporateName:            string(reading[indexes[2]:indexes[3]]),
		LegalNature:              string(reading[indexes[4]:indexes[5]]),
		ResponsibleQualification: string(reading[indexes[6]:indexes[7]]),
		SocialCapital:            string(reading[indexes[8]:indexes[9]]),
		CompanySize:              string(reading[indexes[10]:indexes[11]]),
		FederativeEntity:         string(reading[indexes[12]:indexes[13]]),
	}
}

func (d *EnterpriseData) Persist(data map[uint64]any) error {
	return d.persist(data)
}
