package main

import (
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  "fmt"
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

//func main() {
//  signingMethodStr := env.GetDefaulted("TEMP_JWT_SIGNING_METHOD", "HS512")
//  signingMethod := jwt.GetSigningMethod(signingMethodStr)
//  secret, ok := os.LookupEnv("TEMP_JWT_SECRET_KEY")
//  if !ok {
//    fmt.Println("ERROR:TEMP_JWT_SECRET_KEY environment variable not set")
//    return
//  }
//
//  token, err := GenerateTemporaryToken(signingMethod, secret)
//  if err != nil {
//    fmt.Printf("ERROR:%s", err)
//    return
//  }
//
//  fmt.Println(token)
//}

func main() {
  token := GenerateTemporaryToken(jwt.SigningMethodRS256)
  privKey := OpenPrivKeyPem()
  publicKey := OpenPublicKeyPem()
  signedString, err := token.SignedString(privKey)
  if err != nil {
    log.Fatalln(err)
  }
  fmt.Println(signedString)

  tokenParsed, err := jwt.ParseWithClaims(signedString, &sharedJwt.TemporaryClaims{}, func(token *jwt.Token) (interface{}, error) {
    return publicKey, nil
  })
  if err != nil {
    log.Fatalln(err)
  }

  sharedUtil.DoNothing(tokenParsed)
}

func OpenPrivKeyPem() *rsa.PrivateKey {
  file, err := os.ReadFile("privkey.pem")
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
