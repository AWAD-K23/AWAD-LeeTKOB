#include <unistd.h>
#include <sys/wait.h>
#include <iostream>
int main() {
    int children = 0;
    for (int i = 0; i < 64; ++i) {
        pid_t pid = fork();
        if (pid < 0) { std::cout << "BLOCKED"; break; }
        if (pid == 0) { sleep(2); _exit(0); }
        ++children;
    }
    while (children-- > 0) wait(nullptr);
}
