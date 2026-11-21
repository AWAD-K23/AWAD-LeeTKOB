// Resolve actual C++ call targets. Comments and strings never become CallExprs.
#include <clang-c/Index.h>
#include <iostream>
#include <string>

static std::string text(CXString value) {
    std::string out = clang_getCString(value) ? clang_getCString(value) : "";
    clang_disposeString(value);
    return out;
}

static std::string qualified(CXCursor cursor) {
    std::string name = text(clang_getCursorSpelling(cursor));
    CXCursor parent = clang_getCursorSemanticParent(cursor);
    while (!clang_Cursor_isNull(parent) && clang_getCursorKind(parent) != CXCursor_TranslationUnit) {
        std::string part = text(clang_getCursorSpelling(parent));
        if (!part.empty() && part.rfind("__", 0) != 0) name = part + "::" + name;
        parent = clang_getCursorSemanticParent(parent);
    }
    return name;
}

static CXChildVisitResult visit(CXCursor cursor, CXCursor, CXClientData) {
    CXSourceLocation location = clang_getCursorLocation(cursor);
    CXFile file = nullptr;
    unsigned line = 0;
    clang_getExpansionLocation(location, &file, &line, nullptr, nullptr);
    auto kind = clang_getCursorKind(cursor);
    if ((kind == CXCursor_CallExpr || kind == CXCursor_DeclRefExpr) && file && text(clang_getFileName(file)) == "/work/main.cpp") {
        CXCursor target = clang_getCursorReferenced(cursor);
        if (!clang_Cursor_isNull(target)) {
            // A tab-separated name and line avoid serializing arbitrary student strings.
            std::cout << qualified(target) << '\t' << line << '\n';
        }
    }
    return CXChildVisit_Recurse;
}

int main(int argc, char** argv) {
    if (argc != 2) return 2;
    CXIndex index = clang_createIndex(0, 0);
    const char* args[] = {"-std=c++17", "-xc++"};
    CXTranslationUnit unit = clang_parseTranslationUnit(index, argv[1], args, 2, nullptr, 0, CXTranslationUnit_None);
    if (!unit) return 2;
    bool invalid = false;
    for (unsigned i = 0; i < clang_getNumDiagnostics(unit); ++i) {
        CXDiagnostic diagnostic = clang_getDiagnostic(unit, i);
        invalid |= clang_getDiagnosticSeverity(diagnostic) >= CXDiagnostic_Error;
        clang_disposeDiagnostic(diagnostic);
    }
    if (!invalid) clang_visitChildren(clang_getTranslationUnitCursor(unit), visit, nullptr);
    clang_disposeTranslationUnit(unit);
    clang_disposeIndex(index);
    return invalid ? 2 : 0;
}
