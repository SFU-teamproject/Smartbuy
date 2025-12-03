import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import { LanguageProvider } from './context/LanguageContext';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { SmartphoneList } from './components/products/SmartphoneList';
import { SmartphoneDetail } from './components/products/SmartphoneDetail';
import { UsersList } from './components/admin/UsersList';
import { CartView } from './components/cart/CartView';
import LoginForm from './components/auth/LoginForm';
import SignupForm from './components/auth/SignupForm';
import { Layout } from './components/Layout';
import { NotFoundPage } from './components/NotFoundPage';
import { ProtectedCartRoute } from './components/cart/ProtectedCartRoute';
import { PaymentView } from './components/cart/PaymentView';
import { BankView } from './components/cart/BankView';
import { SuccessView } from './components/cart/SuccessView';
function App() {
  return (
      <LanguageProvider>
    <ThemeProvider>
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

             <Route path="/payment" element={
              <ProtectedRoute>
                <ProtectedCartRoute>
                  <PaymentView />
                </ProtectedCartRoute>
              </ProtectedRoute>
            } />

            <Route path="/bank" element={
              <ProtectedRoute>
                <ProtectedCartRoute>
                  <BankView />
                </ProtectedCartRoute>
              </ProtectedRoute>
            } />

            <Route path="/success" element={
              <ProtectedRoute>
                <SuccessView />
              </ProtectedRoute>
            } />

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </Layout>
      </Router>
    </AuthProvider>
    </ThemeProvider>
    </LanguageProvider>
  );
}

export default App;