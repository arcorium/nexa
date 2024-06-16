package main

import (
  "context"
  "encoding/json"
  "log"
  "net/http"
  "os"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
  st := r.URL.Query().Get("state")
  if st != state {
    log.Println("state token response has different value")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  code := r.URL.Query().Get("code")
  if len(code) == 0 {
    code = r.FormValue("code")
    if len(code) == 0 {
      log.Println("Failed to obtain code")
      w.WriteHeader(http.StatusBadRequest)
      return
    }
  }

  tok, err := cfgs.Exchange(context.Background(), code)
  // Save token
  if err != nil {
    log.Println("Failed to exchange token")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  f, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
  if err != nil {
    log.Println("Failed to create token file")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  defer f.Close()

  err = json.NewEncoder(f).Encode(tok)
  if err != nil {
    log.Println("Failed to encode token to file")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.Write([]byte(`
  <!DOCTYPE html>
  <html>
  <head>
  </head>
  <body>
  
  <h1>You can Close this Browser</h1>
  
  </body>
  </html> 
  `))
  w.WriteHeader(http.StatusOK)
}
