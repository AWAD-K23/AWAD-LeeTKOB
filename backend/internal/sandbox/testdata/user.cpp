#include <unistd.h>
#include <iostream>
int main() { std::cout << (getuid() != 0 && setuid(0) != 0 ? "BLOCKED" : "ESCAPED"); }
