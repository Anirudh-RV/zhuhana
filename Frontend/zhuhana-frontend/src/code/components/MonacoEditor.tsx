import React from "react";
import Editor from "@monaco-editor/react";

type MonacoEditorProps = {
  code: string;
  onChange: (value: string | undefined) => void;
};

const MonacoEditor: React.FC<MonacoEditorProps> = ({ code, onChange }) => {
  return (
    <div style={{ height: "80vh", border: "1px solid #ccc" }}>
      <Editor
        height="100vw"
        defaultLanguage="python"
        defaultValue={code}
        theme="vs-dark"
        onChange={onChange}
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          automaticLayout: true,
        }}
      />
    </div>
  );
};

export default MonacoEditor;
