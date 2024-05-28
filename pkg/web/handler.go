package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ametow/getblock/pkg"

	"log"
)

type Handler struct {
	service *pkg.Service
	context context.Context
}

func NewHandlerContext(s *pkg.Service, ctx context.Context) *Handler {
	return &Handler{s, ctx}
}

func (c *Handler) Get(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.Run(c.context)
	if err != nil {
		log.Println("error while running:", err)
		fmt.Fprintln(w, "Server error ((")
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}
