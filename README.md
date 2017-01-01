# seccomp-compile

This tiny utility compiles seccomp rules in [gosecco](https://github.com/twtiger/gosecco) format and spits them out on stdout. This makes it convenient to add seccomp policies to [bubblewrap](https://github.com/projectatomic/bubblewrap) using simply a shell script.

## Usage
Suppose you want to sandbox `ls` for some reason. You can write your seccomp rules and save them in `~/seccomp/ls-amd64.seccomp` and then just run:
```sh
seccomp-compile -rules ~/seccomp/ls-amd64.seccomp | bwrap \
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
