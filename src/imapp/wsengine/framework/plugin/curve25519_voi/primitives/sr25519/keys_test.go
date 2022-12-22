// Copyright (c) 2017-2020 isis agora lovecruft. All rights reserved.
// Copyright (c) 2019-2020 Web 3 Foundation. All rights reserved.
// Copyright (c) 2021 Oasis Labs Inc. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// 1. Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
// TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package sr25519

import (
	"bytes"
	"crypto/rand"
	"testing"

	"ws/framework/plugin/curve25519_voi/curve/scalar"
	"ws/framework/plugin/curve25519_voi/pkg/testhelpers"
)

func isAllZeros(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func TestScalarDivideByCofactor(t *testing.T) {
	var b [scalar.ScalarSize]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("failed to read random scalar bytes: %v", err)
	}

	// Apply the "clamp", that always will be applied to the input
	// to the routine we are testing.
	b[0] &= 248
	b[31] &= 63
	b[31] |= 64

	sExpected, err := scalar.NewFromBytesModOrder(b[:])
	if err != nil {
		t.Fatalf("failed to deserialize random scalar: %v", err)
	}

	s, err := scalarDivideByCofactor(b[:])
	if err != nil {
		t.Fatalf("failed to divide by cofactor: %v", err)
	}

	s.Mul(s, scalar.NewFromUint64(8)) // Multiply by 8

	if sExpected.Equal(s) != 1 {
		t.Fatalf("(b / 8) * 8 != b (Got %v expected %v)", s, sExpected)
	}
}

func TestMiniSecretKey(t *testing.T) {
	t.Run("s11n", func(t *testing.T) {
		msk, err := GenerateMiniSecretKey(nil)
		if err != nil {
			t.Fatalf("GenerateMiniSecretKey: %v", err)
		}

		b, err := msk.MarshalBinary()
		if err != nil {
			t.Fatalf("msk.MarshalBinary: %v", err)
		}

		if !bytes.Equal(b, msk[:]) {
			t.Fatalf("msk.MarshalBinary() != msk (Got %v, %v)", b, msk)
		}

		var msk2 MiniSecretKey
		if err = msk2.UnmarshalBinary(b); err != nil {
			t.Fatalf("msk2.UnmarshalBinary(b): %v", err)
		}

		if msk2 != *msk {
			t.Fatalf("msk != msk2 (Got %v, %v)", *msk, msk2)
		}
	})
	t.Run("ExpandUniform", func(t *testing.T) {
		expected := testhelpers.MustUnhex(t, "04f0557e7f35e00df0824f458868915368bd5e41fd91f85b177f5907383ac50bdd0660b091e0ec47ecaf1f6ce73e7168fef267770f5030d5c524a49615163471063b66cc8b77aa24f694d073ad72c21a9f296be0fd4ee953d8e58d5d627d435b")
		var kp KeyPair
		if err := kp.UnmarshalBinary(expected); err != nil {
			t.Fatalf("kp.UnmarshalBinary: %v", err)
		}

		var msk MiniSecretKey
		sk := msk.ExpandUniform()

		if !kp.SecretKey().Equal(sk) {
			t.Fatalf("sk != expected (Got %v)", sk)
		}
		if pk := sk.PublicKey(); !kp.PublicKey().Equal(pk) {
			t.Fatalf("sk.PublicKey() != expected (Got %v)", pk)
		}
	})
	t.Run("ExpandEd25519", func(t *testing.T) {
		expected := testhelpers.MustUnhex(t, "caa835781b15c7706f65b71f7a58c807ab360faed6440fb23e0f4c52e930de0a0a6a85eaa642dac835424b5d7c8d637c00408c7a73da672b7f498521420b6dd3def12e42f3e487e9b14095aa8d5cc16a33491f1b50dadcf8811d1480f3fa8627")
		var kp KeyPair
		if err := kp.UnmarshalBinary(expected); err != nil {
			t.Fatalf("kp.UnmarshalBinary: %v", err)
		}

		var msk MiniSecretKey
		sk := msk.ExpandEd25519()

		if !kp.SecretKey().Equal(sk) {
			t.Fatalf("sk != expected (Got %v)", sk)
		}
		if pk := sk.PublicKey(); !kp.PublicKey().Equal(pk) {
			t.Fatalf("sk.PublicKey() != expected (Got %v)", pk)
		}
	})
}

