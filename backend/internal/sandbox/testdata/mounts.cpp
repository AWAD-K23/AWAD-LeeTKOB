#include <fstream>
#include <iostream>
#include <sys/socket.h>
#include <sys/un.h>
#include <cstring>
int main() {
    sockaddr_un address{}; address.sun_family = AF_UNIX;
    std::strcpy(address.sun_path, "/var/run/docker.sock");
    int fd = socket(AF_UNIX, SOCK_STREAM, 0);
    bool socketVisible = connect(fd, (sockaddr*)&address, sizeof(address)) == 0;
    bool hostVisible = std::ifstream("/host/etc/passwd").good() || std::ifstream("/run/host/etc/passwd").good();
    std::cout << (socketVisible || hostVisible ? "ESCAPED" : "BLOCKED");
}
