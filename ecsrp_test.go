package ecsrp

import (
	"bytes"
	"testing"
)

func TestAll(t *testing.T) {
	password := []byte("hunter2")

	s, Vx, Vy := Register(password)

	a, Ax, Ay := CreateChallenge()
	b, Bx, By := CreateResponse(Vx, Vy)

	CSx, CSy := ProcessResponseClient(s, password, a, Ax, Ay, Bx, By)
	SSx, SSy := ProcessResponseServer(b, Ax, Ay, Bx, By, Vx, Vy)

	if Cx, Sx := CSx.Bytes(), SSx.Bytes(); !bytes.Equal(Cx, Sx) {
		t.Fatalf("Disagreement in X, [%x] != [%x]", Cx, Sx)
	}

	if Cy, Sy := CSy.Bytes(), SSy.Bytes(); !bytes.Equal(Cy, Sy) {
		t.Fatalf("Disagreement in Y, [%x] != [%x]", Cy, Sy)
	}
}
