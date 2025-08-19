import {
  DataGrid,
  GridColDef,
  GridRowsProp,
  GridRowParams,
  GridRenderCellParams,
} from "@mui/x-data-grid";
import { useState, Fragment } from "react";
import { Box } from "@mui/material";

interface CustomizedDataGridProps {
  rows: GridRowsProp;
  columns: GridColDef[];
  onRowClick?: (params: GridRowParams) => void;
  onRowDoubleClick?: (params: GridRowParams) => void;
  renderRowDetails?: (rowId: string) => React.ReactNode;
}

export default function CustomizedDataGrid({
  rows,
  columns,
  onRowClick,
  onRowDoubleClick,
  renderRowDetails,
}: CustomizedDataGridProps) {
  const [expandedRowId, setExpandedRowId] = useState<string | null>(null);

  const handleRowClick = (params: GridRowParams) => {
    // Toggle expansion state
    setExpandedRowId((prev) => (prev === params.row.id ? null : params.row.id));
    onRowClick?.(params); // Preserve parent click logic
  };

  const handleRowDoubleClick = (params: GridRowParams) => {
    setExpandedRowId((prev) => (prev === params.row.id ? null : params.row.id));
    onRowDoubleClick?.(params); // preserve original double click behavior
  };

  // Augment rows with detail panel manually
  const displayedRows: GridRowsProp = rows.flatMap((row) => {
    const isExpanded = expandedRowId === row.id;
    return isExpanded
      ? [
          row,
          {
            id: `${row.id}-details`,
            isDetailRow: true,
            parentId: row.id,
          } as any,
        ]
      : [row];
  });

  // Add a "renderCell" to display detail panel in a full-width cell
  const augmentedColumns: GridColDef[] = [
    ...columns,
    {
      field: "__details__",
      headerName: "",
      sortable: false,
      filterable: false,
      disableColumnMenu: true,
      renderCell: (params: GridRenderCellParams) => {
        if (!params.row.isDetailRow || !renderRowDetails) return null;

        return (
          <Box sx={{ width: "100%" }}>
            {renderRowDetails(params.row.parentId)}
          </Box>
        );
      },
      flex: 1,
      minWidth: 300,
    },
  ];

  return (
    <DataGrid
      rows={displayedRows}
      columns={augmentedColumns}
      getRowId={(row) => row.id}
      getRowHeight={(params) => (params.model.isDetailRow ? "auto" : null)}
      onRowClick={handleRowClick}
      onRowDoubleClick={handleRowDoubleClick}
      getRowClassName={(params) =>
        params.indexRelativeToCurrentPage % 2 === 0 ? "even" : "odd"
      }
      initialState={{
        pagination: { paginationModel: { pageSize: 20 } },
      }}
      pageSizeOptions={[10, 20, 50]}
      disableColumnResize
      density="compact"
      localeText={{
        noRowsLabel: "No Algorithms Created",
      }}
      getCellClassName={(params) => {
        if (params.row.isDetailRow && params.field === "__details__") {
          return "detail-cell";
        }
        return "";
      }}
      sx={{
        "& .MuiDataGrid-row:hover": {
          cursor: "pointer",
        },
        "& .MuiDataGrid-cell": {
          borderBottom: "none",
        },
        "& .MuiDataGrid-row.Mui-odd": {
          backgroundColor: "#f9f9f9",
        },

        // 👇 Ensures that detail rows don't have extra padding
        "& .MuiDataGrid-row .MuiDataGrid-cell.detail-cell": {
          padding: 0,
          border: "none",
        },

        "& .subtable-wrapper": {
          width: "100%",
          padding: "16px 24px",
          boxSizing: "border-box",
          backgroundColor: "#fff", // optional: to visually separate
        },
      }}
    />
  );
}
