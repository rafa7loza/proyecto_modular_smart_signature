import { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Outlet } from 'react-router-dom';
import { styled, useTheme } from '@mui/material/styles';
import { Alert, AlertTitle, AppBar, IconButton, Box, CssBaseline, Dialog, Grid, Typography, Toolbar, useMediaQuery } from '@mui/material';
import Breadcrumbs from '../UI/Breadcrumbs';
import Header from '../Header/Header';
import Sidebar from '../Sidebar/Sidebar';
import navigation from '../../menu-items';
import { drawerWidth } from '../../store/constant';
import { SET_MENU } from '../../store/actions';
import { IconChevronRight } from '@tabler/icons-react';
import MainCard from '../UI/MainCard';
import SubCard from '../UI/SubCard';
import { FilePond, registerPlugin } from 'react-filepond';
import FilePondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';
import { Document, Page, Image, StyleSheet } from '@react-pdf/renderer';
import CloseIcon from '@mui/icons-material/Close';
import PDFViewer from '../UI/PDFViewer';
import 'filepond/dist/filepond.min.css';
import MainTable from '../UI/MainTable';
import { REST_API } from '../../endpoint_urls';
import API from '../../Api';

registerPlugin(FilePondPluginImageExifOrientation, FilePondPluginImagePreview);

// styles
const Main = styled('main', { shouldForwardProp: (prop) => prop !== 'open' })(({ theme, open }) => ({
  ...theme.typography.mainContent,
  borderBottomLeftRadius: 0,
  borderBottomRightRadius: 0,
  transition: theme.transitions.create(
    'margin',
    open
      ? {
        easing: theme.transitions.easing.easeOut,
        duration: theme.transitions.duration.enteringScreen
      }
      : {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen
      }
  ),
  [theme.breakpoints.up('md')]: {
    marginLeft: open ? 0 : -(drawerWidth - 20),
    width: `calc(100% - ${drawerWidth}px)`
  },
  [theme.breakpoints.down('md')]: {
    marginLeft: '20px',
    width: `calc(100% - ${drawerWidth}px)`,
    padding: '16px'
  },
  [theme.breakpoints.down('sm')]: {
    marginLeft: '10px',
    width: `calc(100% - ${drawerWidth}px)`,
    padding: '16px',
    marginRight: '10px'
  }
}));

const Scan = () => {
  const theme = useTheme();
  const matchDownMd = useMediaQuery(theme.breakpoints.down('md'));
  // Handle left drawer
  const leftDrawerOpened = useSelector((state) => state.customization.opened);
  const dispatch = useDispatch();
  const [tableData, setTableData] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [base64Image, setBase64Image] = useState('');

  const closeModal = () => {
    setIsModalOpen(false);
  };

  const handleLeftDrawerToggle = () => {
    dispatch({ type: SET_MENU, opened: !leftDrawerOpened });
  };

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    API.get(REST_API.endpoints.userDocuments)
      .then(response => {
        setTableData(response?.data?.documents);
        setIsLoading(false);
      })
      .catch(error => {
        console.error('Error fetching data:', error);

        setIsLoading(false);
      });
  };

  const [files, setFiles] = useState([]);
  const [pond, setPond] = useState(null);

  const tableColumns = [
    { header: 'ID', cell: (row) => row.id },
    { header: 'Nombre', cell: (row) => row.file_name },
  ];

  const handleView = (id) => {
    const document = tableData.find(doc => doc.id === id);
    setBase64Image(document.encoded_content);
    setIsModalOpen(true);
  };


  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      {/* header */}
      <AppBar
        enableColorOnDark
        position="fixed"
        color="inherit"
        elevation={0}
        sx={{
          bgcolor: theme.palette.background.default,
          transition: leftDrawerOpened ? theme.transitions.create('width') : 'none'
        }}
      >
        <Toolbar>
          <Header handleLeftDrawerToggle={handleLeftDrawerToggle} />
        </Toolbar>
      </AppBar>

      {/* drawer */}
      <Sidebar drawerOpen={!matchDownMd ? leftDrawerOpened : !leftDrawerOpened} drawerToggle={handleLeftDrawerToggle} />

      {/* main content */}
      <Main theme={theme} open={leftDrawerOpened}>
        {/* breadcrumb */}
        <Breadcrumbs separator={IconChevronRight} navigation={navigation} icon title rightAlign />
        <MainCard title="Captura tus documento">
          <Dialog
            fullScreen
            open={isModalOpen}
            onClose={closeModal}
          >
            <AppBar position="relative" color="default">
              <IconButton
                  edge="start"
                  color="inherit"
                  onClick={closeModal}
                  aria-label="close"
                >
                  <CloseIcon />
                </IconButton>

        </AppBar>
              <img src={`data:image/jpeg;base64,${base64Image}`} />

          </Dialog>
          <FilePond
            files={files}
            ref={(ref) => (setPond(ref))}
            onupdatefiles={setFiles}
            allowMultiple={false}
            allowRemove={false}
            allowRevert={false}
            onprocessfiles={() => { pond.removeFiles(); }}
            server={{
              process: (fieldName, file, metadata, load, error, progress, abort, transfer, options) => {
                const controller = new AbortController();
                const formData = new FormData();
                formData.append(fieldName, file, file.name);
                console.log('Uploading file:', file.name);
                API.post(REST_API.endpoints.uploadFile, formData, {
                  signal: controller.signal
                }).then(res => {
                  console.log(formData);
                  console.log(res);
                  load(res.data);
                })
                  .catch(err => {
                    console.log(formData);
                    console.log(err);
                    error(err);
                  });
                return {
                  abort: () => {
                    controller.abort();
                    // Let FilePond know the request has been cancelled
                    abort();
                  },
                };  
              },
            }}
            credits={false}
            name="files" /* sets the file input name, it's filepond by default */
            labelIdle='Arrastra tus documentos aquí o <span class="filepond--label-action">Súbelo desde tu computadora</span>'
          />
          <Grid item xs={12}>
            <SubCard title="Documentos previos">
              {isLoading ? (
                <p>Loading...</p>
              ) : (
                <MainTable label="Tabla de documentos" data={tableData} columns={tableColumns} handleView={handleView}></MainTable>
              )}
            </SubCard>
          </Grid>
        </MainCard>
      </Main>
    </Box>
  );
};

export default Scan;
