import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { highlightCpp } from "@/lib/highlight";

// TA content supports Markdown without executing embedded HTML or scripts.
export function LabMarkdown({ children }: { children: string }) {
  return (
    <div className="markdown-content">
      <Markdown
        remarkPlugins={[remarkGfm]}
        skipHtml
        components={{
          code({ className, children }) {
            // Only generated, escaped highlighter markup is rendered as HTML.
            if (/language-(cpp|c\+\+|cxx)\b/.test(className ?? "")) {
              return (
                <code
                  className={className}
                  dangerouslySetInnerHTML={{
                    __html: highlightCpp(String(children)),
                  }}
                />
              );
            }

            return <code className={className}>{children}</code>;
          },
        }}
      >
        {children}
      </Markdown>
    </div>
  );
}
