import { TextField, Typography, Tooltip } from "@mui/material";
import {
  useState,
  useRef,
  useEffect,
  useImperativeHandle,
  forwardRef,
} from "react";

export interface EditableFileNameHandle {
  focusEditMode: () => void;
}

const EditableFileName = forwardRef<
  EditableFileNameHandle,
  {
    name: string;
    onRename: (newName: string) => void;
  }
>(({ name, onRename }, ref) => {
  const [isEditing, setIsEditing] = useState(false);
  const [draftName, setDraftName] = useState(name);
  const [showTooltip, setShowTooltip] = useState(false); // 👈 manually control tooltip
  const inputRef = useRef<HTMLInputElement>(null);
  const hasShownTooltip = useRef(false);
  const [hasAlgorithmId, setHasAlgorithmId] = useState(true);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const hasId = urlParams.has("algorithm_id");

    setHasAlgorithmId(hasId); // store it in state
    console.log("ALGORITHM_ID: ", hasId);

    if (!hasId && !hasShownTooltip.current) {
      hasShownTooltip.current = true;
      setShowTooltip(true);
      const timer = setTimeout(() => setShowTooltip(false), 10000);
      return () => clearTimeout(timer);
    }
  }, []);

  const handleBlurOrSubmit = () => {
    setIsEditing(false);
    if (draftName.trim() && draftName !== name) {
      onRename(draftName.trim());
    } else {
      setDraftName(name);
    }
  };

  useImperativeHandle(ref, () => ({
    focusEditMode: () => setIsEditing(true),
  }));

  return isEditing ? (
    <Tooltip
      title="Save the algorithm with a new name"
      placement="bottom"
      arrow
      open={showTooltip} // 👈 control manually
      slotProps={{
        tooltip: {
          sx: {
            fontSize: "1rem",
            padding: "8px 12px",
            maxWidth: "none",
          },
        },
      }}
    >
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
    </Tooltip>
  ) : (
    <Tooltip
      title="Double-click to rename"
      arrow
      open={showTooltip}
      disableFocusListener
      disableHoverListener
      disableTouchListener
    >
      <Typography
        variant="subtitle1"
        fontWeight="bold"
        onDoubleClick={() => setIsEditing(true)}
        onMouseEnter={() => setShowTooltip(true)}
        onMouseLeave={() => setShowTooltip(false)}
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
    </Tooltip>
  );
});

export default EditableFileName;
