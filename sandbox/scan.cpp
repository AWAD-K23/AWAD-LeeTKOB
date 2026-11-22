// Resolve actual C++ call targets. Comments and strings never become CallExprs.
#include <clang-c/Index.h>
#include <iostream>
#include <string>

// Both the worker sandbox and standalone checker supply their own source path.
static std::string mainPath;

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
    bool submitted = file && text(clang_getFileName(file)) == mainPath;

    // Statement cursors distinguish real repetition from comments and strings.
    if (submitted) {
        const char* statement = nullptr;
        switch (kind) {
            case CXCursor_ForStmt:
            case CXCursor_CXXForRangeStmt: statement = "for"; break;
            case CXCursor_WhileStmt: statement = "while"; break;
            case CXCursor_DoStmt: statement = "do"; break;
            case CXCursor_GotoStmt: statement = "goto"; break;
            default: break;
        }
        if (statement) std::cout << statement << '\t' << line << '\n';
    }

    if (submitted && (kind == CXCursor_CallExpr || kind == CXCursor_DeclRefExpr ||
                      kind == CXCursor_TypeRef || kind == CXCursor_TemplateRef)) {
        CXCursor target = clang_getCursorReferenced(cursor);
        if (!clang_Cursor_isNull(target)) {
            // A tab-separated name and line avoid serializing arbitrary student strings.
            std::cout << qualified(target) << '\t' << line << '\n';
        }
    }

    // Canonical types expose containers hidden behind typedefs, using aliases,
    // or auto. Only declarations in submitted code count, never header bodies.
    if (submitted && (kind == CXCursor_VarDecl || kind == CXCursor_FieldDecl ||
                      kind == CXCursor_ParmDecl)) {
        CXType type = clang_getCanonicalType(clang_getCursorType(cursor));
        while (type.kind == CXType_Pointer || type.kind == CXType_LValueReference ||
               type.kind == CXType_RValueReference) {
            type = clang_getCanonicalType(clang_getPointeeType(type));
        }
        CXCursor declaration = clang_getTypeDeclaration(type);
        if (!clang_Cursor_isNull(declaration)) {
            std::cout << qualified(declaration) << '\t' << line << '\n';
        }
    }
    return CXChildVisit_Recurse;
}

int main(int argc, char** argv) {
    if (argc != 2) return 2;
    mainPath = argv[1];
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
