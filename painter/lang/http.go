package lang

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/DmytroHalai/kpi-3/painter"
)

// HttpHandler конструює обробник HTTP запитів, який дані з запиту віддає у Parser, а потім відправляє отриманий список
// операцій у painter.Loop.
func HttpHandler(loop *painter.Loop, p *Parser, scene *painter.Scene) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			in = strings.NewReader(r.URL.Query().Get("cmd"))
		}

		cmds, err := p.Parse(in, scene)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		loop.Post(painter.OperationList(cmds))
		rw.WriteHeader(http.StatusOK)
	})
}
