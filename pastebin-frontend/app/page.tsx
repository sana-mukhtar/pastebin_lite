"use client";

import { log } from "console";
import { useState } from "react";

export default function Home() {
  const [content, setContent] = useState("");
  const [pasteURL, setPasteURL] = useState("");

  const createPaste = async () => {
    if (!content) return;
    const res = await fetch("http://localhost:8080/api/pastes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ content }),
    });
    const data = await res.json();
    console.log(data);
    setPasteURL(data.url);
  };

  return (
    <div style={{ maxWidth: 600, margin: "50px auto", textAlign: "center", border: "1px solid purple", padding: "20px" }}>
      <h1>Pastebin</h1>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={10}
        cols={50}
        placeholder="Enter your text here"
      />
      <br />
      <button onClick={createPaste} style={{ marginTop: 10, backgroundColor: "purple", color: "white", padding: "10px" }}>
        Create Paste
      </button>
      {pasteURL && (
        <div style={{ marginTop: 20 }}>
          <p>Paste created!</p>
          <a href={pasteURL} target="_blank" rel="noreferrer">
            {pasteURL}
          </a>
        </div>
      )}
    </div>
  );
}
