def main():
    rawbytes = bytes([0xb7, 0x02, 0x40, 0x20, 0x67, 0x80, 0x02, 0x00])
    with open("bootloader", "w+b", buffering=0) as f:
        f.write(rawbytes)


if __name__ == '__main__':
    main()
