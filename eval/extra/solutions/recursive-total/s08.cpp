#include <iostream>
#define STEP while

int main() {
    int n, answer = 0;
    std::cin >> n;
    STEP (n > 0) answer += n--;
    std::cout << answer << '\n';
}
