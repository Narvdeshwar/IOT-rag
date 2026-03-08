import { useState, useRef, useEffect } from "react";
import { Send, Bot, User, Sparkles, Loader2 } from "lucide-react";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { Card } from "../../components/ui/Card";
import { cn } from "../../lib/utils";
import { motion, AnimatePresence } from "framer-motion";

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

        // Simulated API call (User will implement backend)
        setTimeout(() => {
            const assistantMessage: Message = {
                id: (Date.now() + 1).toString(),
                role: "assistant",
                content: `I've analyzed the query: "${input}". Based on the current IoT sensor data, I recommend checking the power consumption patterns in Zone B. Current predictive threat index is stable at 0.12.`,
                timestamp: new Date(),
            };
            setMessages((prev) => [...prev, assistantMessage]);
            setIsLoading(false);
        }, 1500);
    };

    return (
        <div className="flex flex-col h-[calc(100vh-10rem)] max-w-4xl mx-auto">
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight">RAG Engine Explorer</h2>
                    <p className="text-muted-foreground">Context-aware intelligence fusion for IoT telemetry.</p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" size="sm">
                        <Sparkles className="w-4 h-4 mr-2 text-primary" />
                        Optimize Model
                    </Button>
                </div>
            </div>

            <Card className="flex-1 flex flex-col overflow-hidden bg-card/30 backdrop-blur-sm border-border/50">
                <div
                    ref={scrollRef}
                    className="flex-1 overflow-y-auto p-6 space-y-6"
                >
                    <AnimatePresence initial={false}>
                        {messages.map((message) => (
                            <motion.div
                                key={message.id}
                                initial={{ opacity: 0, y: 10, scale: 0.95 }}
                                animate={{ opacity: 1, y: 0, scale: 1 }}
                                className={cn(
                                    "flex items-start gap-4 max-w-[85%]",
                                    message.role === "user" ? "ml-auto flex-row-reverse" : ""
                                )}
                            >
                                <div className={cn(
                                    "w-8 h-8 rounded-lg flex items-center justify-center shrink-0",
                                    message.role === "user" ? "bg-primary" : "bg-muted border border-border"
                                )}>
                                    {message.role === "user" ? <User className="w-5 h-5 text-primary-foreground" /> : <Bot className="w-5 h-5 text-primary" />}
                                </div>
                                <div className={cn(
                                    "p-4 rounded-2xl text-sm leading-relaxed",
                                    message.role === "user"
                                        ? "bg-primary text-primary-foreground rounded-tr-none"
                                        : "bg-muted/50 border border-border/50 rounded-tl-none"
                                )}>
                                    {message.content}
                                </div>
                            </motion.div>
                        ))}
                    </AnimatePresence>
                    {isLoading && (
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            className="flex items-center gap-2 text-muted-foreground text-xs"
                        >
                            <Loader2 className="w-3 h-3 animate-spin" />
                            Engine is fusing intelligence...
                        </motion.div>
                    )}
                </div>

                <div className="p-4 border-t border-border/50 bg-background/50">
                    <form
                        onSubmit={(e) => { e.preventDefault(); handleSend(); }}
                        className="flex gap-2"
                    >
                        <Input
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder="Ask anything about your IoT data..."
                            className="flex-1 bg-background"
                        />
                        <Button type="submit" disabled={isLoading}>
                            {isLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
                        </Button>
                    </form>
                    <p className="text-[10px] text-center text-muted-foreground mt-2">
                        AI-generated insights may require verification. Connected to IOT-RAG v1.0.
                    </p>
                </div>
            </Card>
        </div>
    );
}
