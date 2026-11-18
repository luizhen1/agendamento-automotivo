package application

import (
	"context"
	"strings"

	"agendamento-automotivo/backend/internal/domain"
)

type ClienteService struct{ repo ClienteRepository }

func NewClienteService(repo ClienteRepository) *ClienteService { return &ClienteService{repo: repo} }

func (s *ClienteService) Listar(ctx context.Context) ([]domain.Cliente, error) {
	return s.repo.Listar(ctx)
}

func (s *ClienteService) Criar(ctx context.Context, cliente domain.Cliente) (domain.Cliente, error) {
	cliente.Nome = strings.TrimSpace(cliente.Nome)
	cliente.Telefone = strings.TrimSpace(cliente.Telefone)
	cliente.VeiculoModelo = strings.TrimSpace(cliente.VeiculoModelo)
	cliente.Placa = strings.TrimSpace(cliente.Placa)
	if cliente.Nome == "" || cliente.Telefone == "" || cliente.VeiculoModelo == "" || cliente.Placa == "" {
		return domain.Cliente{}, ErrDadosInvalidos
	}
	if err := s.repo.Criar(ctx, &cliente); err != nil {
		return domain.Cliente{}, err
	}
	return cliente, nil
}
