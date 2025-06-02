package lang

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// HttpHandler конструює обробник HTTP запитів, який дані з запиту віддає у Parser, а потім відправляє отриманий список
// операцій у painter.Loop.
func HttpHandler(loop *painter.Loop, p *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			cmd := r.URL.Query().Get("cmd")
			log.Printf("Received command: %s", cmd)
			in = strings.NewReader(cmd)
		}

		cmds, err := p.Parse(in)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("Posting commands to loop")
		for _, cmd := range cmds {
			loop.Post(cmd)
		}

		rw.WriteHeader(http.StatusOK)
	})
}
