import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { AuthProvider, RequireAuth } from "./components/Auth";
import './App.css';
import SignIn from './components/SignIn';
import SignUp from './components/SignUp';
import Dashboard from './components/Dashboard/Dashboard';
import Scan from './components/Scan/Scan';
import { useSelector } from 'react-redux';
import { ThemeProvider } from '@mui/material/styles';
import themes from './themes/theme';
import Signature from './components/Signature/Signature';
import Signatures from './components/Signatures/Signatures';
import Users from './components/Users/Users';
import Documents from './components/Documents/Documents';

function App() {
  const customization = useSelector((state) => state.customization);

  return (
    <div className="App">
      <BrowserRouter>
      <ThemeProvider theme={themes(customization)}>
      <AuthProvider>
        <Routes>
          <Route 
            protected
            path='dashboard'
            element={
              <RequireAuth>
                <Dashboard />
              </RequireAuth>
            } 
          />
         <Route 
            protected
            path='scan'
            element={
              <RequireAuth>
                <Scan />
              </RequireAuth>
            } 
          />
          <Route 
            protected
            path='signature'
            element={
              <RequireAuth>
                <Signature />
              </RequireAuth>
            } 
          />
          <Route 
            protected
            path='users'
            element={
              <RequireAuth>
                <Users />
              </RequireAuth>
            } 
          />
          <Route 
            protected
            path='signatures'
            element={
              <RequireAuth>
                <Signatures />
              </RequireAuth>
            } 
          />
          <Route 
            protected
            path='documents'
            element={
              <RequireAuth>
                <Documents />
              </RequireAuth>
            } 
          />
          <Route path='login' element={<SignIn />} />
          <Route path='signup' element={<SignUp />} />
        </Routes>
      </AuthProvider>
      </ThemeProvider>
      </BrowserRouter>
    </div>
  );
}

export default App;
