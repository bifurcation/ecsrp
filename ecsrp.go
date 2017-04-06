package ecsrp

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// XXX: Basically all errors are ignored
// XXX: The following are hard-wired, could be turned into an Options struct
//	 curve is P-256
//	 hash is SHA-256
//	 salt is 32 bytes long
//	 Multi-input hashes use concatenation
//	 hashes are interpreted as big-endian
//	 k = 3
func hash(vals ...[]byte) []byte {
	h := sha256.New()
	for _, val := range vals {
		h.Write(val)
	}
	return h.Sum(nil)
}

func randval() (out []byte) {
	out = make([]byte, 32)
	rand.Read(out)
	return
}

func bytes2int(x []byte) *big.Int {
	return big.NewInt(0).SetBytes(x)
}

// Returns 32 or more bytes
func int2bytes(x *big.Int) []byte {
	b := x.Bytes()
	if len(b) >= 32 {
		return b
	}

	n := 32 - len(b)
	fmt.Printf("n = %d\n", n)
	pad := make([]byte, n)
	return append(pad, b...)
}

var (
	curve = elliptic.P256()
	k = []byte{3}
)

func Register(p []byte) (s []byte, Vx, Vy *big.Int) {
	// salt = random
	s = randval()

	// x = H(s || p)
	x := hash(s, p)

	// V = x * G
	Vx, Vy = curve.ScalarBaseMult(x)

	return
}

// Returns a, A = a*G
func CreateChallenge() (a []byte, Ax, Ay *big.Int) {
	a = randval()
	Ax, Ay = curve.ScalarBaseMult(a)
	return
}

// Returns b, B = k*V + b*G
func CreateResponse(Vx, Vy *big.Int) (b []byte, Bx, By *big.Int) {
	b = randval()
	bGx, bGy := curve.ScalarBaseMult(b)
	kVx, kVy := curve.ScalarMult(Vx, Vy, k)
	Bx, By = curve.Add(kVx, kVy, bGx, bGy)
	return
}

// Returns S = (a + u*x) * (B - k * x * G)
func ProcessResponseClient(s, p, a []byte, Ax, Ay, Bx, By *big.Int) (Sx, Sy *big.Int) {
	x := hash(s, p)
	u := hash(int2bytes(Ax), int2bytes(Ay), int2bytes(Bx), int2bytes(By))

	// a + u * x
	ai := bytes2int(a)
	xi := bytes2int(x)
	ui := bytes2int(u)
	coeff := int2bytes(ai.Add(ai, ui.Mul(ui, xi)))

	// B - k * x * G
	xGx, xGy := curve.ScalarBaseMult(x)
	kxGx, kxGy := curve.ScalarMult(xGx, xGy, k)
	kxGy = kxGy.Neg(kxGy)
	Px, Py := curve.Add(Bx, By, kxGx, kxGy)

	Sx, Sy = curve.ScalarMult(Px, Py, coeff)
	return // TODO
}

// Returns S = b * (A + u * V)
func ProcessResponseServer(b []byte, Ax, Ay, Bx, By, Vx, Vy *big.Int) (Sx, Sy *big.Int) {
	u := hash(int2bytes(Ax), int2bytes(Ay), int2bytes(Bx), int2bytes(By))
	uVx, uVy := curve.ScalarMult(Vx, Vy, u)

	Px, Py := curve.Add(Ax, Ay, uVx, uVy)
	Sx, Sy = curve.ScalarMult(Px, Py, b)
	return
}
