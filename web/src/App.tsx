import React from "react";
import { Routes, Route } from "react-router-dom";
import { AnimatePresence } from "framer-motion";

// Layout components
import Layout from "./components/layout/Layout";

// Page components
import Dashboard from "./pages/Dashboard";
import Books from "./pages/Books";
import BookDetails from "./pages/BookDetails";
import Members from "./pages/Members";
import MemberDetails from "./pages/MemberDetails";
import Checkout from "./pages/Checkout";
import Returns from "./pages/Returns";
import Settings from "./pages/Settings";
import NotFound from "./pages/NotFound";
import RoomEntry from "./pages/RoomEntry";
import RoomOverview from "./pages/RoomOverview";

function App() {
  return (
    <AnimatePresence mode="wait">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="books" element={<Books />} />
          <Route path="books/:id" element={<BookDetails />} />
          <Route path="members" element={<Members />} />
          <Route path="members/:id" element={<MemberDetails />} />
          <Route path="checkout" element={<Checkout />} />
          <Route path="returns" element={<Returns />} />
          <Route path="/room-entry" element={<RoomEntry />} />
          <Route path="/room-overview" element={<RoomOverview />} />
          <Route path="settings" element={<Settings />} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </AnimatePresence>
  );
}

export default App;
