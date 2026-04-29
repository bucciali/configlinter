package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"configlinter/internal/engine"
	"configlinter/internal/parser"
	"configlinter/internal/server/response"
)

type Handler struct {
	Registry *parser.Registry
	Engine   *engine.Engine
}

func New(reg *parser.Registry, eng *engine.Engine) *Handler {
	return &Handler{
		Registry: reg,
		Engine:   eng,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

type LintRequest struct {
	Content string `json:"content"`
	Format  string `json:"format"`
}

func (h *Handler) Lint(w http.ResponseWriter, r *http.Request) {
	var req LintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "cannot parse request body")
		return
	}
	defer r.Body.Close()

	if req.Content == "" {
		response.WriteError(w, http.StatusBadRequest, "empty content")
		return
	}

	format := req.Format
	if format == "" {
		format = "yaml"
	}

	p, err := h.Registry.GetByFormat(format)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "unsupported format")
		return
	}

	root, err := p.Parse([]byte(req.Content))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("parse error: %s", err))
		return
	}

	findings := h.Engine.Analyze(root)

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"issues_count": len(findings),
		"issues":       findings,
	})
}

func detectFormat(r *http.Request, body []byte) string {
	ct := r.Header.Get("Content-Type")
	switch {
	case strings.Contains(ct, "yaml"):
		return "yaml"
	case strings.Contains(ct, "toml"):
		return "toml"
	case strings.Contains(ct, "json"):
		return "json"
	}

	trimmed := strings.TrimSpace(string(body))
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return "json"
	}
	return "yaml"
}
