package dh

import "big"
import "os"
import "io"

type DHData struct {
  P *big.Int
  G *big.Int // generally, actually a very small number.
  s *big.Int // "Secret" integer
  public *big.Int
}

/* NB: L is in BYTES not BITS.  You really don't want to wait for a
2048 BYTE DH keyex.  */
func NewDH(r io.Reader, l int, g *big.Int, p *big.Int)(dhd *DHData, err os.Error){
  buff := make([]byte, l)
  n, err := r.Read(buff)
  if n < l && err == nil { err = os.NewError("Random source provided insufficent data") }
  if err != nil { return }
  dhd = &DHData {P: p, G: g, s: big.NewInt(0) }
  // Use random data for our secret
  dhd.s.SetBytes(buff)
  // But it might be negative, so we Abs it.
  dhd.s.Abs(dhd.s)
  return
}

// Based on the secret generated at creation (as well as
// G&P)
func (self *DHData)ComputePublic()(q *big.Int){
  if self.public == nil {
    self.public = big.NewInt(0)
    self.public.Exp(self.G, self.s, self.P) 
  }
  q = self.public
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
  out.Exp(in, self.s, self.P)
  return
}
