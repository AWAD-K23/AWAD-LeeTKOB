#include <sys/socket.h>
#include <arpa/inet.h>
#include <iostream>
int main() {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    sockaddr_in address{};
    address.sin_family = AF_INET;
    address.sin_port = htons(443);
    inet_pton(AF_INET, "1.1.1.1", &address.sin_addr);
    std::cout << (connect(fd, (sockaddr*)&address, sizeof(address)) == 0 ? "ESCAPED" : "BLOCKED");
}
