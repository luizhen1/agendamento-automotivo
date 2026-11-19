package application

import (
	"context"
	"strings"

	"agendamento-automotivo/backend/internal/domain"
)

type AgendamentoService struct{ repo AgendamentoRepository }

func NewAgendamentoService(repo AgendamentoRepository) *AgendamentoService {
	return &AgendamentoService{repo: repo}
}

func (s *AgendamentoService) Listar(ctx context.Context) ([]domain.Agendamento, error) {
	return s.repo.Listar(ctx)
}

func (s *AgendamentoService) Criar(ctx context.Context, clienteID, servicoID uint, dataHora, observacoes string) (domain.Agendamento, error) {
	dataHora, observacoes = strings.TrimSpace(dataHora), strings.TrimSpace(observacoes)
	if clienteID == 0 || servicoID == 0 || dataHora == "" {
		return domain.Agendamento{}, ErrDadosInvalidos
	}
	item := domain.Agendamento{ClienteID: clienteID, ServicoID: servicoID, DataHora: dataHora, Observacoes: observacoes, Status: domain.StatusAgendado}
	if err := s.repo.CriarSeHorarioLivre(ctx, &item); err != nil {
		return domain.Agendamento{}, err
	}
	return item, nil
}

func (s *AgendamentoService) Atualizar(ctx context.Context, id uint, status *string, dataHora *string, observacoes *string) (domain.Agendamento, error) {
	if id == 0 || (status == nil && dataHora == nil && observacoes == nil) {
		return domain.Agendamento{}, ErrDadosInvalidos
	}
	changes := domain.AtualizacaoAgendamento{}
	if status != nil {
		value := domain.StatusAgendamento(strings.ToUpper(strings.TrimSpace(*status)))
		if !value.Valido() {
			return domain.Agendamento{}, ErrDadosInvalidos
		}
		changes.Status = &value
	}
	if dataHora != nil {
		value := strings.TrimSpace(*dataHora)
		if value == "" {
			return domain.Agendamento{}, ErrDadosInvalidos
		}
		changes.DataHora = &value
	}
	if observacoes != nil {
		value := strings.TrimSpace(*observacoes)
		changes.Observacoes = &value
	}
	return s.repo.Atualizar(ctx, id, changes)
}

func (s *AgendamentoService) Excluir(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrDadosInvalidos
	}
	return s.repo.Excluir(ctx, id)
}
