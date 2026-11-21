#include <unistd.h>
int main() { char buffer[4096]{}; for (;;) write(1, buffer, sizeof(buffer)); }
