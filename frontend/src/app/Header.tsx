"use client";
import React from "react";
import { useRouter, usePathname } from "next/navigation";

// TODO: account情報へのリンク

export default function Header() {
    const router = useRouter();
    const pathname = usePathname();
    const handleLogout = async () => {
        const api_base_url = process.env.NEXT_PUBLIC_API_BASE_URL;
        await fetch(`${api_base_url}/session`, {
            method: "DELETE",
            credentials: "include",
        });
        router.push("/login");
    };

    return (
        <header className="w-full flex items-center justify-between px-6 py-4 bg-blue-600 text-white shadow">
            <div className="text-lg font-bold">Shift WebApp</div>
            {pathname !== "/login" && (
                <button
                    onClick={handleLogout}
                    className="bg-white text-blue-600 px-4 py-2 rounded font-semibold hover:bg-blue-100 transition"
                >
                    ログアウト
                </button>
            )}
        </header>
    );
}
