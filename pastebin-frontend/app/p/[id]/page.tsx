"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function PastePage() {
    const params = useParams(); // gets { id: string }
    const id = params?.id;

    const [paste, setPaste] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!id) return;

        async function fetchPaste() {
            try {
                const res = await fetch(`http://localhost:8080/api/pastes/${id}`, {
                    cache: "no-store",
                });
                if (!res.ok) {
                    setPaste(null);
                } else {
                    const data = await res.json();
                    setPaste(data);
                }
            } catch {
                setPaste(null);
            } finally {
                setLoading(false);
            }
        }

        fetchPaste();
    }, [id]);

    if (loading) {
        return <main style={{ padding: "20px" }}>Loading...</main>;
    }

    if (!paste) {
        return (
            <main style={{ padding: "20px" }}>
                <h2>Paste unavailable</h2>
                <p>This paste has expired or reached its view limit.</p>
            </main>
        );
    }

    return (
        <main style={{ padding: "20px" }}>
            <h1>Paste</h1>
            <pre
                style={{
                    whiteSpace: "pre-wrap",
                    border: "3px solid #444",
                    padding: "12px",
                    marginTop: "12px",
                }}
            >
                {paste.content}
            </pre>

            {paste.remaining_views !== null && (
                <p>Remaining views: {paste.remaining_views}</p>
            )}

            {paste.expires_at && (
                <p>Expires at: {new Date(paste.expires_at).toLocaleString()}</p>
            )}
        </main>
    );
}
