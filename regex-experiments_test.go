package commuter

import (
  "testing"
  "math"
)

func TestRewriteGPX(*testing.T) {
  var v float64
  v = math.Average([]float64{1,2})
  if v != 1.5 {
    t.Error("expected 1.5, got ", v)
  }
}
