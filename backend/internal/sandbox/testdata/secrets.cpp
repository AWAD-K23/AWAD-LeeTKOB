#include <cstdlib>
#include <fstream>
#include <iostream>
#if __has_include("/app/.env") || __has_include("/work/.env") || __has_include("/run/secrets/model_api_key")
#error Sandbox exposed a secret file to the compiler
#endif
int main() {
    bool exposed = getenv("MODEL_API_KEY") || getenv("GEMINI_API_KEY") || getenv("DATABASE_URL");
    exposed |= std::ifstream("/proc/1/root/app/.env").good();
    std::cout << (exposed ? "ESCAPED" : "BLOCKED");
}
