    .global main
    .text
main:
    mov     message(%rip), %rdi
    call    puts
    ret

message:
    .asciz  "Hello, world"
