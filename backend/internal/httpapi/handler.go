package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"agendamento-automotivo/backend/internal/application"
	"agendamento-automotivo/backend/internal/domain"
)

type Handler struct {
	clientes     *application.ClienteService
	servicos     *application.ServicoService
	agendamentos *application.AgendamentoService
}

func NewHandler(clientes *application.ClienteService, servicos *application.ServicoService, agendamentos *application.AgendamentoService) *Handler {
	return &Handler{clientes: clientes, servicos: servicos, agendamentos: agendamentos}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/clientes", h.clientesHandler)
	mux.HandleFunc("/servicos", h.servicosHandler)
	mux.HandleFunc("/servicos/", h.servicoIDHandler)
	mux.HandleFunc("/agendamentos", h.agendamentosHandler)
	mux.HandleFunc("/agendamentos/", h.agendamentoIDHandler)
	return cors(mux)
}

type clienteRequest struct {
	Nome          string `json:"nome"`
	Telefone      string `json:"telefone"`
	VeiculoModelo string `json:"veiculo_modelo"`
	Placa         string `json:"placa"`
}

type clienteResponse struct {
	ID            uint      `json:"id"`
	Nome          string    `json:"nome"`
	Telefone      string    `json:"telefone"`
	VeiculoModelo string    `json:"veiculo_modelo"`
	Placa         string    `json:"placa"`
	CreatedAt     time.Time `json:"created_at"`
}

type servicoRequest struct {
	Nome           string  `json:"nome"`
	Descricao      string  `json:"descricao"`
	Valor          float64 `json:"valor"`
	DuracaoMinutos int     `json:"duracao_minutos"`
}

type servicoResponse struct {
	ID             uint      `json:"id"`
	Nome           string    `json:"nome"`
	Descricao      string    `json:"descricao"`
	Valor          float64   `json:"valor"`
	DuracaoMinutos int       `json:"duracao_minutos"`
	CreatedAt      time.Time `json:"created_at"`
}

type agendamentoRequest struct {
	ClienteID   uint   `json:"cliente_id"`
	ServicoID   uint   `json:"servico_id"`
	DataHora    string `json:"data_hora"`
	Observacoes string `json:"observacoes"`
}

type agendamentoPatchRequest struct {
	Status      *string `json:"status"`
	DataHora    *string `json:"data_hora"`
	Observacoes *string `json:"observacoes"`
}

