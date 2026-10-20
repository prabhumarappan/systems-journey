import sys

unhex = {
    "0": 0x0,
    "1": 0x1,
    "2": 0x2,
    "3": 0x3,
    "4": 0x4,
    "5": 0x5,
    "6": 0x6,
    "7": 0x7,
    "8": 0x8,
    "9": 0x9,
    "a": 0xA,
    "b": 0xB,
    "c": 0xC,
    "d": 0xD,
    "e": 0xE,
    "f": 0xF,
    "A": 0xA,
    "B": 0xB,
    "C": 0xC,
    "D": 0xD,
    "E": 0xE,
    "F": 0xF,
}


def main():
    args = sys.argv
    if len(args) == 1:
        src_file = sys.stdin
    elif len(args) == 2:
        src_file = open(sys.argv[1], "r")
    else:
        print(f"usage: {sys.argv[0]} <filename>")
        sys.exit(1)

    with src_file as f:
        src = f.read().split()
        out = []
        for x in src:
            print(x)
            high = unhex[x[0]]
            low = unhex[x[1]]
            out.append(high << 4 | low)
        sys.stdout.buffer.write(bytes(out))


if __name__ == "__main__":
    main()
