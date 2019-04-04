#!/usr/bin/env python2

from pwn import *

context(log_level="warning")

for i in range(100):
    shell = remote("2018shell.picoctf.com", 34490)

    shell.sendafter(": ", "A"*i + "\n")

    enc = shell.recvline()

    print i, len(enc), enc

    shell.close()
