"use client";
import React from "react";
import { useRouter, usePathname } from "next/navigation";
import { useSession } from "./hooks";
import { Button } from "./components/ui";
import { del } from "./lib/api";

// TODO: headerの名前残る

export default function Header() {
    const router = useRouter();
    const pathname = usePathname();
    const { user, isLoading } = useSession();

    // ログインページでない場合のみ表示
    const showUserInfo = pathname !== "/login";

    const handleLogout = async () => {
        const api_base_url = process.env.NEXT_PUBLIC_API_BASE_URL;
        await del(`${api_base_url}/session`);
        router.push("/login");
    };

    return (
        <header className="w-full flex items-center justify-between px-6 py-4 bg-blue-600 text-white shadow">
            <div className="text-lg font-bold">Shift WebApp</div>
            {showUserInfo && (
                <div className="flex items-center gap-4">
                    {!isLoading && user && (
                        <div className="text-sm font-medium">{user.name}さん</div>
                    )}
                    <Button
                        onClick={handleLogout}
                        variant="outline"
                        className="bg-white text-blue-600 hover:bg-blue-100 transition"
                    >
                        ログアウト
                    </Button>
                </div>
            )}
        </header>
    );
}
