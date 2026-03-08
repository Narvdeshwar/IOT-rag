import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MainLayout } from "./components/layout/main-layout";
import { DashboardView } from "./modules/dashboard/dashboard-view";
import { RagInterface } from "./modules/rag/rag-interface";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <MainLayout>
          <Routes>
            <Route path="/" element={<DashboardView />} />
            <Route path="/rag" element={<RagInterface />} />
            {/* Fallback routes */}
            <Route path="*" element={<DashboardView />} />
          </Routes>
        </MainLayout>
      </Router>
    </QueryClientProvider>
  );
}

export default App;
