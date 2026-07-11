package handlers

import (
	"encoding/json"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

type APIIntegrationHandler struct {
	serviceFuncionario *services.FuncionarioService
	serviceVisitante   *services.VisitanteService
	serviceTerceiro    *services.TerceiroService
}

func NewAPIIntegrationHandler(
	sf *services.FuncionarioService,
	sv *services.VisitanteService,
	st *services.TerceiroService,
) *APIIntegrationHandler {
	return &APIIntegrationHandler{
		serviceFuncionario: sf,
		serviceVisitante:   sv,
		serviceTerceiro:    st,
	}

}

// FuncionarioCadastroDTO define o formato limpo para o desenvolvedor externo
type FuncionarioCadastroDTO struct {
	Matricula       string `json:"matricula" example:"15420"`
	Nome            string `json:"nome" example:"Jefferson de Souza Martins"`
	Cargo           string `json:"cargo" example:"Analista de Infraestrutura"`
	Departamento    string `json:"departamento" example:"TI"`
	CentroDeCusto   string `json:"centro_de_custo" example:"TI-002"`
	Documento       string `json:"documento" example:"44411122299"`
	Credencial      string `json:"credencial" example:"998822"`
	Ativo           string `json:"ativo" example:"Sim"`
	GrupoRefeicaoID uint   `json:"grupo_refeicao_id" example:"2"`
	EquipamentosID  []uint `json:"equipamentos_id" example:"[1, 2]"`
}

// RespostaPadrao define o retorno JSON da sua API
type RespostaPadrao struct {
	Sucesso  bool   `json:"sucesso"`
	Mensagem string `json:"mensagem"`
	ID       uint   `json:"id,omitempty"`
}

// @Summary      Integrar Funcionário
// @Description  Recebe os dados cadastrais de um funcionário vindo do controle de acesso e grava no refeitório
// @Tags         Integração
// @Accept       json
// @Produce      json
// @Param        funcionario body models.Funcionario true "Dados do Funcionário"
// @Success      201  {object}  RespostaPadrao
// @Failure      400  {object}  RespostaPadrao
// @Failure      500  {object}  RespostaPadrao
// @Router       /api/integracao/funcionario [post]

// ReceberFuncionario trata o POST enviado pelo controle de acesso
func (h *APIIntegrationHandler) ReceberFuncionario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "Método não permitido."})
		return
	}

	// 1. Decodifica o formato limpo (DTO)
	var dto FuncionarioCadastroDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "JSON inválido: " + err.Error()})
		return
	}

	// 2. Mapeia os dados recebidos para a Struct real do Banco (Models)
	var empID uint = 1
	funcionario := models.Funcionario{
		Matricula:       dto.Matricula,
		Nome:            &dto.Nome,
		Cargo:           &dto.Cargo,
		Departamento:    &dto.Departamento,
		CentroDeCusto:   &dto.CentroDeCusto,
		Documento:       &dto.Documento,
		Credencial:      &dto.Credencial,
		Ativo:           &dto.Ativo,
		GrupoRefeicaoID: dto.GrupoRefeicaoID,
		EmpresaID:       &empID,
	}

	// 3. Monta a lista de ponteiros/structs virtuais de equipamentos com base nos IDs que vieram no array simples
	for _, idEquip := range dto.EquipamentosID {
		funcionario.Equipamentos = append(funcionario.Equipamentos, models.Equipamento{
			IDEquip: idEquip,
		})
	}

	// 4. Executa o service
	err := h.serviceFuncionario.Incluir(&funcionario)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "Erro ao salvar no refeitório: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: true, Mensagem: "Funcionário integrado com sucesso!", ID: funcionario.ID})
}

// ReceberVisitante trata a integração de visitantes
// @Summary      Integrar Visitante
// @Description  Recebe os dados de um visitante vindo do controle de acesso e grava no refeitório
// @Tags         Integração
// @Accept       json
// @Produce      json
// @Param        visitante body models.Visitante true "Dados do Visitante"
// @Success      201  {object}  RespostaPadrao
// @Failure      400  {object}  RespostaPadrao
// @Failure      500  {object}  RespostaPadrao
// @Router       /api/integracao/visitante [post]
// ReceberVisitante trata a integracao de visitantes
func (h *APIIntegrationHandler) ReceberVisitante(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "Metodo nao permitido."})
		return
	}

	var visitante models.Visitante

	if err := json.NewDecoder(r.Body).Decode(&visitante); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "JSON invalido: " + err.Error()})
		return
	}

	if visitante.EmpresaID == nil {
		var empID uint = 1
		visitante.EmpresaID = &empID
	}

	err := h.serviceVisitante.Incluir(&visitante)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "Erro ao salvar visitante: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: true, Mensagem: "Visitante integrado com sucesso!", ID: visitante.ID})

}

// ReceberTerceiro trata a integração de terceiros
// @Summary      Integrar Terceiro
// @Description  Recebe os dados de um terceiro vindo do controle de acesso e grava no refeitório
// @Tags         Integração
// @Accept       json
// @Produce      json
// @Param        terceiro body models.Terceiro true "Dados do Terceiro"
// @Success      201  {object}  RespostaPadrao
// @Failure      400  {object}  RespostaPadrao
// @Failure      500  {object}  RespostaPadrao
// @Router       /api/integracao/terceiro [post]
// ReceberTerceiro trata a integracao de terceiros
func (h *APIIntegrationHandler) ReceberTerceiro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "Metodo nao permitido."})
		return
	}

	var terceiro models.Terceiro
	if err := json.NewDecoder(r.Body).Decode(&terceiro); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "JSON invalido: " + err.Error()})
		return
	}

	if terceiro.EmpresaID == nil {
		var empID uint = 1
		terceiro.EmpresaID = &empID
	}

	err := h.serviceTerceiro.Incluir(&terceiro)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: false, Mensagem: "Erro ao salvar terceiro: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RespostaPadrao{Sucesso: true, Mensagem: "Terceiro integrado com sucesso!", ID: terceiro.ID})

}
