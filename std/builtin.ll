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
