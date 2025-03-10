import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from 'store';

interface PrivateRouteProps {
    requiredRole?: string; // Например, 'admin'
    // если children переданы напрямую
    children?: React.ReactNode;
}

export const PrivateRoute: React.FC<PrivateRouteProps> = ({ requiredRole, children }) => {
    const { token, user } = useSelector((state: RootState) => state.auth);

    // Если токена или данных пользователя нет, перенаправляем на /login
    if (!token || !user) {
        return <Navigate to="/login" replace />;
    }

    // Если требуется проверка роли и роль не соответствует, показываем сообщение
    if (requiredRole && user.role.toLowerCase() !== requiredRole.toLowerCase()) {
        return <div>You do not have sufficient privileges to access this page.</div>;
    }

    return children ? <>{children}</> : <Outlet />;
};
