ECSRP
=====

This repo has a prototype of a variant of the
[SRP](https://en.wikipedia.org/wiki/Secure_Remote_Password_protocol) PAKE
protocol that allows the use of elliptic curve groups instead of a finite-field
group.

The original SRP protocol has the following outline:

* Registration:
  * Client chooses password, salt
  * `x = H(salt, password)`
  * `V = g^x`
  * `V` is registered with the server
* Authentication:
  * Client generates random `a` and `A = g^a`
  * Server generates random `b` and `B = k * V + g^b`, for some agreed `k`
  * Client and server exchange `A` and `B`, compute `u = H(A, B)`
  * Client computes `K = (B - k * g^x)^(a + u*x)`
  * Server computes `K = (A * V^u)^b

To turn this into an EC-based protocol, we start by replacing the
exponentiations with elliptic curve multiplications.  In this new framing `A`,
`B`, and `V` are elliptic curve points; other values are scalars in the
relevant finite field.

```
A = a * G
B = k * V + b * G
K = (a + u*x) * (B - k * x * G)
  = b * (A + u * V)
```

The confirmation steps of SRP can then proceed without modification.
