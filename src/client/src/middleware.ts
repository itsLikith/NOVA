import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Run middleware for all non-api and non-static paths
export const config = {
  matcher: ['/((?!api|_next|static|favicon.ico).*)'],
}

export async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl

  // allow landing page and auth pages to proceed without validation
  if (pathname === '/' || pathname.startsWith('/auth')) return NextResponse.next()

  // validate via auth service using the request cookies
  try {
    const validateUrl = new URL('/api/v1/auth/validate', req.url)
    const resp = await fetch(validateUrl.toString(), {
      method: 'GET',
      // forward cookie header so auth service can read the token cookie
      headers: { cookie: req.headers.get('cookie') || '' },
    })

    if (resp.ok) {
      // token valid -> if user is on auth pages, redirect to dashboard; otherwise allow
      if (pathname.startsWith('/auth')) {
        return NextResponse.redirect(new URL('/dashboard', req.url))
      }

      return NextResponse.next()
    }

    // invalid -> go to auth (unless already on auth)
    if (pathname.startsWith('/auth')) return NextResponse.next()
    return NextResponse.redirect(new URL('/auth', req.url))
  } catch (err) {
    // on error, redirect to auth (unless already on auth)
    if (pathname.startsWith('/auth')) return NextResponse.next()
    return NextResponse.redirect(new URL('/auth', req.url))
  }
}
