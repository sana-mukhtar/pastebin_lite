"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";

export default function PastePage() {
    const router = useRouter();
    const params = useParams();
    const id = params.id as string;

    const [data, setData] = useState<{
        content: string;
        remaining_views: number | null;
        expires_at: string | null;
    } | null>(null);

    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchPaste = async () => {
            try {
                const res = await fetch(`http://localhost:8080/api/pastes/${id}`);
                const json = await res.json();

                if (!res.ok) {
                    setError(json.error || "Paste unavailable");
                    return;
                }

                setData(json);
            } catch {
                setError("Failed to reach backend");
            }
        };

        fetchPaste();
    }, [id]);

    if (error) {
        return (
            <div style={{ margin: 40 }}>
                <p style={{ color: "red" }}>{error}</p>
                <button onClick={() => router.push("/")}>Go Home</button>
            </div>
        );
    }

    if (!data) {
        return <p style={{ margin: 40 }}>Loading...</p>;
    }

    return (
        <div style={{ maxWidth: 600, margin: "40px auto", fontFamily: "sans-serif" }}>
            <h2>Paste</h2>

            <pre
                style={{
                    border: "3px solid #999",
                    padding: "12px",
                    whiteSpace: "pre-wrap",
                    wordBreak: "break-word",
                }}
            >
                {data.content}
            </pre>

            <p>
                <strong>Remaining views:</strong>{" "}
                {data.remaining_views === null ? "Unlimited" : data.remaining_views}
            </p>

            <p>
                <strong>Expires at:</strong>{" "}
                {data.expires_at
                    ? new Date(data.expires_at).toLocaleString()
                    : "Never"}
            </p>
        </div>
    );
}
