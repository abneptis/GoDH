package dh

// Copyright(c) 2011 Abneptis LLC, all rights reserved.
// Original Author: James D. Nurmi <james@abneptis.com>

import "big"
import "container/vector"
import "io"
import "os"
import "time"

type SimpleDiffieHellmanCache interface {
  GenerateNew()(*DHData, int64, os.Error)
  GetCurrent(int64)(*DHData, int64, os.Error)
}

type cachedRecord struct {
  Expires int64
  data *DHData
}

type dhCache struct {
  DefaultLength int
  DefaultPrime *big.Int
  DefaultGroup *big.Int
  DefaultReader io.Reader
  DefaultExpiration int64
  cache *vector.Vector
}

func NewDHCache(l int, g *big.Int, p *big.Int, exp int64, r io.Reader)(*dhCache){
  return &dhCache {
    DefaultLength: l,
    DefaultPrime: p,
    DefaultGroup: g,
    DefaultReader: r,
    DefaultExpiration: exp,
    cache: &vector.Vector{},
  }
}

func (self *dhCache)GenerateNew()(dhdata *DHData, exp int64, err os.Error){
  dhdata, err = NewDH(self.DefaultReader, self.DefaultLength, self.DefaultGroup, self.DefaultPrime)
  if err == nil {
    cr := cachedRecord {
      Expires: time.Nanoseconds() + self.DefaultExpiration,
      data: dhdata,
    }
    dhdata.ComputePublic()
    self.cache.Push(cr)
    exp = cr.Expires
  }
  return
}

func (self *dhCache)GetCurrent(min int64)(d *DHData, exp int64, err os.Error){
  for i := range(*self.cache){
    now := time.Nanoseconds()
    cr := ((*self.cache)[i]).(cachedRecord)
    if cr.Expires > now + min {
      d = cr.data
      exp = cr.Expires
      return
    }
    if cr.Expires < now {
      self.cache.Delete(i)
    }
  }
  d, exp, err = self.GenerateNew()
  return
}
