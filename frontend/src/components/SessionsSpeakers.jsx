import { useState, useEffect } from 'react';
import { getSessions } from '../services/api';

function SessionsSpeakers() {
  const [sessions, setSessions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchSessions();
  }, []);

  const fetchSessions = async () => {
    try {
      setLoading(true);
      const response = await getSessions();
      setSessions(response.data || []);
      setError(null);
    } catch (err) {
      console.error('Failed to load sessions:', err);
      setError('Failed to load sessions');
      setSessions([]); // Ensure sessions is an array
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <section className="sessions-section">
        <div className="container">
          <h2 className="section-title">Sessions & Speakers</h2>
          <div className="loading">Loading sessions...</div>
        </div>
      </section>
    );
  }

  if (error) {
    return (
      <section className="sessions-section">
        <div className="container">
          <h2 className="section-title">Sessions & Speakers</h2>
          <div className="error">{error}</div>
        </div>
      </section>
    );
  }

  return (
    <section className="sessions-section">
      <div className="container">
        <h2 className="section-title">Sessions & Speakers</h2>
        <div className="sessions-grid">
          {!sessions || sessions.length === 0 ? (
            <div className="no-sessions">No sessions available yet. Check back soon!</div>
          ) : (
            sessions.map((session) => (
              <div key={session.id} className="session-card">
                <div className="session-header">
                  <h3 className="session-title">{session.title}</h3>
                  <span className="session-time">{session.time}</span>
                </div>
                <p className="session-description">{session.description}</p>
                {session.speaker && (
                  <div className="speaker-info">
                    {session.speaker.photoUrl && (
                      <img
                        src={session.speaker.photoUrl}
                        alt={session.speaker.name}
                        className="speaker-photo"
                      />
                    )}
                    <div className="speaker-details">
                      <h4 className="speaker-name">{session.speaker.name}</h4>
                      {session.speaker.bio && (
                        <p className="speaker-bio">{session.speaker.bio}</p>
                      )}
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </section>
  );
}

export default SessionsSpeakers;
