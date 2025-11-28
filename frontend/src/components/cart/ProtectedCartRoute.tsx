import { useAuth } from '../../context/AuthContext';
import { FC, ReactNode } from "react";
import { Navigate } from "react-router-dom";
interface ProtectedCartRoutePropse {
    children : ReactNode
}
export const ProtectedCartRoute = ({children} : ProtectedCartRoutePropse) => {
    const {cart} = useAuth();
    console.log(cart)
    if (cart) {
        return <>{children}</>;
    } 
    return <Navigate to="/cart" replace></Navigate>
}