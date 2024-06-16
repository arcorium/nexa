package main

import (
  "fmt"
  "github.com/caarlos0/env/v10"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/drive/v3"
  "log"
  "net/http"
  "nexa/services/authentication/config"
  "nexa/services/file_storage/util"
  sharedEnv "nexa/shared/env"
  "strconv"
  "sync"
)

var state = "state-token"
var cfgs *oauth2.Config
var wg sync.WaitGroup

func main() {
  err := sharedEnv.LoadEnvs("dev.env")
  if err != nil {
    log.Fatalln(err)
  }

  var svConf config.GoogleOauthConfig
  err = env.Parse(&svConf)
  if err != nil {
    log.Fatalln(err)
  }
  if len(svConf.ServerAddress) == 0 {
    // random port
    port, err := util.GetAvailablePort()
    if err != nil {
      log.Fatalln("Failed to obtain port: ", err)
    }
    svConf.ServerAddress = "localhost:" + strconv.Itoa(port)
  }

  mux := http.NewServeMux()
  mux.HandleFunc(svConf.RedirectEndpoint, CallbackHandler)

  server := &http.Server{Handler: mux, Addr: svConf.ServerAddress}
  wg.Add(1)
  go func() {
    defer wg.Done()
    err := server.ListenAndServe()
    if err != nil {
      log.Println(err)
    }
  }()

  var cfg = &oauth2.Config{
    ClientID:     svConf.ClientID,
    ClientSecret: svConf.ClientSecret,
    Endpoint:     google.Endpoint,
    RedirectURL:  fmt.Sprintf("http://%s/%s", svConf.ServerAddress, svConf.RedirectEndpoint),
    Scopes:       []string{drive.DriveScope},
  }
  cfgs = cfg

  authUrl := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

  fmt.Printf("Go to following link:\n%s", authUrl)

  wg.Wait()
}
