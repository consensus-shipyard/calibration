package http

import (
	"errors"
	"html/template"
	"net/http"
	"path"

	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"

	"github.com/consensus-shipyard/calibration/faucet/internal/data"
	"github.com/consensus-shipyard/calibration/faucet/internal/faucet"
	"github.com/consensus-shipyard/calibration/faucet/internal/platform/web"
)

type FaucetWebService struct {
	log            *logging.ZapEventLogger
	faucet         *faucet.Service
	backendAddress string
}

func NewWebService(log *logging.ZapEventLogger, faucet *faucet.Service, backendAddress string) *FaucetWebService {
	return &FaucetWebService{
		log:            log,
		faucet:         faucet,
		backendAddress: backendAddress,
	}
}

func (h *FaucetWebService) handleFunds(w http.ResponseWriter, r *http.Request) {
	var req data.FundRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Address == "" {
		web.RespondError(w, http.StatusBadRequest, errors.New("empty address"))
		return
	}

	h.log.Infof(">>> %s -> {%s}\n", r.RemoteAddr, req.Address)

	targetAddr := common.HexToAddress(req.Address)

	err := h.faucet.FundAddress(r.Context(), targetAddr)
	if err != nil {
		h.log.Errorw("Failed to fund address", "addr", targetAddr, "err", err)
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FaucetWebService) handleHome(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./static/index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

func (h *FaucetWebService) handleScript(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./static/js/scripts.js")
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	if err = tmpl.Execute(w, h.backendAddress); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-type", "text/javascript")
}
