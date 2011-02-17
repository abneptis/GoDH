package dh

// Copyright(c) 2011 Abneptis LLC, all rights reserved.
// Original Author: James D. Nurmi <james@abneptis.com>

import "container/vector"
import "os"
import "time"

type SimpleDiffieHellmanCache interface {
  GetCurrent(int64)(*DHData, int64, os.Error)
}

type cachedRecord struct {
  Expires int64
  data *DHData
}


type dhCache struct {
  generator DiffieHellmanGenerator
  DefaultExpiration int64
  cache *vector.Vector
}

func NewDHCache(g DiffieHellmanGenerator, exp int64)(*dhCache){
  return &dhCache {
    generator: g,
    DefaultExpiration: exp,
    cache: &vector.Vector{},
  }
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
  exp = time.Nanoseconds() + self.DefaultExpiration
  cr := cachedRecord {
    Expires: exp,
  }
  d, err = self.generator.GenerateNew()
  if err == nil {
    cr.data = d
    self.cache.Push(cr)
  }
  return
}
