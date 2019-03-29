# RISC-V "Bedrock Bootstrapping"
This is a project to bootstrap a RISC-V system which can do _something_ useful, from the very bottom up.

"From the bottom up" means:
- Start with machine code, which can be assembled by hand
- Create a simple assembler which supports jumps to labels
- Implement some sort of higher level language (or a subset of one)
- Use the higher level language to achieve some useful task

There's probably enough scope there that it's fine to leave "some task" to be defined. An example might be
successfully serving an HTTP response over the network, or maybe encrypting, signing or hashing some value.

## Project Motivation
On most modern systems, the bootstrapping process is hidden from users by virtue of the fact that they install
an OS from some other medium (e.g. a USB stick) which was compiled on someone _else's_ computer. That means that
if you're installing Windows, you're trusting a lot of people at Microsoft. If you install Linux, you're trusing
whoever provides your distro of choice which is, again, a group of people. You're also trusting whoever wrote the
compiler that was used, and anyone who wrote any of the software which is installed, such as Notepad on Windows
or Vim on a Linux box... and in turn, whatever compiler was used to compile those binaries.

On Windows, "trust" is the end of the road - you don't get to see most of the code that's running. On Linux you
generally _do_ get that right, but in practise _nobody_ on the planet can feasibly vet every single line of code
they're running, whether they have access to it or not. On top of that, modern systems are often based on proprietary
hardware, with proprietary firmware. Availability of the source for these devices is vanishingly rare, and even where
it can be replaced, it often requires specialist knowledge and tools to achieve.

In practise, the statements above aren't generally a problem - doing modern computing requires a lot of stuff just to
get off the ground, and if you use Linux you can get most of the way to a pretty secure environment. The dedicated
enthusiast might choose to use [Linux from Scratch](http://www.linuxfromscratch.org/lfs/) to have more control over
their OS and the rate at which updates are merged (albeit at the risk of missing a security update and actually
decreasing the overall security of their system). The proprietary, closed hardware is a pain point, but RISC-V and
the ecosystem building around it present an opportunity to improve the market. The dream of an entirely open, usable
system seems to be getting close.

TL;DR: Why bother with this bootstrapping project then?

Well, frankly, it'll be fun to poke at the limits of what's possible. By implementing a whole system - or at least,
the software for a whole system - from as close to "scratch" as possible, and making that implementation easy enough
to follow, we can potentially allow for people to develop special purpose systems entirely from scratch where they
_genuinely have_ audited the whole system and where they really can account for every line of code. For some crypto
work that's a boast that could really mean something. For other non-crypto stuff: it'll be cool, so why not?

## Scope
As the scope of this project is immense, limiting scope is important. We'll explicitly not attempt the following:

### No Linux / Other Kernels
If we introduce Linux or a different kernel, we introduce potentially millions of lines of code and potentially also
license obligations. We want to account for every line of code and allow everyone to be able to reproduce the results
on their own hardware with any modifications they like.

### "Bedrock" means Machine Code
By "Bedrock means Machine Code" we mean that the lowest layer we'll consider is a RISC-V processor running machine
code that we got onto it _somehow_. For now we won't try to burn a ROM, or do anything in hardware. This means that
there will be _some_ dependence on using another machine (e.g. to copy the code to a development board or to run an
emulator or simulator).

Being able to hook up some device and write raw machine code directly into a ROM and thereby avoid any dependence on
other full-featured operating systems sounds great and would be fun to do, but it's out of scope here.

### RISC-V > ARM > x86\_64
RISC-V wins by default since it's open source. x86\_64 loses by default since the instruction set is so complex and
that would make writing machine code or the basic assembler much much harder.

ARM is more interesting since hardware is very widely available whereas RISC-V hardware isn't at the time of writing.
A parallel project to do the same "bedrock bootstrap" on ARM would be great if only because Raspberry Pi boards are so
prevalent and cheap. That said, RISC-V is open source so that's where we'll start.

In any case if we lack a good emulation solution running on commodity (x86\_64) hardware, the project will be nearly
impossibly hard and that takes the sting out of a lack of physical hardware.

## Useful Links
- [ocd SPI flash driver](https://github.com/riscv/riscv-openocd/blob/riscv/src/flash/nor/fespi.c)
- [OSDev Bare Bones](https://wiki.osdev.org/HiFive-1_Bare_Bones#The_Bare_Bones)
