#include <iostream>

int total(int n) {
    return n == 0 ? 0 : n + total(n - 1);
}

int main() {
    int n;
    std::cin >> n;
    std::cout << total(n) << '\n';
}
