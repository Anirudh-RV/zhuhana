// MonacoEditor.tsx
import Editor, { useMonaco } from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react";
import { useEffect, useRef } from "react";
import * as monacoEditor from "monaco-editor";

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

  const handleMount: OnMount = (editor, monacoInstance) => {
    editorRef.current = editor;

    // Configure Python language and theme if needed
    monacoInstance.languages.register({ id: "python" });

    monacoInstance.languages.setMonarchTokensProvider("python", {
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

    monacoInstance.editor.defineTheme("python-dark", {
      base: "vs-dark",
      inherit: true,
      rules: [
        { token: "keyword", foreground: "C586C0", fontStyle: "bold" },
        { token: "string", foreground: "CE9178" },
        { token: "number", foreground: "B5CEA8" },
        { token: "comment", foreground: "6A9955" },
        { token: "identifier", foreground: "9CDCFE" },
      ],
      colors: {},
    });

    monacoInstance.languages.registerCompletionItemProvider("python", {
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
              kind: monacoInstance.languages.CompletionItemKind.Function,
              insertText: "print()",
              documentation: "Print to stdout",
              range,
            },
            {
              label: "def",
              kind: monacoInstance.languages.CompletionItemKind.Keyword,
              insertText: "def ",
              range,
            },
            {
              label: "class",
              kind: monacoInstance.languages.CompletionItemKind.Keyword,
              insertText: "class ",
              range,
            },
          ],
        };
      },
    });

    monacoInstance.editor.setTheme("python-dark");

    onMount?.(editor, monacoInstance);
  };

  return (
    <div style={{ height: "100%" }}>
      <Editor
        height="100%"
        defaultLanguage="python"
        value={code}
        onChange={onChange}
        onMount={handleMount}
        theme="python-dark"
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
