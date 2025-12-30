"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

export default function PastePage() {
    const { id } = useParams();
    const [content, setContent] = useState<string | null>(null);
    const [error, setError] = useState("");

    useEffect(() => {
        fetch(`http://localhost:8080/api/pastes/${id}`)
            .then((res) => {
                if (!res.ok) throw new Error("Paste unavailable");
                return res.json();
            })
            .then((data) => setContent(data.content))
            .catch((e) => setError(e.message));
    }, [id]);

    if (error) return <p>{error}</p>;

    return (
        <pre style={{ maxWidth: 600, margin: "40px auto" }}>
            {content}
        </pre>
    );
}
