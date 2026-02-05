import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { AuthProvider } from 'react-oidc-context'
import type { AuthProviderProps } from 'react-oidc-context'
import './index.css'
import App from './App.tsx'

const oidcConfig: AuthProviderProps = {
  authority: import.meta.env.VITE_ZITADEL_AUTHORITY || "https://auth.homelab.chapar.tech",
  client_id: import.meta.env.VITE_ZITADEL_CLIENT_ID || "358141840966287363",
  redirect_uri: window.location.origin,
  post_logout_redirect_uri: window.location.origin,
  scope: "openid profile email",
  onSigninCallback: () => {
    window.history.replaceState({}, document.title, window.location.pathname);
  },
};

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider {...oidcConfig}>
      <App />
    </AuthProvider>
  </StrictMode>,
)
