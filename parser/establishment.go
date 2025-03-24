package parser

type EstablishmentData struct {
	persist func(map[uint64]any) error

	CNPJBasico            string // Número base de inscrição no CNPJ (oito primeiros dígitos do CNPJ)
	CNPJOrdem             string // Número do estabelecimento de inscrição no CNPJ (dígitos 9 a 12)
	CNPJDV                string // Dígito verificador do CNPJ (dois últimos dígitos)
	IdentificadorMatriz   string // Código identificador: 1 - Matriz, 2 - Filial
	NomeFantasia          string // Nome fantasia do estabelecimento
	SituacaoCadastral     string // Código da situação cadastral: 01 - Nula, 2 - Ativa, etc.
	DataSituacaoCadastral string // Data do evento da situação cadastral
	MotivoSituacao        string // Código do motivo da situação cadastral
	CidadeExterior        string // Nome da cidade no exterior, se aplicável
	Pais                  string // Código do país
	DataInicioAtividade   string // Data de início da atividade
	CNAEPrincipal         string // Código da atividade econômica principal do estabelecimento
	CNAESecundaria        string // Códigos das atividades econômicas secundárias
	TipoLogradouro        string // Descrição do tipo de logradouro
	Logradouro            string // Nome do logradouro onde está o estabelecimento
	Numero                string // Número do endereço, ou "S/N" se não houver
	Complemento           string // Complemento do endereço
	Bairro                string // Bairro onde está localizado o estabelecimento
	CEP                   string // Código de endereço postal (CEP)
	UF                    string // Unidade da federação (estado)
	Municipio             string // Código do município
	DDD1                  string // Código DDD do telefone 1
	Telefone1             string // Número do telefone 1
	DDD2                  string // Código DDD do telefone 2
	Telefone2             string // Número do telefone 2
	DDDFax                string // Código DDD do fax
	Fax                   string // Número do fax
	Email                 string // E-mail do contribuinte
	SituacaoEspecial      string // Situação especial da empresa
	DataSituacaoEspecial  string // Data em que a empresa entrou em situação especial
}

var _ Parser = (*EstablishmentData)(nil)

func NewEstablishment(persist func(map[uint64]any) error) Parser {
	return &EstablishmentData{persist: persist}
}

func (d *EstablishmentData) ID(reading []byte, indexes []int) uint64 {
	return hash(reading[indexes[0]:indexes[1]])
}

func (d *EstablishmentData) Size() int {
	return 60
}

func (d *EstablishmentData) Data(reading []byte, indexes []int) Parser {
	return &EstablishmentData{
		CNPJBasico:            string(reading[indexes[0]:indexes[1]]),
		CNPJOrdem:             string(reading[indexes[2]:indexes[3]]),
		CNPJDV:                string(reading[indexes[4]:indexes[5]]),
		IdentificadorMatriz:   string(reading[indexes[6]:indexes[7]]),
		NomeFantasia:          string(reading[indexes[8]:indexes[9]]),
		SituacaoCadastral:     string(reading[indexes[10]:indexes[11]]),
		DataSituacaoCadastral: string(reading[indexes[12]:indexes[13]]),
		MotivoSituacao:        string(reading[indexes[14]:indexes[15]]),
		CidadeExterior:        string(reading[indexes[16]:indexes[17]]),
		Pais:                  string(reading[indexes[18]:indexes[19]]),
		DataInicioAtividade:   string(reading[indexes[20]:indexes[21]]),
		CNAEPrincipal:         string(reading[indexes[22]:indexes[23]]),
		CNAESecundaria:        string(reading[indexes[24]:indexes[25]]),
		TipoLogradouro:        string(reading[indexes[26]:indexes[27]]),
		Logradouro:            string(reading[indexes[28]:indexes[29]]),
		Numero:                string(reading[indexes[30]:indexes[31]]),
		Complemento:           string(reading[indexes[32]:indexes[33]]),
		Bairro:                string(reading[indexes[34]:indexes[35]]),
		CEP:                   string(reading[indexes[36]:indexes[37]]),
		UF:                    string(reading[indexes[38]:indexes[39]]),
		Municipio:             string(reading[indexes[40]:indexes[41]]),
		DDD1:                  string(reading[indexes[42]:indexes[43]]),
		Telefone1:             string(reading[indexes[44]:indexes[45]]),
		DDD2:                  string(reading[indexes[46]:indexes[47]]),
		Telefone2:             string(reading[indexes[48]:indexes[49]]),
		DDDFax:                string(reading[indexes[50]:indexes[51]]),
		Fax:                   string(reading[indexes[52]:indexes[53]]),
		Email:                 string(reading[indexes[54]:indexes[55]]),
		SituacaoEspecial:      string(reading[indexes[56]:indexes[57]]),
		DataSituacaoEspecial:  string(reading[indexes[58]:indexes[59]]),
	}
}

func (d *EstablishmentData) Persist(data map[uint64]any) error {
	return d.persist(data)
}
