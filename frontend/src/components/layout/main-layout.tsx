import {
    Cpu,
} from "lucide-react";

export function MainLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="flex h-screen bg-background text-foreground overflow-hidden relative">
            {/* Ambient Background Glows */}
            <div className="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
                <div className="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-primary/20 blur-[120px] rounded-full" />
                <div className="absolute top-[20%] -right-[10%] w-[30%] h-[30%] bg-purple-500/10 blur-[100px] rounded-full" />
                <div className="absolute -bottom-[10%] left-[20%] w-[35%] h-[35%] bg-blue-500/10 blur-[110px] rounded-full" />
            </div>

            {/* Main Content */}
            <main className="flex-1 flex flex-col min-w-0 overflow-hidden relative z-10">
                <header className="h-16 border-b border-border/50 flex items-center justify-between px-8 bg-background/30 backdrop-blur-xl sticky top-0 z-10">
                    <div className="flex items-center gap-4">
                        <div className="flex items-center gap-3">
                            <div className="w-7 h-7 rounded-lg bg-primary/20 border border-primary/30 flex items-center justify-center">
                                <Cpu className="text-primary w-4 h-4" />
                            </div>
                            <span className="font-bold text-lg tracking-tight bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">IoT-RAG Engine</span>
                        </div>
                        <div className="h-4 w-[1px] bg-border/50 mx-2" />
                        <h1 className="text-xs font-medium text-muted-foreground uppercase tracking-[0.2em]">
                            Intelligence Fusion Console
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