type agendamentoResponse struct {
	ID          uint            `json:"id"`
	ClienteID   uint            `json:"cliente_id"`
	Cliente     clienteResponse `json:"cliente"`
	ServicoID   uint            `json:"servico_id"`
	Servico     servicoResponse `json:"servico"`
	DataHora    string          `json:"data_hora"`
	Status      string          `json:"status"`
	Observacoes string          `json:"observacoes"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (h *Handler) clientesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.clientes.Listar(r.Context())
		if err != nil {
			writeInternalError(w, err)
			return
		}
		response := make([]clienteResponse, 0, len(items))
		for _, item := range items {
			response = append(response, toClienteResponse(item))
		}
		writeJSON(w, http.StatusOK, response)
	case http.MethodPost:
		var input clienteRequest
		if !decode(w, r, &input) {
			return
		}
		item, err := h.clientes.Criar(r.Context(), domain.Cliente{Nome: input.Nome, Telefone: input.Telefone, VeiculoModelo: input.VeiculoModelo, Placa: input.Placa})
		if err != nil {
			writeApplicationError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, toClienteResponse(item))
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) servicosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.servicos.Listar(r.Context())
		if err != nil {
			writeInternalError(w, err)
			return
		}
		response := make([]servicoResponse, 0, len(items))
		for _, item := range items {
			response = append(response, toServicoResponse(item))
		}
		writeJSON(w, http.StatusOK, response)
	case http.MethodPost:
		var input servicoRequest
		if !decode(w, r, &input) {
			return
		}
		item, err := h.servicos.Criar(r.Context(), domain.Servico{Nome: input.Nome, Descricao: input.Descricao, Valor: input.Valor, DuracaoMinutos: input.DuracaoMinutos})
		if err != nil {
			writeApplicationError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, toServicoResponse(item))
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) servicoIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		methodNotAllowed(w)
		return
	}
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/servicos/"), "/")
	id64, err := strconv.ParseUint(path, 10, 64)
	if err != nil || id64 == 0 {
		writeError(w, http.StatusBadRequest, "id de serviço inválido")
		return
	}
	if err := h.servicos.Excluir(r.Context(), uint(id64)); err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensagem": "Serviço excluído com sucesso"})
}

func (h *Handler) agendamentosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.agendamentos.Listar(r.Context())
		if err != nil {
			writeInternalError(w, err)
			return
		}
		response := make([]agendamentoResponse, 0, len(items))
		for _, item := range items {
			response = append(response, toAgendamentoResponse(item))
		}
		writeJSON(w, http.StatusOK, response)
	case http.MethodPost:
		var input agendamentoRequest
		if !decode(w, r, &input) {
			return
		}
		item, err := h.agendamentos.Criar(r.Context(), input.ClienteID, input.ServicoID, input.DataHora, input.Observacoes)
		if err != nil {
			writeApplicationError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, toAgendamentoResponse(item))
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) agendamentoIDHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/agendamentos/"), "/")
	id64, err := strconv.ParseUint(path, 10, 64)
	if err != nil || id64 == 0 {
		writeError(w, http.StatusBadRequest, "id de agendamento inválido")
		return
	}
	id := uint(id64)
	switch r.Method {
	case http.MethodPatch:
		var input agendamentoPatchRequest
		if !decode(w, r, &input) {
			return
		}
		item, err := h.agendamentos.Atualizar(r.Context(), id, input.Status, input.DataHora, input.Observacoes)
		if err != nil {
			writeApplicationError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, toAgendamentoResponse(item))
	case http.MethodDelete:
		if err := h.agendamentos.Excluir(r.Context(), id); err != nil {
			writeApplicationError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"mensagem": "Agendamento excluído com sucesso"})
	default:
		methodNotAllowed(w)
	}
}

func toClienteResponse(item domain.Cliente) clienteResponse {
	return clienteResponse{ID: item.ID, Nome: item.Nome, Telefone: item.Telefone, VeiculoModelo: item.VeiculoModelo, Placa: item.Placa, CreatedAt: item.CreatedAt}
}

func toServicoResponse(item domain.Servico) servicoResponse {
	return servicoResponse{ID: item.ID, Nome: item.Nome, Descricao: item.Descricao, Valor: item.Valor, DuracaoMinutos: item.DuracaoMinutos, CreatedAt: item.CreatedAt}
}

func toAgendamentoResponse(item domain.Agendamento) agendamentoResponse {
	return agendamentoResponse{ID: item.ID, ClienteID: item.ClienteID, Cliente: toClienteResponse(item.Cliente), ServicoID: item.ServicoID, Servico: toServicoResponse(item.Servico), DataHora: item.DataHora, Status: string(item.Status), Observacoes: item.Observacoes, CreatedAt: item.CreatedAt}
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return false
	}
	return true
}

func writeApplicationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrConflitoHorario):
		writeError(w, http.StatusConflict, application.ErrConflitoHorario.Error())
	case errors.Is(err, application.ErrNaoEncontrado):
		writeError(w, http.StatusNotFound, "agendamento não encontrado")
	case errors.Is(err, application.ErrDadosInvalidos):
		writeError(w, http.StatusBadRequest, "dados obrigatórios inválidos")
	case errors.Is(err, application.ErrRegistroEmUso):
		writeError(w, http.StatusConflict, "não é possível excluir o serviço porque há agendamentos vinculados")
	default:
		writeInternalError(w, err)
	}
}

func writeInternalError(w http.ResponseWriter, err error) {
	log.Printf("erro interno da API: %v", err)
	writeError(w, http.StatusInternalServerError, "erro interno do servidor")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"erro": message})
}

func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "método não permitido")
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
