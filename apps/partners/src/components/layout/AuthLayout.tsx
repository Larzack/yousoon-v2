import { Outlet } from 'react-router-dom'

export function AuthLayout() {
  return (
    <div className="min-h-screen flex">
      {/* Left side - Brand */}
      <div className="hidden lg:flex lg:w-1/2 bg-yousoon-black items-center justify-center p-12">
        <div className="max-w-md text-center">
          <div className="mb-8">
            <svg
              className="w-20 h-20 mx-auto text-yousoon-gold"
              viewBox="0 0 100 100"
              fill="currentColor"
            >
              <circle cx="50" cy="50" r="45" stroke="currentColor" strokeWidth="2" fill="none" />
              <text
                x="50"
                y="55"
                textAnchor="middle"
                fontSize="24"
                fontWeight="bold"
                fill="currentColor"
              >
                YS
              </text>
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-white mb-4">
            Yousoon Business
          </h1>
          <p className="text-gray-400 text-lg">
            Gérez vos offres, suivez vos performances et développez votre activité
            avec Yousoon.
          </p>
        </div>
      </div>

      {/* Right side - Auth form */}
      <div className="flex-1 flex items-center justify-center p-8 bg-background">
        <div className="w-full max-w-md">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
