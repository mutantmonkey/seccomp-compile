# seccomp-compile

As of November 2024, this package is no longer maintained. There have been no updates to gosecco since 2017, and I haven't used this utility in years. I recommend using the [libseccomp bindings for Go](https://github.com/seccomp/libseccomp-golang) to write your rules and export them in BPF format instead.

This tiny utility compiles seccomp rules in [gosecco](https://github.com/twtiger/gosecco) format and spits them out on stdout. This makes it convenient to add seccomp policies to [bubblewrap](https://github.com/projectatomic/bubblewrap) using simply a shell script.

## Usage
Suppose you want to sandbox `ls` for some reason. You can write your seccomp rules and save them in `~/seccomp/x86_64/ls.seccomp` and then just run:
```sh
seccomp-compile -rules ~/seccomp/x86_64/ls.seccomp | bwrap \
    --unshare-ipc \
    --unshare-pid \
    --unshare-net \
    --unshare-uts \
    --ro-bind /usr /usr \
    --ro-bind /lib /lib \
    --ro-bind /lib64 /lib64 \
    --seccomp 0 \
    ls /
```

Perhaps you need to sandbox a program that still needs normal access to stdin. You can use parameterized file descriptors in combination with process substitution to do this in a clean way without creating any intermediate files on disk. These features are supported by many shells, including bash and zsh. Just write your seccomp rules as usual and run:
```sh
integer seccomp
exec {seccomp}< <(seccomp-compile -rules ~/seccomp/x86_64/less.seccomp)
cat /etc/passwd /etc/hosts /etc/resolv.conf | bwrap \
    --unshare-ipc \
    --unshare-pid \
    --unshare-net \
    --unshare-uts \
    --ro-bind /usr /usr \
    --ro-bind /lib /lib \
    --ro-bind /lib64 /lib64 \
    --seccomp $seccomp \
    less
exec {seccomp}<&-
```
