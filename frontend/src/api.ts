const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface Attendee {
  id: string;
  name: string;
  email: string;
  designation: string;
  registeredAt: string;
}

export interface AttendeesResponse {
  count: number;
  attendees: Attendee[];
}

export interface Speaker {
  id: string;
  name: string;
  bio: string;
  photoUrl?: string;
}

export interface Session {
  id: string;
  title: string;
  description: string;
  time: string;
  speakerIds: string[];
}

export interface SessionWithSpeakers extends Session {
  speakers: Speaker[];
}

export interface StatsData {
  designation: string;
  count: number;
}

// Public API
export const getAttendees = async (): Promise<AttendeesResponse> => {
  const response = await fetch(`${API_URL}/api/attendees`);
  if (!response.ok) throw new Error('Failed to fetch attendees');
  return response.json();
};

export const registerAttendee = async (data: {
  name: string;
  email: string;
  designation: string;
}): Promise<Attendee> => {
  const response = await fetch(`${API_URL}/api/attendees`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to register');
  }
  return response.json();
};

export const getSessions = async (): Promise<SessionWithSpeakers[]> => {
  const response = await fetch(`${API_URL}/api/sessions`);
  if (!response.ok) throw new Error('Failed to fetch sessions');
  return response.json();
};

export const getSpeakers = async (): Promise<Speaker[]> => {
  const response = await fetch(`${API_URL}/api/speakers`);
  if (!response.ok) throw new Error('Failed to fetch speakers');
  return response.json();
};

// Admin API
const getAdminHeaders = (password: string) => ({
  'Content-Type': 'application/json',
  'X-Admin-Password': password,
});

export const adminLogin = async (password: string): Promise<void> => {
  const response = await fetch(`${API_URL}/api/admin/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });
  if (!response.ok) throw new Error('Invalid password');
};

export const getAdminAttendees = async (password: string): Promise<Attendee[]> => {
  const response = await fetch(`${API_URL}/api/admin/attendees`, {
    headers: getAdminHeaders(password),
  });
  if (!response.ok) throw new Error('Failed to fetch attendees');
  const data: AttendeesResponse = await response.json();
  return data.attendees;
};

export const deleteAttendee = async (id: string, password: string): Promise<void> => {
  const response = await fetch(`${API_URL}/api/admin/attendees/${id}`, {
    method: 'DELETE',
    headers: getAdminHeaders(password),
  });
  if (!response.ok) throw new Error('Failed to delete attendee');
};

export const addOrUpdateSession = async (
  session: Partial<Session> & { title: string; speakerIds: string[] },
  password: string
): Promise<Session> => {
  const response = await fetch(`${API_URL}/api/admin/sessions`, {
    method: 'POST',
    headers: getAdminHeaders(password),
    body: JSON.stringify(session),
  });
  if (!response.ok) throw new Error('Failed to save session');
  return response.json();
};

export const deleteSession = async (id: string, password: string): Promise<void> => {
  const response = await fetch(`${API_URL}/api/admin/sessions/${id}`, {
    method: 'DELETE',
    headers: getAdminHeaders(password),
  });
  if (!response.ok) throw new Error('Failed to delete session');
};

export const addOrUpdateSpeaker = async (
  speaker: Partial<Speaker> & { name: string },
  password: string
): Promise<Speaker> => {
  const response = await fetch(`${API_URL}/api/admin/speakers`, {
    method: 'POST',
    headers: getAdminHeaders(password),
    body: JSON.stringify(speaker),
  });
  if (!response.ok) throw new Error('Failed to save speaker');
  return response.json();
};

export const deleteSpeaker = async (id: string, password: string): Promise<void> => {
  const response = await fetch(`${API_URL}/api/admin/speakers/${id}`, {
    method: 'DELETE',
    headers: getAdminHeaders(password),
  });
  if (!response.ok) throw new Error('Failed to delete speaker');
};

export const getStats = async (password: string): Promise<StatsData[]> => {
  const response = await fetch(`${API_URL}/api/admin/stats`, {
    headers: getAdminHeaders(password),
  });
  if (!response.ok) throw new Error('Failed to fetch stats');
  return response.json();
};

