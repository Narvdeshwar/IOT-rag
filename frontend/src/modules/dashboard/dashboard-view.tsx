import {
    Activity,
    Zap,
    ShieldAlert,
    Waves,
    TrendingUp,
    AlertCircle
} from "lucide-react";
import { Button } from "../../components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../../components/ui/Card";
import { motion } from "framer-motion";
import { cn } from "../../lib/utils";

const stats = [
    { label: "Active Nodes", value: "1,284", icon: Zap, color: "text-yellow-500", bg: "bg-yellow-500/10" },
    { label: "Throughput", value: "48.2 GB/s", icon: Activity, color: "text-blue-500", bg: "bg-blue-500/10" },
    { label: "Threat Index", value: "Low", valueColor: "text-green-500", icon: ShieldAlert, color: "text-green-500", bg: "bg-green-500/10" },
    { label: "Fusion Accuracy", value: "99.4%", icon: TrendingUp, color: "text-purple-500", bg: "bg-purple-500/10" },
];

export function DashboardView() {
    return (
        <div className="space-y-8">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
                        Operational Dashboard
                    </h2>
                    <p className="text-muted-foreground mt-1">Real-time telemetry and predictive analysis overview.</p>
                </div>
                <div className="flex gap-3">
                    <div className="px-4 py-2 bg-primary/10 border border-primary/20 rounded-full flex items-center gap-2">
                        <div className="w-2 h-2 rounded-full bg-primary animate-ping" />
                        <span className="text-xs font-bold text-primary">LIVE MONITORING</span>
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {stats.map((stat, i) => (
                    <motion.div
                        key={stat.label}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: i * 0.1 }}
                    >
                        <Card className="hover:border-primary/50 transition-colors border-border/50 bg-card/50">
                            <CardContent className="p-6">
                                <div className="flex items-center justify-between mb-4">
                                    <div className={cn("p-2 rounded-lg", stat.bg)}>
                                        <stat.icon className={cn("w-5 h-5", stat.color)} />
                                    </div>
                                    <span className="text-xs font-medium text-muted-foreground">Updated Now</span>
                                </div>
                                <h4 className="text-sm font-medium text-muted-foreground">{stat.label}</h4>
                                <div className="flex items-baseline gap-2 mt-1">
                                    <span className={cn("text-2xl font-bold", stat.valueColor || "text-foreground")}>
                                        {stat.value}
                                    </span>
                                </div>
                            </CardContent>
                        </Card>
                    </motion.div>
                ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <Card className="lg:col-span-2 border-border/50 bg-card/50 backdrop-blur-sm">
                    <CardHeader>
                        <CardTitle>System Telemetry</CardTitle>
                        <CardDescription>Visualizing packet distribution across nodes.</CardDescription>
                    </CardHeader>
                    <CardContent className="h-[300px] flex items-center justify-center border-t border-border/50">
                        <div className="flex flex-col items-center gap-4 text-center">
                            <div className="relative w-48 h-48">
                                <motion.div
                                    animate={{ rotate: 360 }}
                                    transition={{ duration: 10, repeat: Infinity, ease: "linear" }}
                                    className="absolute inset-0 border-2 border-dashed border-primary/30 rounded-full"
                                />
                                <motion.div
                                    animate={{ rotate: -360 }}
                                    transition={{ duration: 15, repeat: Infinity, ease: "linear" }}
                                    className="absolute inset-4 border-2 border-dashed border-blue-500/20 rounded-full"
                                />
                                <div className="absolute inset-0 flex items-center justify-center">
                                    <Waves className="w-12 h-12 text-primary opacity-50" />
                                </div>
                            </div>
                            <p className="text-sm text-muted-foreground max-w-[200px]">
                                Real-time data visualization module initializing...
                            </p>
                        </div>
                    </CardContent>
                </Card>

                <Card className="border-border/50 bg-card/50 backdrop-blur-sm">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <AlertCircle className="w-5 h-5 text-yellow-500" />
                            Active Alerts
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        {[1, 2, 3].map((_, i) => (
                            <div key={i} className="flex gap-4 p-3 rounded-lg bg-accent/30 border border-border/50">
                                <div className="w-1 h-10 bg-yellow-500 rounded-full shrink-0" />
                                <div className="space-y-1 overflow-hidden">
                                    <p className="text-sm font-medium truncate">Anomaly detected in Node-7{i}</p>
                                    <p className="text-xs text-muted-foreground">37 seconds ago • Zone C</p>
                                </div>
                            </div>
                        ))}
                        <Button variant="outline" className="w-full text-xs" size="sm">
                            View All Alerts
                        </Button>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}

