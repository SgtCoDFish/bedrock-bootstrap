|  x# | ABI name | Class      | Preserved? | Typical use           |
| :-: | -------- | ---------- | ---------- | --------------------- |
|  x0 | zero     | Constant   | —          | Hard-wired 0          |
|  x1 | ra       | Link       | No         | Return address        |
|  x2 | sp       | Stack ptr  | Yes        | Stack pointer         |
|  x3 | gp       | Global ptr | Yes        | Global data access    |
|  x4 | tp       | Thread ptr | Yes        | TLS / hart-local data |
|  x5 | t0       | Temporary  | No         | Scratch               |
|  x6 | t1       | Temporary  | No         | Scratch               |
|  x7 | t2       | Temporary  | No         | Scratch               |
|  x8 | s0 / fp  | Saved      | Yes        | Saved reg / frame ptr |
|  x9 | s1       | Saved      | Yes        | Saved reg             |
| x10 | a0       | Argument   | No         | Arg 0 / return value  |
| x11 | a1       | Argument   | No         | Arg 1 / return value  |
| x12 | a2       | Argument   | No         | Arg 2                 |
| x13 | a3       | Argument   | No         | Arg 3                 |
| x14 | a4       | Argument   | No         | Arg 4                 |
| x15 | a5       | Argument   | No         | Arg 5                 |
| x16 | a6       | Argument   | No         | Arg 6                 |
| x17 | a7       | Argument   | No         | Arg 7 / syscall #     |
| x18 | s2       | Saved      | Yes        | Saved reg             |
| x19 | s3       | Saved      | Yes        | Saved reg             |
| x20 | s4       | Saved      | Yes        | Saved reg             |
| x21 | s5       | Saved      | Yes        | Saved reg             |
| x22 | s6       | Saved      | Yes        | Saved reg             |
| x23 | s7       | Saved      | Yes        | Saved reg             |
| x24 | s8       | Saved      | Yes        | Saved reg             |
| x25 | s9       | Saved      | Yes        | Saved reg             |
| x26 | s10      | Saved      | Yes        | Saved reg             |
| x27 | s11      | Saved      | Yes        | Saved reg             |
| x28 | t3       | Temporary  | No         | Scratch               |
| x29 | t4       | Temporary  | No         | Scratch               |
| x30 | t5       | Temporary  | No         | Scratch               |
| x31 | t6       | Temporary  | No         | Scratch               |

