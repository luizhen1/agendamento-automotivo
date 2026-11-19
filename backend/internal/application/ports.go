package application

import (
	"context"

	"agendamento-automotivo/backend/internal/domain"
)

type ClienteRepository interface {
	Listar(context.Context) ([]domain.Cliente, error)
	Criar(context.Context, *domain.Cliente) error
}

type ServicoRepository interface {
	Listar(context.Context) ([]domain.Servico, error)
	Criar(context.Context, *domain.Servico) error
	Excluir(context.Context, uint) error
}

type AgendamentoRepository interface {
	Listar(context.Context) ([]domain.Agendamento, error)
	CriarSeHorarioLivre(context.Context, *domain.Agendamento) error
	Atualizar(context.Context, uint, domain.AtualizacaoAgendamento) (domain.Agendamento, error)
	Excluir(context.Context, uint) error
}
