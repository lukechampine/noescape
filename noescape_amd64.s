// +build amd64,!noasm

#include "textflag.h"

// func callRead(r io.Reader, b []byte) (int, error)
TEXT ·callRead(SB), $64-64
	// Before calling the Read method, we need to put all the data it needs on
	// the stack. The Read method obviously needs b, but less obviously, it also
	// needs r_data -- the receiver object. When you call a method, the receiver
	// is passed as the first argument. So we'll start by copying r_data onto the
	// stack, via the AX register.
	MOVQ r_data+8(FP), AX
	MOVQ AX, (SP)

	// Next, we'll copy b, which consists of three words:
	MOVQ b_base+16(FP), BX
	MOVQ b_len+24(FP), CX
	MOVQ b_cap+32(FP), DX
	MOVQ BX, 8(SP)
	MOVQ CX, 16(SP)
	MOVQ DX, 24(SP)

	// Now we can invoke the Read method. To do so, we need to load the actual
	// function pointer. It is stored in the 'itable' of r, in the 'fun' array;
	// see runtime.itab. The offset of 'fun' within 'itab' is 24 bytes.
	MOVQ r_itable+0(FP), AX
	CALL 24(AX)

	// Before the method returned, it placed its return values on the stack. We
	// now need to copy those values into the "return value slots" of our own
	// stack. We're returning three words in total: one for the int, and two for
	// the error interface.
	MOVQ 32(SP), AX
	MOVQ 40(SP), BX
	MOVQ 48(SP), CX
	MOVQ AX, ret+40(FP)
	MOVQ BX, ret1_itable+48(FP)
	MOVQ CX, ret1_data+56(FP)
	RET

// func callWrite(w io.Writer, b []byte) (int, error)
TEXT ·callWrite(SB), $64-64
	// identical to callRead
	MOVQ r_data+8(FP), AX
	MOVQ b_base+16(FP), BX
	MOVQ b_len+24(FP), CX
	MOVQ b_cap+32(FP), DX
	MOVQ AX, (SP)
	MOVQ BX, 8(SP)
	MOVQ CX, 16(SP)
	MOVQ DX, 24(SP)
	MOVQ r_itable+0(FP), AX
	CALL 24(AX)
	MOVQ 32(SP), AX
	MOVQ 40(SP), BX
	MOVQ 48(SP), CX
	MOVQ AX, ret+40(FP)
	MOVQ BX, ret1_itable+48(FP)
	MOVQ CX, ret1_data+56(FP)
	RET
