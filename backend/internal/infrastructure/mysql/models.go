package mysql

import (
	"time"

	"agendamento-automotivo/backend/internal/domain"
)

type clienteModel struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Nome          string `gorm:"size:120;not null"`
	Telefone      string `gorm:"size:20;not null"`
	VeiculoModelo string `gorm:"size:80;not null"`
	Placa         string `gorm:"size:15;not null"`
	CreatedAt     time.Time
}

func (clienteModel) TableName() string { return "clientes" }

type servicoModel struct {
	ID             uint    `gorm:"primaryKey;autoIncrement"`
	Nome           string  `gorm:"size:100;not null"`
	Descricao      string  `gorm:"size:255"`
	Valor          float64 `gorm:"type:decimal(10,2);not null"`
	DuracaoMinutos int     `gorm:"not null"`
	CreatedAt      time.Time
}

func (servicoModel) TableName() string { return "servicos" }

type agendamentoModel struct {
	ID          uint         `gorm:"primaryKey;autoIncrement"`
	ClienteID   uint         `gorm:"not null;index"`
	Cliente     clienteModel `gorm:"foreignKey:ClienteID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ServicoID   uint         `gorm:"not null;index"`
	Servico     servicoModel `gorm:"foreignKey:ServicoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DataHora    string       `gorm:"size:30;not null;index"`
	Status      string       `gorm:"size:20;not null;default:AGENDADO"`
	Observacoes string       `gorm:"size:255"`
	CreatedAt   time.Time
}

func (agendamentoModel) TableName() string { return "agendamentos" }

func clienteToDomain(item clienteModel) domain.Cliente {
	return domain.Cliente{ID: item.ID, Nome: item.Nome, Telefone: item.Telefone, VeiculoModelo: item.VeiculoModelo, Placa: item.Placa, CreatedAt: item.CreatedAt}
}

func servicoToDomain(item servicoModel) domain.Servico {
	return domain.Servico{ID: item.ID, Nome: item.Nome, Descricao: item.Descricao, Valor: item.Valor, DuracaoMinutos: item.DuracaoMinutos, CreatedAt: item.CreatedAt}
}

func agendamentoToDomain(item agendamentoModel) domain.Agendamento {
	return domain.Agendamento{ID: item.ID, ClienteID: item.ClienteID, Cliente: clienteToDomain(item.Cliente), ServicoID: item.ServicoID, Servico: servicoToDomain(item.Servico), DataHora: item.DataHora, Status: domain.StatusAgendamento(item.Status), Observacoes: item.Observacoes, CreatedAt: item.CreatedAt}
}
