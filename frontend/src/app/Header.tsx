"use client";
import React, { useState, useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { get, del } from "./lib/api";

export default function Header() {
    const router = useRouter();
    const pathname = usePathname();
    const [userName, setUserName] = useState<string>("");
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        // ログインページでない場合のみユーザー情報を取得
        if (pathname !== "/login") {
            fetchUserInfo();
        }
    }, [pathname]);

    const fetchUserInfo = async () => {
        setIsLoading(true);
        try {
            const response = await get(`/session`);

            // ログインページへのリダイレクトは get() で自動的に処理される
            if (response && response.ok) {
                const data = await response.json();
                setUserName(data.user.name);
            } else if (response) {
                console.error("Failed to fetch user info");
            }
        } catch (error) {
            console.error("Error fetching user info:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleLogout = async () => {
        const api_base_url = process.env.NEXT_PUBLIC_API_BASE_URL;
        await del(`${api_base_url}/session`);
        router.push("/login");
    };

    return (
        <header className="w-full flex items-center justify-between px-6 py-4 bg-blue-600 text-white shadow">
            <div className="text-lg font-bold">Shift WebApp</div>
            {pathname !== "/login" && (
                <div className="flex items-center gap-4">
                    {!isLoading && userName && (
                        <div className="text-sm font-medium">{userName}さん</div>
                    )}
                    <button
                        onClick={handleLogout}
                        className="bg-white text-blue-600 px-4 py-2 rounded font-semibold hover:bg-blue-100 transition"
                    >
                        ログアウト
                    </button>
                </div>
            )}
        </header>
    );
}
