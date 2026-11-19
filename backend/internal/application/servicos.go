package application

import (
	"context"
	"strings"

	"agendamento-automotivo/backend/internal/domain"
)

type ServicoService struct{ repo ServicoRepository }

func NewServicoService(repo ServicoRepository) *ServicoService { return &ServicoService{repo: repo} }

func (s *ServicoService) Listar(ctx context.Context) ([]domain.Servico, error) {
	return s.repo.Listar(ctx)
}

func (s *ServicoService) Criar(ctx context.Context, servico domain.Servico) (domain.Servico, error) {
	servico.Nome = strings.TrimSpace(servico.Nome)
	servico.Descricao = strings.TrimSpace(servico.Descricao)
	if servico.Nome == "" || servico.Valor <= 0 || servico.DuracaoMinutos <= 0 {
		return domain.Servico{}, ErrDadosInvalidos
	}
	if err := s.repo.Criar(ctx, &servico); err != nil {
		return domain.Servico{}, err
	}
	return servico, nil
}

func (s *ServicoService) Excluir(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrDadosInvalidos
	}
	return s.repo.Excluir(ctx, id)
}
