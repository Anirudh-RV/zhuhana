import { enUS } from "@mui/x-charts/locales";
import {
  DataGrid,
  GridColDef,
  GridRowsProp,
  GridRowParams,
} from "@mui/x-data-grid";

interface CustomizedDataGridProps {
  rows: GridRowsProp;
  columns: GridColDef[];
  onRowDoubleClick?: (params: GridRowParams) => void;
}

export default function CustomizedDataGrid({
  rows,
  columns,
  onRowDoubleClick,
}: CustomizedDataGridProps) {
  return (
    <DataGrid
      checkboxSelection
      rows={rows}
      columns={columns}
      getRowId={(row) => row.id}
      onRowDoubleClick={onRowDoubleClick}
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
        ...enUS.components.MuiChartsLocalizationProvider.defaultProps
          .localeText,
        noRowsLabel: "No Algorithms Created",
      }}
      sx={{
        "& .MuiDataGrid-row:hover": {
          cursor: "pointer",
        },
      }}
    />
  );
}
