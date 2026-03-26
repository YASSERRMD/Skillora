import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Auth routes that users must be logged in to access.
const protectedRoutes = ["/dashboard", "/workspace", "/marketplace", "/onboarding"];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Check if current route is protected.
  const isProtected = protectedRoutes.some((route) => pathname.startsWith(route));

  if (isProtected) {
    // Check for the JWT cookie. The robust verification happens on the Go backend,
    // but the presence of the cookie is enough for edge middleware to decide on routing.
    const token = request.cookies.get("skillora_token");

    if (!token?.value) {
      // Redirect to login, preserving the original URL to redirect back after auth (optional enhancement).
      const url = request.nextUrl.clone();
      url.pathname = "/login";
      return NextResponse.redirect(url);
    }
  }

  return NextResponse.next();
}

// Config ensures middleware only runs on necessary paths, matching protected routes
// and excluding static files, API routes, Next.js internals, etc.
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
};
