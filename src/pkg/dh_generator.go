package dh

// Copyright(c) 2011 Abneptis LLC, all rights reserved.
// Original Author: James D. Nurmi <james@abneptis.com>

import "big"
import "crypto/rand"
import "io"
import "os"

type DiffieHellmanGenerator interface {
  GenerateNew()(*DHData, os.Error)
}

type dhGenerator struct {
  Prime *big.Int
  Group *big.Int
  Rand io.Reader
  SecretSize int
}

func NewGenerator(g, prime *big.Int, bits int)(DiffieHellmanGenerator){
  return dhGenerator {
    Prime: prime,
    Group: g,
    Rand: rand.Reader,
    SecretSize: (bits + 7)/8,
  }
}

func (self dhGenerator)GenerateNew()(*DHData, os.Error){
  return  NewDH(self.Rand, self.SecretSize, self.Group, self.Prime)
}
