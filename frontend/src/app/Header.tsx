"use client";
import React from "react";
import { useRouter, usePathname } from "next/navigation";

export default function Header() {
    const router = useRouter();
    const pathname = usePathname();
    const handleLogout = async () => {
        await fetch("/api/logout", {
            method: "POST",
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
