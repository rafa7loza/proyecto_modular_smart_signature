import { useContext, createContext } from 'react';
import { useLocation, Navigate } from 'react-router-dom'
import { useJwt } from 'react-jwt';
import { REST_API } from '../endpoint_urls';
import API from '../Api'


export let AuthContext = createContext();
const DEFAULT_TIMEOUT = 100;

const fakeAuthProvider = {
  signin(callback) {
    setTimeout(callback, DEFAULT_TIMEOUT);
  }
}

export function AuthProvider({children}) {
  const signin = (loginForm, callback) => {
    return fakeAuthProvider.signin(() => {
      API.post(REST_API.endpoints.login, loginForm)
      .then(res => {
        const token = res.data.token;
        localStorage.setItem('token', token);
        callback();
      })
      .catch(res => {
        console.log('error al entrar', res);
      })
    })
  }

  const auth = {
    signin
  };

  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>
}

export function useAuth() {
  let auth = useContext(AuthContext);
  return auth;  
}

export function useAuthToken() {
  const token = localStorage.getItem('token');
  const { decodedToken, isExpired } = useJwt(token);
  return token ? {decodedToken, isExpired} : {decodedToken: null, isExpired: true};
}

export function RequireAuth({children}) {
  let location = useLocation();
  let authToken = useAuthToken();

  if (authToken.isExpired) {
    return <Navigate to='/login' state={{ from: location }}></Navigate>
  }

  return children;
}
