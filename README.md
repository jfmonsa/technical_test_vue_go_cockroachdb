# 📈 Stock Recommender App

Una aplicación simple construida con Go, Vue y CockroachDB que permite obtener, almacenar y visualizar recomendaciones de acciones de bolsa. Implementa un sistema de recomendación basado en reglas para sugerir las mejores acciones en las que invertir hoy.

<div style="display: flex; gap: 16px; flex-wrap: wrap;">
    <img src="./docs/images/Recommendations.png" alt="Screenshot of recommendations view" style="max-width: 48%; height: auto; border: 1px solid #ddd; border-radius: 4px;">
    <img src="./docs/images/homeView.png" alt="Screenshot of home view" style="max-width: 48%; height: auto; border: 1px solid #ddd; border-radius: 4px;">
</div>

---

## 🚀 Características

- ETL automático desde una API externa con paginación.
- Backend en Go con endpoints RESTful.
- Base de datos relacional (CockroachDB).
- Frontend en Vue.js para búsqueda, filtrado y visualización.
- Algoritmo de recomendación de acciones basado en potencial de ganancia, acción del bróker, cambio de rating y actualidad.

---

## 🛠️ Tecnologías

- **Backend:** Go + Chi + CockroachDB
- **Frontend:** Vue.js + Pinia + Tailwind CSS
- **Base de datos:** CockroachDB

---

## 🧪 Requisitos

- Go ≥ 1.19
- Docker (para correr CockroachDB localmente)
- curl o Postman para probar la API
- Node.js ≥ 16 (para el frontend, si aplica)

---

## Requerimientos: Como fueron resueltos y sus retos

### ⚙️ Parte 1: Implementación del Proceso ETL (_Connect to the API and store the data_)

Para resolver el requerimiento de descargar la información de stocks del api externa proporcionada que implica un proceso de **_Extraer, Transformar y Cargar (ETL)_** datos desde una API externa hacia una base de datos en la nube (CockroachDB), se desarrolló un proceso automático en Go que ejecuta los siguientes pasos:

#### **_✅ Extracción_**

Se consume una API que retorna información de acciones de manera paginada. Para cada página, se realiza una solicitud `GET` autenticada. El resultado se deserializa en estructuras tipo `APIRawItem`.

- **_Reto_**: La API no siempre responde con datos consistentes. En ocasiones, ciertos campos vienen vacíos o mal formateados. Para solucionarlo, se agregaron validaciones robustas antes de transformar los datos.

#### **_🔄 Transformación_**

Cada item recibido se transforma en un objeto Stock, validando campos obligatorios (ticker, time) y convirtiendo strings numéricos con símbolos como `$` y , a `float64`. También se descartan registros con valores inválidos en campos como `rating_from` o `rating_to`, que podrían sesgar análisis posteriores.

- **_Reto_**: Algunos registros venían con ratings vacíos o inconsistentes. Decidí ignorarlos durante la transformación y registrar estos fallos en una tabla aparte (failed_items), para poder analizarlos sin afectar la calidad del dataset principal.

#### **_💾 Carga_**

Una vez transformado el dato, se inserta en la tabla stocks. Se utiliza la estrategia ON CONFLICT DO NOTHING para evitar duplicados basados en la clave (`ticker`, `time`).

- **_Reto_**: A pesar de las validaciones, podrían presentarse errores al insertar (por ejemplo, por campos nulos no controlados). En ese caso, también se guarda el item fallido junto con el mensaje de error en la tabla `failed_items`.

#### **_🧾 Registro de errores_**

Para asegurar la trazabilidad, todo el proceso genera logs en archivos con timestamps, ubicados en una carpeta `logs/`. Además, se implementó una tabla en la base de datos para guardar los registros que fallaron en las fases de transformación o carga, con sus respectivos mensajes de error y la fase en la que ocurrió el problema.

---

### 🛠️ Parte 2.1: API RESTful

Una vez los datos fueron almacenados correctamente en la base de datos, el siguiente paso fue crear una **API RESTful** para exponer esta información a través de endpoints organizados, seguros y bien estructurados.

#### **_🔌 Endpoints Implementados_**

Se definieron 3 endpoints que permiten consultar información de acciones almacenadas:

##### 🔍 `GET /stocks`

```shell
curl "http://localhost:8080/stocks?search=pharma&sort_by=target_to&order=desc&page=1&limit=10"
```

response:

```json
{
  "items": [
    {
      "ticker": "REGN",
      "company": "Regeneron Pharmaceuticals",
      "brokerage": "BMO Capital Markets",
      "action": "target lowered by",
      "rating_from": "Outperform",
      "rating_to": "Outperform",
      "target_from": 800,
      "target_to": 600,
      "time": "2025-06-03T00:30:05.843101Z"
    }
    // etc...
  ],
  "limit": 10,
  "page": 1,
  "total": 83,
  "totalPages": 9
}
```

##### ⭐ `GET /stocks/recommendations`

Devuelve las mejores acciones para invertir hoy según el algoritmo.

```shell
curl http://localhost:8080/stocks/recommendations
```

response

```json
[
  {
    "ticker": "SPXC",
    "company": "SPX Technologies",
    "brokerage": "UBS Group",
    "action": "upgraded by",
    "rating_from": "Neutral",
    "rating_to": "Buy",
    "target_from": 160,
    "target_to": 182,
    "time": "2025-06-03T00:30:06.138894Z",
    "Score": 14.4375
  },
  {
    "ticker": "PATH",
    "company": "UiPath",
    "brokerage": "Morgan Stanley",
    "action": "target raised by",
    "rating_from": "Equal Weight",
    "rating_to": "Equal Weight",
    "target_from": 12,
    "target_to": 15,
    "time": "2025-06-03T00:30:06.134828Z",
    "Score": 8.25
  }
  // etc...
]
```

