package metrics

import (
	"github.com/arl/statsviz"
	"net/http"
)

// Server 可视化实时监控 /debug/statsviz
func Server(addr string) error {
	mux := http.NewServeMux()
	if err := statsviz.Register(mux); err != nil {
		return err
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		return err
	}

	return nil
}
