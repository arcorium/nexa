package main

import (
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "github.com/arcorium/nexa/shared/env"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/golang-jwt/jwt/v5"
  "log"
  "os"
  "time"
)

const TOKEN_TTL = time.Minute * 10
const TOKEN_ISSUER = "nexa-token_generator"
const TOKEN_SUBJECT = "setup"

func GenerateTemporaryToken(method jwt.SigningMethod) *jwt.Token {
  token := jwt.NewWithClaims(method, &sharedJwt.TemporaryClaims{
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer:    TOKEN_ISSUER,
      Subject:   TOKEN_SUBJECT,
      Audience:  nil,
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_TTL)),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
      ID:        sharedUtil.RandomString(32),
    },
  })

  return token
}

func main() {
  signingMethodStr, ok := os.LookupEnv("TEMP_JWT_SIGNING_METHOD")
  signingMethod := jwt.GetSigningMethod(signingMethodStr)
  if !ok {
    signingMethod = sharedJwt.DefaultSigningMethod
  }

  privKeyPath := env.GetDefaulted("PRIVATE_KEY_PATH", "privkey.pem")
  token := GenerateTemporaryToken(signingMethod)
  privKey := OpenPrivKeyPem(privKeyPath)

  signedToken, err := token.SignedString(privKey)
  if err != nil {
    fmt.Printf("ERROR:%s", err)
    return
  }

  fmt.Printf(signedToken)
}

func OpenPrivKeyPem(path string) *rsa.PrivateKey {
  file, err := os.ReadFile(path)
  if err != nil {
    log.Fatalln(err)
  }

  p, rest := pem.Decode(file)
  sharedUtil.DoNothing(p, rest)
  privateKey, err := x509.ParsePKCS8PrivateKey(p.Bytes)
  if err != nil {
    log.Fatalln(err)
  }

  key := privateKey.(*rsa.PrivateKey)
  err = key.Validate()
  if err != nil {
    log.Fatalln(err)
  }

  return key
}

func OpenPublicKeyPem() *rsa.PublicKey {
  file, err := os.ReadFile("pubkey.pem")
  if err != nil {
    log.Fatalln(err)
  }

  p, rest := pem.Decode(file)
  sharedUtil.DoNothing(p, rest)
  publicKey, err := x509.ParsePKIXPublicKey(p.Bytes)
  if err != nil {
    log.Fatalln(err)
  }

  key := publicKey.(*rsa.PublicKey)
  sharedUtil.DoNothing(key)
  return key
}