##### 📌 `GET /stocks/:ticker`

Devuelve detalles de una acción específica.

```shell
curl "http://localhost:8080/stocks/AKBA"
```

response:

```json
{
  "ticker": "AKBA",
  "company": "Akebia Therapeutics",
  "brokerage": "Leerink Partners",
  "action": "initiated by",
  "rating_from": "Outperform",
  "rating_to": "Outperform",
  "target_from": 7,
  "target_to": 7,
  "time": "2025-04-29T00:30:06.253903Z"
}
```

#### 🧱 Organización: Handler, Service y Repository

Se siguió una arquitectura de 3 capas:

```
handler.go  ⟶  service.go  ⟶  repository_interface.go  ⟶  repository_cockroachdb.go
```

- handler.go: expone los endpoints HTTP y maneja parámetros de la petición.
- service.go: contiene la lógica de negocio (ej. cálculo del puntaje).
- repository_interface.go: define las operaciones esperadas sobre los datos.
- repository_cockroachdb.go: implementación de las operaciones usando SQL con pgx.

🧩 Ventaja: esta separación facilita testeo, mantenimiento y futura migración de la capa de datos.

---

### 🛠️ Parte 2.2: Frontend

El frontend fue desarrollado en Vue 3 con TypeScript, siguiendo buenas prácticas como el uso de la API de composición, componentes modulares y separación clara entre vistas, componentes y lógica de estado.

Se utilizó Pinia para el manejo del estado global, especialmente para:

- Cachear peticiones a la API (evitando llamadas innecesarias).
- Almacenar filtros de búsqueda como el score mínimo y el límite de resultados.
- Gestionar tablas y estadísticas derivadas (top 3 recomendaciones, promedio, distribución, etc.).

Gracias a Pinia y la reactividad de Vue, el sistema actualiza automáticamente los datos al modificar los filtros, mejorando la experiencia del usuario.

- Además, se aplicaron buenas prácticas como:
- Organización por carpetas (components, stores, models, views).
- Componentes reutilizables y tipados.
- Uso de Tailwind CSS para un diseño responsivo y limpio.
- Integración eficiente con el backend mediante variables de entorno y peticiones dinámicas.

<!--
- posibles mejoras si la app escala
TODO: expliación de como funciona el algoritmo y posibles mejores
- Arquitectura de eventos, como se podría escalar este sistema para trabajar con datos "vivos" cron-jobs, lambdas, microservicios, colas


TODO: despliegue con terraform
TODO: tests e2e y algunos tests


TODO: explicar como se maneja todo el SDLC (testing, ci, deploy etc...), herramietnas de desarrollo: linters, format, testing, git-hooks, pipelines
---

TODO: Get started variable de entorno etc.
 -->

<!--
TODO: incluir imágenes de la app funcionando
TODO: script o docker-compose para levantar toda el app por primera vez

-->

---

### 🧠 Parte 3: Algoritmo de recomendación (Cálculo del Score)

Se implementó un **algoritmo de scoring** que evalúa cada registro de la base de datos y le asigna un puntaje numérico basado en múltiples factores relevantes para la toma de decisiones de inversión.

El objetivo es priorizar aquellas acciones con mayor potencial y menor riesgo, según la información disponible.

#### **¿Cómo funciona el algoritmo de score?**

El algoritmo toma en cuenta los siguientes aspectos para cada acción:

1. **Potencial de ganancia (`target_from` vs `target_to`):**

   - Calcula el porcentaje de ganancia potencial entre el precio objetivo anterior y el nuevo.
   - Si el potencial es positivo, suma más puntos (más peso si la subida es grande).
   - Si es negativo, penaliza el score (pero con menor peso).

2. **Acción recomendada por el bróker (`action`):**

   - Acciones como "upgraded", "initiated", "target raised" suman puntos.
   - Acciones como "downgraded" o "target lowered" restan puntos.
   - Se utiliza una tabla de palabras clave y su peso relativo.

3. **Cambio de rating (`rating_from` y `rating_to`):**

   - Se asigna un valor numérico a cada rating (por ejemplo, "Strong Buy" = 10, "Sell" = 1).
   - Si el rating mejora, se suman puntos proporcionales al cambio.
   - Si el rating empeora, se penaliza el score.

4. **Actualidad de la recomendación (`time`):**
   - Las recomendaciones más recientes reciben un bono extra.
   - Si la recomendación es del mismo día, suma más puntos; si es de la última semana, suma menos.

#### **Ejemplo de fórmula simplificada:**

```go
score := 0.0
// 1. Potencial de ganancia
if s.TargetFrom > 0 {
    potential := ((s.TargetTo - s.TargetFrom) / s.TargetFrom) * 100
    if potential > 0 {
        score += potential / 4
    } else {
        score += potential / 10
    }
}
// 2. Acción recomendada
// ... (ver tabla de acciones y pesos)
// 3. Cambio de rating
// ... (ver tabla de ratings y pesos)
// 4. Actualidad
// ... (bono por fecha reciente)
```

#### **Ventajas del enfoque:**

- Permite comparar acciones de diferentes sectores y precios en una misma escala.
- Es flexible: se pueden ajustar los pesos según la experiencia o feedback de usuarios.
- Facilita la visualización y el ranking de las mejores oportunidades de inversión.

#### **Notas adicionales:**

- El score se calcula automáticamente durante el proceso ETL y se almacena en la base de datos para eficiencia.
- El endpoint `/stocks/recommendations` filtra por un score mínimo y ordena por el score más alto, devolviendo las mejores opciones para invertir hoy.

---
