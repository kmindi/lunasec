package controller

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/refinery-labs/loq/util/auth"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/refinery-labs/loq/service"
	"github.com/refinery-labs/loq/types"
	"github.com/refinery-labs/loq/types/event"
	"github.com/refinery-labs/loq/util"
)

type metaController struct {
	meta          service.MetadataService
	jwtVerifier service.JwtVerifier
	grant         service.GrantService
}

// MetaController ...
type MetaController interface {
	GetMetadata(w http.ResponseWriter, req *http.Request)
	SetMetadata(w http.ResponseWriter, req *http.Request)
}

// NewMetaController ...
func NewMetaController(meta service.MetadataService, jwtVerifier service.JwtVerifier, grant service.GrantService) MetaController {
	return &metaController{
		meta:          meta,
		jwtVerifier: jwtVerifier,
		grant:                     grant,
	}
}

// GetMetadata ...
func (s *metaController) GetMetadata(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received GetMetadata request")

	input := event.MetadataGetRequest{}
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(b, &input); err != nil {
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	meta, err := s.meta.GetMetadata(types.Token(input.TokenID))
	if err != nil {
		statusCode := 500
		if err.Error() == "unable to locate metadata for token" {
			statusCode = 404
		}
		util.RespondError(w, statusCode, err)
		return
	}

	resp := event.MetadataGetResponse{
		Metadata: meta.CustomMetadata,
	}

	util.Respond(w, resp)
}

// SetMetadata ...
func (s *metaController) SetMetadata(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received SetMetadata request")

	input := event.MetadataSetRequest{}
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(b, &input); err != nil {
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	claims, err := auth.GetRequestClaims(s.jwtVerifier, r)
	if err != nil {
		err = errors.Wrap(err, "unable to verify token jwt with claims")
		return
	}

	if err := s.meta.SetMetadata(types.Token(input.TokenID), claims, input.Metadata); err != nil {
		util.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	resp := event.MetadataSetResponse{}

	util.Respond(w, resp)
}
