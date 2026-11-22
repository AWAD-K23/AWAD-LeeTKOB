#include <iostream>

// A text search for "for", "while", or "goto" would wrongly reject this.
int total(int n) {
    const char* example = "for (;;) while (true) goto done;";
    (void)example;
    return n == 0 ? 0 : n + total(n - 1);
}

int main() {
    int n;
    std::cin >> n;
    std::cout << total(n) << '\n';
}
