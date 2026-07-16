# Systems programming & what a program is (eblog pt 1, ch 1–3)

Session S2 · Source: https://eblog.fly.dev/startingsystems1.html · Notes: `notes/os/eblog-programs.md`

```mermaid
mindmap
  root((systems programming))
    what counts - any ONE of
      talks to the OS or hardware
      tight performance constraints
      low-level = individual bytes and registers
    mindset
      computer is a physical machine you can fully understand
      tear it apart, put it back together
    what is a program
      executable file the OS runs as machine instructions
      just bytes
        code bytes - the instructions
        data bytes - the hello world string
      loaded into memory to run
    program shapes
      run once with input, then exit
      long-running
        daemon - background, no terminal
        shell - long-running AND interactive, not either/or
    everything is bytes
      input and output are just bytes
```

**One-line takeaway:** a program is just bytes — some are code, some are data — and a systems programmer is someone who isn't afraid to look at them directly.
