package restricted

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/repository"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services"
	core "grovehub/libs/core/gateway"
)

const RenamePattern = "PATCH /school/restricted/registry/schools/{schoolId}/name"
const RenameScope = "internal.school.registry.rename"
const RenameIssuer = "internal-tool-school-rename"

func NewRename(db *sql.DB, raw, readRaw string) (http.HandlerFunc, error) {
	if strings.TrimSpace(raw) == "" {
		return func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(403) }, nil
	}
	if len(raw) > 256*1024 {
		return nil, errors.New("invalid rename grants")
	}
	var keys []core.WorkloadKeyConfig
	var extra any
	d := json.NewDecoder(strings.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&keys) != nil || d.Decode(&extra) != io.EOF || len(keys) == 0 || len(keys) > 100 {
		return nil, errors.New("invalid rename grants")
	}
	var readers []Grant
	if strings.TrimSpace(readRaw) != "" && json.Unmarshal([]byte(readRaw), &readers) != nil {
		return nil, errors.New("invalid directory grants")
	}
	public := map[string]bool{}
	ids := map[string]bool{}
	for _, r := range readers {
		public[r.Key.PublicKey] = true
		ids[r.Key.KeyID] = true
	}
	for _, k := range keys {
		if k.Issuer != RenameIssuer || !strings.HasPrefix(k.Subject, RenameIssuer+":") || !services.ValidID(strings.TrimPrefix(k.Subject, RenameIssuer+":")) || k.AllowDelegation || len(k.Scopes) != 1 || k.Scopes[0] != RenameScope || public[k.PublicKey] || ids[k.KeyID] {
			return nil, errors.New("invalid or reused rename credential")
		}
		public[k.PublicKey] = true
		ids[k.KeyID] = true
	}
	credentials, err := core.NewEd25519Credentials(keys)
	if err != nil {
		return nil, err
	}
	verifier := core.ScopedWorkloadVerifier{Credentials: credentials, Audience: "school", Routes: map[string]core.RoutePolicy{RenamePattern: {Scope: RenameScope}}}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		school := r.PathValue("schoolId")
		if !services.ValidID(school) || r.Method != "PATCH" || r.URL.Path != "/school/restricted/registry/schools/"+school+"/name" || r.URL.RawPath != "" || r.URL.RawQuery != "" || r.URL.ForceQuery || len(r.Header.Values("Cookie")) != 0 || len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
			w.WriteHeader(403)
			return
		}
		ctx, err := verifier.VerifyWorkload(r.Context(), r)
		if err != nil {
			w.WriteHeader(403)
			return
		}
		claims, ok := core.WorkloadFrom(ctx)
		if !ok || claims.Subject != RenameIssuer+":"+school || len(claims.Scopes) != 1 || claims.Scopes[0] != RenameScope || claims.ActorGID != "" || claims.TenantID != "" || claims.SessionID != "" {
			w.WriteHeader(403)
			return
		}
		keys := r.Header.Values("Idempotency-Key")
		if len(keys) != 1 || !services.ValidID(keys[0]) {
			w.WriteHeader(400)
			return
		}
		var in struct {
			CorrelationID   *string `json:"correlationId"`
			ExpectedVersion int64   `json:"expectedVersion"`
			ExpectedName    string  `json:"expectedName"`
			Name            string  `json:"name"`
		}
		d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
		d.DisallowUnknownFields()
		var extra any
		if d.Decode(&in) != nil || d.Decode(&extra) != io.EOF || (in.CorrelationID != nil && !services.ValidID(*in.CorrelationID)) || in.ExpectedVersion < 1 || !services.ValidName(in.ExpectedName) || !services.ValidName(in.Name) || in.ExpectedName == in.Name {
			w.WriteHeader(400)
			return
		}
		probe, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = repository.Rename(probe, db, school, keys[0], claims.Subject, in.ExpectedName, in.Name, in.ExpectedVersion, in.CorrelationID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			w.WriteHeader(404)
		case errors.Is(err, repository.ErrInactive):
			w.WriteHeader(403)
		case errors.Is(err, repository.ErrConflict):
			w.WriteHeader(409)
		case err != nil:
			w.WriteHeader(503)
		default:
			w.WriteHeader(204)
		}
	}, nil
}
