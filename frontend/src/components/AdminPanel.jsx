import { useState, useEffect } from 'react';
import { adminLogin, getAttendees, getAttendeeStats, createOrUpdateSpeaker, createOrUpdateSession, getSpeakers, getSessions } from '../services/api';
import PieChart from './PieChart';

function AdminPanel({ onClose }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [password, setPassword] = useState('');
  const [loginError, setLoginError] = useState('');
  const [loading, setLoading] = useState(false);
  const [attendees, setAttendees] = useState([]);
  const [stats, setStats] = useState({});
  const [speakers, setSpeakers] = useState([]);
  const [sessions, setSessions] = useState([]);
  const [activeTab, setActiveTab] = useState('attendees');
  
  // Speaker form
  const [speakerForm, setSpeakerForm] = useState({ id: '', name: '', bio: '', photoUrl: '' });
  const [speakerFormError, setSpeakerFormError] = useState('');
  
  // Session form
  const [sessionForm, setSessionForm] = useState({ id: '', title: '', description: '', time: '', speakerId: '' });
  const [sessionFormError, setSessionFormError] = useState('');

  useEffect(() => {
    if (isAuthenticated) {
      fetchData();
    }
  }, [isAuthenticated]);

  useEffect(() => {
    if (isAuthenticated && activeTab === 'speakers') {
      fetchSpeakers();
    }
    if (isAuthenticated && activeTab === 'sessions') {
      fetchSessions();
      fetchSpeakers(); // Load speakers for dropdown
    }
  }, [isAuthenticated, activeTab]);

  const fetchSpeakers = async () => {
    try {
      const response = await getSpeakers();
      setSpeakers(response.data);
    } catch (err) {
      console.error('Failed to fetch speakers:', err);
    }
  };

  const fetchSessions = async () => {
    try {
      const response = await getSessions();
      setSessions(response.data);
    } catch (err) {
      console.error('Failed to fetch sessions:', err);
    }
  };

  const fetchData = async () => {
    try {
      const [attendeesRes, statsRes] = await Promise.all([
        getAttendees(),
        getAttendeeStats(),
      ]);
      setAttendees(attendeesRes.data);
      setStats(statsRes.data);
    } catch (err) {
      console.error('Failed to fetch admin data:', err);
    }
  };

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoginError('');
    setLoading(true);

    try {
      const response = await adminLogin(password);
      if (response.data.success) {
        setIsAuthenticated(true);
      } else {
        setLoginError(response.data.message || 'Invalid password');
      }
    } catch (err) {
      setLoginError('Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleSpeakerSubmit = async (e) => {
    e.preventDefault();
    setSpeakerFormError('');
    setLoading(true);

    try {
      await createOrUpdateSpeaker(speakerForm);
      setSpeakerForm({ id: '', name: '', bio: '', photoUrl: '' });
      fetchSpeakers();
      fetchData();
    } catch (err) {
      setSpeakerFormError('Failed to save speaker');
    } finally {
      setLoading(false);
    }
  };

  const handleSessionSubmit = async (e) => {
    e.preventDefault();
    setSessionFormError('');
    setLoading(true);

    try {
      await createOrUpdateSession(sessionForm);
      setSessionForm({ id: '', title: '', description: '', time: '', speakerId: '' });
      fetchSessions();
      fetchData();
    } catch (err) {
      setSessionFormError('Failed to save session');
    } finally {
      setLoading(false);
    }
  };

  const editSpeaker = (speaker) => {
    setSpeakerForm({
      id: speaker.id,
      name: speaker.name,
      bio: speaker.bio || '',
      photoUrl: speaker.photoUrl || '',
    });
    setActiveTab('speakers');
  };

  const editSession = (session) => {
    setSessionForm({
      id: session.id,
      title: session.title,
      description: session.description,
      time: session.time,
      speakerId: session.speakerID || '',
    });
    setActiveTab('sessions');
  };

  if (!isAuthenticated) {
    return (
      <div className="admin-modal-overlay" onClick={onClose}>
        <div className="admin-modal-content" onClick={(e) => e.stopPropagation()}>
          <div className="admin-modal-header">
            <h2>Admin Login</h2>
            <button className="modal-close" onClick={onClose}>×</button>
          </div>
          <form onSubmit={handleLogin} className="admin-login-form">
            <div className="form-group">
              <label htmlFor="admin-password">Password</label>
              <input
                type="password"
                id="admin-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            {loginError && <div className="error-message">{loginError}</div>}
            <button type="submit" className="admin-login-button" disabled={loading}>
              {loading ? 'Logging in...' : 'Login'}
            </button>
          </form>
        </div>
      </div>
    );
  }

  const chartData = Object.entries(stats).map(([name, value]) => ({
    name,
    value,
  }));

  return (
    <div className="admin-modal-overlay" onClick={onClose}>
      <div className="admin-modal-content admin-panel" onClick={(e) => e.stopPropagation()}>
        <div className="admin-modal-header">
          <h2>Admin Panel</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="admin-tabs">
          <button
            className={activeTab === 'attendees' ? 'active' : ''}
            onClick={() => setActiveTab('attendees')}
          >
            Attendees
          </button>
          <button
            className={activeTab === 'speakers' ? 'active' : ''}
            onClick={() => setActiveTab('speakers')}
          >
            Speakers
          </button>
          <button
            className={activeTab === 'sessions' ? 'active' : ''}
            onClick={() => setActiveTab('sessions')}
          >
            Sessions
          </button>
          <button
            className={activeTab === 'stats' ? 'active' : ''}
            onClick={() => setActiveTab('stats')}
          >
            Statistics
          </button>
        </div>
        <div className="admin-content">
          {activeTab === 'attendees' && (
            <div className="admin-section">
              <h3>Attendee Details</h3>
              <div className="table-container">
                <table className="admin-table">
                  <thead>
                    <tr>
                      <th>Full Name</th>
                      <th>Email</th>
                      <th>Designation</th>
                      <th>Registered At</th>
                    </tr>
                  </thead>
                  <tbody>
                    {attendees.length === 0 ? (
                      <tr>
                        <td colSpan="4" className="no-data">No attendees yet</td>
                      </tr>
                    ) : (
                      attendees.map((attendee) => (
                        <tr key={attendee.id}>
                          <td>{attendee.fullName}</td>
                          <td>{attendee.email}</td>
                          <td>{attendee.designation}</td>
                          <td>{new Date(attendee.registeredAt?.seconds * 1000 || attendee.registeredAt).toLocaleString()}</td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          )}
          {activeTab === 'speakers' && (
            <div className="admin-section">
              <h3>Add/Update Speaker</h3>
              <form onSubmit={handleSpeakerSubmit} className="admin-form">
                <div className="form-group">
                  <label>Name *</label>
                  <input
                    type="text"
                    value={speakerForm.name}
                    onChange={(e) => setSpeakerForm({ ...speakerForm, name: e.target.value })}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Bio</label>
                  <textarea
                    value={speakerForm.bio}
                    onChange={(e) => setSpeakerForm({ ...speakerForm, bio: e.target.value })}
                    rows="3"
                  />
                </div>
                <div className="form-group">
                  <label>Photo URL</label>
                  <input
                    type="url"
                    value={speakerForm.photoUrl}
                    onChange={(e) => setSpeakerForm({ ...speakerForm, photoUrl: e.target.value })}
                  />
                </div>
                {speakerFormError && <div className="error-message">{speakerFormError}</div>}
                <button type="submit" className="admin-submit-button" disabled={loading}>
                  {speakerForm.id ? 'Update Speaker' : 'Add Speaker'}
                </button>
              </form>
              {speakers.length > 0 && (
                <div className="existing-items">
                  <h4>Existing Speakers</h4>
                  {speakers.map((speaker) => (
                    <div key={speaker.id} className="existing-item">
                      <span>{speaker.name}</span>
                      <button onClick={() => editSpeaker(speaker)}>Edit</button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
          {activeTab === 'sessions' && (
            <div className="admin-section">
              <h3>Add/Update Session</h3>
              <form onSubmit={handleSessionSubmit} className="admin-form">
                <div className="form-group">
                  <label>Title *</label>
                  <input
                    type="text"
                    value={sessionForm.title}
                    onChange={(e) => setSessionForm({ ...sessionForm, title: e.target.value })}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Description *</label>
                  <textarea
                    value={sessionForm.description}
                    onChange={(e) => setSessionForm({ ...sessionForm, description: e.target.value })}
                    rows="3"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Time *</label>
                  <input
                    type="text"
                    value={sessionForm.time}
                    onChange={(e) => setSessionForm({ ...sessionForm, time: e.target.value })}
                    placeholder="e.g., 10:00 AM - 11:00 AM"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Speaker ID *</label>
                  <select
                    value={sessionForm.speakerId}
                    onChange={(e) => setSessionForm({ ...sessionForm, speakerId: e.target.value })}
                    required
                  >
                    <option value="">Select Speaker</option>
                    {speakers.map((speaker) => (
                      <option key={speaker.id} value={speaker.id}>
                        {speaker.name}
                      </option>
                    ))}
                  </select>
                </div>
                {sessionFormError && <div className="error-message">{sessionFormError}</div>}
                <button type="submit" className="admin-submit-button" disabled={loading}>
                  {sessionForm.id ? 'Update Session' : 'Add Session'}
                </button>
              </form>
              {sessions.length > 0 && (
                <div className="existing-items">
                  <h4>Existing Sessions</h4>
                  {sessions.map((session) => (
                    <div key={session.id} className="existing-item">
                      <span>{session.title} - {session.time}</span>
                      <button onClick={() => editSession(session)}>Edit</button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
          {activeTab === 'stats' && (
            <div className="admin-section">
              <h3>Attendee Breakdown by Designation</h3>
              {chartData.length > 0 ? (
                <PieChart data={chartData} />
              ) : (
                <div className="no-data">No data available</div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default AdminPanel;

