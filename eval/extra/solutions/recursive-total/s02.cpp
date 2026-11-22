#include <iostream>

int main() {
    int n, result = 0;
    std::cin >> n;
    for (int i = 1; i <= n; ++i) result += i;
    std::cout << result << '\n';
}
