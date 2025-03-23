package main

type EnterpriseData struct {
	BasicCNPJ                string
	CorporateName            string
	LegalNature              string
	ResponsibleQualification string
	SocialCapital            string
	CompanySize              string
	FederativeEntity         string
}

func (d *EnterpriseData) ID(reading []byte, indexes []int) uint64 {
	return hash(reading[indexes[0]:indexes[1]])
}

func (d *EnterpriseData) Size() int {
	return 14
}

func (d *EnterpriseData) Data(reading []byte, indexes []int) *EnterpriseData {
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
