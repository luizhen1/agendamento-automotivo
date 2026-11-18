package domain

import "time"

type Cliente struct {
	ID            uint
	Nome          string
	Telefone      string
	VeiculoModelo string
	Placa         string
	CreatedAt     time.Time
}

type Servico struct {
	ID             uint
	Nome           string
	Descricao      string
	Valor          float64
	DuracaoMinutos int
	CreatedAt      time.Time
}

type StatusAgendamento string

const (
	StatusAgendado  StatusAgendamento = "AGENDADO"
	StatusConcluido StatusAgendamento = "CONCLUIDO"
	StatusCancelado StatusAgendamento = "CANCELADO"
)

func (s StatusAgendamento) Valido() bool {
	return s == StatusAgendado || s == StatusConcluido || s == StatusCancelado
}

type Agendamento struct {
	ID          uint
	ClienteID   uint
	Cliente     Cliente
	ServicoID   uint
	Servico     Servico
	DataHora    string
	Status      StatusAgendamento
	Observacoes string
	CreatedAt   time.Time
}

type AtualizacaoAgendamento struct {
	Status   *StatusAgendamento
	DataHora *string
}
