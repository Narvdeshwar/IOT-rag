import React from "react";
import { Link, useLocation } from "react-router-dom";
import {
    LayoutDashboard,
    MessageSquare,
    Settings,
    Cpu,
    Database,
    ChevronRight
} from "lucide-react";
import { cn } from "../../lib/utils";

const sidebarItems = [
    { name: "Dashboard", icon: LayoutDashboard, path: "/" },
    { name: "IoT Devices", icon: Cpu, path: "/devices" },
    { name: "RAG Explorer", icon: MessageSquare, path: "/rag" },
    { name: "Knowledge Base", icon: Database, path: "/kb" },
    { name: "Settings", icon: Settings, path: "/settings" },
];

export function MainLayout({ children }: { children: React.ReactNode }) {
    const location = useLocation();

    return (
        <div className="flex h-screen bg-background text-foreground overflow-hidden">
            {/* Sidebar */}
            <aside className="w-64 border-r border-border bg-card/50 backdrop-blur-xl flex flex-col">
                <div className="p-6 flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
                        <Cpu className="text-primary-foreground w-5 h-5" />
                    </div>
                    <span className="font-bold text-xl tracking-tight">IoT-RAG</span>
                </div>

                <nav className="flex-1 px-4 space-y-2">
                    {sidebarItems.map((item) => {
                        const isActive = location.pathname === item.path;
                        const Icon = item.icon;
                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                className={cn(
                                    "flex items-center gap-3 px-3 py-2 rounded-lg transition-all group",
                                    isActive
                                        ? "bg-primary text-primary-foreground shadow-lg shadow-primary/20"
                                        : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                                )}
                            >
                                <Icon className={cn("w-5 h-5", isActive ? "text-primary-foreground" : "text-muted-foreground group-hover:text-accent-foreground")} />
                                <span className="flex-1 font-medium">{item.name}</span>
                                {isActive && <ChevronRight className="w-4 h-4 ml-auto" />}
                            </Link>
                        );
                    })}
                </nav>

                <div className="p-4 border-t border-border">
                    <div className="flex items-center gap-3 p-2 rounded-lg bg-accent/50">
                        <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-primary to-purple-500" />
                        <div className="flex-1 overflow-hidden">
                            <p className="text-sm font-medium truncate">Senior Architect</p>
                            <p className="text-xs text-muted-foreground truncate">Admin System</p>
                        </div>
                    </div>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 flex flex-col min-w-0 overflow-hidden">
                <header className="h-16 border-b border-border flex items-center justify-between px-8 bg-background/50 backdrop-blur-md sticky top-0 z-10">
                    <div className="flex items-center gap-4">
                        <h1 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
                            System Console / {sidebarItems.find(i => i.path === location.pathname)?.name || "Live"}
                        </h1>
                    </div>
                    <div className="flex items-center gap-4">
                        <div className="flex h-2 w-2 rounded-full bg-green-500 animate-pulse" />
                        <span className="text-xs font-mono text-muted-foreground">ENGINE_STATUS: OPTIMAL</span>
                    </div>
                </header>

                <div className="flex-1 overflow-y-auto p-8">
                    <div className="max-w-7xl mx-auto h-full">
                        {children}
                    </div>
                </div>
            </main>
        </div>
    );
}
