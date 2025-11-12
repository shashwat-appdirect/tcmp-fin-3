import { useState, useEffect } from 'react'
import Sessions from './components/Sessions'
import Register from './components/Register'
import Location from './components/Location'
import Footer from './components/Footer'
import AdminPanel from './components/AdminPanel'
import { getSessions } from './api'

function App() {
  const [sessions, setSessions] = useState<any[]>([])
  const [isAdmin, setIsAdmin] = useState(false)
  const [adminPassword, setAdminPassword] = useState<string | null>(
    localStorage.getItem('adminPassword')
  )

  useEffect(() => {
    const loadSessions = async () => {
      try {
        const data = await getSessions()
        setSessions(data)
      } catch (error) {
        console.error('Failed to load sessions:', error)
      }
    }
    loadSessions()
  }, [])

  useEffect(() => {
    if (adminPassword) {
      setIsAdmin(true)
    }
  }, [adminPassword])

  const handleAdminLogin = (password: string) => {
    setAdminPassword(password)
    localStorage.setItem('adminPassword', password)
    setIsAdmin(true)
  }

  const handleAdminLogout = () => {
    setAdminPassword(null)
    localStorage.removeItem('adminPassword')
    setIsAdmin(false)
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <header className="bg-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <h1 className="text-4xl font-bold text-gray-900 text-center animate-fade-in">
            AppDirect India AI Workshop
          </h1>
          <p className="text-center text-gray-600 mt-2 animate-slide-up">
            Join us for an exciting day of AI innovation and learning
          </p>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Sessions sessions={sessions} />
        <Register />
        <Location />
        {isAdmin && adminPassword && (
          <AdminPanel password={adminPassword} onLogout={handleAdminLogout} />
        )}
      </main>

      <Footer onAdminLogin={handleAdminLogin} isAdmin={isAdmin} />
    </div>
  )
}

export default App

