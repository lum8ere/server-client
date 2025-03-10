import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { jwtDecode } from 'jwt-decode';
import instance from 'service/api';

export interface DecodedToken {
    user_id: string;
    username: string;
    role: 'OBSERVER' | 'OBSERVER_PLUS' | 'ADMIN';
    exp: number;
    iat: number;
}

export interface AuthState {
    token: string | null;
    user: DecodedToken | null;
    loading: boolean;
    error: string | null;
}

const initialState: AuthState = {
    token: localStorage.getItem('token') || null,
    user: localStorage.getItem('token')
        ? jwtDecode<DecodedToken>(localStorage.getItem('token') as string)
        : null,
    loading: false,
    error: null
};

interface LoginCredentials {
    email: string;
    password: string;
}

// Async thunk для логина
export const loginUser = createAsyncThunk<string, LoginCredentials>(
    'auth/loginUser',
    async (credentials, { rejectWithValue }) => {
        try {
            const response = await instance.post('/api/auth/login', credentials);
            const token: string = response.data.token;
            return token;
        } catch (err: any) {
            return rejectWithValue(err.response?.data || 'Login failed');
        }
    }
);

export const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        logout(state) {
            state.token = null;
            state.user = null;
            state.error = null;
            localStorage.removeItem('token');
        },
        setToken(state, action: PayloadAction<string>) {
            state.token = action.payload;
            try {
                state.user = jwtDecode<DecodedToken>(action.payload);
            } catch (error) {
                state.user = null;
            }
            localStorage.setItem('token', action.payload);
        },
        loadToken(state) {
            const token = localStorage.getItem('token');
            if (token) {
                state.token = token;
                try {
                    state.user = jwtDecode<DecodedToken>(token);
                } catch (error) {
                    state.user = null;
                }
            }
        }
    },
    extraReducers: (builder) => {
        builder.addCase(loginUser.pending, (state) => {
            state.loading = true;
            state.error = null;
        });
        builder.addCase(loginUser.fulfilled, (state, action) => {
            state.loading = false;
            state.token = action.payload;
            try {
                state.user = jwtDecode<DecodedToken>(action.payload);
            } catch (error) {
                state.user = null;
            }
            localStorage.setItem('token', action.payload);
        });
        builder.addCase(loginUser.rejected, (state, action) => {
            state.loading = false;
            state.error = action.payload as string;
        });
    }
});

export const { logout, setToken, loadToken } = authSlice.actions;
export default authSlice.reducer;
