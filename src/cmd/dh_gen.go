package main

import "dh"
import "logger"

import "big"
import "crypto/rand"
import "flag"
import "os"

var dhLen = flag.Int("dh-len", 2048, "Length of secret to generate")
var dhSecret = flag.String("dh-secret", "", "Secret to use (numeric)")
var dhGroup = flag.Int64("dh-group", 5, "DH Group to use")
var dhPrime = flag.String("dh-prime", "0xDCF93A0B883972EC0E19989AC5A2CE310E1D37717E8D9571BB7623731866E61EF75A2E27898B057F9891C2E27A639C3F29B60814581CD3B2CA3986D2683705577D45C2E7E52DC81C7A171876E5CEA74B1448BFDFAF18828EFD2519F14E45E3826634AF1949E5B535CC829A483B8A76223E5D490A257F05BDFF16F2FB22C583AB", "DH Prime to use")
var dhLHS = flag.String("lhs", "", "LHS")


func main(){
  flag.Parse()
  prime := big.NewInt(0)
  prime, _ = prime.SetString(*dhPrime, 0)
  
  g := big.NewInt(*dhGroup)
  logger.Log("DHGenStart:Prime", prime.String())
  logger.Log("DHGenStart:Group", g.String())
  mydh, err := dh.NewDH(rand.Reader, (*dhLen + 7)/8, g, prime)
  if *dhSecret != "" {
    mydh.S.SetString(*dhSecret, 0)
  }
  logger.Log("DHGenStart:Secret", mydh.S.String())
  if err != nil {
    logger.Log("DHGenFailure", err)
    os.Exit(1)
  }
  logger.Log("DHPublic", mydh.ComputePublic().String())
  if *dhLHS != "" {
    c, ok := big.NewInt(0).SetString(*dhLHS, 0)
    if !ok {
      logger.Log("Invalid LHS Key", nil)
      os.Exit(1)
    }
    shared, err :=mydh.ComputeShared(c)
    if err != nil {
      logger.Log("Invalid LHS Computation", err)
      os.Exit(1)
    }
    logger.Log("DHShared", shared.String())
  }
}
