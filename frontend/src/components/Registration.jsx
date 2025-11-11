import { useState, useEffect } from 'react';
import { getAttendeeCount, registerAttendee } from '../services/api';
import ConfirmationModal from './ConfirmationModal';

const DESIGNATIONS = [
  'Software Engineer',
  'Product Manager',
  'Designer',
  'QA Engineer',
  'DevOps Engineer',
  'Data Scientist',
  'Business Analyst',
  'Other',
];

function Registration() {
  const [attendeeCount, setAttendeeCount] = useState(0);
  const [formData, setFormData] = useState({
    fullName: '',
    email: '',
    designation: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showConfirmation, setShowConfirmation] = useState(false);

  useEffect(() => {
    fetchCount();
    const interval = setInterval(fetchCount, 3000); // Poll every 3 seconds
    return () => clearInterval(interval);
  }, []);

  const fetchCount = async () => {
    try {
      const response = await getAttendeeCount();
      setAttendeeCount(response.data?.count || 0);
    } catch (err) {
      console.error('Failed to fetch attendee count:', err);
      // Don't set error state, just log it
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    // Validation
    if (!formData.fullName || !formData.email || !formData.designation) {
      setError('Please fill in all fields');
      setLoading(false);
      return;
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setError('Please enter a valid email address');
      setLoading(false);
      return;
    }

    try {
      await registerAttendee(formData);
      setShowConfirmation(true);
      setFormData({ fullName: '', email: '', designation: '' });
      fetchCount(); // Refresh count immediately
    } catch (err) {
      if (err.response?.status === 409) {
        setError('This email is already registered');
      } else {
        setError('Registration failed. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="registration-section">
      <div className="container">
        <h2 className="section-title">Register for the Event</h2>
        <div className="registration-content">
          <div className="attendee-count-box">
            <div className="count-label">Live Attendee Count</div>
            <div className="count-value">{attendeeCount}</div>
          </div>
          <div className="registration-form-container">
            <form onSubmit={handleSubmit} className="registration-form">
              <div className="form-group">
                <label htmlFor="fullName">Full Name</label>
                <input
                  type="text"
                  id="fullName"
                  value={formData.fullName}
                  onChange={(e) =>
                    setFormData({ ...formData, fullName: e.target.value })
                  }
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="email">Email</label>
                <input
                  type="email"
                  id="email"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="designation">Designation</label>
                <select
                  id="designation"
                  value={formData.designation}
                  onChange={(e) =>
                    setFormData({ ...formData, designation: e.target.value })
                  }
                  required
                >
                  <option value="">Select Designation</option>
                  {DESIGNATIONS.map((designation) => (
                    <option key={designation} value={designation}>
                      {designation}
                    </option>
                  ))}
                </select>
              </div>
              {error && <div className="error-message">{error}</div>}
              <button type="submit" className="register-button" disabled={loading}>
                {loading ? 'Registering...' : 'Register'}
              </button>
            </form>
          </div>
        </div>
      </div>
      {showConfirmation && (
        <ConfirmationModal onClose={() => setShowConfirmation(false)} />
      )}
    </section>
  );
}

export default Registration;

