%bool = type i1
%int = type i64
%int8 = type i8
%int16 = type i16
%int32 = type i32
%int64 = type i64
%uint = type i64
%uint8 = type i8
%uint16 = type i16
%uint32 = type i32
%uint64 = type i64
%uintptr = type i64
%float32 = type float
%float64 = type double
%complex64 = type { %float32, %float32 }
%complex128 = type { %float64, %float64 }
%string = type { i8*, %int }
%unsafe.Pointer = type i8*

@builtin.newline = global [1 x i8] c"\0A"

; ssize_t write(int fildes, const void *buf, size_t nbyte)
declare i64 @write(i64 %fd, i8* %buf, i64 %n)

; func println(args ...Type)
;
;    The println built-in function formats its arguments in an
;    implementation-specific way and writes the result to standard error. Spaces
;    are always added between arguments and a newline is appended. Println is
;    useful for bootstrapping and debugging; it is not guaranteed to stay in the
;    language.
define void @println(%string %s, ...) {
	; TODO: handle variadic arguments and include space between each argument.

	; print string to standard output
	%data = extractvalue %string %s, 0
	%len = extractvalue %string %s, 1
	call i64 @write(i64 0, i8* %data, i64 %len)
	; print newline to standard output
	%newline = getelementptr [1 x i8], [1 x i8]* @builtin.newline, i64 0, i64 0
	call i64 @write(i64 0, i8* %newline, i64 1)
	ret void
}

; func print(args ...Type)
;
;    The print built-in function formats its arguments in an
;    implementation-specific way and writes the result to standard error. Print
;    is useful for bootstrapping and debugging; it is not guaranteed to stay in
;    the language.
define void @print(%string %s, ...) {
	; TODO: handle variadic arguments.

	; print string to standard output
	%data = extractvalue %string %s, 0
	%len = extractvalue %string %s, 1
	call i64 @write(i64 0, i8* %data, i64 %len)
	ret void
}

; wrapnilchk returns ptr if non-nil, panics otherwise.
; (For use in indirection wrappers.)
;
;    func ssa:wrapnilchk(ptr *T, recvType, methodName string) *T
define i8* @"ssa:wrapnilchk"(i8* %ptr, %string %recvType, %string %methodName) {
	%ptr_val = ptrtoint i8* %ptr to %uintptr
	%is_null = icmp eq %uintptr %ptr_val, 0
	br i1 %is_null, label %success, label %fail

success:
	ret i8* %ptr

fail:
	; TODO: add panic message with "recvType.methodName".
	unreachable
}

; func cmp.string(x, y string) int
;
;    cmp.string compares x with y lexically byte-wise and returns:
;
;        -1 if x <  y
;         0 if x == y
;        +1 if x >  y
define %int @cmp.string(%string %x, %string %y) {
entry:
	%x_data = extractvalue %string %x, 0 ; %uint8*
	%x_len = extractvalue %string %x, 1  ; %int
	%y_data = extractvalue %string %y, 0 ; %uint8*
	%y_len = extractvalue %string %y, 1  ; %int
	; l = min(len(x), len(y))
	%l = call %int @internal.min(%int %x_len, %int %y_len)
	br label %loop.pre

	; for (i := 0; i < l; i++)
loop.pre:
	%i.ptr = alloca %int
	store %int 0, %int* %i.ptr
	br label %loop.cond

loop.cond:
	%i = load %int, %int* %i.ptr
	%cond = icmp slt %int %i, %l
	br i1 %cond, label %loop.body, label %loop.exit

loop.body:
	%xp = getelementptr %uint8, %uint8* %x_data, %int %i
	%yp = getelementptr %uint8, %uint8* %y_data, %int %i
	%xc = load %uint8, %uint8* %xp
	%yc = load %uint8, %uint8* %yp
	br label %check_byte_less

check_byte_less:
	%byte_less = icmp ult %uint8 %xc, %yc
	br i1 %byte_less, label %ret_less, label %check_byte_greater

check_byte_greater:
	%byte_greater = icmp ugt %uint8 %xc, %yc
	br i1 %byte_greater, label %ret_greater, label %loop.post

loop.post:
	%i.inc = add %int %i, 1
	store %int %i.inc, %int* %i.ptr
	br label %loop.cond

loop.exit:
	br label %check_len_less

check_len_less:
	%len_less = icmp slt %int %x_len, %y_len
	br i1 %len_less, label %ret_less, label %check_len_greater

check_len_greater:
	%len_greater = icmp sgt %int %x_len, %y_len
	br i1 %len_greater, label %ret_greater, label %ret_equal

ret_less:
	ret %int -1

ret_equal:
	ret %int 0

ret_greater:
	ret %int 1
}

; min returns the smaller of x and y.
define %int @internal.min(%int %x, %int %y) {
entry:
	%cond = icmp slt %int %x, %y
	br i1 %cond, label %x_min, label %y_min

x_min:
	ret %int %x

y_min:
	ret %int %y
}
