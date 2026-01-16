# Game Core - Backend Server

Backend del juego Pong construido en Go. Maneja toda la lógica del juego, física, detección de colisiones y comunicación en tiempo real con los clientes vía WebSocket.

## Stack Tecnológico

- **Lenguaje**: Go 1.21+
- **Framework Web**: Gorilla Mux
- **WebSockets**: Gorilla WebSocket
- **Gestión de IDs**: Google UUID

## Arquitectura

El servidor implementa un game loop deterministico a 60 TPS (Ticks Per Second) que:
- Procesa inputs de jugadores
- Actualiza física y posiciones
- Detecta colisiones
- Gestiona power-ups y obstáculos
- Sincroniza estado con clientes

### Estructura del Proyecto

## Características

### Game Loop
- Actualización fija a 60 TPS
- Simulación determinística
- Sincronización precisa de estado

### Física
- Colisiones AABB (Axis-Aligned Bounding Box)
- Respuesta dinámica según punto de impacto
- Incremento progresivo de velocidad
- Efectos de spin en la bola

### Networking
- Comunicación bidireccional vía WebSocket
- Protocol buffers o JSON para mensajes
- Reconciliación de estado cliente-servidor
- Manejo robusto de desconexiones




## Configuración

### Variables de Entorno

```bash
# Puerto del servidor
GAME_PORT=8080

# Configuración del game loop
GAME_TICK_RATE=60
GAME_FIELD_WIDTH=800
GAME_FIELD_HEIGHT=600

# Gestión de salas
GAME_MAX_ROOMS=50
GAME_ROOM_TIMEOUT=300

# Logging
GAME_LOG_LEVEL=info

# CORS (para desarrollo)
GAME_ALLOW_ORIGINS=http://localhost:3000
```

