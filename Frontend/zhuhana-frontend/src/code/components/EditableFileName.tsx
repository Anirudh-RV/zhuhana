import { TextField, Typography } from "@mui/material";
import { useState, useRef, useEffect } from "react";

export default function EditableFileName({
  name,
  onRename,
}: {
  name: string;
  onRename: (newName: string) => void;
}) {
  const [isEditing, setIsEditing] = useState(false);
  const [draftName, setDraftName] = useState(name);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  const handleBlurOrSubmit = () => {
    setIsEditing(false);
    if (draftName.trim() && draftName !== name) {
      onRename(draftName.trim());
    } else {
      setDraftName(name); // reset if unchanged
    }
  };

  return isEditing ? (
    <TextField
      inputRef={inputRef}
      variant="standard"
      value={draftName}
      onChange={(e) => setDraftName(e.target.value)}
      onBlur={handleBlurOrSubmit}
      onKeyDown={(e) => {
        if (e.key === "Enter") handleBlurOrSubmit();
        if (e.key === "Escape") {
          setIsEditing(false);
          setDraftName(name);
        }
      }}
      size="small"
      sx={{
        fontWeight: "bold",
        input: { textAlign: "center" },
        width: "200px",
      }}
    />
  ) : (
    <Typography
      variant="subtitle1"
      fontWeight="bold"
      onDoubleClick={() => setIsEditing(true)}
      sx={{
        cursor: "pointer",
        textAlign: "center",
        whiteSpace: "nowrap",
        overflow: "hidden",
        textOverflow: "ellipsis",
        maxWidth: "200px",
      }}
    >
      {name}
    </Typography>
  );
}
