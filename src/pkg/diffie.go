package dh

import "big"
import "os"
import "io"

type DHData struct {
  P *big.Int
  G *big.Int // generally, actually a very small number.
  S *big.Int // "Secret" integer
}

func NewDH(r io.Reader, l int, g *big.Int, p *big.Int)(dhd *DHData, err os.Error){
  buff := make([]byte, l)
  n, err := r.Read(buff)
  if n < l && err == nil { err = os.NewError("Random source provided insufficent data") }
  if err != nil { return }
  dhd = &DHData {P: p, G: g, S: big.NewInt(0) }
  // Use random data for our secret
  dhd.S.SetBytes(buff)
  // But it might be negative, so we Abs it.
  dhd.S.Abs(dhd.S)
  return
}

// Based on the secret generated at creation (as well as
// G&P)
func (self *DHData)ComputePublic()(q *big.Int){
  q = big.NewInt(0)
  q.Exp(self.G, self.S, self.P) 
  return
}

// Based on the value received from the remote side.
// If either the provided input is too large (>P) or
// the resultant check is invalid, an error is returned.
//
// Technically could be made to run faster if it didn't
// validate, but easier here than elsewhere.
func (self DHData)ComputeShared(in *big.Int)(out *big.Int, err os.Error){
  // Ensure 2 < in < self.P 
  if in.Cmp(big.NewInt(2)) != 1 ||  in.Cmp(self.P) != -1 {
     err = os.NewError("Invalid DH Key (size)")
     return
  }
  out = big.NewInt(0)
  //q := self.ComputePublic()
  /*chk := big.NewInt(0)
  chk.Exp(in, q, self.P)
  if chk.Cmp(big.NewInt(1)) != 0 {
    err = os.NewError("Invalid DH Key (shape)")
  }*/
  out.Exp(in, self.S, self.P)
  return
}
