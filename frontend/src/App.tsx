import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { SmartphoneList } from './components/products/SmartphoneList';
import { SmartphoneDetail } from './components/products/SmartphoneDetail';
import { UsersList } from './components/admin/UsersList';
import { CartView } from './components/cart/CartView';
import { LoginForm } from './components/auth/LoginForm';
import { SignupForm } from './components/auth/SignupForm';
import { Layout } from './components/Layout';
import { NotFoundPage } from './components/NotFoundPage';
function App() {
  return (
    <AuthProvider>
      <Router>
        <Layout>
          <Routes>
            <Route path="/login" element={<LoginForm />} />
            <Route path="/signup" element={<SignupForm />} />
            
            <Route path="/" element={
              <ProtectedRoute>
                <SmartphoneList />
              </ProtectedRoute>
            } />
            
            <Route path="/smartphones/:id" element={
              <ProtectedRoute>
                <SmartphoneDetail />
              </ProtectedRoute>
            } />
            
            <Route path="/users" element={
              <ProtectedRoute adminOnly>
                <UsersList />
              </ProtectedRoute>
            } />
            
            <Route path="/cart" element={
              <ProtectedRoute>
                <CartView />
              </ProtectedRoute>
            } />

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </Layout>
      </Router>
    </AuthProvider>
  );
}

export default App;