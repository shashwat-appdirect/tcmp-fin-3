import { useState, useEffect } from 'react'
import {
  Attendee,
  Speaker,
  Session,
  getAdminAttendees,
  deleteAttendee,
  getSpeakers,
  addOrUpdateSpeaker,
  deleteSpeaker,
  getSessions,
  addOrUpdateSession,
  deleteSession,
  getStats,
  StatsData,
} from '../api'
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts'

interface AdminPanelProps {
  password: string
  onLogout: () => void
}

const COLORS = ['#4F46E5', '#7C3AED', '#EC4899', '#F59E0B', '#10B981', '#3B82F6', '#EF4444']

export default function AdminPanel({ password, onLogout }: AdminPanelProps) {
  const [attendees, setAttendees] = useState<Attendee[]>([])
  const [speakers, setSpeakers] = useState<Speaker[]>([])
  const [sessions, setSessions] = useState<Session[]>([])
  const [stats, setStats] = useState<StatsData[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'attendees' | 'speakers' | 'sessions' | 'stats'>(
    'attendees'
  )

  // Form states
  const [speakerForm, setSpeakerForm] = useState<Partial<Speaker>>({})
  const [sessionForm, setSessionForm] = useState<Partial<Session>>({ speakerIds: [] })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [attendeesData, speakersData, sessionsData, statsData] = await Promise.all([
        getAdminAttendees(password),
        getSpeakers(),
        getSessions(),
        getStats(password),
      ])
      setAttendees(attendeesData)
      setSpeakers(speakersData)
      setSessions(sessionsData)
      setStats(statsData)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteAttendee = async (id: string) => {
    if (!confirm('Are you sure you want to delete this attendee?')) return
    try {
      await deleteAttendee(id, password)
      await loadData()
    } catch (error) {
      alert('Failed to delete attendee')
    }
  }

  const handleSaveSpeaker = async () => {
    if (!speakerForm.name) {
      alert('Name is required')
      return
    }
    try {
      await addOrUpdateSpeaker(speakerForm as Speaker, password)
      setSpeakerForm({})
      await loadData()
    } catch (error) {
      alert('Failed to save speaker')
    }
  }

  const handleDeleteSpeaker = async (id: string) => {
    if (!confirm('Are you sure you want to delete this speaker?')) return
    try {
      await deleteSpeaker(id, password)
      await loadData()
    } catch (error) {
      alert('Failed to delete speaker')
    }
  }

  const handleSaveSession = async () => {
    if (!sessionForm.title) {
      alert('Title is required')
      return
    }
    try {
      await addOrUpdateSession(sessionForm as Session, password)
      setSessionForm({ speakerIds: [] })
      await loadData()
    } catch (error) {
      alert('Failed to save session')
    }
  }

  const handleDeleteSession = async (id: string) => {
    if (!confirm('Are you sure you want to delete this session?')) return
    try {
      await deleteSession(id, password)
      await loadData()
    } catch (error) {
      alert('Failed to delete session')
    }
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600">Loading admin panel...</p>
      </div>
    )
  }

  return (
    <div className="mb-16 bg-white rounded-lg shadow-lg p-8 animate-fade-in">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-3xl font-bold text-gray-900">Admin Panel</h2>
        <button
          onClick={onLogout}
          className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors"
        >
          Logout
        </button>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-gray-200 mb-6">
        {(['attendees', 'speakers', 'sessions', 'stats'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-6 py-3 font-medium capitalize ${
              activeTab === tab
                ? 'border-b-2 border-indigo-600 text-indigo-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab}
          </button>
        ))}
      </div>

      {/* Attendees Tab */}
      {activeTab === 'attendees' && (
        <div>
          <h3 className="text-xl font-semibold mb-4">Attendees ({attendees.length})</h3>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Email
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Designation
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Registered At
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {attendees.map((attendee) => (
                  <tr key={attendee.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {attendee.name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {attendee.email}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {attendee.designation}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {new Date(attendee.registeredAt).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <button
                        onClick={() => handleDeleteAttendee(attendee.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Speakers Tab */}
      {activeTab === 'speakers' && (
        <div>
          <h3 className="text-xl font-semibold mb-4">Speakers</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div className="bg-gray-50 p-4 rounded-lg">
              <h4 className="font-semibold mb-3">Add/Update Speaker</h4>
              <div className="space-y-3">
                <input
                  type="text"
                  placeholder="Name *"
                  value={speakerForm.name || ''}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, name: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                />
                <textarea
                  placeholder="Bio"
                  value={speakerForm.bio || ''}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, bio: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                  rows={3}
                />
                <input
                  type="text"
                  placeholder="Photo URL (optional)"
                  value={speakerForm.photoUrl || ''}
                  onChange={(e) => setSpeakerForm({ ...speakerForm, photoUrl: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                />
                {speakerForm.id && (
                  <p className="text-sm text-indigo-600">
                    Editing: {speakerForm.name}
                  </p>
                )}
                <div className="flex gap-2">
                  <button
                    onClick={handleSaveSpeaker}
                    className="flex-1 bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700"
                  >
                    {speakerForm.id ? 'Update' : 'Add'} Speaker
                  </button>
                  {speakerForm.id && (
                    <button
                      onClick={() => setSpeakerForm({})}
                      className="bg-gray-300 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-400"
                    >
                      Clear
                    </button>
                  )}
                </div>
              </div>
            </div>
            <div>
              <h4 className="font-semibold mb-3">Existing Speakers</h4>
              <div className="space-y-2 max-h-96 overflow-y-auto">
                {speakers.map((speaker) => (
                  <div
                    key={speaker.id}
                    className="flex justify-between items-center p-3 bg-white border border-gray-200 rounded-lg"
                  >
                    <div>
                      <p className="font-medium">{speaker.name}</p>
                      {speaker.bio && <p className="text-sm text-gray-600">{speaker.bio}</p>}
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => setSpeakerForm(speaker)}
                        className="text-indigo-600 hover:text-indigo-800 text-sm"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDeleteSpeaker(speaker.id)}
                        className="text-red-600 hover:text-red-800 text-sm"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Sessions Tab */}
      {activeTab === 'sessions' && (
        <div>
          <h3 className="text-xl font-semibold mb-4">Sessions</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div className="bg-gray-50 p-4 rounded-lg">
              <h4 className="font-semibold mb-3">Add/Update Session</h4>
              <div className="space-y-3">
                <input
                  type="text"
                  placeholder="Title *"
                  value={sessionForm.title || ''}
                  onChange={(e) => setSessionForm({ ...sessionForm, title: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                />
                <textarea
                  placeholder="Description"
                  value={sessionForm.description || ''}
                  onChange={(e) => setSessionForm({ ...sessionForm, description: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                  rows={3}
                />
                <input
                  type="text"
                  placeholder="Time"
                  value={sessionForm.time || ''}
                  onChange={(e) => setSessionForm({ ...sessionForm, time: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                />
                <div>
                  <label className="block text-sm font-medium mb-2">Speakers (select multiple)</label>
                  <div className="max-h-32 overflow-y-auto border border-gray-300 rounded-lg p-2">
                    {speakers.map((speaker) => (
                      <label key={speaker.id} className="flex items-center space-x-2 py-1">
                        <input
                          type="checkbox"
                          checked={sessionForm.speakerIds?.includes(speaker.id) || false}
                          onChange={(e) => {
                            const ids = sessionForm.speakerIds || []
                            if (e.target.checked) {
                              setSessionForm({ ...sessionForm, speakerIds: [...ids, speaker.id] })
                            } else {
                              setSessionForm({
                                ...sessionForm,
                                speakerIds: ids.filter((id) => id !== speaker.id),
                              })
                            }
                          }}
                          className="rounded"
                        />
                        <span className="text-sm">{speaker.name}</span>
                      </label>
                    ))}
                  </div>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={handleSaveSession}
                    className="flex-1 bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700"
                  >
                    {sessionForm.id ? 'Update' : 'Add'} Session
                  </button>
                  {sessionForm.id && (
                    <button
                      onClick={() => setSessionForm({ speakerIds: [] })}
                      className="bg-gray-300 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-400"
                    >
                      Clear
                    </button>
                  )}
                </div>
              </div>
            </div>
            <div>
              <h4 className="font-semibold mb-3">Existing Sessions</h4>
              <div className="space-y-2 max-h-96 overflow-y-auto">
                {sessions.map((session) => (
                  <div
                    key={session.id}
                    className="p-3 bg-white border border-gray-200 rounded-lg"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <p className="font-medium">{session.title}</p>
                        {session.time && <p className="text-sm text-gray-600">{session.time}</p>}
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => setSessionForm(session)}
                          className="text-indigo-600 hover:text-indigo-800 text-sm"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDeleteSession(session.id)}
                          className="text-red-600 hover:text-red-800 text-sm"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                    {session.description && (
                      <p className="text-sm text-gray-600">{session.description}</p>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Stats Tab */}
      {activeTab === 'stats' && (
        <div>
          <h3 className="text-xl font-semibold mb-4">Attendee Statistics by Designation</h3>
          <div className="bg-gray-50 p-6 rounded-lg">
            <ResponsiveContainer width="100%" height={400}>
              <PieChart>
                <Pie
                  data={stats}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ designation, count, percent }) =>
                    `${designation}: ${count} (${(percent * 100).toFixed(0)}%)`
                  }
                  outerRadius={120}
                  fill="#8884d8"
                  dataKey="count"
                >
                  {stats.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}
    </div>
  )
}

