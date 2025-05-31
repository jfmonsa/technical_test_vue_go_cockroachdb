# 📈 Stock Recommender App

Una aplicación simple construida con Go, Vue y CockroachDB que permite obtener, almacenar y visualizar recomendaciones de acciones de bolsa. Implementa un sistema de recomendación basado en reglas para sugerir las mejores acciones en las que invertir hoy.

## 🚀 Características

- ETL automático desde una API externa con paginación.
- Backend en Go con endpoints RESTful.
- Base de datos relacional (CockroachDB).
- Frontend en Vue.js para búsqueda, filtrado y visualización.
- Algoritmo de recomendación de acciones basado en potencial de ganancia, acción del bróker, cambio de rating y actualidad.

---

## 🛠️ Tecnologías

- **Backend:** Go + Chi + CockroachDB
- **Frontend:** Vue.js
- **Infraestructura:** Terraform (AWS + Lambdas)
- **Base de datos:** CockroachDB

---

## 🧪 Requisitos

- Go ≥ 1.19
- Docker (para correr CockroachDB localmente)
- curl o Postman para probar la API
- Node.js ≥ 16 (para el frontend, si aplica)

---
