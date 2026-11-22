import hljs from "highlight.js/lib/core";
import cpp from "highlight.js/lib/languages/cpp";

hljs.registerLanguage("cpp", cpp);

// highlight.js escapes source text before adding its own token markup.
export function highlightCpp(source: string) {
  return hljs.highlight(source, { language: "cpp" }).value;
}
