# RISC-V "Bedrock Bootstrapping"

This is a project to bootstrap a RISC-V system which can do _something_ useful, from the very bottom up with minimal external dependencies.

"From the bottom up" means:

- Start with machine code, which can be assembled by hand
- Create a simple assembler which supports jumps to labels
- Implement some sort of higher level language (or a subset of one)
- Use the higher level language to achieve some useful task

Ideally, it should be possible to do most of this from scratch on the board itself, requiring only the minimal possible interaction with other computers.

There's probably enough scope that it's reasonable to leave "some task" to be defined later. An example might be successfully serving an HTTP request, or maybe encrypting or signing some value.

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

RISC-V wins by default since it's open source. x86\_64 loses by default since the instruction set is so complex and
that would make writing machine code or the basic assembler much much harder.

ARM is more interesting since hardware is very widely available whereas RISC-V hardware isn't at the time of writing.
A parallel project to do the same "bedrock bootstrap" on ARM would be great if only because Raspberry Pi boards are so
prevalent and cheap. That said, RISC-V is open source so that's where we'll start.

In any case if we lack a good emulation solution running on commodity (x86\_64) hardware, the project will be nearly
impossibly hard and that takes the sting out of a lack of physical hardware.

### Use the HiFive1 Rev A Only (to start with)

The [HiFive1](https://www.sifive.com/boards/hifive1) is a RISC-V board by SiFive; it's low-cost, simple and very open indeed.

By settling on one hardware board to run on (alongside Qemu) we'll limit the amount of time we need to spend worrying about harware differences. We should try to keep our options open in terms of using other boards in the future, but we'll focus on the HiFive1. Specifically we're using Revision A, not Revision B or later - although the differences aren't huge.

## Project Motivation

TL;DR: Why bother with this bootstrapping project?

Well, frankly, it'll be fun. By implementing a whole system - or at least, the software for part of a whole system - from as close to "scratch" as possible, and making that implementation easy enough to follow, we can potentially allow for people to develop special purpose systems entirely from scratch where they _genuinely have_ audited the _whole_ system and where they really can account for _every_ line of code. For some crypto work that's a boast that could really mean something. For other stuff at least it'll be cool, so why not?

### Longer Motivation

On most modern systems, the bootstrapping process is hidden from users by virtue of the fact that they install an OS from some other medium (e.g. a USB stick) which was compiled on someone _else's_ computer. That means that if you're installing Windows, you're trusting a lot of people at Microsoft. If you install Linux, you're trusing whoever provides your distro of choice which is, again, a group of people. You're also trusting whoever wrote the compiler that was used, and anyone who wrote any of the software which is installed, such as Notepad on Windows or Vim on a Linux box... and in turn, whatever compiler was used to compile those binaries, and the computer which was used to compile them, and so on.

On Windows, "trust" is the end of the road - you don't get to see most of the code that's running and check for yourself if you're happy with it.

On Linux you generally _do_ get that right, but in practise _nobody_ on the planet can feasibly vet every single line of code they're running, whether they have access to the source or not. Even with the source, one has to trust the binaries provided by their package manager unless they're reproducible - which most aren't. Compiling everything yourself is possible, but the problem of not being able to feasibly audit the code remains.

On top of this, modern systems usually run on proprietary hardware, with proprietary firmware. Availability of the source for these devices is vanishingly rare, and even where it can be replaced it often requires specialist knowledge and tools to achieve.

In practise, the issues above aren't generally problems - doing modern computing requires a lot of stuff just to get off the ground, and if you use Linux you can get most of the way to a pretty secure and trustable environment. The dedicated enthusiast might choose to use [Linux from Scratch](http://www.linuxfromscratch.org/lfs/) to have more control over
their OS and the rate at which updates are merged. The proprietary, closed hardware is a pain point, but RISC-V and the ecosystem building around it present an opportunity to improve the market. The dream of an entirely open usable system seems to be getting close, and by having the ability to bootstrap everything on one system as far as is possible, we're gaining more openness.