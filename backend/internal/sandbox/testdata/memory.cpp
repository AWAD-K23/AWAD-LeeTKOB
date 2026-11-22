#include <cstdlib>
#include <iostream>
int main() {
    // Volatile writes and reads force real pages, even with optimizations.
    volatile unsigned char* memory = (volatile unsigned char*)malloc(1024ULL * 1024 * 1024);
    if (!memory) return 1;
    for (unsigned long long i = 0; i < 1024ULL * 1024 * 1024; ++i) memory[i] = 123;
    unsigned long long total = 0;
    for (unsigned long long i = 0; i < 1024ULL * 1024 * 1024; ++i) total += memory[i];
    std::cout << total;
}
