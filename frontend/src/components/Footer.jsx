import { useState } from 'react';
import AdminPanel from './AdminPanel';

function Footer() {
  const [showAdminPanel, setShowAdminPanel] = useState(false);

  return (
    <>
      <footer className="footer">
        <div className="container">
          <p className="footer-text">
            © 2025 AppDirect India Tech Meetup
          </p>
          <button
            className="admin-login-link"
            onClick={() => setShowAdminPanel(true)}
          >
            Admin Login
          </button>
        </div>
      </footer>
      {showAdminPanel && (
        <AdminPanel onClose={() => setShowAdminPanel(false)} />
      )}
    </>
  );
}

export default Footer;

