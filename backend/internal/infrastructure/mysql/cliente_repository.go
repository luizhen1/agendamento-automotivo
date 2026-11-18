package mysql

import (
	"context"

	"agendamento-automotivo/backend/internal/application"
	"agendamento-automotivo/backend/internal/domain"
	"gorm.io/gorm"
)

type ClienteRepository struct{ db *gorm.DB }

func NewClienteRepository(db *gorm.DB) *ClienteRepository { return &ClienteRepository{db: db} }

func (r *ClienteRepository) Listar(ctx context.Context) ([]domain.Cliente, error) {
	var rows []clienteModel
	if err := r.db.WithContext(ctx).Order("nome asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Cliente, 0, len(rows))
	for _, row := range rows {
		items = append(items, clienteToDomain(row))
	}
	return items, nil
}

func (r *ClienteRepository) Criar(ctx context.Context, item *domain.Cliente) error {
	row := clienteModel{Nome: item.Nome, Telefone: item.Telefone, VeiculoModelo: item.VeiculoModelo, Placa: item.Placa}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	*item = clienteToDomain(row)
	return nil
}

var _ application.ClienteRepository = (*ClienteRepository)(nil)
