package application

import "errors"

var (
	ErrNaoEncontrado   = errors.New("registro não encontrado")
	ErrConflitoHorario = errors.New("Conflito de horário: já existe um veículo agendado para esta data e hora.")
	ErrDadosInvalidos  = errors.New("dados inválidos")
)