func TestSecretKey(t *testing.T) {
	t.Run("s11n", func(t *testing.T) {
		sk, err := GenerateSecretKey(nil)
		if err != nil {
			t.Fatalf("GenerateSecretKey: %v", err)
		}

		b, err := sk.MarshalBinary()
		if err != nil {
			t.Fatalf("sk.MarshalBinary: %v", err)
		}

		var sk2 SecretKey
		if err = sk2.UnmarshalBinary(b); err != nil {
			t.Fatalf("sk2.UnmarshalBinary(b): %v", err)
		}

		if sk.key.Equal(sk2.key) != 1 {
			t.Fatalf("sk.key != sk2.key (Got %v, %v)", sk.key, sk2.key)
		}
		if sk.nonce != sk2.nonce {
			t.Fatalf("sk.nonce != sk2.nonce (Got %v, %v)", sk.nonce, sk2.nonce)
		}

		// Ensure uninitialized keys serialize.
		var skUninit SecretKey

		b, err = skUninit.MarshalBinary()
		if err != nil {
			t.Fatalf("skUninit.MarshalBinary: %v", err)
		}

		if l := len(b); l != SecretKeySize {
			t.Fatalf("invalid serialized skUninit lenght: %v", l)
		}
		if !isAllZeros(b) {
			t.Fatalf("invalid serialized skUnint (Got %v)", b)
		}
	})
	t.Run("Equal", func(t *testing.T) {
		sk, err := GenerateSecretKey(nil)
		if err != nil {
			t.Fatalf("GenerateSecretKey: %v", err)
		}

		if !sk.Equal(sk) {
			t.Fatalf("sk != sk")
		}

		sk2, err := GenerateSecretKey(nil)
		if err != nil {
			t.Fatalf("GenerateSecretKey: %v", err)
		}

		if sk.Equal(sk2) {
			t.Fatalf("sk == sk2")
		}
	})
}

func TestPublicKey(t *testing.T) {
	t.Run("s11n", func(t *testing.T) {
		kp, err := GenerateKeyPair(nil)
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}

		pk := kp.PublicKey()
		b, err := pk.MarshalBinary()
		if err != nil {
			t.Fatalf("pk.MarshalBinary: %v", err)
		}

		var pk2 PublicKey
		if err = pk2.UnmarshalBinary(b); err != nil {
			t.Fatalf("pk2.UnmarshalBinary(b): %v", err)
		}

		if pk.compressed.Equal(&pk2.compressed) != 1 {
			t.Fatalf("pk.compressed != pk2.compressed (Got %v, %v)", pk.compressed, pk2.compressed)
		}
		if pk.point.Equal(pk2.point) != 1 {
			t.Fatalf("pk.point != pk2.point (Got %v, %v)", pk.point, pk2.point)
		}

		// Ensure uninitialized keys serialize.
		var pkUninit PublicKey

		b, err = pkUninit.MarshalBinary()
		if err != nil {
			t.Fatalf("pkUninit.MarshalBinary: %v", err)
		}

		if l := len(b); l != PublicKeySize {
			t.Fatalf("invalid serialized pkUninit lenght: %v", l)
		}
		if !isAllZeros(b) {
			t.Fatalf("invalid serialized pkUnint (Got %v)", b)
		}
	})
	t.Run("Equal", func(t *testing.T) {
		kp, err := GenerateKeyPair(nil)
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}

		pk := kp.PublicKey()
		if !pk.Equal(pk) {
			t.Fatalf("pk != pk")
		}

		kp2, err := GenerateKeyPair(nil)
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}

		pk2 := kp2.PublicKey()
		if pk.Equal(pk2) {
			t.Fatalf("pk == pk2")
		}
	})
}

func TestKeyPair(t *testing.T) {
	t.Run("s11n", func(t *testing.T) {
		kp, err := GenerateKeyPair(nil)
		if err != nil {
			t.Fatalf("GenerateKeyPair: %v", err)
		}

		b, err := kp.MarshalBinary()
		if err != nil {
			t.Fatalf("kp.MarshalBinary: %v", err)
		}

		var kp2 KeyPair
		if err = kp2.UnmarshalBinary(b); err != nil {
			t.Fatalf("kp2.UnmarshalBinary: %v", err)
		}

		if !kp.sk.Equal(kp2.sk) {
			t.Fatalf("kp.sk != kp2.sk (Got %v, %v)", kp.sk, kp2.sk)
		}
		if !kp.pk.Equal(kp2.pk) {
			t.Fatalf("kp.pk != kp2.pk (Got %v, %v)", kp.pk, kp2.pk)
		}

		// Ensure uninitialized key pairs serialize.
		var kpUninit KeyPair
		b, err = kpUninit.MarshalBinary()
		if err != nil {
			t.Fatalf("kpUninit.MarshalBinary: %v", err)
		}

		if l := len(b); l != KeyPairSize {
			t.Fatalf("invalid serialized kpUninit lenght: %v", l)
		}
		if !isAllZeros(b) {
			t.Fatalf("invalid serialized kpUnint (Got %v)", b)
		}
	})
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := GenerateKeyPair(nil); err != nil {
			b.Fatalf("GenerateKeyPair: %v", err)
		}
	}
}
