import React from 'react';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Button } from '@mui/material';
import VisibilityIcon from '@mui/icons-material/Visibility';

const TableComponent = ({ label, data, columns, handleView, handleDownload}) => {
  return (
    <TableContainer component={Paper}>
      <Table aria-label={label}>
        <TableHead>
          <TableRow>
          {handleView ?    <TableCell>Ver</TableCell> : <></>}
            {columns.map((column) => (
              <TableCell key={column.header}>{column.header}</TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {data.length > 0 ? data.map((row) => (
            <TableRow key={row.id}>
              {handleView ?    <TableCell> <Button variant="outlined" color="primary" onClick={() => handleView(row.id)}><VisibilityIcon></VisibilityIcon></Button></TableCell> : <></>}
              {columns.map((column) => (
                <TableCell key={column.header}>{column.cell(row)}</TableCell>
              ))}
            </TableRow>
          )) : <></>}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default TableComponent;
