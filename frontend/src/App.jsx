import './App.css';
import Header from './components/Header';
import SessionsSpeakers from './components/SessionsSpeakers';
import Registration from './components/Registration';
import Location from './components/Location';
import Footer from './components/Footer';

function App() {
  return (
    <div className="app">
      <Header />
      <main className="main-content">
        <SessionsSpeakers />
        <Registration />
        <Location />
      </main>
      <Footer />
    </div>
  );
}

export default App;
