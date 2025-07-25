import { TextField, Typography, Tooltip, Box } from "@mui/material";
import {
  useState,
  useRef,
  useEffect,
  useImperativeHandle,
  forwardRef,
} from "react";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";

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
  const [showTooltip, setShowTooltip] = useState(false);
  const [hovering, setHovering] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const hasShownTooltip = useRef(false);
  const [hasAlgorithmId, setHasAlgorithmId] = useState(true);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const hasId = urlParams.has("algorithm_id");
    setHasAlgorithmId(hasId);

    if (!hasId && !hasShownTooltip.current) {
      hasShownTooltip.current = true;
      setShowTooltip(true);
      const timer = setTimeout(() => setShowTooltip(false), 10000);
      return () => clearTimeout(timer);
    }
  }, []);

  useImperativeHandle(ref, () => ({
    focusEditMode: () => setIsEditing(true),
  }));

  useEffect(() => {
    setDraftName(name);
  }, [name]);

  const handleBlurOrSubmit = () => {
    setIsEditing(false);
    setHovering(false); // 👈 Reset hover state

    if (draftName.trim() && draftName !== name) {
      onRename(draftName.trim());
    } else {
      setDraftName(name);
    }
  };

  return isEditing ? (
    <Tooltip
      title="Save the algorithm with a new name"
      placement="bottom"
      arrow
      open={showTooltip}
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
            setHovering(false); // 👈 Reset hover state
          }
        }}
        size="small"
        autoFocus
        sx={{
          fontWeight: "bold",
          input: {
            textAlign: "center",
            fontSize: "1rem",
            fontWeight: "bold",
          },
          width: "200px",
        }}
      />
    </Tooltip>
  ) : (
    <Box
      onMouseEnter={() => setHovering(true)}
      onMouseLeave={() => setHovering(false)}
      onClick={() => setIsEditing(true)}
      sx={{
        cursor: "text", // I-beam
        borderRadius: "4px",
        padding: "2px 6px",
        border: hovering ? "1px solid #1976d2" : "1px solid transparent",
        transition: "border 0.2s",
        display: "inline-block",
        maxWidth: "200px",
      }}
    >
      <Typography
        variant="subtitle1"
        fontWeight="bold"
        sx={{
          textAlign: "center",
          whiteSpace: "nowrap",
          overflow: "hidden",
          textOverflow: "ellipsis",
        }}
      >
        {name}
      </Typography>
    </Box>
  );
});

export default EditableFileName;
