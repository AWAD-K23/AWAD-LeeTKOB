#include <iostream>

struct Counter {
    static int total(int n, int accumulator = 0) {
        return n == 0 ? accumulator : total(n - 1, accumulator + n);
    }
};

int main() {
    int n;
    std::cin >> n;
    std::cout << Counter::total(n) << '\n';
}
