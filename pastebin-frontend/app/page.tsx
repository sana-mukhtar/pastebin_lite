"use client";

import { useState } from "react";

export default function Home() {
  const [content, setContent] = useState("");
  const [maxViews, setMaxViews] = useState(3);
  const [ttlSeconds, setTtlSeconds] = useState(300);
  const [pasteURL, setPasteURL] = useState("");
  const [error, setError] = useState("");

  const createPaste = async () => {
    setError("");
    setPasteURL("");

    try {
      const res = await fetch(process.env.NEXT_PUBLIC_API_URL), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content, max_views: maxViews, ttl_seconds: ttlSeconds }),
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Failed to create paste");
      }

      const data = await res.json();
      setPasteURL(data.url);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const inputStyle = { width: "100%", padding: 8, marginBottom: 10, border: "2px solid #0070f3", borderRadius: 4 };
  const smallInputStyle = { width: 80, padding: 6, border: "2px solid #0070f3", borderRadius: 4 };

  return (
    <div style={{ padding: 20, fontFamily: "sans-serif", maxWidth: 600, margin: "40px auto" }}>
      <h1 style={{ textAlign: "center", marginBottom: 20 }}>Pastebin</h1>
      {error && <p style={{ color: "red", textAlign: "center" }}>{error}</p>}
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={10}
        style={inputStyle}
        placeholder="Enter your text here"
      />
      <div style={{ marginBottom: 10 }}>
        Max Views: <input type="number" value={maxViews} onChange={(e) => setMaxViews(Number(e.target.value))} style={{ ...smallInputStyle, width: 50 }} />
      </div>
      <div style={{ marginBottom: 10 }}>
        TTL (seconds): <input type="number" value={ttlSeconds} onChange={(e) => setTtlSeconds(Number(e.target.value))} style={smallInputStyle} />
      </div>
      <button onClick={createPaste} style={{ width: "100%", padding: 10, background: "#0070f3", color: "#fff", border: "none", borderRadius: 4 }}>
        Create Paste
      </button>

      {pasteURL && (
        <div style={{ marginTop: 20, textAlign: "center" }}>
          <p style={{ color: "green" }}>Paste created!</p>
          <a href={pasteURL} target="_blank" rel="noreferrer" style={{ color: "#0070f3" }}>
            {pasteURL}
          </a>
        </div>
      )}
    </div>
  );
}
