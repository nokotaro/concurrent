package association

import (
    "io"
    "fmt"
    "log"
    "bytes"
    "net/http"
    "encoding/json"
)

type AssociationHandler struct {
    service AssociationService
}

func NewAssociationHandler(service AssociationService) AssociationHandler {
    return AssociationHandler{service: service}
}

func (h AssociationHandler) Handle(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Headers", "*")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
    switch r.Method {
        case http.MethodGet:
            var author = r.URL.Query().Get("author")
            associations := h.service.GetOwn(author)
            jsonstr, err := json.Marshal(associations)
            if err != nil {
                log.Fatalf("getCharacters json.Marshal error:%v", err)
            }
            fmt.Fprint(w, string(jsonstr))
        case http.MethodPost:
            body := r.Body
            defer body.Close()

            buf := new(bytes.Buffer)
            io.Copy(buf, body)

            var association Association
            json.Unmarshal(buf.Bytes(), &association)

            h.service.PostAssociation(association)
            fmt.Fprintf(w, "{\"message\": \"accept\"}")
        case http.MethodDelete:
            body := r.Body
            defer body.Close()

            buf := new(bytes.Buffer)
            io.Copy(buf, body)

            var request deleteQuery
            json.Unmarshal(buf.Bytes(), &request)

            h.service.Delete(request.Id)
            fmt.Fprintf(w, "{\"message\": \"accept\"}")
        case http.MethodOptions:
            return
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            fmt.Fprint(w, "Method not allowed.")
    }
}

