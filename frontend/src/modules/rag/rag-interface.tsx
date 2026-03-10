import { useState, useRef, useEffect } from "react";
import { Send, Bot, User, Sparkles, Loader2 } from "lucide-react";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { Card } from "../../components/ui/Card";
import { cn } from "../../lib/utils";
import { motion, AnimatePresence } from "framer-motion";
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

interface Message {
    id: string;
    role: "user" | "assistant";
    content: string;
    timestamp: Date;
}

export function RagInterface() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            role: "assistant",
            content: "Hello! I am the IoT-RAG Engine. How can I help you analyze your telemetry data today?",
            timestamp: new Date(),
        },
    ]);
    const [input, setInput] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [messages]);

    const handleSend = async () => {
        if (!input.trim() || isLoading) return;

        const userMessage: Message = {
            id: Date.now().toString(),
            role: "user",
            content: input,
            timestamp: new Date(),
        };

        setMessages((prev) => [...prev, userMessage]);
        setInput("");
        setIsLoading(true);

        try {
            const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"}/query`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ query: input }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const contentType = response.headers.get("Content-Type");

            // Placeholder for assistant message that we'll update
            const assistantMessageId = (Date.now() + 1).toString();
            const assistantMessage: Message = {
                id: assistantMessageId,
                role: "assistant",
                content: "",
                timestamp: new Date(),
            };
            setMessages((prev) => [...prev, assistantMessage]);

            if (contentType?.includes("application/json")) {
                // Handle cached response
                const result = await response.json();
                setMessages((prev) =>
                    prev.map(msg => msg.id === assistantMessageId ? { ...msg, content: result.answer } : msg)
                );
            } else if (contentType?.includes("text/event-stream")) {
                // Handle streaming response
                const reader = response.body?.getReader();
                if (!reader) throw new Error("No reader available");

                const decoder = new TextDecoder();
                let accumulatedContent = "";

                while (true) {
                    const { done, value } = await reader.read();
                    if (done) break;

                    const chunk = decoder.decode(value, { stream: true });
                    const lines = chunk.split("\n");

                    for (const line of lines) {
                        if (line.startsWith("data: ")) {
                            try {
                                const data = JSON.parse(line.slice(6));
                                accumulatedContent += data.token;
                                
                                setMessages((prev) =>
                                    prev.map(msg => msg.id === assistantMessageId ? { ...msg, content: accumulatedContent } : msg)
                                );
                            } catch (e) {
                                console.error("Error parsing JSON token:", e);
                            }
                        }
                    }
                }
            }
        } catch (error) {
            console.error("Failed to fetch query:", error);
            const errorMessage: Message = {
                id: (Date.now() + 1).toString(),
                role: "assistant",
                content: "Sorry, I encountered an error while processing your request. Please try again later.",
                timestamp: new Date(),
            };
            setMessages((prev) => [...prev, errorMessage]);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="flex flex-col h-[calc(100vh-8rem)] max-w-7xl mx-auto px-4">
            <div className="flex items-center justify-between mb-8">
                <div className="space-y-1">
                    <h2 className="text-3xl font-extrabold tracking-tight bg-gradient-to-r from-white via-primary/80 to-gray-500 bg-clip-text text-transparent">
                        Intelligence Fusion Console
                    </h2>
                    <p className="text-muted-foreground/80 text-sm font-medium">
                        Real-time neural RAG explorer for distributed IoT telemetry streams.
                    </p>
                </div>
                <div className="flex gap-3">
                    <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-primary/10 border border-primary/20 text-[10px] font-mono text-primary animate-pulse">
                        <div className="w-1.5 h-1.5 rounded-full bg-primary" />
                        OLLAMA: ACTIVE
                    </div>
                </div>
            </div>

            <Card className="flex-1 flex flex-col overflow-hidden bg-black/40 backdrop-blur-2xl border-white/5 shadow-2xl relative group">
                {/* Decorative corner accents */}
                <div className="absolute top-0 left-0 w-20 h-20 border-t-2 border-l-2 border-primary/20 rounded-tl-3xl pointer-events-none" />
                <div className="absolute bottom-0 right-0 w-20 h-20 border-b-2 border-r-2 border-primary/20 rounded-br-3xl pointer-events-none" />

                <div
                    ref={scrollRef}
                    className="flex-1 overflow-y-auto p-8 space-y-8 scroll-smooth"
                >
                    <AnimatePresence initial={false}>
                        {messages.map((message) => (
                            <motion.div
                                key={message.id}
                                initial={{ opacity: 0, y: 20, filter: "blur(10px)" }}
                                animate={{ opacity: 1, y: 0, filter: "blur(0px)" }}
                                className={cn(
                                    "flex items-start gap-6 max-w-[90%]",
                                    message.role === "user" ? "ml-auto flex-row-reverse" : ""
                                )}
                            >
                                <div className={cn(
                                    "w-10 h-10 rounded-xl flex items-center justify-center shrink-0 shadow-lg",
                                    message.role === "user" 
                                        ? "bg-gradient-to-br from-primary to-blue-600 shadow-primary/20" 
                                        : "bg-white/5 border border-white/10 backdrop-blur-md"
                                )}>
                                    {message.role === "user" ? (
                                        <User className="w-5 h-5 text-primary-foreground" />
                                    ) : (
                                        <Bot className="w-5 h-5 text-primary" />
                                    )}
                                </div>
                                <div className={cn(
                                    "p-6 rounded-3xl text-sm leading-relaxed overflow-hidden relative group/msg",
                                    message.role === "user"
                                        ? "bg-primary/20 border border-primary/30 text-foreground rounded-tr-none shadow-[0_0_20px_rgba(59,130,246,0.1)]"
                                        : "bg-white/5 border border-white/10 backdrop-blur-xl rounded-tl-none shadow-[0_0_30px_rgba(255,255,255,0.02)] prose prose-invert prose-sm max-w-none"
                                )}>
                                    {message.role === "assistant" ? (
                                        <ReactMarkdown 
                                            remarkPlugins={[remarkGfm]}
                                            components={{
                                                table: ({node, ...props}) => (
                                                    <div className="my-6 rounded-2xl border border-white/10 overflow-hidden bg-white/5 backdrop-blur-sm">
                                                        <div className="overflow-x-auto">
                                                            <table className="min-w-full divide-y divide-white/10" {...props} />
                                                        </div>
                                                    </div>
                                                ),
                                                thead: ({node, ...props}) => <thead className="bg-white/5" {...props} />,
                                                th: ({node, ...props}) => (
                                                    <th className="px-5 py-3 text-left font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground" {...props} />
                                                ),
                                                td: ({node, ...props}) => (
                                                    <td className="px-5 py-3 border-t border-white/5 font-medium text-white/90" {...props} />
                                                ),
                                                p: ({node, ...props}) => <p className="mb-4 last:mb-0 leading-relaxed text-white/80" {...props} />,
                                                strong: ({node, ...props}) => <strong className="text-white font-semibold" {...props} />,
                                            }}
                                        >
                                            {message.content}
                                        </ReactMarkdown>
                                    ) : (
                                        <div className="text-white/90">{message.content}</div>
                                    )}
                                    <span className="absolute bottom-2 right-4 text-[8px] font-mono text-muted-foreground opacity-0 group-hover/msg:opacity-100 transition-opacity uppercase tracking-widest">
                                        {message.timestamp.toLocaleTimeString()}
                                    </span>
                                </div>
                            </motion.div>
                        ))}
                    </AnimatePresence>
                    {isLoading && (
                        <motion.div
                            initial={{ opacity: 0, x: -10 }}
                            animate={{ opacity: 1, x: 0 }}
                            className="flex items-center gap-4 text-primary text-[10px] font-mono tracking-widest uppercase ml-16 bg-primary/5 px-4 py-2 rounded-full border border-primary/20 w-fit"
                        >
                            <Loader2 className="w-3 h-3 animate-spin" />
                            Synchronizing Neural Buffers...
                        </motion.div>
                    )}
                </div>

                <div className="p-8 border-t border-white/5 bg-black/40 backdrop-blur-xl relative">
                    {/* Input Glow */}
                    <div className="absolute inset-0 bg-primary/2 h-[1px] -top-[1px] blur-sm" />
                    
                    <form
                        onSubmit={(e) => { e.preventDefault(); handleSend(); }}
                        className="flex gap-4 max-w-5xl mx-auto relative group"
                    >
                        <div className="relative flex-1">
                            <Input
                                value={input}
                                onChange={(e) => setInput(e.target.value)}
                                placeholder="Query neural database for IoT telemetry analysis..."
                                className="flex-1 bg-white/[0.03] border-white/10 h-14 pl-6 pr-12 rounded-2xl transition-all focus:bg-white/[0.05] focus:border-primary/50 focus:ring-0 placeholder:text-muted-foreground/40 font-medium"
                            />
                            <div className="absolute right-4 top-1/2 -translate-y-1/2 w-6 h-6 rounded-md bg-white/5 border border-white/10 flex items-center justify-center pointer-events-none group-focus-within:border-primary/30 group-focus-within:bg-primary/5 transition-colors">
                                <Sparkles className="w-3 h-3 text-muted-foreground group-focus-within:text-primary transition-colors" />
                            </div>
                        </div>
                        <Button 
                            type="submit" 
                            disabled={isLoading}
                            className="h-14 px-8 rounded-2xl bg-primary hover:bg-primary/90 shadow-lg shadow-primary/20 transition-all active:scale-95 disabled:opacity-50"
                        >
                            {isLoading ? (
                                <Loader2 className="w-5 h-5 animate-spin" />
                            ) : (
                                <div className="flex items-center gap-2">
                                    <span className="font-bold tracking-wide">EXECUTE</span>
                                    <Send className="w-4 h-4" />
                                </div>
                            )}
                        </Button>
                    </form>
                    <div className="flex items-center justify-center gap-6 mt-6">
                        <p className="text-[10px] font-mono text-muted-foreground/60 uppercase tracking-[0.3em]">
                            System: IOT-RAG v1.0.0-Stable
                        </p>
                        <div className="h-1 w-1 rounded-full bg-white/20" />
                        <p className="text-[10px] font-mono text-muted-foreground/60 uppercase tracking-[0.3em]">
                            Latent Dimension: 768
                        </p>
                        <div className="h-1 w-1 rounded-full bg-white/20" />
                        <p className="text-[10px] font-mono text-muted-foreground/60 uppercase tracking-[0.3em]">
                            Quantum Encrypted
                        </p>
                    </div>
                </div>
            </Card>
        </div>
    );
}
