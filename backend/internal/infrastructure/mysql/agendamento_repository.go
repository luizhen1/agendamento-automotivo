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
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		releaseLock, err := acquireSlotLock(tx, item.DataHora)
		if err != nil {
			return err
		}
		defer releaseLock()
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
	var updated agendamentoModel
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item agendamentoModel
		if err := tx.First(&item, id).Error; err != nil {
			return mapNotFound(err)
		}

		targetTime := item.DataHora
		if changes.DataHora != nil {
			targetTime = *changes.DataHora
		}
		targetStatus := item.Status
		if changes.Status != nil {
			targetStatus = string(*changes.Status)
		}

		if targetStatus != string(domain.StatusCancelado) {
			releaseLock, err := acquireSlotLock(tx, targetTime)
			if err != nil {
				return err
			}
			defer releaseLock()
			var count int64
			if err := tx.Model(&agendamentoModel{}).
				Where("id <> ? AND data_hora = ? AND status <> ?", id, targetTime, string(domain.StatusCancelado)).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return application.ErrConflitoHorario
			}
		}

		updates := make(map[string]interface{}, 3)
		if changes.Status != nil {
			updates["status"] = string(*changes.Status)
		}
		if changes.DataHora != nil {
			updates["data_hora"] = *changes.DataHora
		}
		if changes.Observacoes != nil {
			updates["observacoes"] = *changes.Observacoes
		}
		if err := tx.Model(&item).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Preload("Cliente").Preload("Servico").First(&updated, id).Error
	})
	if err != nil {
		return domain.Agendamento{}, err
	}
	return agendamentoToDomain(updated), nil
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

func acquireSlotLock(tx *gorm.DB, dataHora string) (func(), error) {
	lockName := "agendamento:" + dataHora
	var acquired int
	if err := tx.Raw("SELECT GET_LOCK(?, 5)", lockName).Scan(&acquired).Error; err != nil {
		return nil, err
	}
	if acquired != 1 {
		return nil, errors.New("não foi possível obter bloqueio do horário")
	}
	return func() {
		_ = tx.Exec("SELECT RELEASE_LOCK(?)", lockName).Error
	}, nil
}

func mapNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return application.ErrNaoEncontrado
	}
	return err
}

var _ application.AgendamentoRepository = (*AgendamentoRepository)(nil)
