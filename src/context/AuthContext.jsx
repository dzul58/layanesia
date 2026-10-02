import React, { createContext, useContext, useState, useEffect } from 'react';
import { api } from '../api';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem('layanesia_token') || null);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const logout = () => {
    localStorage.removeItem('layanesia_token');
    setToken(null);
    setUser(null);
  };

  const fetchMe = async (authToken) => {
    const activeToken = authToken || token;
    if (!activeToken) {
      setUser(null);
      setLoading(false);
      return;
    }

    try {
      const data = await api('/api/v1/users/me', { token: activeToken });
      setUser(data.user || null);
    } catch (err) {
      if (err.status === 401) {
        logout();
      } else {
        console.error('Error fetching user profile:', err);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMe(token);
  }, [token]);

  const login = async (email, password) => {
    const data = await api('/api/v1/auth/login', {
      method: 'POST',
      body: { email, password },
    });
    localStorage.setItem('layanesia_token', data.token);
    setToken(data.token);
    setUser(data.user);
    return data;
  };

  const register = async (registerData) => {
    const data = await api('/api/v1/auth/register', {
      method: 'POST',
      body: registerData,
    });
    localStorage.setItem('layanesia_token', data.token);
    setToken(data.token);
    setUser(data.user);
    return data;
  };

  const switchMode = async (newMode) => {
    if (!token) return;
    const data = await api('/api/v1/users/switch-mode', {
      method: 'PATCH',
      token,
      body: { mode: newMode },
    });
    if (data.token) {
      localStorage.setItem('layanesia_token', data.token);
      setToken(data.token);
    }
    setUser(data.user);
    return data;
  };

  const updateProfile = async (profileData) => {
    if (!token) return;
    const data = await api('/api/v1/users/profile', {
      method: 'PUT',
      token,
      body: profileData,
    });
    setUser(data.user);
    return data;
  };

  const verifyKTP = async (ktpNumber, ktpImageURL) => {
    if (!token) return;
    const data = await api('/api/v1/users/verify-ktp', {
      method: 'POST',
      token,
      body: { ktp_number: ktpNumber, ktp_image_url: ktpImageURL },
    });
    setUser(data.user);
    return data;
  };

  const uploadResume = async (resumeURL) => {
    if (!token) return;
    const data = await api('/api/v1/users/upload-resume', {
      method: 'POST',
      token,
      body: { resume_url: resumeURL },
    });
    setUser(data.user);
    return data;
  };

  const uploadFile = async (file, purpose) => {
    if (!token) return;
    const formData = new FormData();
    formData.append('file', file);
    if (purpose) formData.append('purpose', purpose);
    const data = await api('/api/v1/users/upload-file', {
      method: 'POST',
      token,
      formData,
    });
    return data.file_url;
  };

  return (
    <AuthContext.Provider
      value={{
        token,
        user,
        loading,
        login,
        register,
        logout,
        switchMode,
        updateProfile,
        verifyKTP,
        uploadResume,
        uploadFile,
        refreshUser: () => fetchMe(token),
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
