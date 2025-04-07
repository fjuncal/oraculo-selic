# 📡 Messaging Oracle

A distributed system simulation that processes and validates messages exchanged between services via messaging queues like **ActiveMQ** and **IBM MQ**. The application supports multi-database integration, message tracking, and a modern UI for visualizing message flows and statuses.

---

## 🚀 Features

- 🔄 Message sending and receiving via ActiveMQ or IBM MQ
- 🧠 Data validation and correlation
- 🗃️ Multiple PostgreSQL database support
- 📊 Message status tracking with custom correlation ID
- 🖥️ Frontend for real-time message visualization
- ⚙️ Easily switch between environments (local/dev/cert)

---

## 🧱 Tech Stack

- Golang
- PostgreSQL (3 instances)
- ActiveMQ / IBM MQ
- React (Next.js)
- REST API
- UUID correlation for messages

---

## ⚙️ Environment Variables

Before running the project, make sure to configure the following variables:

- DATABASE_URL_1=<your_primary_database_url>
- DATABASE_URL_2=<your_secondary_database_url>
- DATABASE_URL_3=<your_tertiary_database_url>


- MESSAGING_TYPE=activemq             # or ibmmq
- QUEUE_URL=localhost:61616           # ActiveMQ example URL
You can define these variables in a .env file or via command line when running the project.

📦 Running the Project
# Clone the repository
git clone https://github.com/fjuncal/oraculo-selic.git
cd oraculo-selic

# Set your environment variables or export them in terminal

# Run the backend
go run main.go
Make sure all PostgreSQL databases are running and accessible.

🌐 Frontend
This project comes with a modern frontend (Next.js) to visualize message flows and statuses.

cd frontend
npm install
npm run dev
Access at: http://localhost:3000

📄 License
MIT License

👤 Author


Fellipe Juncal


LinkedIn • GitHub