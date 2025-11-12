export default function Location() {
  return (
    <section className="mb-16 animate-fade-in">
      <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Event Location</h2>
      <div className="bg-white rounded-lg shadow-lg p-8 max-w-6xl mx-auto">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div>
            <h3 className="text-2xl font-semibold text-gray-900 mb-4">AppDirect India</h3>
            <div className="space-y-3 text-gray-700">
              <p className="flex items-start">
                <svg
                  className="w-5 h-5 text-indigo-600 mr-2 mt-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
                <span>Pune, Maharashtra, India</span>
              </p>
              <p className="flex items-start">
                <svg
                  className="w-5 h-5 text-indigo-600 mr-2 mt-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
                <span>Date: TBD</span>
              </p>
              <p className="flex items-start">
                <svg
                  className="w-5 h-5 text-indigo-600 mr-2 mt-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <span>Time: TBD</span>
              </p>
            </div>
            <div className="mt-6 p-4 bg-indigo-50 rounded-lg">
              <p className="text-sm text-indigo-900">
                Join us for an exciting day of AI innovation, networking, and learning. Connect
                with industry experts and fellow professionals.
              </p>
            </div>
          </div>
          <div className="w-full h-full min-h-[400px]">
            <iframe
              src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3930.310034205593!2d73.92600257972352!3d18.515585566381823!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x3bc2c18cf4eaad8d%3A0xc5835f1d9e3a91d3!2sAppDirect%20India!5e0!3m2!1sen!2sin!4v1762854087901!5m2!1sen!2sin"
              width="100%"
              height="100%"
              style={{ border: 0 }}
              allowFullScreen
              loading="lazy"
              referrerPolicy="no-referrer-when-downgrade"
              className="rounded-lg"
            />
          </div>
        </div>
      </div>
    </section>
  )
}

