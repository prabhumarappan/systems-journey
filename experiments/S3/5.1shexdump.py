import sys


def main():
    args = sys.argv
    if len(args) == 1:
        src = sys.stdin.buffer
    elif len(args) == 2:
        src = open(sys.argv[1], "rb")
    else:
        print(f"usage: {sys.argv[0]} <filename>", file=sys.stderr)
        sys.exit(1)

    with src as f:
        nibble = "0123456789abcdef"

        while len(data := f.read(16)) > 0:
            encoded = []
            for b in data:
                out = ""
                out += nibble[b >> 4]
                out += nibble[b & 15]
                encoded.append(out)
            print(" ".join(encoded), file=sys.stdout)


if __name__ == "__main__":
    main()
