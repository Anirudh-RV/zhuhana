import CodeMirror from "@uiw/react-codemirror";
import { python } from "@codemirror/lang-python";
import { githubDark, githubLight } from "@uiw/codemirror-themes-all";
import { EditorView } from "@codemirror/view";
import type { Extension } from "@codemirror/state";
import { completionKeymap } from "@codemirror/autocomplete";
import { keymap } from "@codemirror/view";
import { lineNumbers, highlightActiveLineGutter } from "@codemirror/view";
import { useColorScheme } from "@mui/material/styles";

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

  const { mode, systemMode } = useColorScheme();
  const resolvedMode = mode === "system" ? systemMode : mode;

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
        theme={resolvedMode === "dark" ? githubDark : githubLight}
        extensions={[...baseExtensions, ...extraExtensions]}
        onChange={onChange}
        onCreateEditor={onCreateEditor}
        style={{
          flex: 1,
          height: "100%",
        }}
        basicSetup={{
          lineNumbers: false, // optional: disable if redundant
        }}
      />
    </div>
  );
};

export default CodeMirrorEditor;
