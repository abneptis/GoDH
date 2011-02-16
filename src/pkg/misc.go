package dh


import "big"
import "bytes"
import "encoding/base64"
import "os"


// Misc helper functions for creating B2 padded Bigs (often useful w/ DH exchanges)

/* Returns the appropriate base-two padded
bytes (assuming the underlying big representation
remains b2c) */
func BigBytes(i *big.Int)(buff []byte){
  ib := i.Bytes()
  shift := 0
  shiftbyte := byte(0)
  switch i.Cmp(big.NewInt(0)) {
    case 1:
      // Positive must be padded if high-bit is 1
      if ib[0] & 0x80 == 0x80 {
        shift = 1
      }
    case -1:
      // Negative numbers with a leading high-bit will also need
      // to be padded, but with a single bit tagging its 'negativity'
      if ib[0] & 0x80 == 0x80 {
        shift = 1
        shiftbyte = 0x80 
        
      }
  }
  buff = make([]byte, len(ib)+shift)
  if shift == 1 { buff[0] = shiftbyte }
  copy(buff[shift:], ib)
  
  return
}

/* Encode a BigInt to a base64 number;
  Also does appropriate b2 padding of the
  integer (if it's not negative)
*/
func BigIntToB64(i *big.Int)(string){
  b := BigBytes(i)
  buff := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
  base64.StdEncoding.Encode(buff, b)
  return string(buff)
}

func B64ToBigInt(in string, b *big.Int)(err os.Error){
  bsize := base64.StdEncoding.DecodedLen(len(in))
  buff := make([]byte, bsize)
  n, err := base64.StdEncoding.Decode(buff, bytes.NewBufferString(in).Bytes())
  neg := false
  if err == nil {
    buff = buff[0:n]
    if buff[0] & 0x80 == 0x80 {
      neg = true
      buff[0] &= 0x7f
    }
    b.SetBytes(buff)
    // In case the passed in big was negative...
    b.Abs(b)
    if neg { b.Neg(b) }
  }
  return
}
