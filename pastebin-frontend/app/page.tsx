"use client";

import { useState } from "react";

export default function Home() {
  const [content, setContent] = useState("");
  const [maxViews, setMaxViews] = useState("");
  const [ttl, setTtl] = useState("");
  const [url, setUrl] = useState("");

  const createPaste = async () => {
    const body: any = { content };
    if (maxViews) body.max_views = Number(maxViews);
    if (ttl) body.ttl_seconds = Number(ttl);

    const res = await fetch("http://localhost:8080/api/pastes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      alert("Failed to create paste");
      return;
    }

    const data = await res.json();
    setUrl(data.url);
  };

  return (
    <div style={{ maxWidth: 600, margin: "40px auto" }}>
      <h2>Create Paste</h2>

      <textarea
        rows={6}
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="Paste content"
        style={{ width: "100%" }}
      />

      <input
        placeholder="Max views (optional)"
        value={maxViews}
        onChange={(e) => setMaxViews(e.target.value)}
        style={{ border: "1", borderColor: "blue" }}
      />

      <input
        placeholder="TTL seconds (optional)"
        value={ttl}
        onChange={(e) => setTtl(e.target.value)}
        style={{ border: "5px", borderColor: "blue" }}

      />

      <button onClick={createPaste} style={{ backgroundColor: "lightskyblue" }}>Create</button>

      {url && (
        <p>
          Share: <a href={url}>{url}</a>
        </p>
      )}
    </div>
  );
}
