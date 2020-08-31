# RISC-V "Bedrock Bootstrapping"

A project to bootstrap a RISC-V system which can do _something_ useful, from the very bottom up with minimal external dependencies.

"From the bottom up" means:

- Start with machine code, which could be written by hand
- Create a simple assembler which supports jumps to labels
- Implement some sort of higher level language (or a subset of one)
- Use the higher level language to achieve some useful task

It should ideally be possible to do most of this from scratch on the board itself, requiring only minimal interaction with other computer systems.

The "something useful" is for now left undefined - the scope here is large enough that it'll take a while to get there.

## Scope

As the scope of this project could be absolutely immense, limiting what we're going to do is important.

### No Linux / Other Kernels

If we introduce Linux or a different kernel, we introduce potentially millions of lines of code and potentially also license obligations. We want to account for every line of code and allow everyone to be able to reproduce the results on their own hardware with any modifications they like.

### "Bedrock" means Machine Code

This means that the lowest layer we'll consider is RISC-V machine
code that we got onto a machine _somehow_. For now we won't try to burn anything into one-time-programmable memory, or do anything in hardware apart from writing our code into SPI Flash where it can be run, and maybe connecting some peripherals.

This means that there will be _some_ dependence on using another machine (e.g. to copy the code to a development board or to run an emulator or simulator).

Being able to hook up some device and write raw machine code directly into flash and thereby avoid any dependence on other full-featured operating systems sounds great but could be a step too far. In any case we'll do our best to not restrict this possibility.

### RISC-V > ARM > x86\_64

RISC-V wins by default since it's open source! Further to that, RISC-V has a beautiful simplicity and elegance in its design, which makes it attractive as an ISA.

x86\_64 loses by default since the instruction set is so complex. It's possible to construct a single x86 instruction which is 15 bytes in total, and likewise a single instruction which is only 1 byte. That variety makes writing, parsing and generating machine code much, much harder. Likewise, parsing x86 is vastly more complicated than RISC-V or ARM where we can (choose to) have every instruction take the same length.

ARM is more interesting than x86 since hardware is very widely available whereas RISC-V hardware isn't at the time of writing.
A parallel project to do the same "bedrock bootstrap" on ARM would be great if only because ARM development boards like the Raspberry Pi are so
prevalent and cheap. Still, the open design of RISC-V makes it slightly more attractive for use here.

In any case if we lack a good emulation solution running on commodity (x86\_64) hardware, the project will be nearly
impossibly hard. That means that we need an emulator to be available, which in turn takes the sting out of a lack of physical hardware.

### Use the HiFive1 Rev A (to start with)

The [HiFive1](https://www.sifive.com/boards/hifive1) is a RISC-V board by SiFive; it's low-cost, simple and very open.

By settling on one hardware board to run on (alongside QEMU) we'll limit the amount of time we need to spend worrying about hardware differences. We should try to keep our options open in terms of using other boards in the future, but we'll focus on the HiFive1. Specifically we're using Revision A, not Revision B or later - although the differences aren't huge.

### RV32I

RISC-V has the notion of "extensions". The base instruction set, known as `I` is relatively small compared to the base instruction set of many other ISAs. The base set is supplemented by optional extensions, which a vendor may or may not choose to implement on a per-chip basis. These include `M` for multiplication and division, `A` for atomic operations and `C` for compressed instructions which map to base instructions and take less memory to store. There are other extensions too.

Together, `RV32IMAC` describes the extensions available on the HiFive1 board, but we're only going to use `I`. That's because every RISC-V processor must support `I` and that's something we can rely on for the future - so by targeting only the base set, we can maximise the portability of the code we write. When RISC-V based Linux-capable machines become more common, we'll confidently be able to target them. (Of course, the boot process differs between machines, and that would have to be tweaked - but we'd have to do that anyway.)

In addition, RV32I is very simple, which reduces the "surface area" of knowledge we need to get off the ground and start developing and understanding. When we're dealing with machine code (and in fact, more generally) simplicity is king.

## Project Motivation

I.e. "Why bother?"

Frankly, it'll be fun. By implementing a whole system - or at least, the software for a whole system - from as close to "scratch" as possible, and making that implementation easy enough to follow, we can potentially allow for people to develop special purpose systems entirely from scratch where they _genuinely have_ audited the _whole_ system and where they really can account for _every_ line of code. For some crypto work that's a boast that could really mean something. For other stuff at least it'll be cool, so why not?

### Longer Motivation

On most modern systems, the bootstrapping process is hidden from users by virtue of the fact that they install an OS from some other medium (e.g. a USB stick) which was compiled on someone _else's_ computer. That means that if you're installing Windows, you're trusting a lot of people at Microsoft. If you install Linux, you're trusting whoever provides your distro of choice which is, again, a group of people. You're also trusting whoever wrote the compiler that was used, and anyone who wrote any of the software which is installed, such as Notepad on Windows or Vim on a Linux box; in turn you trust whatever compiler was used to compile those binaries, and the computer which was used to compile them, and so on.

On Windows, "trust" is the end of the road - you don't get to see most of the code that's running and check it.

On Linux you generally _do_ get that right, but in practise _nobody_ can feasibly vet every single line of code they're running in a Linux system, whether they have access to the source or not. Even with the source, one has to trust the binaries provided by their package manager unless they're reproducible - which most aren't. Compiling everything yourself is possible, but the problem of not being able to feasibly audit the code remains.

On top of this, modern systems usually run on proprietary hardware, with proprietary firmware. Source availability for firmware is vanishingly rare, and even where it can be replaced it often requires specialist knowledge and tools.

In practise, the issues above aren't generally problems - doing modern computing requires a lot of stuff just to get off the ground, and if you use Linux you can get most of the way to a pretty secure and trustable environment. The dedicated enthusiast might choose to use [Linux from Scratch](http://www.linuxfromscratch.org/lfs/) to have more control over
their OS and the rate at which updates are merged. The proprietary, closed hardware is a pain point, but RISC-V and the ecosystem building around it present an opportunity to improve the market. The dream of an entirely open usable system seems to be getting close, and by having the ability to bootstrap everything on one system as far as is possible, we're gaining more openness.
