import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react/";

import { useRef } from "react";
import * as monacoEditor from "monaco-editor";
import githubDark from "monaco-themes/themes/GitHub Dark.json";

import * as monaco from "monaco-editor";

type MonacoEditorProps = {
  code: string;
  onChange: (value: string | undefined) => void;
  onMount?: OnMount;
  errorLines?: number[];
};

const MonacoEditor: React.FC<MonacoEditorProps> = ({
  code,
  onChange,
  onMount,
  errorLines,
}) => {
  const editorRef = useRef<monacoEditor.editor.IStandaloneCodeEditor | null>(
    null
  );

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor;

    // Register the Python language (basic tokenizer) - language server will provide rich features
    monaco.languages.register({ id: "python" });
    monaco.languages.setMonarchTokensProvider("python", {
      tokenizer: {
        root: [
          [
            /\b(class|def|return|if|else|elif|for|while|import|from|as|pass|break|continue|lambda|try|except|raise|with|yield|assert|in|is|not|and|or)\b/,
            "keyword",
          ],
          [/"[^"]*"|'[^']*'/, "string"],
          [/[0-9]+/, "number"],
          [/[a-zA-Z_]\w*/, "identifier"],
          [/#.*/, "comment"],
        ],
      },
    });

    // Register the GitHub Dark theme
    monaco.editor.defineTheme(
      "github-dark",
      githubDark as monacoEditor.editor.IStandaloneThemeData
    );
    monaco.editor.setTheme("github-dark"); // Set the theme you defined

    // Optional: You can keep this basic autocomplete if you want, but the LSP will offer more.
    // Consider removing it if LSP provides all necessary completions to avoid conflicts/redundancy.
    monaco.languages.registerCompletionItemProvider("python", {
      provideCompletionItems: (model, position) => {
        const word = model.getWordUntilPosition(position);
        const range = {
          startLineNumber: position.lineNumber,
          endLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endColumn: word.endColumn,
        };
        return {
          suggestions: [
            {
              label: "print",
              kind: monaco.languages.CompletionItemKind.Function,
              insertText: "print()",
              documentation: "Print to stdout",
              range,
            },
            {
              label: "def",
              kind: monaco.languages.CompletionItemKind.Keyword,
              insertText: "def ",
              range,
            },
            {
              label: "class",
              kind: monaco.languages.CompletionItemKind.Keyword,
              insertText: "class ",
              range,
            },
          ],
        };
      },
    });

    onMount?.(editor, monaco);
  };

  return (
    <div style={{ height: "100%" }}>
      <Editor
        height="100%"
        defaultLanguage="python"
        value={code}
        onChange={onChange}
        onMount={handleMount}
        theme="github-dark" // Use the theme name you defined and set
        options={{
          lineNumbers: (lineNumber: number) =>
            errorLines?.includes(lineNumber) ? "❗" : String(lineNumber),
          fontSize: 14,
          glyphMargin: true,
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          automaticLayout: true,
        }}
      />
    </div>
  );
};

export default MonacoEditor;
