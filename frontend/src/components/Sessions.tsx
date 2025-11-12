import { SessionWithSpeakers } from '../api'

interface SessionsProps {
  sessions: SessionWithSpeakers[]
}

export default function Sessions({ sessions }: SessionsProps) {
  return (
    <section className="mb-16 animate-fade-in">
      <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Sessions & Speakers</h2>
      {sessions.length === 0 ? (
        <div className="text-center text-gray-600 py-12">
          <p>No sessions scheduled yet. Check back soon!</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {sessions.map((session) => (
            <div
              key={session.id}
              className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition-shadow duration-300 animate-slide-up"
            >
              <h3 className="text-xl font-semibold text-gray-900 mb-2">{session.title}</h3>
              <p className="text-sm text-indigo-600 mb-4 font-medium">{session.time}</p>
              <p className="text-gray-700 mb-4">{session.description}</p>
              {session.speakers.length > 0 && (
                <div className="mt-4 pt-4 border-t border-gray-200">
                  <p className="text-sm font-semibold text-gray-900 mb-2">Speakers:</p>
                  <div className="space-y-2">
                    {session.speakers.map((speaker) => (
                      <div key={speaker.id} className="flex items-start">
                        {speaker.photoUrl && (
                          <img
                            src={speaker.photoUrl}
                            alt={speaker.name}
                            className="w-10 h-10 rounded-full mr-3 object-cover"
                          />
                        )}
                        <div>
                          <p className="text-sm font-medium text-gray-900">{speaker.name}</p>
                          {speaker.bio && (
                            <p className="text-xs text-gray-600">{speaker.bio}</p>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </section>
  )
}

