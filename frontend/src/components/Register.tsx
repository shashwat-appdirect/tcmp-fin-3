import { useState, useEffect } from 'react'
import { getAttendees, registerAttendee, AttendeesResponse } from '../api'

const DESIGNATIONS = [
  'Software Engineer',
  'Senior Software Engineer',
  'Tech Lead',
  'Engineering Manager',
  'Product Manager',
  'Product Designer',
  'UX Designer',
  'Data Scientist',
  'DevOps Engineer',
  'QA Engineer',
  'Full Stack Developer',
  'Frontend Developer',
  'Backend Developer',
  'Mobile Developer',
  'Architect',
  'CTO',
  'Other',
]

export default function Register() {
  const [attendeeCount, setAttendeeCount] = useState(0)
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [designation, setDesignation] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const loadCount = async () => {
      try {
        const data: AttendeesResponse = await getAttendees()
        setAttendeeCount(data.count)
      } catch (err) {
        console.error('Failed to load attendee count:', err)
      }
    }
    loadCount()
    const interval = setInterval(loadCount, 5000) // Refresh every 5 seconds
    return () => clearInterval(interval)
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess(false)

    try {
      await registerAttendee({ name, email, designation })
      setSuccess(true)
      setName('')
      setEmail('')
      setDesignation('')
      // Refresh count
      const data: AttendeesResponse = await getAttendees()
      setAttendeeCount(data.count)
      setTimeout(() => setSuccess(false), 3000)
    } catch (err: any) {
      setError(err.message || 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="mb-16 animate-fade-in">
      <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Register Now</h2>
      <div className="bg-white rounded-lg shadow-lg p-8 max-w-4xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 items-center">
          {/* Live Count */}
          <div className="text-center md:text-left">
            <div className="bg-indigo-100 rounded-lg p-6">
              <p className="text-sm text-indigo-700 font-medium mb-2">Live Attendee Count</p>
              <p className="text-4xl font-bold text-indigo-900">{attendeeCount}</p>
            </div>
          </div>

          {/* Registration Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                Name
              </label>
              <input
                type="text"
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                placeholder="Enter your name"
              />
            </div>
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                placeholder="Enter your email"
              />
            </div>
            <div>
              <label
                htmlFor="designation"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Designation
              </label>
              <select
                id="designation"
                value={designation}
                onChange={(e) => setDesignation(e.target.value)}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              >
                <option value="">Select designation</option>
                {DESIGNATIONS.map((d) => (
                  <option key={d} value={d}>
                    {d}
                  </option>
                ))}
              </select>
            </div>
            {error && <p className="text-red-600 text-sm">{error}</p>}
            <div className="flex justify-center md:justify-end pt-2">
              <button
                type="submit"
                disabled={loading || !name || !email || !designation}
                className="bg-indigo-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-indigo-700 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-md hover:shadow-lg"
              >
                {loading ? 'Registering...' : 'Register'}
              </button>
            </div>
          </form>
        </div>
      </div>

      {/* Success Popup */}
      {success && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 animate-fade-in">
          <div className="bg-white rounded-lg p-8 max-w-md mx-4 animate-slide-up">
            <div className="text-center">
              <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4">
                <svg
                  className="h-6 w-6 text-green-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Registration Successful!</h3>
              <p className="text-sm text-gray-600 mb-4">
                Thank you for registering. We look forward to seeing you at the event!
              </p>
              <button
                onClick={() => setSuccess(false)}
                className="bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 transition-colors"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </section>
  )
}

