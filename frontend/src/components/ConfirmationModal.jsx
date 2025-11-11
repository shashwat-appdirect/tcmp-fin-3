function ConfirmationModal({ onClose }) {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Registration Successful!</h3>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          <p>Thank you for registering. We look forward to seeing you at the event!</p>
        </div>
        <div className="modal-footer">
          <button className="modal-button" onClick={onClose}>Close</button>
        </div>
      </div>
    </div>
  );
}

export default ConfirmationModal;

