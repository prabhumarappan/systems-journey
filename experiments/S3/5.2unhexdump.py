import sys


def main():
    args = sys.argv
    if len(args) == 1:
        src = sys.stdin
    elif len(args) == 2:
        src = open(sys.argv[1], "rb")
    else:
        print(f"usage: {sys.argv[0]} <filename>")
        sys.exit(1)

    with src as f:
        while len(data := f.read(3)) > 0:
            sys.stdout.buffer.write(bytes([int(data, 16)]))


if __name__ == "__main__":
    main()
