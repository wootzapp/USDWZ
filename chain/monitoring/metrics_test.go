package monitoring

import "testing"

func TestSetCollateral(t *testing.T) {
	SetCollateral(10)
	if collateralGauge.Desc() == nil {
		t.Fatal("gauge not initialized")
	}
}
