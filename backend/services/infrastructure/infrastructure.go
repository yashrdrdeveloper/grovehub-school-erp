package infrastructure

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	infrastructurerepositoryimpl "grovehub/services/infrastructure/repository/impl"
	infrastructureservices "grovehub/services/infrastructure/services"
	infrastructureserviceimpl "grovehub/services/infrastructure/services/impl"
)

const ensureDKIMPath = "/internal/infrastructure/mail/dkim/ensure"

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
	Client() MailClient
}

type module struct {
	logger      *slog.Logger
	client      MailClient
	dkimService infrastructureservices.DKIMService
}

type ensureDKIMRequest struct {
	Domain   string `json:"domain"`
	Selector string `json:"selector"`
}

type apiResponse struct {
	Success   bool                                 `json:"success"`
	Data      *infrastructureservices.DKIMIdentity `json:"data,omitempty"`
	ErrorCode string                               `json:"errorCode,omitempty"`
	Message   string                               `json:"message,omitempty"`
}

func New(logger *slog.Logger, db *gorm.DB) Module {
	dkimRepository := infrastructurerepositoryimpl.NewDKIMRepository(db)
	dkimService := infrastructureserviceimpl.NewDKIMService(dkimRepository)
	smtpService := infrastructureserviceimpl.NewSMTPServiceFromEnvironment()

	return &module{
		logger:      logger,
		client:      NewMailClient(smtpService),
		dkimService: dkimService,
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST "+ensureDKIMPath, service.ensureDKIM)
	service.logger.Info("infrastructure service registered")
}

func (service *module) Client() MailClient { return service.client }

func (service *module) ensureDKIM(writer http.ResponseWriter, request *http.Request) {
	var input ensureDKIMRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, apiResponse{
			Success: false, ErrorCode: "INFRASTRUCTURE_INVALID_DKIM_REQUEST", Message: "invalid DKIM request",
		})
		return
	}

	identity, err := service.dkimService.Ensure(request.Context(), input.Domain, input.Selector)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INFRASTRUCTURE_DKIM_PREPARE_FAILED"
		message := "unable to prepare DKIM identity"
		switch {
		case errors.Is(err, infrastructureservices.ErrInvalidDKIMIdentity):
			status = http.StatusBadRequest
			code = "INFRASTRUCTURE_INVALID_DKIM_IDENTITY"
			message = err.Error()
		case errors.Is(err, infrastructureservices.ErrDKIMSelectorConflict):
			status = http.StatusConflict
			code = "INFRASTRUCTURE_DKIM_SELECTOR_CONFLICT"
			message = err.Error()
		}
		writeJSON(writer, status, apiResponse{Success: false, ErrorCode: code, Message: message})
		return
	}

	writeJSON(writer, http.StatusOK, apiResponse{Success: true, Data: &identity})
}

func writeJSON(writer http.ResponseWriter, status int, response apiResponse) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(response)
}
