import CodeMirror from "@uiw/react-codemirror";
import { python } from "@codemirror/lang-python";
import { githubDark, githubLight } from "@uiw/codemirror-themes-all";
import { EditorView } from "@codemirror/view";
import type { Extension } from "@codemirror/state";
import { completionKeymap } from "@codemirror/autocomplete";
import { keymap } from "@codemirror/view";
import { lineNumbers, highlightActiveLineGutter } from "@codemirror/view";
import { useColorScheme } from "@mui/material/styles";
import { autocompletion, acceptCompletion } from "@codemirror/autocomplete";
import { indentMore, indentWithTab } from "@codemirror/commands";

type CodeMirrorEditorProps = {
  code: string;
  onChange: (value: string) => void;
  onCreateEditor: (view: EditorView) => void;
  extraExtensions?: Extension[];
};

const customTabKey = (view: EditorView) => {
  return acceptCompletion(view) || indentWithTab.run?.(view) || false;
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
    autocompletion({ activateOnTyping: true }),
    keymap.of([
      {
        key: "Tab",
        run: (view) => {
          // If a completion is active, accept it
          if (acceptCompletion(view)) return true;
          // Otherwise, indent normally
          return indentWithTab.run?.(view) ?? false;
        },
        preventDefault: true,
      },
      {
        key: "Enter",
        run: acceptCompletion,
        preventDefault: true,
      },
    ]),
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
