import { useState } from 'react'
import AdminLogin from './AdminLogin'

interface FooterProps {
  onAdminLogin: (password: string) => void
  isAdmin: boolean
}

export default function Footer({ onAdminLogin, isAdmin }: FooterProps) {
  const [showLogin, setShowLogin] = useState(false)

  return (
    <footer className="bg-gray-900 text-white mt-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-4 md:mb-0">
            <p className="text-gray-400">© 2025 AppDirect India AI Workshop. All rights reserved.</p>
          </div>
          <div>
            {!isAdmin ? (
              <button
                onClick={() => setShowLogin(true)}
                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg transition-colors duration-200"
              >
                Admin Login
              </button>
            ) : (
              <p className="text-green-400">Admin Mode Active</p>
            )}
          </div>
        </div>
      </div>
      {showLogin && (
        <AdminLogin
          onLogin={onAdminLogin}
          onClose={() => setShowLogin(false)}
        />
      )}
    </footer>
  )
}

