package mysql

import (
	"context"
	"errors"

	"agendamento-automotivo/backend/internal/application"
	"agendamento-automotivo/backend/internal/domain"
	drivermysql "github.com/go-sql-driver/mysql"
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

func (r *ServicoRepository) Excluir(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&servicoModel{}, id)
	if result.Error != nil {
		var mysqlErr *drivermysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1451 {
			return application.ErrRegistroEmUso
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return application.ErrNaoEncontrado
	}
	return nil
}

var _ application.ServicoRepository = (*ServicoRepository)(nil)
