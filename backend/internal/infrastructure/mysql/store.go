package mysql

import (
	"context"

	"gorm.io/gorm"
)

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&clienteModel{}, &servicoModel{}, &agendamentoModel{})
}

func (s *Store) Seed(ctx context.Context) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&clienteModel{}).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			items := []clienteModel{
				{Nome: "Ana Souza", Telefone: "(11) 99999-1001", VeiculoModelo: "Honda Civic", Placa: "ABC1D23"},
				{Nome: "Bruno Lima", Telefone: "(11) 99999-1002", VeiculoModelo: "Toyota Corolla", Placa: "DEF4G56"},
				{Nome: "Carla Mendes", Telefone: "(11) 99999-1003", VeiculoModelo: "Jeep Renegade", Placa: "GHI7J89"},
			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&servicoModel{}).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			items := []servicoModel{
				{Nome: "Insulfilm Nano Cerâmica", Descricao: "Aplicação de película nano cerâmica", Valor: 650, DuracaoMinutos: 180},
				{Nome: "Polimento Cristalizado", Descricao: "Polimento técnico com proteção", Valor: 450, DuracaoMinutos: 240},
				{Nome: "Higienização Interna", Descricao: "Limpeza profunda de bancos e interior", Valor: 320, DuracaoMinutos: 180},
			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
