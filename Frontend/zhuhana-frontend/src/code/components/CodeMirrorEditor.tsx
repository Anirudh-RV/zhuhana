import CodeMirror from "@uiw/react-codemirror";
import { python } from "@codemirror/lang-python";
import { githubDark } from "@uiw/codemirror-themes-all";
import { EditorView } from "@codemirror/view";
import type { Extension } from "@codemirror/state";
import { completionKeymap } from "@codemirror/autocomplete";
import { keymap } from "@codemirror/view";
import { lineNumbers, highlightActiveLineGutter } from "@codemirror/view";

type CodeMirrorEditorProps = {
  code: string;
  onChange: (value: string) => void;
  onCreateEditor: (view: EditorView) => void;
  extraExtensions?: Extension[];
};

const CodeMirrorEditor: React.FC<CodeMirrorEditorProps> = ({
  code,
  onChange,
  onCreateEditor,
  extraExtensions = [],
}) => {
  const baseExtensions: Extension[] = [
    python(),
    lineNumbers(),
    highlightActiveLineGutter(),
    EditorView.lineWrapping,
    keymap.of(completionKeymap), // Standard keybindings for autocompletion
  ];

  return (
    <div
      style={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        flex: 1,
        overflow: "hidden",
      }}
    >
      <CodeMirror
        value={code}
        height="100%"
        theme={githubDark}
        extensions={[...baseExtensions, ...extraExtensions]}
        onChange={onChange}
        onCreateEditor={onCreateEditor}
        style={{ flex: 1, height: "100%" }} // ✅ Ensures the inner editor stretches
        basicSetup={{
          lineNumbers: false, // optional: disable if redundant
        }}
      />
    </div>
  );
};

export default CodeMirrorEditor;
