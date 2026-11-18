package mysql

import (
	"context"
	"errors"

	"agendamento-automotivo/backend/internal/application"
	"agendamento-automotivo/backend/internal/domain"
	"gorm.io/gorm"
)

type AgendamentoRepository struct{ db *gorm.DB }

func NewAgendamentoRepository(db *gorm.DB) *AgendamentoRepository {
	return &AgendamentoRepository{db: db}
}

func (r *AgendamentoRepository) Listar(ctx context.Context) ([]domain.Agendamento, error) {
	var rows []agendamentoModel
	if err := r.db.WithContext(ctx).Preload("Cliente").Preload("Servico").Order("data_hora asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Agendamento, 0, len(rows))
	for _, row := range rows {
		items = append(items, agendamentoToDomain(row))
	}
	return items, nil
}

func (r *AgendamentoRepository) CriarSeHorarioLivre(ctx context.Context, item *domain.Agendamento) error {
	row := agendamentoModel{ClienteID: item.ClienteID, ServicoID: item.ServicoID, DataHora: item.DataHora, Status: string(item.Status), Observacoes: item.Observacoes}
	lockName := "agendamento:" + item.DataHora
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var acquired int
		if err := tx.Raw("SELECT GET_LOCK(?, 5)", lockName).Scan(&acquired).Error; err != nil {
			return err
		}
		if acquired != 1 {
			return errors.New("não foi possível obter bloqueio do horário")
		}
		defer tx.Exec("SELECT RELEASE_LOCK(?)", lockName)
		var count int64
		if err := tx.Model(&agendamentoModel{}).Where("data_hora = ? AND status <> ?", item.DataHora, string(domain.StatusCancelado)).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return application.ErrConflitoHorario
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		return tx.Preload("Cliente").Preload("Servico").First(&row, row.ID).Error
	})
	if err != nil {
		return err
	}
	*item = agendamentoToDomain(row)
	return nil
}

func (r *AgendamentoRepository) Atualizar(ctx context.Context, id uint, changes domain.AtualizacaoAgendamento) (domain.Agendamento, error) {
	var item agendamentoModel
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return domain.Agendamento{}, mapNotFound(err)
	}
	if changes.DataHora != nil && item.Status != string(domain.StatusCancelado) {
		ocupado, err := r.existeHorarioOcupado(ctx, *changes.DataHora, id)
		if err != nil {
			return domain.Agendamento{}, err
		}
		if ocupado {
			return domain.Agendamento{}, application.ErrConflitoHorario
		}
	}
	updates := make(map[string]interface{}, 2)
	if changes.Status != nil {
		updates["status"] = string(*changes.Status)
	}
	if changes.DataHora != nil {
		updates["data_hora"] = *changes.DataHora
	}
	if err := r.db.WithContext(ctx).Model(&item).Updates(updates).Error; err != nil {
		return domain.Agendamento{}, err
	}
	if err := r.db.WithContext(ctx).Preload("Cliente").Preload("Servico").First(&item, id).Error; err != nil {
		return domain.Agendamento{}, mapNotFound(err)
	}
	return agendamentoToDomain(item), nil
}

func (r *AgendamentoRepository) Excluir(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&agendamentoModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return application.ErrNaoEncontrado
	}
	return nil
}

func (r *AgendamentoRepository) existeHorarioOcupado(ctx context.Context, dataHora string, ignorarID uint) (bool, error) {
	query := r.db.WithContext(ctx).Model(&agendamentoModel{}).Where("data_hora = ? AND status <> ?", dataHora, string(domain.StatusCancelado))
	if ignorarID != 0 {
		query = query.Where("id <> ?", ignorarID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func mapNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return application.ErrNaoEncontrado
	}
	return err
}

var _ application.AgendamentoRepository = (*AgendamentoRepository)(nil)
