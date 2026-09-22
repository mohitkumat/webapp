import { Routes, Route, Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { analytics } from "./lib/analytics";

function HomePage() {
  return (
    <main>
      <h1>Project Template</h1>
      <p>Frontend is up. See the health check page to verify the API connection.</p>
      <nav>
        <Link to="/health">Health check</Link>
      </nav>
    </main>
  );
}

function HealthPage() {
  return <HealthCheck />;
}

type Status = "loading" | "ok" | "error";

function HealthCheck() {
  const [status, setStatus] = useState<Status>("loading");

  useEffect(() => {
    analytics.track("health_check_viewed");

    const apiUrl = import.meta.env.VITE_API_URL ?? "http://127.0.0.1:8080";
    fetch(`${apiUrl}/health`)
      .then((res) => {
        const next = res.ok ? "ok" : "error";
        setStatus(next);
        analytics.track("health_check_result", { status: next });
      })
      .catch(() => {
        setStatus("error");
        analytics.track("health_check_result", { status: "error" });
      });
  }, []);

  if (status === "loading") return <p role="status">Checking API…</p>;
  if (status === "error")
    return <p role="alert">Could not reach the API. Is it running?</p>;
  return <p role="status">API is healthy.</p>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/health" element={<HealthPage />} />
    </Routes>
  );
}
