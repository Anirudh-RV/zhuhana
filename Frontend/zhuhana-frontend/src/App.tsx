import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import HomePage from "./home/HomePage";
import Blog from "./home/Blog";
import CompanyPage from "./home/Company";
import DemoPage from "./home/Demo";
import FeaturesPage from "./home/Features";
import PricingPage from "./home/Pricing";
import PrivacyPage from "./home/Privacy";
import TermsPage from "./home/Terms";
import ContactPage from "./home/Contact";

import Login from "./login/login";
import SignUp from "./signup/signup";
import Code from "./code/Code";
import Dashboard from "./dashboard/Dashboard";
import { AuthProvider } from "./AuthContext";
import ProtectedRoute from "./ProtectedRoute";
import PublicRoute from "./PublicRoute";

import "./App.css";
import "highlight.js/styles/github-dark.css";

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public routes — redirect to dashboard if user is already logged in */}
          <Route
            path="/"
            element={
              <PublicRoute>
                <HomePage />
              </PublicRoute>
            }
          />
          <Route
            path="/login"
            element={
              <PublicRoute>
                <Login />
              </PublicRoute>
            }
          />
          <Route
            path="/signup"
            element={
              <PublicRoute>
                <SignUp />
              </PublicRoute>
            }
          />
          <Route
            path="/blog"
            element={
              <PublicRoute>
                <Blog />
              </PublicRoute>
            }
          />
          <Route
            path="/company"
            element={
              <PublicRoute>
                <CompanyPage />
              </PublicRoute>
            }
          />
          <Route
            path="/contact"
            element={
              <PublicRoute>
                <ContactPage />
              </PublicRoute>
            }
          />
          <Route
            path="/demo"
            element={
              <PublicRoute>
                <DemoPage />
              </PublicRoute>
            }
          />
          <Route
            path="/features"
            element={
              <PublicRoute>
                <FeaturesPage />
              </PublicRoute>
            }
          />
          <Route
            path="/pricing"
            element={
              <PublicRoute>
                <PricingPage />
              </PublicRoute>
            }
          />
          <Route
            path="/privacy"
            element={
              <PublicRoute>
                <PrivacyPage />
              </PublicRoute>
            }
          />
          <Route
            path="/terms"
            element={
              <PublicRoute>
                <TermsPage />
              </PublicRoute>
            }
          />

          {/* Protected routes — only for authenticated users */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/code"
            element={
              <ProtectedRoute>
                <Code />
              </ProtectedRoute>
            }
          />

          {/* Catch-all fallback */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
