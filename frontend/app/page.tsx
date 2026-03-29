"use client";

import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { motion } from "framer-motion";
import { ArrowRight, Sparkles, Zap, Shield, Globe } from "lucide-react";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-zinc-50 dark:bg-black overflow-hidden font-sans">
      {/* Navbar Minimal */}
      <nav className="flex items-center justify-between p-6 w-full max-w-7xl mx-auto z-10">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-indigo-600 dark:bg-indigo-500 flex items-center justify-center text-white font-bold shadow-lg shadow-indigo-500/20">
            S
          </div>
          <span className="text-xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">Skillora</span>
        </div>
        <div className="flex items-center gap-4 relative z-50">
          <Link href="/login" className={cn(buttonVariants({ variant: "ghost" }), "hidden sm:inline-flex text-zinc-600 dark:text-zinc-300 hover:text-zinc-900 dark:hover:text-white transition-colors")}>
            Sign In
          </Link>
          <Link href="/login" className={cn(buttonVariants({ variant: "default" }), "bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600 text-white rounded-full shadow-md shadow-indigo-500/20 transition-all font-medium")}>
            Get Started
          </Link>
        </div>
      </nav>

      {/* Hero Section */}
      <main className="flex-1 flex flex-col items-center justify-center relative px-6 z-10 -mt-20">
        
        {/* Decorative background glows */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[600px] opacity-30 dark:opacity-20 pointer-events-none blur-[120px] rounded-full bg-gradient-to-tr from-indigo-500 via-emerald-400 to-indigo-600 mix-blend-multiply dark:mix-blend-screen" />
        
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, ease: "easeOut" }}
          className="max-w-4xl text-center space-y-8 relative z-10"
        >
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 text-sm font-medium border border-indigo-200 dark:border-indigo-800/50 backdrop-blur-sm mb-4 mx-auto">
            <Sparkles className="w-4 h-4" />
            <span>The AI-Driven Skill Exchange Platform</span>
          </div>
          
          <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight text-zinc-900 dark:text-white leading-tight md:leading-[1.1]">
            Trade Real Skills,
            <span className="block mt-2 pb-2 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-indigo-500 to-emerald-500 dark:from-indigo-400 dark:via-indigo-300 dark:to-emerald-400">Not Just Talk.</span>
          </h1>
          
          <p className="max-w-2xl mx-auto text-lg md:text-xl text-zinc-600 dark:text-zinc-400 leading-relaxed">
            Skillora is the ultimate marketplace where your knowledge equals currency. Connect, exchange expertise, and let our AI appraise your skills with double-entry ledger precision.
          </p>
          
          <motion.div 
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3, duration: 0.5 }}
            className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4 relative z-50"
          >
            <Link href="/login" className={cn(buttonVariants({ size: "lg" }), "w-full sm:w-auto text-base h-14 px-8 rounded-full bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600 text-white shadow-lg shadow-indigo-500/25 transition-all hover:scale-105 active:scale-95 group font-medium sm:flex items-center justify-center")}>
              Enter the Marketplace
              <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
            </Link>
            <Link href="/marketplace" className={cn(buttonVariants({ size: "lg", variant: "outline" }), "w-full sm:w-auto text-base h-14 px-8 rounded-full border-zinc-200 dark:border-zinc-800 bg-white/50 dark:bg-zinc-900/50 backdrop-blur-sm hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-900 dark:text-white transition-all font-medium sm:inline-flex")}>
              Explore Skills
            </Link>
          </motion.div>
        </motion.div>
        
        {/* Feature Grid Mini */}
        <motion.div 
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.8 }}
          className="grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-5xl mx-auto mt-24 relative z-10 w-full"
        >
          {[
            { icon: Zap, title: "AI Skill Appraiser", desc: "Automated vector-based valuation of your expertise." },
            { icon: Shield, title: "ACID Escrow", desc: "True double-entry ledger ensures safe skill exchanges." },
            { icon: Globe, title: "Global Marketplace", desc: "Match seamlessly using semantic vector search." },
          ].map((feature, i) => (
            <div key={i} className="p-6 rounded-3xl bg-white/60 dark:bg-zinc-900/40 backdrop-blur-md border border-zinc-200 dark:border-zinc-800 shadow-sm hover:shadow-md transition-all group hover:-translate-y-1">
              <div className="w-12 h-12 rounded-xl bg-indigo-100 dark:bg-indigo-900/50 flex items-center justify-center text-indigo-600 dark:text-indigo-400 mb-4 group-hover:scale-110 transition-transform">
                <feature.icon className="w-6 h-6" />
              </div>
              <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100 mb-2">{feature.title}</h3>
              <p className="text-zinc-600 dark:text-zinc-400 text-sm leading-relaxed">{feature.desc}</p>
            </div>
          ))}
        </motion.div>
      </main>
    </div>
  );
}
