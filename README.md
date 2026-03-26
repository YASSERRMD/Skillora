<div align="center">
  <img src="assets/logo.png" alt="Skillora Logo" width="250" />
  <h1>Skillora: The AI-Driven Skill Exchange & Barter Platform</h1>
  <p>Learn, Connect, and Grow — by exchanging what you know for what you want to learn.</p>
</div>

<br/>

> **The Vision:** This project was inspired by and blossomed from the [On-Demand Skill Exchange Platform](https://www.linkedin.com/pulse/on-demand-skill-exchange-platform-mohamed-yasser-ednqc/) article authored by **Mohamed Yasser**. In today’s fast-paced world, individuals are continuously seeking new ways to enhance their skills without financial barriers. Skillora is the technical manifestation of that vision: a modern Web3-inspired, AI-driven barter economy strictly for the exchange of knowledge.

---

## 🚀 Overview

**Skillora** is a revolutionary web-based marketplace that replaces traditional financial transactions with **Skill Credits**. Instead of paying money for services, users trade their expertise. 

A web developer might offer coding mentorship in exchange for marketing advice; a yogi could provide virtual classes in return for guitar lessons. Skillora automates this exchange securely using a **double-entry ledger system**, while leveraging artificial intelligence to **appraise skills**, generate **vector embeddings**, and **semantically match** user demands across the platform.

## ✨ Core Features

- 🧠 **AI-Powered Skill Appraisal:** Users don't just self-report skills. The platform integrates an intelligent LLM Oracle that appraises skill proposals, estimating a relative "Credit Value" and confirming proficiency before committing it to a user's verified portfolio.
- 💳 **Double-Entry Barter Ledger:** Every knowledge exchange operates on a strict, ACID-compliant double-entry accounting ledger ensuring that credits are perfectly conserved, created, or transferred when barters are proposed, accepted, and completed.
- 🎯 **Semantic Vector Matching Engine:** Built on Postgres `pgvector`, the matching engine converts requested topics into semantic embeddings via AI, searching millions of skill nodes to find the best possible human matches for what you want to learn.
- 🛡 **State-of-the-Art Architecture:** The backbone of Skillora is built with a highly secure, parallelized **Golang API** and a stunning, responsive **Next.js** frontend.

## 🛠️ Technology Stack

Skillora was designed for immense scale, leveraging industry-standard modern tooling:

### Backend (Core Eng)
* **Golang 1.25** & **Gin Framework** for ultra-fast API routing and dependency injection.
* **PostgreSQL + pgxpool** for native ACID transactions.
* **pgvector** for High-Dimensional Nearest Neighbor (HNSW) vector search.
* **Redis** for OAuth caching, rate-limiting, and ephemeral sessions.
* **Server-Sent Events (SSE)** for real-time notification dispatching.

### Frontend
* **Next.js 15 (React 19)** via App Router.
* **Zustand** for lightweight global state (Authentication & Theming).
* **TanStack React Query** for fetching, caching, and optimistic mutations.
* **Tailwind CSS & Shadcn UI** for gorgeous, accessible component design.
* **Lucide Icons & Sonner** for rich micro-interactions and toast flows.

## 🐳 Quick Start (Docker)

Skillora is completely Dockerized, running its infrastructure and applications entirely isolated from your host machine.

### Prerequisites
* [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running.
* GitHub OAuth keys (if testing Google login, create credentials in GCP Console).

```bash
# 1. Clone the repository
git clone https://github.com/YASSERRMD/Skillora.git
cd Skillora

# 2. Setup environment variables
cp .env.example .env

# 3. Fire it all up
docker-compose up -d --build
```

**Access Points:**
* 🎨 **Frontend Application:** `http://localhost:3000`
* 🔌 **Backend API Router:** `http://localhost:8080/api/v1`

*(PostgreSQL and Redis are hidden completely inside the internal Docker bridge network for maximum security).*

## 💡 How It Works

1. **Verify Your Skills**: Head to your Dashboard and use the "Add Skill" wizard. Write what you know. Skillora's AI evaluates your write-up, determines a fair market credit value, and issues it a proficiency badge.
2. **Find a Match**: Enter the Marketplace. Type exactly what you need (e.g., *"I need someone to help me migrate to React 19"*). The semantic embedding engine calculates cosine similarities across all users.
3. **Propose a Barter**: Propose an exchange. You offer your credits (or your own verified skills) in exchange for theirs.
4. **Learn and Grow**: Once accepted, the transaction goes into "Active" status. Upon completion, the ledger processes the credit movement securely using double-entry logic.

## 🤝 Contributing

We welcome contributions to Skillora. If you have an idea for a feature, notice a bug, or want to enhance the matching algorithms, please open an Issue or submit a Pull Request following conventional commits.

## 📄 License

Skillora is open-source software licensed under the [MIT License](LICENSE). 
