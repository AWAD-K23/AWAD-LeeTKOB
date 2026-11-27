#include <iostream>

int unused(int n) {
    return n == 0 ? 0 : n + unused(n - 1);
}

int main() {
    int n;
    std::cin >> n;
    std::cout << (n * n + n) / 2 << '\n';
}
