package mysql

import (
	"context"

	"agendamento-automotivo/backend/internal/domain"
	"gorm.io/gorm"
)

type ServicoRepository struct{ db *gorm.DB }

func NewServicoRepository(db *gorm.DB) *ServicoRepository { return &ServicoRepository{db: db} }

func (r *ServicoRepository) Listar(ctx context.Context) ([]domain.Servico, error) {
	var rows []servicoModel
	if err := r.db.WithContext(ctx).Order("nome asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Servico, 0, len(rows))
	for _, row := range rows {
		items = append(items, servicoToDomain(row))
	}
	return items, nil
}

func (r *ServicoRepository) Criar(ctx context.Context, item *domain.Servico) error {
	row := servicoModel{Nome: item.Nome, Descricao: item.Descricao, Valor: item.Valor, DuracaoMinutos: item.DuracaoMinutos}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	*item = servicoToDomain(row)
	return nil
}
