import { useEffect } from 'react'
import { useAuth } from 'react-oidc-context'
import { Dashboard } from '@/components/dashboard/Dashboard'
import { setAuthToken, clearAuthToken } from '@/api/client'

function App() {
  const auth = useAuth();

  useEffect(() => {
    if (!auth.isLoading && !auth.isAuthenticated && !auth.error) {
      console.log('No active session, redirecting to login...');
      auth.signinRedirect();
    }
  }, [auth.isLoading, auth.isAuthenticated, auth.error, auth.signinRedirect]);

  useEffect(() => {
    if (auth.isAuthenticated && auth.user?.access_token) {
      setAuthToken(auth.user.access_token);
    } else {
      clearAuthToken();
    }
  }, [auth.isAuthenticated, auth.user]);

  if (auth.isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading authentication...</div>
      </div>
    );
  }

  if (auth.error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg text-red-500">Authentication error: {auth.error.message}</div>
      </div>
    );
  }

  if (!auth.isAuthenticated) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen gap-4">
        <h1 className="text-2xl font-bold">Redirecting to login...</h1>
      </div>
    );
  }

  // Set the token synchronously during render if we're authenticated.
  // This ensures that when <Dashboard /> mounts and immediately calls the API,
  // the token is already present in the client middleware.
  if (auth.user?.access_token) {
    setAuthToken(auth.user.access_token);
  }

  return <Dashboard />
}

export default App
