import Editor from "@monaco-editor/react";
import type { OnMount } from "@monaco-editor/react/";
import * as monacoEditor from "monaco-editor";
import githubDark from "monaco-themes/themes/GitHub Dark.json";
import { useRef, useEffect } from "react";
import { createPythonLanguageClient } from "./lspClient";
import { Parser, Language } from "web-tree-sitter";

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

  // ✅ Tree-sitter setup: runs once after mount
  useEffect(() => {
    const setupTreeSitter = async () => {
      // Init WASM runtime
      await Parser.init({ locateFile: () => "/tree-sitter.wasm" });

      // Load Python grammar
      const parser = new Parser();
      const lang = await Language.load("/tree-sitter-python.wasm");
      parser.setLanguage(lang);

      // Get editor content and parse
      const model = editorRef.current?.getModel();
      if (!model) return;

      const code = model.getValue();
      const tree = parser.parse(code);
      const rootNode = tree?.rootNode;

      // Example: log function definitions
      rootNode?.namedChildren
        .filter((node) => node?.type === "function_definition")
        .forEach((node) => {
          const nameNode = node?.namedChildren.find(
            (n) => n?.type === "identifier"
          );
          console.log("Function:", nameNode?.text);
        });
    };

    if (editorRef.current) {
      setupTreeSitter().catch(console.error);
    }
  }, []);

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor;

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

    monaco.editor.defineTheme(
      "github-dark",
      githubDark as monacoEditor.editor.IStandaloneThemeData
    );
    monaco.editor.setTheme("github-dark");

    // ✅ Start LSP client
    const ws = new WebSocket("ws://localhost:3001");

    ws.onopen = () => {
      const client = createPythonLanguageClient(); // assumes messageTransports setup inside
      client.start();
    };

    const model = editor.getModel();
    if (model) monaco.editor.setModelLanguage(model, "python");

    onMount?.(editor, monaco);
  };

  return (
    <div style={{ height: "100%" }}>
      <Editor
        height="100%"
        defaultLanguage="python"
        defaultValue={code}
        onChange={onChange}
        onMount={handleMount}
        theme="github-dark"
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
