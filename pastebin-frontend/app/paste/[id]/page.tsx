"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

export default function Paste() {
    const { id } = useParams();
    const [content, setContent] = useState<string | null>(null);
    const [error, setError] = useState("");

    useEffect(() => {
        if (!id) return;
        fetch(`http://localhost:8080/api/pastes/${id}`)
            .then((res) => {
                if (!res.ok) throw new Error("Paste not found");
                return res.json();
            })
            .then((data) => setContent(data.content))
            .catch((err) => setError(err.message));
    }, [id]);

    return (
        <div style={{ maxWidth: 600, margin: "50px auto", textAlign: "center" }}>
            <h1>Paste Content</h1>
            {error && <p>{error}</p>}
            {content && (
                <pre style={{ textAlign: "left", padding: 10 }}>
                    {content}
                </pre>
            )}
        </div>
    );
}
