"use client";

import { Button } from "@/components/ui/button";
import { Icons } from "@/components/icons";

export default function LoginPage() {
  const handleGoogleLogin = () => {
    // Redirect to the Go backend OAuth2 login endpoint.
    window.location.href = "http://localhost:8080/api/v1/auth/google/login";
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-950 via-indigo-900 to-indigo-950 px-4">
      <div className="w-full max-w-md bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-2xl flex flex-col items-center">
        {/* Brand Icon or Logo Placeholder */}
        <div className="w-16 h-16 bg-indigo-500 rounded-full flex items-center justify-center mb-6 shadow-lg shadow-indigo-500/30">
          <span className="text-2xl text-white font-bold tracking-tighter">S</span>
        </div>

        <h1 className="text-3xl font-bold text-white mb-2 tracking-tight">
          Welcome to Skillora
        </h1>
        <p className="text-indigo-200 mb-8 text-center text-sm">
          Join the AI-driven barter economy for knowledge. Exchange your skills for credits and learn something new.
        </p>

        <Button
          variant="outline"
          className="w-full h-12 flex flex-row items-center gap-3 text-base text-gray-900 bg-white hover:bg-gray-50 border-gray-200 transition-all font-medium rounded-xl mb-4"
          onClick={handleGoogleLogin}
        >
          <Icons.google className="w-5 h-5" />
          Continue with Google
        </Button>

        <p className="text-xs text-indigo-300 text-center mt-6">
          By signing in, you agree to our Terms of Service and Privacy Policy.
        </p>
      </div>
    </div>
  );
}
